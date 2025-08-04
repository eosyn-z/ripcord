package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

const (
	ProtocolVersion = "1.0"
	
	MessageTypeHeartbeat = "heartbeat"
	MessageTypeJoin      = "join"
	MessageTypeLeave     = "leave"
	MessageTypeChat      = "chat"
	MessageTypeSync      = "sync"
	MessageTypePing      = "ping"
	MessageTypePong      = "pong"
	MessageTypeInvite    = "invite"
	MessageTypeDM        = "dm"
	MessageTypeRoomInfo  = "room_info"
	MessageTypeUserInfo  = "user_info"
	MessageTypeBlock     = "block"
	MessageTypeUnblock   = "unblock"
)

type ProtocolMessage struct {
	Version   string      `json:"version"`
	Type      string      `json:"type"`
	MessageID string      `json:"message_id"`
	From      string      `json:"from"`
	To        string      `json:"to,omitempty"`
	RoomID    string      `json:"room_id,omitempty"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
	Signature string      `json:"signature"`
}

type HeartbeatPayload struct {
	Nickname    string   `json:"nickname"`
	PublicKey   string   `json:"public_key"`
	I2PAddress  string   `json:"i2p_address"`
	ActiveRooms []string `json:"active_rooms"`
}

type ChatPayload struct {
	Content   string `json:"content"`
	IsCommand bool   `json:"is_command,omitempty"`
}

type JoinPayload struct {
	RoomID     string `json:"room_id"`
	InviteCode string `json:"invite_code,omitempty"`
	Nickname   string `json:"nickname"`
	PublicKey  string `json:"public_key"`
}

type LeavePayload struct {
	RoomID string `json:"room_id"`
	Reason string `json:"reason,omitempty"`
}

type InvitePayload struct {
	RoomID      string `json:"room_id"`
	RoomName    string `json:"room_name"`
	InviteCode  string `json:"invite_code"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private"`
}

type DMPayload struct {
	Content     string `json:"content"`
	IsEncrypted bool   `json:"is_encrypted"`
}

type RoomInfoPayload struct {
	RoomID      string   `json:"room_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsPrivate   bool     `json:"is_private"`
	Moderators  []string `json:"moderators"`
	Participants []string `json:"participants"`
}

type UserInfoPayload struct {
	Nickname    string `json:"nickname"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	IsBlocked   bool   `json:"is_blocked,omitempty"`
}

type SyncPayload struct {
	RoomID       string        `json:"room_id"`
	LastSyncTime int64         `json:"last_sync_time"`
	Messages     []SyncMessage `json:"messages,omitempty"`
}

type SyncMessage struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

func NewProtocolMessage(msgType, from, messageID string) *ProtocolMessage {
	return &ProtocolMessage{
		Version:   ProtocolVersion,
		Type:      msgType,
		MessageID: messageID,
		From:      from,
		Timestamp: time.Now().Unix(),
	}
}

func (pm *ProtocolMessage) SetPayload(payload interface{}) {
	pm.Payload = payload
}

func (pm *ProtocolMessage) ToJSON() ([]byte, error) {
	return json.Marshal(pm)
}

func (pm *ProtocolMessage) Sign(privateKey ed25519.PrivateKey) error {
	data, err := pm.GetSignableData()
	if err != nil {
		return err
	}
	signature := ed25519.Sign(privateKey, data)
	pm.Signature = hex.EncodeToString(signature)
	return nil
}

func (pm *ProtocolMessage) GetSignableData() ([]byte, error) {
	tempMsg := *pm
	tempMsg.Signature = ""
	return json.Marshal(tempMsg)
}

func (pm *ProtocolMessage) VerifySignature(publicKey ed25519.PublicKey) bool {
	if pm.Signature == "" {
		return false
	}
	
	signature, err := hex.DecodeString(pm.Signature)
	if err != nil {
		return false
	}
	
	data, err := pm.GetSignableData()
	if err != nil {
		return false
	}
	
	return ed25519.Verify(publicKey, data, signature)
}

func (pm *ProtocolMessage) IsValid() error {
	if pm.Version != ProtocolVersion {
		return errors.New("unsupported protocol version")
	}
	
	if pm.Type == "" {
		return errors.New("message type is required")
	}
	
	if pm.From == "" {
		return errors.New("from field is required")
	}
	
	if pm.MessageID == "" {
		return errors.New("message ID is required")
	}
	
	if pm.Timestamp == 0 {
		return errors.New("timestamp is required")
	}
	
	return nil
}

func ParseProtocolMessage(data []byte) (*ProtocolMessage, error) {
	var msg ProtocolMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	
	if err := msg.IsValid(); err != nil {
		return nil, err
	}
	
	return &msg, nil
}

func (pm *ProtocolMessage) GetTypedPayload() (interface{}, error) {
	if pm.Payload == nil {
		return nil, nil
	}
	
	payloadBytes, err := json.Marshal(pm.Payload)
	if err != nil {
		return nil, err
	}
	
	switch pm.Type {
	case MessageTypeHeartbeat:
		var payload HeartbeatPayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeChat:
		var payload ChatPayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeJoin:
		var payload JoinPayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeLeave:
		var payload LeavePayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeInvite:
		var payload InvitePayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeDM:
		var payload DMPayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeRoomInfo:
		var payload RoomInfoPayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeUserInfo:
		var payload UserInfoPayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	case MessageTypeSync:
		var payload SyncPayload
		err = json.Unmarshal(payloadBytes, &payload)
		return payload, err
	default:
		return pm.Payload, nil
	}
} 