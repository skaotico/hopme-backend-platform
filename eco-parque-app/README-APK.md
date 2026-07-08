# Guía de Construcción y Distribución (APK)

Este documento detalla los pasos para compilar la aplicación móvil **Eco Parque App** en un archivo instalable (`.apk`) para pruebas locales o distribución interna.

## 1. Requisitos Previos

Antes de comenzar, asegúrate de tener instalado:
- **Node.js** (versión recomendada LTS)
- Una cuenta gratuita en [Expo](https://expo.dev/) (para compilaciones en la nube)
- La interfaz de línea de comandos de Expo Application Services (EAS CLI)

Para instalar el CLI de EAS a nivel global, ejecuta en tu terminal:
```bash
npm install -g eas-cli
```

Inicia sesión en tu cuenta de Expo:
```bash
eas login
```

---

## 2. Configuración del Proyecto

El proyecto ya cuenta con el archivo `eas.json` configurado. Este archivo le dice a EAS cómo debe compilar la aplicación.

Para generar un APK instalable (en lugar de un AAB para la Play Store), utilizamos el perfil `preview` que incluye la directiva:
```json
"android": {
  "buildType": "apk"
}
```

---

## 3. Generar el APK (Compilación en la Nube)

La forma más sencilla de compilar la aplicación es usando los servidores gratuitos de Expo. Así no necesitas configurar Android Studio ni Java en tu computadora.

Desde la carpeta raíz del proyecto de React Native (`eco-parque-app`), ejecuta:

```bash
eas build -p android --profile preview
```

### Proceso:
1. EAS subirá tu código fuente a la nube.
2. Comenzará a compilar la app. Puedes ver el progreso en tu terminal o en el panel web de Expo.
3. Una vez finalizado, la terminal te entregará un **enlace directo** y un **código QR**.

---

## 4. Alternativa: Compilación Local

Si prefieres no usar la nube de Expo y quieres que el procesamiento se haga en tu CPU, puedes forzar la compilación local.
*⚠️ **Nota:** Requiere tener instalado y configurado el Android SDK, NDK y Java JDK.*

Ejecuta:
```bash
eas build -p android --profile preview --local
```
Al terminar, el archivo `build-XXXXX.apk` quedará guardado directamente en la carpeta de tu proyecto.

---

## 5. Distribución e Instalación

Una vez que tengas el enlace (si compilaste en la nube) o el archivo físico (si compilaste localmente), puedes distribuirlo a tu equipo de pruebas.

### Para Emuladores (Android Studio)
Si estás usando un emulador en tu PC, simplemente arrastra el archivo `.apk` descargado hacia la ventana del emulador. Se instalará automáticamente.

### Para Dispositivos Físicos (Teléfonos Android)
1. **Opción A (Código QR):** Con tu teléfono, escanea el código QR que apareció en la terminal al terminar el build en la nube.
2. **Opción B (Transferencia directa):** Envía el archivo `.apk` por Google Drive, correo electrónico, Telegram, o transfiérelo por cable USB.
3. En el dispositivo, asegúrate de tener habilitada la opción **"Instalar aplicaciones de orígenes desconocidos"** (en Ajustes > Seguridad).
4. Toca el archivo `.apk` desde el explorador de archivos del teléfono e instálalo.

---

## Solución de Problemas Comunes

- **Error de "No credentials found":** Si EAS te pregunta por un *keystore*, permite que Expo genere uno automáticamente o sigue los pasos en pantalla para proveer el tuyo.
- **Fallos por dependencias:** Asegúrate de que `node_modules` esté actualizado (borra la carpeta y corre `npm install` si tienes dudas) antes de lanzar el build local. En builds de nube, EAS instala las dependencias desde cero usando tu `package-lock.json`.
