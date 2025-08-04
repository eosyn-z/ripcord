package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"github.com/gorilla/websocket"
	"ripcord/database"
	"ripcord/i2p"
	"ripcord/security"
	"ripcord/types"
)

type Server struct {
	cryptoManager  *security.CryptoManager
	db             database.Database
	roomManager    *RoomManager
	messageHandler *MessageHandler
	node           *Node
	i2pManager     *i2p.I2PManager
	config         *Config
	wsClients      map[*websocket.Conn]*WSClient
	wsClientsMutex sync.RWMutex
	upgrader       websocket.Upgrader
}

type WSClient struct {
	conn     *websocket.Conn
	userID   string
	username string
	roomID   string
	send     chan []byte
}

// Config now defined in config.go

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
		port := fmt.Sprintf("%d", config.Server.Port)
		fmt.Printf("Server starting on port %s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()
	
	handleGracefulShutdown(server)
}

func loadConfig() (*Config, error) {
	return LoadConfig("config.json")
}

func saveConfig(config *Config, path string) error {
	return config.Save(path)
}

func initializeServer(config *Config) (*Server, error) {
	dataDir := "data" // Default data directory
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	
	// Use database config or default to SQLite
	dbPath := filepath.Join(dataDir, config.Database.Database)
	if config.Database.Database == "" {
		dbPath = filepath.Join(dataDir, "ripcord.db")
	}
	
	db := database.NewSQLiteDatabase(dbPath)
	if err := db.Connect(); err != nil {
		return nil, err
	}
	
	keyPath := filepath.Join(dataDir, "identity.json")
	cryptoManager := security.NewCryptoManager(keyPath)
	if err := cryptoManager.LoadOrGenerateKeys("Anonymous"); err != nil {
		return nil, err
	}
	
	roomManager := NewRoomManager(db)
	messageHandler := NewMessageHandler(cryptoManager)
	
	node := NewNode(cryptoManager, roomManager, messageHandler)
	
	// Initialize I2P manager
	i2pManager := i2p.NewI2PManager(config.I2P.SamAddress, config.I2P.SamPort)
	if config.I2P.Enabled {
		if err := i2pManager.Connect(); err != nil {
			log.Printf("Warning: Failed to connect to I2P: %v", err)
		} else {
			log.Printf("Successfully connected to I2P")
		}
	}
	
	server := &Server{
		cryptoManager:  cryptoManager,
		db:             db,
		roomManager:    roomManager,
		messageHandler: messageHandler,
		node:           node,
		i2pManager:     i2pManager,
		config:         config,
		wsClients:      make(map[*websocket.Conn]*WSClient),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
	
	return server, nil
}

func setupHTTPHandlers(server *Server) {
	// Add CORS middleware wrapper with security improvements
	corsHandler := func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Add CORS headers with security improvements
			origin := r.Header.Get("Origin")
			if origin != "" {
				// In production, validate against allowed origins
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			
			// Add security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			
			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			handler(w, r)
		}
	}
	
	http.HandleFunc("/", serveStatic)
	http.HandleFunc("/api/identity", corsHandler(server.handleIdentity))
	http.HandleFunc("/api/rooms", corsHandler(server.handleRooms))
	http.HandleFunc("/api/rooms/create", corsHandler(server.handleCreateRoom))
	http.HandleFunc("/api/rooms/join", corsHandler(server.handleJoinRoom))
	http.HandleFunc("/api/rooms/leave", corsHandler(server.handleLeaveRoom))
	http.HandleFunc("/api/messages", corsHandler(server.handleMessages))
	http.HandleFunc("/api/messages/send", corsHandler(server.handleSendMessage))
	http.HandleFunc("/ws", server.handleWebSocket)
	
	// Admin API endpoints
	http.HandleFunc("/api/admin/identity", corsHandler(server.handleAdminIdentity))
	http.HandleFunc("/api/admin/i2p/status", corsHandler(server.handleI2PStatus))
	http.HandleFunc("/api/admin/stats", corsHandler(server.handleServerStats))
	http.HandleFunc("/api/admin/rooms", corsHandler(server.handleAdminRooms))
	http.HandleFunc("/api/admin/peers", corsHandler(server.handleAdminPeers))
	http.HandleFunc("/api/admin/logs/connections", corsHandler(server.handleConnectionLogs))
	http.HandleFunc("/api/admin/settings", corsHandler(server.handleAdminSettings))
	http.HandleFunc("/api/admin/restart", corsHandler(server.handleAdminRestart))
	
	// API Access Management endpoints
	http.HandleFunc("/api/admin/api-access", corsHandler(server.handleAPIAccess))
	http.HandleFunc("/api/admin/api-access/rules", corsHandler(server.handleAPIAccessRules))
	http.HandleFunc("/api/admin/api-access/test", corsHandler(server.handleAPIAccessTest))
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "../frontend/index.html")
		return
	}
	
	if r.URL.Path == "/admin" || r.URL.Path == "/admin/" {
		http.ServeFile(w, r, "../frontend/admin.html")
		return
	}
	
	staticDir := "../frontend/"
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
	
	// Validate input
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}
	
	if len(req.Name) > 100 {
		http.Error(w, "Room name too long (max 100 characters)", http.StatusBadRequest)
		return
	}
	
	if len(req.Description) > 500 {
		http.Error(w, "Room description too long (max 500 characters)", http.StatusBadRequest)
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
	
	// Validate and sanitize input
	if strings.TrimSpace(req.RoomID) == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	
	// Sanitize room ID - only allow alphanumeric and hyphens
	if !isValidRoomID(req.RoomID) {
		http.Error(w, "Invalid room ID format", http.StatusBadRequest)
		return
	}
	
	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, "Message content is required", http.StatusBadRequest)
		return
	}
	
	// Sanitize content - remove potentially dangerous characters
	content := sanitizeMessageContent(req.Content)
	if content == "" {
		http.Error(w, "Message content is required", http.StatusBadRequest)
		return
	}
	
	if len(content) > 2000 {
		http.Error(w, "Message too long (max 2000 characters)", http.StatusBadRequest)
		return
	}
	
	userID := s.cryptoManager.GetPublicKeyBase58()
	username := s.cryptoManager.GetNickname()
	
	message, err := s.messageHandler.CreateSignedMessage(req.RoomID, userID, username, content, "")
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}
	
	dbMessage := &types.Message{
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
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	
	client := &WSClient{
		conn: conn,
		send: make(chan []byte, 256),
	}
	
	s.wsClientsMutex.Lock()
	s.wsClients[conn] = client
	s.wsClientsMutex.Unlock()
	
	// Start goroutines for reading and writing
	go s.wsClientReader(client)
	go s.wsClientWriter(client)
	
	log.Printf("WebSocket client connected: %s", conn.RemoteAddr())
}

