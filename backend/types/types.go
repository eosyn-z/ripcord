package types

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Message represents a chat message with all necessary fields
type Message struct {
	ID        string    `json:"id" db:"id"`
	RoomID    string    `json:"room_id" db:"room_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	Type      string    `json:"type" db:"type"`
	Encrypted bool      `json:"encrypted" db:"encrypted"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Signature string    `json:"signature,omitempty" db:"signature"`
}

// Room represents a chat room
type Room struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	InviteCode  string    `json:"invite_code" db:"invite_code"`
	IsPrivate   bool      `json:"is_private" db:"is_private"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Participants []string  `json:"participants,omitempty"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	PublicKey string    `json:"public_key" db:"public_key"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	LastSeen  time.Time `json:"last_seen" db:"last_seen"`
	IsBlocked bool      `json:"is_blocked" db:"is_blocked"`
}

// Message types constants
const (
	MessageTypeText    = "text"
	MessageTypeCommand = "command"
	MessageTypeSystem  = "system"
	MessageTypeDM      = "dm"
	MessageTypeFile    = "file"
)

// Protocol message types
const (
	ProtocolMessageTypeJoin     = "join"
	ProtocolMessageTypeLeave    = "leave"
	ProtocolMessageTypeInvite   = "invite"
	ProtocolMessageTypeBlock    = "block"
	ProtocolMessageTypeUnblock  = "unblock"
	ProtocolMessageTypeDM       = "dm"
	ProtocolMessageTypeHeartbeat = "heartbeat"
)

// SlashCommand represents a parsed slash command
type SlashCommand struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// ProtocolMessage represents a protocol-level message
type ProtocolMessage struct {
	Type      string    `json:"type"`
	FromID    string    `json:"from_id"`
	MessageID string    `json:"message_id"`
	Data      string    `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// I2PAddress represents an I2P destination address
type I2PAddress struct {
	Destination string `json:"destination"`
	Base32      string `json:"base32"`
	Base64      string `json:"base64"`
}

// ConnectionStatus represents the status of a connection
type ConnectionStatus struct {
	Connected bool   `json:"connected"`
	Address   string `json:"address,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ServerStats represents server statistics
type ServerStats struct {
	Uptime   float64 `json:"uptime"`
	Rooms    int     `json:"rooms"`
	Messages int     `json:"messages"`
	Peers    int     `json:"peers"`
}

// AdminSettings represents admin-configurable settings
type AdminSettings struct {
	I2P struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		TunnelLength int    `json:"tunnel_length"`
	} `json:"i2p"`
	Server struct {
		Port                  int `json:"port"`
		MaxPeers             int `json:"max_peers"`
		MessageRetentionDays int `json:"message_retention_days"`
	} `json:"server"`
	Security struct {
		AutoBlockMalicious          bool `json:"auto_block_malicious"`
		RequireSignatureVerification bool `json:"require_signature_verification"`
		RateLimitPerMinute          int  `json:"rate_limit_per_minute"`
	} `json:"security"`
}

// Message methods for cryptographic operations and parsing

func (m *Message) Sign(privateKey ed25519.PrivateKey) error {
	if privateKey == nil {
		return errors.New("private key is nil")
	}
	
	signableData, err := m.getSignableData()
	if err != nil {
		return err
	}
	
	signature := ed25519.Sign(privateKey, signableData)
	m.Signature = hex.EncodeToString(signature)
	return nil
}

func (m *Message) getSignableData() ([]byte, error) {
	temp := *m
	temp.Signature = ""
	return json.Marshal(temp)
}

func (m *Message) VerifySignature(publicKey ed25519.PublicKey) bool {
	if m.Signature == "" || publicKey == nil {
		return false
	}
	
	signature, err := hex.DecodeString(m.Signature)
	if err != nil {
		return false
	}
	
	signableData, err := m.getSignableData()
	if err != nil {
		return false
	}
	
	return ed25519.Verify(publicKey, signableData, signature)
}

func (m *Message) IsSlashCommand() bool {
	return strings.HasPrefix(m.Content, "/")
}

func (m *Message) ParseSlashCommand() (*SlashCommand, error) {
	if !m.IsSlashCommand() {
		return nil, errors.New("not a slash command")
	}
	
	parts := strings.Fields(m.Content)
	if len(parts) == 0 {
		return nil, errors.New("empty command")
	}
	
	command := strings.TrimPrefix(parts[0], "/")
	args := parts[1:]
	
	return &SlashCommand{
		Command: command,
		Args:    args,
	}, nil
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
} 