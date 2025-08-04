package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"ripcord/database"
	"ripcord/security"
)

type Server struct {
	cryptoManager  *security.CryptoManager
	db             database.Database
	roomManager    *RoomManager
	messageHandler *MessageHandler
	node           *Node
	config         *Config
}

type Config struct {
	Port       string `json:"port"`
	DataDir    string `json:"data_dir"`
	I2PEnabled bool   `json:"i2p_enabled"`
	Nickname   string `json:"nickname"`
}

func main() {
	fmt.Println("Starting Ripcord - Decentralized Secure Chat Platform")
	
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	
	server, err := initializeServer(config)
	if err != nil {
		log.Fatal("Failed to initialize server:", err)
	}
	
	setupHTTPHandlers(server)
	
	go func() {
		fmt.Printf("Server starting on port %s\n", config.Port)
		if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()
	
	handleGracefulShutdown(server)
}

func loadConfig() (*Config, error) {
	configPath := "config.json"
	
	config := &Config{
		Port:       "8080",
		DataDir:    "./data",
		I2PEnabled: false,
		Nickname:   "Anonymous",
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, saveConfig(config, configPath)
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	
	if err := json.Unmarshal(data, config); err != nil {
		return config, err
	}
	
	return config, nil
}

func saveConfig(config *Config, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

func initializeServer(config *Config) (*Server, error) {
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, err
	}
	
	dbPath := filepath.Join(config.DataDir, "ripcord.db")
	db := database.NewSQLiteDatabase(dbPath)
	if err := db.Connect(); err != nil {
		return nil, err
	}
	
	keyPath := filepath.Join(config.DataDir, "identity.json")
	cryptoManager := security.NewCryptoManager(keyPath)
	if err := cryptoManager.LoadOrGenerateKeys(config.Nickname); err != nil {
		return nil, err
	}
	
	roomManager := NewRoomManager(db)
	messageHandler := NewMessageHandler(cryptoManager)
	
	node := NewNode(cryptoManager, roomManager, messageHandler)
	
	server := &Server{
		cryptoManager:  cryptoManager,
		db:             db,
		roomManager:    roomManager,
		messageHandler: messageHandler,
		node:           node,
		config:         config,
	}
	
	return server, nil
}

func setupHTTPHandlers(server *Server) {
	http.HandleFunc("/", serveStatic)
	http.HandleFunc("/api/identity", server.handleIdentity)
	http.HandleFunc("/api/rooms", server.handleRooms)
	http.HandleFunc("/api/rooms/create", server.handleCreateRoom)
	http.HandleFunc("/api/rooms/join", server.handleJoinRoom)
	http.HandleFunc("/api/rooms/leave", server.handleLeaveRoom)
	http.HandleFunc("/api/messages", server.handleMessages)
	http.HandleFunc("/api/messages/send", server.handleSendMessage)
	http.HandleFunc("/ws", server.handleWebSocket)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./frontend/index.html")
		return
	}
	
	staticDir := "./frontend/"
	http.StripPrefix("/", http.FileServer(http.Dir(staticDir))).ServeHTTP(w, r)
}

func (s *Server) handleIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	identity := map[string]interface{}{
		"nickname":    s.cryptoManager.GetNickname(),
		"public_key":  s.cryptoManager.GetPublicKeyBase58(),
		"fingerprint": s.cryptoManager.GetPublicKeyFingerprint(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(identity)
}

func (s *Server) handleRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	rooms, err := s.db.GetRooms()
	if err != nil {
		http.Error(w, "Failed to get rooms", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPrivate   bool   `json:"is_private"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	userID := s.cryptoManager.GetPublicKeyBase58()
	username := s.cryptoManager.GetNickname()
	publicKey := s.cryptoManager.GetPublicKeyBase58()
	
	room, err := s.roomManager.CreateRoom(req.Name, req.Description, req.IsPrivate, userID, username, publicKey)
	if err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func (s *Server) handleJoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		InviteCode string `json:"invite_code"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	userID := s.cryptoManager.GetPublicKeyBase58()
	username := s.cryptoManager.GetNickname()
	publicKey := s.cryptoManager.GetPublicKeyBase58()
	
	room, err := s.roomManager.JoinRoomByInvite(req.InviteCode, userID, username, publicKey)
	if err != nil {
		http.Error(w, "Failed to join room", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func (s *Server) handleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		RoomID string `json:"room_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	userID := s.cryptoManager.GetPublicKeyBase58()
	
	if err := s.roomManager.LeaveRoom(req.RoomID, userID); err != nil {
		http.Error(w, "Failed to leave room", http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "room_id parameter required", http.StatusBadRequest)
		return
	}
	
	messages, err := s.db.GetMessages(roomID, 50)
	if err != nil {
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		RoomID  string `json:"room_id"`
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	userID := s.cryptoManager.GetPublicKeyBase58()
	username := s.cryptoManager.GetNickname()
	
	message, err := s.messageHandler.CreateSignedMessage(req.RoomID, userID, username, req.Content, "")
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}
	
	dbMessage := &database.Message{
		ID:        message.ID,
		RoomID:    message.RoomID,
		UserID:    message.UserID,
		Username:  message.Username,
		Content:   message.Content,
		Type:      message.Type,
		Encrypted: message.Encrypted,
		Timestamp: message.Timestamp,
		Signature: message.Signature,
	}
	
	if err := s.db.SaveMessage(dbMessage); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement WebSocket for real-time communication
	http.Error(w, "WebSocket not implemented yet", http.StatusNotImplemented)
}

func handleGracefulShutdown(server *Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	<-c
	fmt.Println("\nShutting down gracefully...")
	
	if server.db != nil {
		server.db.Disconnect()
	}
	
	fmt.Println("Server stopped")
	os.Exit(0)
} 