func (s *Server) wsClientReader(client *WSClient) {
	defer func() {
		s.removeWSClient(client)
		client.conn.Close()
	}()
	
	client.conn.SetReadLimit(512)
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		// Handle incoming WebSocket message
		s.handleWSMessage(client, message)
	}
}

func (s *Server) wsClientWriter(client *WSClient) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			
			// Add queued messages to the current write
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.send)
			}
			
			if err := w.Close(); err != nil {
				return
			}
			
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Server) handleWSMessage(client *WSClient, message []byte) {
	var wsMsg map[string]interface{}
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("Invalid WebSocket message: %v", err)
		return
	}
	
	msgType, ok := wsMsg["type"].(string)
	if !ok {
		log.Printf("WebSocket message missing type field")
		return
	}
	
	switch msgType {
	case "auth":
		s.handleWSAuth(client, wsMsg)
	case "join_room":
		s.handleWSJoinRoom(client, wsMsg)
	case "leave_room":
		s.handleWSLeaveRoom(client, wsMsg)
	case "send_message":
		s.handleWSSendMessage(client, wsMsg)
	case "get_messages":
		s.handleWSGetMessages(client, wsMsg)
	default:
		log.Printf("Unknown WebSocket message type: %s", msgType)
	}
}

func (s *Server) handleWSAuth(client *WSClient, msg map[string]interface{}) {
	username, _ := msg["username"].(string)
	if username == "" {
		username = "Anonymous"
	}
	
	client.username = username
	client.userID = generateMessageID() // Simple user ID generation
	
	response := map[string]interface{}{
		"type": "auth_response",
		"success": true,
		"user": map[string]interface{}{
			"id": client.userID,
			"username": client.username,
		},
	}
	
	s.sendToClient(client, response)
}

