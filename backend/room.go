package main

import (
	"crypto/rand"
	"errors"
	"strings"
	"sync"
	"time"
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"ripcord/database"
)

type Room struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	InviteCode  string             `json:"invite_code"`
	IsPrivate   bool               `json:"is_private"`
	Members     map[string]*Member `json:"members"`
	Moderators  map[string]bool    `json:"moderators"`
	Messages    []*Message         `json:"messages,omitempty"` // For testing compatibility
	CreatedAt   time.Time          `json:"created_at"`
	mu          sync.RWMutex       `json:"-"`
}

type Member struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	PublicKey string    `json:"public_key"`
	JoinedAt  time.Time `json:"joined_at"`
	Role      string    `json:"role"`
	IsBlocked bool      `json:"is_blocked"`
}

type RoomManager struct {
	rooms   map[string]*Room
	db      database.Database
	mu      sync.RWMutex
}

const (
	RoleMember    = "member"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

func NewRoomManager(db database.Database) *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
		db:    db,
	}
}

func NewRoom(name, description string, isPrivate bool, creatorID string) *Room {
	roomID := uuid.New().String()
	inviteCode := generateInviteCode()
	
	room := &Room{
		ID:          roomID,
		Name:        name,
		Description: description,
		InviteCode:  inviteCode,
		IsPrivate:   isPrivate,
		Members:     make(map[string]*Member),
		Moderators:  make(map[string]bool),
		Messages:    make([]*Message, 0),
		CreatedAt:   time.Now(),
	}
	
	room.Moderators[creatorID] = true
	return room
}

func generateInviteCode() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	encoded := base58.Encode(bytes)
	if len(encoded) > 16 {
		return encoded[:16]
	}
	return encoded
}

func (rm *RoomManager) CreateRoom(name, description string, isPrivate bool, creatorID, creatorUsername, creatorPublicKey string) (*Room, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	room := NewRoom(name, description, isPrivate, creatorID)
	
	creator := &Member{
		UserID:    creatorID,
		Username:  creatorUsername,
		PublicKey: creatorPublicKey,
		JoinedAt:  time.Now(),
		Role:      RoleAdmin,
		IsBlocked: false,
	}
	
	room.Members[creatorID] = creator
	rm.rooms[room.ID] = room
	
	dbRoom := &database.Room{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		InviteCode:  room.InviteCode,
		IsPrivate:   room.IsPrivate,
		CreatedAt:   room.CreatedAt,
	}
	
	if err := rm.db.SaveRoom(dbRoom); err != nil {
		return nil, err
	}
	
	if err := rm.db.AddRoomParticipant(room.ID, creatorID); err != nil {
		return nil, err
	}
	
	return room, nil
}

func (rm *RoomManager) GetRoom(roomID string) (*Room, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	if room, exists := rm.rooms[roomID]; exists {
		return room, nil
	}
	
	dbRoom, err := rm.db.GetRoom(roomID)
	if err != nil {
		return nil, err
	}
	
	room := &Room{
		ID:          dbRoom.ID,
		Name:        dbRoom.Name,
		Description: dbRoom.Description,
		InviteCode:  dbRoom.InviteCode,
		IsPrivate:   dbRoom.IsPrivate,
		Members:     make(map[string]*Member),
		Moderators:  make(map[string]bool),
		Messages:    make([]*Message, 0),
		CreatedAt:   dbRoom.CreatedAt,
	}
	
	participants, err := rm.db.GetRoomParticipants(roomID)
	if err != nil {
		return nil, err
	}
	
	for _, userID := range participants {
		user, err := rm.db.GetUser(userID)
		if err != nil {
			continue
		}
		
		member := &Member{
			UserID:    user.ID,
			Username:  user.Username,
			PublicKey: user.PublicKey,
			JoinedAt:  user.CreatedAt,
			Role:      RoleMember,
			IsBlocked: user.IsBlocked,
		}
		
		room.Members[userID] = member
	}
	
	rm.rooms[roomID] = room
	return room, nil
}

func (rm *RoomManager) GetRoomByInviteCode(inviteCode string) (*Room, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	for _, room := range rm.rooms {
		if room.InviteCode == inviteCode {
			return room, nil
		}
	}
	
	rooms, err := rm.db.GetRooms()
	if err != nil {
		return nil, err
	}
	
	for _, dbRoom := range rooms {
		if dbRoom.InviteCode == inviteCode {
			return rm.GetRoom(dbRoom.ID)
		}
	}
	
	return nil, errors.New("room not found")
}

func (r *Room) AddMember(userID, username, publicKey string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.Members[userID]; exists {
		return errors.New("user already in room")
	}
	
	member := &Member{
		UserID:    userID,
		Username:  username,
		PublicKey: publicKey,
		JoinedAt:  time.Now(),
		Role:      RoleMember,
		IsBlocked: false,
	}
	
	r.Members[userID] = member
	return nil
}

func (r *Room) RemoveMember(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	delete(r.Members, userID)
	delete(r.Moderators, userID)
	return nil
}

func (r *Room) IsMember(userID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	_, exists := r.Members[userID]
	return exists
}

func (r *Room) IsModerator(userID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.Moderators[userID]
}

func (r *Room) PromoteToModerator(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if member, exists := r.Members[userID]; exists {
		member.Role = RoleModerator
		r.Moderators[userID] = true
		return nil
	}
	
	return errors.New("user not found in room")
}

func (r *Room) DemoteFromModerator(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if member, exists := r.Members[userID]; exists {
		member.Role = RoleMember
		delete(r.Moderators, userID)
		return nil
	}
	
	return errors.New("user not found in room")
}

func (r *Room) BlockUser(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if member, exists := r.Members[userID]; exists {
		member.IsBlocked = true
		return nil
	}
	
	return errors.New("user not found in room")
}

func (r *Room) UnblockUser(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if member, exists := r.Members[userID]; exists {
		member.IsBlocked = false
		return nil
	}
	
	return errors.New("user not found in room")
}

func (r *Room) IsUserBlocked(userID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if member, exists := r.Members[userID]; exists {
		return member.IsBlocked
	}
	
	return false
}

func (r *Room) GetMembersList() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	members := make([]string, 0, len(r.Members))
	for userID := range r.Members {
		members = append(members, userID)
	}
	
	return members
}

func (r *Room) GetMembersInfo() []Member {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	members := make([]Member, 0, len(r.Members))
	for _, member := range r.Members {
		members = append(members, *member)
	}
	
	return members
}

func (rm *RoomManager) JoinRoomByInvite(inviteCode, userID, username, publicKey string) (*Room, error) {
	room, err := rm.GetRoomByInviteCode(inviteCode)
	if err != nil {
		return nil, err
	}
	
	if err := room.AddMember(userID, username, publicKey); err != nil {
		return nil, err
	}
	
	if err := rm.db.AddRoomParticipant(room.ID, userID); err != nil {
		return nil, err
	}
	
	return room, nil
}

func (rm *RoomManager) LeaveRoom(roomID, userID string) error {
	room, err := rm.GetRoom(roomID)
	if err != nil {
		return err
	}
	
	if err := room.RemoveMember(userID); err != nil {
		return err
	}
	
	return rm.db.RemoveRoomParticipant(roomID, userID)
}

// AddMessage adds a message to the room's in-memory message list (for testing compatibility)
func (r *Room) AddMessage(msg *Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.Messages == nil {
		r.Messages = make([]*Message, 0)
	}
	
	r.Messages = append(r.Messages, msg)
	return nil
} 