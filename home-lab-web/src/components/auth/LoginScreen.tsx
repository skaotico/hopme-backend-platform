import React, { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { Loader2 } from 'lucide-react';
import './Auth.css';

interface LoginScreenProps {
  onGoToRegister: () => void;
}

export const LoginScreen: React.FC<LoginScreenProps> = ({ onGoToRegister }) => {
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsSubmitting(true);
    const result = await login(email, password);
    if (!result.success) {
      setError(result.error || 'Error al iniciar sesión');
      setIsSubmitting(false);
    }
  };

  return (
    <div className="auth-container animate-fade-in">
      <div className="card auth-card glass-panel">
        <h1 className="auth-title">EcoParque</h1>
        <p className="auth-subtitle">Inicia sesión en tu cuenta</p>
        
        <form className="auth-form" onSubmit={handleSubmit}>
          <div>
            <label htmlFor="email">Correo Electrónico</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="tu@email.com"
              required
              disabled={isSubmitting}
            />
          </div>
          <div>
            <label htmlFor="password">Contraseña</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              disabled={isSubmitting}
            />
          </div>
          
          {error && <p className="error-text">{error}</p>}
          
          <button 
            type="submit" 
            className="btn btn-primary w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? <Loader2 className="animate-spin" size={20} /> : 'Iniciar Sesión'}
          </button>
        </form>
        
        <div className="auth-footer">
          ¿No tienes una cuenta?
          <span className="auth-link" onClick={onGoToRegister}>Regístrate aquí</span>
        </div>
      </div>
    </div>
  );
};