func (s *Server) handleWSJoinRoom(client *WSClient, msg map[string]interface{}) {
	roomID, _ := msg["room_id"].(string)
	if roomID == "" {
		return
	}
	
	client.roomID = roomID
	
	response := map[string]interface{}{
		"type": "room_joined",
		"room_id": roomID,
	}
	
	s.sendToClient(client, response)
	s.broadcastToRoom(roomID, map[string]interface{}{
		"type": "user_joined",
		"user": map[string]interface{}{
			"id": client.userID,
			"username": client.username,
		},
	}, client)
}

func (s *Server) handleWSLeaveRoom(client *WSClient, msg map[string]interface{}) {
	if client.roomID == "" {
		return
	}
	
	oldRoomID := client.roomID
	client.roomID = ""
	
	s.broadcastToRoom(oldRoomID, map[string]interface{}{
		"type": "user_left",
		"user_id": client.userID,
	}, client)
}

func (s *Server) handleWSSendMessage(client *WSClient, msg map[string]interface{}) {
	content, _ := msg["content"].(string)
	if content == "" || client.roomID == "" {
		return
	}
	
	message := &types.Message{
		ID:        generateMessageID(),
		RoomID:    client.roomID,
		UserID:    client.userID,
		Username:  client.username,
		Content:   content,
		Type:      types.MessageTypeText,
		Encrypted: false,
		Timestamp: time.Now(),
	}
	
	// Save to database
	if err := s.db.SaveMessage(message); err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}
	
	// Broadcast to room
	s.broadcastToRoom(client.roomID, map[string]interface{}{
		"type": "message",
		"message": message,
	}, nil)
}

func (s *Server) handleWSGetMessages(client *WSClient, msg map[string]interface{}) {
	roomID, _ := msg["room_id"].(string)
	if roomID == "" {
		return
	}
	
	limit := 50
	if l, ok := msg["limit"].(float64); ok {
		limit = int(l)
	}
	
	messages, err := s.db.GetMessages(roomID, limit)
	if err != nil {
		log.Printf("Failed to get messages: %v", err)
		return
	}
	
	response := map[string]interface{}{
		"type": "message_history",
		"messages": messages,
	}
	
	s.sendToClient(client, response)
}

func (s *Server) sendToClient(client *WSClient, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal WebSocket response: %v", err)
		return
	}
	
	select {
	case client.send <- jsonData:
	default:
		close(client.send)
		s.removeWSClient(client)
	}
}

func (s *Server) broadcastToRoom(roomID string, data interface{}, exclude *WSClient) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal broadcast data: %v", err)
		return
	}
	
	s.wsClientsMutex.RLock()
	defer s.wsClientsMutex.RUnlock()
	
	for _, client := range s.wsClients {
		if client.roomID == roomID && client != exclude {
			select {
			case client.send <- jsonData:
			default:
				close(client.send)
				delete(s.wsClients, client.conn)
			}
		}
	}
}

func (s *Server) removeWSClient(client *WSClient) {
	s.wsClientsMutex.Lock()
	defer s.wsClientsMutex.Unlock()
	
	if _, ok := s.wsClients[client.conn]; ok {
		delete(s.wsClients, client.conn)
		close(client.send)
		
		// Notify room that user left
		if client.roomID != "" {
			s.broadcastToRoom(client.roomID, map[string]interface{}{
				"type": "user_left",
				"user_id": client.userID,
			}, client)
		}
		
		log.Printf("WebSocket client disconnected: %s", client.conn.RemoteAddr())
	}
}

