# EcoParque App

Bienvenido al proyecto **EcoParque**. Esta aplicación móvil está construida con **React Native y Expo** y forma parte de un ecosistema más amplio orientado a microservicios desarrollados en **Golang**. 

A continuación encontrarás los detalles de cómo está construida la app y los pasos exactos para levantar todo el ecosistema.

---

## 🏗 Arquitectura del Sistema

El ecosistema se divide en una aplicación cliente (Frontend) y múltiples APIs independientes (Backend), que se comunican a través de peticiones HTTP REST.

### 1. Frontend (Mobile App)
- **Tecnología**: React Native gestionado con Expo.
- **Patrón de Navegación**: Enrutamiento gestionado por estado local en `App.tsx` (desacoplado por flujos lógicos).
- **Gestión de Estado y Peticiones**: Custom Hooks (`useAuth`, `useParques`, `useZonas`, `useArboles`).
- **Mapas y Geolocalización**: Uso de `react-native-maps` y `expo-location` para capturar coordenadas automáticas.
- **Estilos**: Cada componente visual (`.tsx`) cuenta con su archivo de estilos segregado (`.styles.ts`) con una estética moderna en tonos oscuros (Deep Teal/Cyan).

### 2. Backend (Microservicios en Go)
El backend sigue un patrón arquitectónico hexagonal/puertos y adaptadores, lo que garantiza bajo acoplamiento.
Actualmente se han definido los siguientes servicios:
- **`api-auth`** (Puerto: 9090): Gestión de Usuarios, JWT y Roles.
- **`api-parque`** (Puerto: 9091): Módulo territorial para creación y listado de EcoParques.
- **`api-zona`** (Puerto: 9092): Subdivisión de parques en Zonas.
- **`api-flora`** (Puerto: 9093): Gestión del catálogo botánico y árboles plantados (con datos biométricos y GPS).

**Base de Datos y Caché:**
Cada microservicio requiere acceso a una base de datos **PostgreSQL** y un cluster de **Redis** para optimización y cacheo.

---

## 🚀 Guía de Instalación y Ejecución

Para correr todo el sistema de manera local, debes levantar primero la capa de persistencia, luego los microservicios, y finalmente la app móvil.

### Paso 1: Configurar la Persistencia de Datos
Asegúrate de que tu instancia de PostgreSQL y Redis en tu servidor local (ej. `192.168.1.28`) estén funcionando y aceptando conexiones. 

### Paso 2: Levantar los Microservicios (Backend)
Debes abrir terminales independientes para cada API y ejecutarlas. 
Revisa previamente el archivo `.env` de cada carpeta para asegurar que los puertos (`PORT=909X`) y las credenciales de BD sean correctas.

```bash
# Terminal 1: Auth
cd "api-auth"
go run cmd/api/main.go

# Terminal 2: Parques
cd "api-parque"
go run cmd/api/main.go

# Terminal 3: Zonas
cd "api-zona"
go run cmd/api/main.go

# Terminal 4: Flora (Cuando esté creada)
cd "api-flora"
go run cmd/api/main.go
```

### Paso 3: Configurar el Frontend
Antes de correr la aplicación móvil, asegúrate de que el archivo que consolida las URLs apunte a la IP de tu computadora dentro de tu red local.

Abre el archivo `eco-parque-app/src/services/api.ts` y verifica que las variables apunten correctamente a los servicios, por ejemplo:
```typescript
export const API_AUTH_URL = "http://192.168.1.27:9090/api/v1";
export const API_PARQUE_URL = "http://192.168.1.27:9091/api/v1";
export const API_ZONA_URL = "http://192.168.1.27:9092/api/v1";
export const API_FLORA_URL = "http://192.168.1.27:9093/api/v1";
```

### Paso 4: Levantar la Aplicación Móvil
Abre una terminal en la carpeta principal de la app y descarga las dependencias.

```bash
cd "eco-parque-app"
npm install
```

Luego, arranca el servidor de **Expo**:

```bash
# Limpiar caché e iniciar por túnel (recomendado si usas dispositivo físico con red estricta)
npx expo start -c --tunnel

# O arranque normal en red local
npx expo start
```

Escanea el código QR desde tu celular utilizando la app **Expo Go** (Android) o la cámara (iOS).

---

## 🛠 Solución de Problemas Frecuentes

1. **Error: `JSON Parse error: Unexpected character: <`**
   - **Causa**: La aplicación móvil está intentando llamar a una API que está apagada, o la ruta no existe. El proxy de Expo devuelve una página HTML de error (como un 404/502).
   - **Solución**: Asegúrate de que el microservicio correspondiente (ej. `api-zona`) esté corriendo en consola sin errores, y que los puertos coincidan.
   
2. **Error de Conexión a Base de Datos (`dial tcp 192.168.1.28:5432`)**
   - **Causa**: El microservicio no encuentra a PostgreSQL en la red.
   - **Solución**: Revisa si la base de datos está encendida. Si estás en Docker, asegúrate de que el contenedor de Postgres no esté pausado. Revisa tus archivos `.env`.

3. **Error en Expo (`Unexpected token` en un archivo `.tsx`)**
   - **Causa**: Desincronización de la caché de Metro Bundler.
   - **Solución**: Detén el proceso en la consola y vuelve a correr con `npx expo start -c`.
