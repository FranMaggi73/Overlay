# Game Overlay - Lectura de Memoria

## Instalación

### Backend (Go)
1. Instalar Go: https://go.dev/dl/
2. Instalar dependencias:
   ```bash
   cd backend
   go mod download
   ```
3. Compilar:
   ```bash
   go main.go
   ```

### Frontend (SvelteKit)
1. Instalar Node.js: https://nodejs.org/
2. Instalar dependencias:
   ```bash
   cd frontend
   npm install
   ```
3. Ejecutar en desarrollo:
   ```bash
   npm run dev
   ```
4. Construir para producción:
   ```bash
   npm run build

### Overlay (Electron)
1. Instalar Node.js: https://nodejs.org/
2. Instalar dependencias:
   ```bash
   cd overlayApp
   npm install
   ```
3. Construir para producción:
   ```bash
   npm start
   ```