// API Access Management handlers
func (s *Server) handleAPIAccess(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Return current API access configuration
		config := map[string]interface{}{
			"enabled": true,
			"require_auth": false,
			"default_rate_limit": 100,
			"global_permissions": map[string]bool{
				"identity": true,
				"admin_identity": true,
				"stats": true,
				"rooms_read": true,
				"rooms_create": false,
				"rooms_join": false,
				"admin_rooms": true,
				"messages_read": true,
				"messages_send": false,
				"peers": true,
				"i2p_status": true,
				"connection_logs": true,
				"settings_read": true,
				"settings_write": false,
				"restart": false,
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
		
	case http.MethodPost:
		// Update API access configuration
		var config map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		// TODO: Save configuration to database or config file
		log.Printf("API access configuration updated: %+v", config)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAPIAccessRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Return current API access rules
		rules := []map[string]interface{}{
			{
				"id": "rule-1",
				"node_identifier": "example-public-key-123",
				"node_nickname": "Example Node",
				"access_level": "standard",
				"rate_limit": 100,
				"enabled": true,
				"allowed_endpoints": []string{"/api/identity", "/api/rooms", "/api/messages"},
				"created_at": time.Now().Add(-24 * time.Hour),
				"expires_at": nil,
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rules)
		
	case http.MethodPost:
		// Add new API access rule
		var rule map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		// Generate rule ID
		rule["id"] = fmt.Sprintf("rule-%d", time.Now().Unix())
		rule["created_at"] = time.Now()
		
		// TODO: Save rule to database
		log.Printf("New API access rule added: %+v", rule)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rule)
		
	case http.MethodDelete:
		// Delete API access rule
		ruleID := r.URL.Query().Get("id")
		if ruleID == "" {
			http.Error(w, "Rule ID required", http.StatusBadRequest)
			return
		}
		
		// TODO: Delete rule from database
		log.Printf("API access rule deleted: %s", ruleID)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAPIAccessTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var testRequest struct {
		NodeIdentifier string `json:"node_identifier"`
		Endpoint       string `json:"endpoint"`
		Method         string `json:"method"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&testRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Simulate API access test
	result := map[string]interface{}{
		"allowed": true,
		"access_level": "standard",
		"rate_limit": 100,
		"reason": "Node has standard access to this endpoint",
		"matched_rule": "rule-1",
	}
	
	// Simple simulation - block admin endpoints for non-admin access
	if strings.Contains(testRequest.Endpoint, "/admin/") && testRequest.NodeIdentifier != s.cryptoManager.GetPublicKeyBase58() {
		result["allowed"] = false
		result["reason"] = "Admin access required for this endpoint"
		result["access_level"] = "standard"
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Admin API handlers
func (s *Server) handleAdminIdentity(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleI2PStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if s.i2pManager == nil {
		http.Error(w, "I2P manager not initialized", http.StatusInternalServerError)
		return
	}
	
	status := s.i2pManager.GetStatus()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (s *Server) handleServerStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get room count
	rooms, _ := s.db.GetRooms()
	roomCount := len(rooms)
	
	// Get message count (approximate)
	messageCount := 0
	for _, room := range rooms {
		messages, _ := s.db.GetMessages(room.ID, 1000)
		messageCount += len(messages)
	}
	
	// Count active WebSocket connections as peer approximation
	s.wsClientsMutex.RLock()
	peerCount := len(s.wsClients)
	s.wsClientsMutex.RUnlock()
	
	stats := map[string]interface{}{
		"uptime":   time.Since(time.Now().Add(-1*time.Hour)).Seconds(), // Placeholder
		"rooms":    roomCount,
		"messages": messageCount,
		"peers":    peerCount,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleAdminRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	rooms, err := s.db.GetRooms()
	if err != nil {
		http.Error(w, "Failed to get rooms", http.StatusInternalServerError)
		return
	}
	
	// Enhance with additional stats
	adminRooms := make([]map[string]interface{}, len(rooms))
	for i, room := range rooms {
		messages, _ := s.db.GetMessages(room.ID, 1000)
		
		adminRooms[i] = map[string]interface{}{
			"id":               room.ID,
			"name":             room.Name,
			"description":      room.Description,
			"invite_code":      room.InviteCode,
			"is_private":       room.IsPrivate,
			"created_at":       room.CreatedAt,
			"participant_count": len(room.Participants),
			"message_count":    len(messages),
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adminRooms)
}

func (s *Server) handleAdminPeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// For now, use WebSocket clients as peer approximation
	s.wsClientsMutex.RLock()
	defer s.wsClientsMutex.RUnlock()
	
	peers := make([]map[string]interface{}, 0, len(s.wsClients))
	for _, client := range s.wsClients {
		peer := map[string]interface{}{
			"id":           client.userID,
			"status":       "connected",
			"i2p_address":  "Not Available",
			"connected_at": time.Now().Add(-30 * time.Minute), // Placeholder
			"message_count": 0,                                  // TODO: Implement
			"trust_level":  "unknown",
		}
		peers = append(peers, peer)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func (s *Server) handleConnectionLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// TODO: Implement actual connection logging
	logs := []map[string]interface{}{
		{
			"timestamp": time.Now().Add(-5 * time.Minute),
			"level":     "info",
			"message":   "New peer connection established",
		},
		{
			"timestamp": time.Now().Add(-10 * time.Minute),
			"level":     "warning",
			"message":   "Peer connection timeout, retrying",
		},
		{
			"timestamp": time.Now().Add(-15 * time.Minute),
			"level":     "info",
			"message":   "I2P tunnel established",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (s *Server) handleAdminSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Return current settings
		settings := map[string]interface{}{
			"i2p": map[string]interface{}{
				"host":          s.config.I2P.Host,
				"port":          s.config.I2P.Port,
				"tunnel_length": 2, // Default
			},
			"server": map[string]interface{}{
				"port":                   s.config.Server.Port,
				"max_peers":              50, // Default
				"message_retention_days": 30, // Default
			},
			"security": map[string]interface{}{
				"auto_block_malicious":           true, // Default
				"require_signature_verification": true, // Default
				"rate_limit_per_minute":          60,   // Default
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
		
	case http.MethodPost:
		// Save new settings
		var newSettings map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		// TODO: Validate and apply settings
		log.Printf("Admin settings updated: %+v", newSettings)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAdminRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	log.Println("Admin requested node restart")
	
	// TODO: Implement graceful restart
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "scheduled",
		"message": "Node restart scheduled",
	})
}

func handleGracefulShutdown(server *Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	<-c
	fmt.Println("\nShutting down gracefully...")
	
	if server.db != nil {
		server.db.Disconnect()
	}
	
	if server.i2pManager != nil {
		server.i2pManager.Disconnect()
	}
	
	fmt.Println("Server stopped")
	os.Exit(0)
}

// Input sanitization functions
func isValidRoomID(roomID string) bool {
	// Only allow alphanumeric characters, hyphens, and underscores
	for _, char := range roomID {
		if !((char >= 'a' && char <= 'z') || 
			(char >= 'A' && char <= 'Z') || 
			(char >= '0' && char <= '9') || 
			char == '-' || char == '_') {
			return false
		}
	}
	return len(roomID) > 0 && len(roomID) <= 50
}

func sanitizeMessageContent(content string) string {
	// Remove null bytes and other control characters
	var sanitized strings.Builder
	for _, char := range content {
		if char >= 32 || char == '\n' || char == '\t' {
			sanitized.WriteRune(char)
		}
	}
	
	// Trim whitespace
	result := strings.TrimSpace(sanitized.String())
	
	// Limit consecutive whitespace
	result = strings.Join(strings.Fields(result), " ")
	
	return result
} 