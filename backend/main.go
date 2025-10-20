package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"syscall"
	"time"
	"unsafe"

	"github.com/gorilla/websocket"
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procReadProcessMemory = kernel32.NewProc("ReadProcessMemory")
	procOpenProcess       = kernel32.NewProc("OpenProcess")
)

type GameData struct {
	RivalName    string    `json:"rivalName"`
	RivalRank    string    `json:"rivalRank"`
	GameMode     string    `json:"gameMode"`
	Timestamp    time.Time `json:"timestamp"`
	IsInGame     bool      `json:"isInGame"`
	MatchStarted bool      `json:"matchStarted"`
}

type MemoryReader struct {
	processHandle uintptr
	baseAddress   uintptr
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir todas las conexiones (cambiar en producción)
	},
}

// OpenProcess abre un proceso por su PID
func OpenProcess(pid uint32) (uintptr, error) {
	const PROCESS_ALL_ACCESS = 0x1F0FFF
	handle, _, err := procOpenProcess.Call(
		uintptr(PROCESS_ALL_ACCESS),
		0,
		uintptr(pid),
	)

	if handle == 0 {
		return 0, fmt.Errorf("no se pudo abrir el proceso: %v", err)
	}

	return handle, nil
}

// ReadMemory lee memoria del proceso
func ReadMemory(handle uintptr, address uintptr, size uint32) ([]byte, error) {
	buffer := make([]byte, size)
	var bytesRead uintptr

	ret, _, err := procReadProcessMemory.Call(
		handle,
		address,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&bytesRead)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("error leyendo memoria: %v", err)
	}

	return buffer, nil
}

// FindProcess busca el proceso del juego por nombre
func FindProcessByName(name string) (uint32, error) {
	// Aquí deberías implementar la búsqueda del proceso
	// Por ahora retornamos un ejemplo
	// En producción, usar golang.org/x/sys/windows para enumerar procesos

	// Ejemplo: buscar "League of Legends.exe", "Teamfight Tactics.exe", etc.
	return 0, fmt.Errorf("implementar búsqueda de proceso")
}

// ReadGameData lee los datos del juego desde memoria
func (mr *MemoryReader) ReadGameData() (*GameData, error) {
	// Aquí defines las direcciones de memoria específicas de tu juego
	// Estas direcciones cambian según el juego y deben ser encontradas con Cheat Engine

	// Ejemplo para TFT/LoL (direcciones ficticias, debes encontrar las reales)
	const (
		RIVAL_NAME_OFFSET = 0x12345678
		RIVAL_RANK_OFFSET = 0x23456789
		GAME_STATE_OFFSET = 0x34567890
	)

	gameData := &GameData{
		Timestamp: time.Now(),
	}

	// Leer nombre del rival (asumiendo que es un string UTF-8)
	nameBytes, err := ReadMemory(mr.processHandle, mr.baseAddress+RIVAL_NAME_OFFSET, 64)
	if err == nil {
		gameData.RivalName = string(nameBytes)
	}

	// Leer rank del rival
	rankBytes, err := ReadMemory(mr.processHandle, mr.baseAddress+RIVAL_RANK_OFFSET, 32)
	if err == nil {
		gameData.RivalRank = string(rankBytes)
	}

	// Leer estado del juego
	stateBytes, err := ReadMemory(mr.processHandle, mr.baseAddress+GAME_STATE_OFFSET, 4)
	if err == nil && len(stateBytes) >= 4 {
		state := uint32(stateBytes[0]) | uint32(stateBytes[1])<<8 |
			uint32(stateBytes[2])<<16 | uint32(stateBytes[3])<<24
		gameData.IsInGame = state == 1
		gameData.MatchStarted = state == 2
	}

	return gameData, nil
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	log.Println("Cliente conectado")

	// Simular lectura de memoria (cambiar por lectura real)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Aquí iría la lectura real de memoria
		// Por ahora enviamos datos de prueba
		gameData := &GameData{
			RivalName:    "ProPlayer123",
			RivalRank:    "Diamond II",
			GameMode:     "Ranked",
			Timestamp:    time.Now(),
			IsInGame:     true,
			MatchStarted: true,
		}

		if err := conn.WriteJSON(gameData); err != nil {
			log.Println("Error enviando datos:", err)
			break
		}
	}
}

// Handler para datos estáticos
func handleGameData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Datos de ejemplo
	gameData := &GameData{
		RivalName:    "ShadowMaster",
		RivalRank:    "Master 250LP",
		GameMode:     "Ranked Solo",
		Timestamp:    time.Now(),
		IsInGame:     true,
		MatchStarted: true,
	}

	json.NewEncoder(w).Encode(gameData)
}

func main() {
	// Configurar rutas
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/api/gamedata", handleGameData)

	// Servir archivos estáticos (para desarrollo)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("Overlay Backend Running"))
	})

	port := ":8080"
	log.Printf("Servidor corriendo en http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Función helper para inicializar el lector de memoria
func InitMemoryReader(processName string) (*MemoryReader, error) {
	pid, err := FindProcessByName(processName)
	if err != nil {
		return nil, err
	}

	handle, err := OpenProcess(pid)
	if err != nil {
		return nil, err
	}

	// La dirección base debe ser encontrada dinámicamente
	// Aquí un ejemplo estático
	return &MemoryReader{
		processHandle: handle,
		baseAddress:   0x400000,
	}, nil
}
