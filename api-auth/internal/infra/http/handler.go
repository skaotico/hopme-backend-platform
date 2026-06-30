package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"c4-auth/internal/domain/model"
	"c4-auth/internal/domain/port"
	"c4-auth/internal/infra/http/middleware"
	"c4-auth/internal/infra/http/response"
)

// RegisterHandler expone el endpoint HTTP para el registro de nuevos usuarios
func RegisterHandler(registerUC port.RegisterUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req response.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Cuerpo de petición JSON inválido", err.Error())
			return
		}

		user, err := registerUC.Execute(r.Context(), req.Username, req.Email, req.Password)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidEmail):
				response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidEmail, err.Error(), "")
			case errors.Is(err, model.ErrInvalidUsername):
				response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidUsername, err.Error(), "")
			case errors.Is(err, model.ErrInvalidPassword):
				response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidPassword, err.Error(), "")
			case errors.Is(err, model.ErrUserAlreadyExists):
				response.Failure(w, http.StatusConflict, response.ErrCodeUserExists, err.Error(), "")
			default:
				response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error interno al procesar el registro", err.Error())
			}
			return
		}

		response.Success(w, http.StatusCreated, user)
	}
}

// LoginHandler expone el endpoint HTTP para el inicio de sesión
func LoginHandler(loginUC port.LoginUseCase, secureCookie bool, rtExpiration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req response.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Cuerpo de petición JSON inválido", err.Error())
			return
		}

		ipAddress := r.Header.Get("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = r.RemoteAddr
		}
		deviceInfo := r.Header.Get("User-Agent")

		token, refreshToken, err := loginUC.Execute(r.Context(), req.Email, req.Password, ipAddress, deviceInfo)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidCredentials):
				response.Failure(w, http.StatusUnauthorized, response.ErrCodeInvalidCredentials, err.Error(), "")
			case errors.Is(err, model.ErrUnauthorized):
				response.Failure(w, http.StatusForbidden, response.ErrCodeUserDisabled, "La cuenta de usuario se encuentra desactivada", "")
			default:
				response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error interno al procesar el ingreso", err.Error())
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/api/v1/auth",
			Expires:  time.Now().Add(rtExpiration),
			HttpOnly: true,
			Secure:   secureCookie,
			SameSite: http.SameSiteStrictMode,
		})

		response.Success(w, http.StatusOK, response.TokenResponse{
			Token: token,
		})
	}
}

// RefreshHandler expone el endpoint HTTP para rotar tokens
func RefreshHandler(refreshUC port.RefreshTokenUseCase, secureCookie bool, rtExpiration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			response.Failure(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "Refresh token ausente", "No se encontró la cookie refresh_token")
			return
		}
		refreshTokenStr := cookie.Value

		ipAddress := r.Header.Get("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = r.RemoteAddr
		}
		deviceInfo := r.Header.Get("User-Agent")

		token, newRefreshToken, err := refreshUC.Execute(r.Context(), refreshTokenStr, ipAddress, deviceInfo)
		if err != nil {
			response.Failure(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "Token de refresco inválido o expirado", err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefreshToken,
			Path:     "/api/v1/auth",
			Expires:  time.Now().Add(rtExpiration),
			HttpOnly: true,
			Secure:   secureCookie,
			SameSite: http.SameSiteStrictMode,
		})

		response.Success(w, http.StatusOK, response.TokenResponse{
			Token: token,
		})
	}
}

// LogoutHandler expone el endpoint HTTP para cerrar sesión
func LogoutHandler(logoutUC port.LogoutUseCase, secureCookie bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshTokenStr string
		if cookie, err := r.Cookie("refresh_token"); err == nil {
			refreshTokenStr = cookie.Value
		}

		// Recuperar los claims del contexto inyectados por el middleware JWT
		claims, ok := middleware.GetUserClaims(r.Context())
		if !ok {
			response.Failure(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "Sesión de usuario no válida o inactiva", "")
			return
		}

		// Calcular el tiempo restante de vida del JWT en segundos
		expTime := claims.ExpiresAt.Time
		ttl := int(time.Until(expTime).Seconds())

		err := logoutUC.Execute(r.Context(), claims.ID, ttl, refreshTokenStr)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error interno al cerrar sesión", err.Error())
			return
		}

		// Invalidar la cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/api/v1/auth",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   secureCookie,
			SameSite: http.SameSiteStrictMode,
		})

		response.Success(w, http.StatusOK, map[string]string{"message": "Sesión cerrada correctamente"})
	}
}


// MeHandler expone el endpoint HTTP para obtener los datos del usuario firmante del JWT actual
func MeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Recuperar los claims del contexto inyectados por el middleware JWT
		claims, ok := middleware.GetUserClaims(r.Context())
		if !ok {
			response.Failure(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "Sesión de usuario no válida o inactiva", "")
			return
		}

		// Responder con los claims básicos del JWT
		response.Success(w, http.StatusOK, map[string]interface{}{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"email":    claims.Email,
			"roles":    claims.Roles,
		})
	}
}

// HealthHandler expone un endpoint público y ligero para comprobaciones de estado de Docker / Portainer
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.3",
		})
	}
}
