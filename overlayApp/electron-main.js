// electron-main.js
const { app, BrowserWindow, ipcMain, globalShortcut } = require('electron');

let overlayWindow;
let clickThroughMode = true; // Por defecto, puedes hacer click al juego

function createOverlay() {
  overlayWindow = new BrowserWindow({
    width: 1920,
    height: 1080,
    transparent: true,
    frame: false,
    alwaysOnTop: true,
    skipTaskbar: true,
    resizable: false,
    focusable: false, // No toma focus por defecto
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      enableRemoteModule: false
    }
  });

  // Cargar tu overlay de localhost
  overlayWindow.loadURL('http://localhost:5173');

  // Click-through ACTIVADO por defecto (puedes hacer click al juego)
  overlayWindow.setIgnoreMouseEvents(true, { forward: true });

  // Mantener siempre visible
  overlayWindow.setAlwaysOnTop(true, 'screen-saver');
  overlayWindow.setVisibleOnAllWorkspaces(true);

  overlayWindow.on('closed', () => {
    overlayWindow = null;
  });

  // Mostrar mensaje de atajos al iniciar
  overlayWindow.webContents.on('did-finish-load', () => {
    console.log(`
╔════════════════════════════════════════════╗
║   GAME OVERLAY - ATAJOS DE TECLADO         ║
╠════════════════════════════════════════════╣
║  Ctrl+Shift+I → Modo INTERACTIVO           ║
║     (Desactiva click-through para usar     ║
║      los botones del overlay)              ║
║                                            ║
║  Ctrl+Shift+G → Modo JUEGO                 ║
║     (Activa click-through para jugar)      ║
║                                            ║
║  Ctrl+Shift+O → Mostrar/Ocultar overlay    ║
║  Ctrl+Shift+R → Recargar overlay           ║
║  Ctrl+Shift+Q → Cerrar aplicación          ║
╠════════════════════════════════════════════╣
║  MODO ACTUAL: JUEGO (click-through ON)     ║
╚════════════════════════════════════════════╝
    `);
  });
}

app.whenReady().then(() => {
  createOverlay();
  
  // Ctrl+Shift+I - MODO INTERACTIVO (puedes usar botones del overlay)
  globalShortcut.register('Ctrl+Shift+I', () => {
    if (overlayWindow) {
      clickThroughMode = false;
      overlayWindow.setIgnoreMouseEvents(false);
      overlayWindow.setFocusable(true);
      overlayWindow.focus();
      console.log('\n🎯 MODO INTERACTIVO: Ahora puedes usar los botones del overlay');
    }
  });

  // Ctrl+Shift+G - MODO JUEGO (vuelve a click-through)
  globalShortcut.register('Ctrl+Shift+G', () => {
    if (overlayWindow) {
      clickThroughMode = true;
      overlayWindow.setIgnoreMouseEvents(true, { forward: true });
      overlayWindow.setFocusable(false);
      console.log('\n🎮 MODO JUEGO: Ahora puedes hacer click al juego debajo');
    }
  });

  // Ctrl+Shift+O - Mostrar/Ocultar overlay
  globalShortcut.register('Ctrl+Shift+O', () => {
    if (overlayWindow) {
      if (overlayWindow.isVisible()) {
        overlayWindow.hide();
        console.log('\n👁️ Overlay OCULTO');
      } else {
        overlayWindow.show();
        console.log('\n👁️ Overlay VISIBLE');
      }
    }
  });

  // Ctrl+Shift+R - Recargar overlay
  globalShortcut.register('Ctrl+Shift+R', () => {
    if (overlayWindow) {
      overlayWindow.reload();
      console.log('\n🔄 Overlay recargado');
    }
  });

  // Ctrl+Shift+Q - Cerrar overlay
  globalShortcut.register('Ctrl+Shift+Q', () => {
    console.log('\n❌ Cerrando overlay...');
    app.quit();
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createOverlay();
  }
});

app.on('will-quit', () => {
  globalShortcut.unregisterAll();
});