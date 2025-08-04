package database

import (
	"database/sql"
	"fmt"
	"time"
	"errors"
	_ "modernc.org/sqlite"
	"ripcord/types"
)

// Room represents a chat room in the database layer
type Room struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	InviteCode  string    `json:"invite_code" db:"invite_code"`
	IsPrivate   bool      `json:"is_private" db:"is_private"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Participants []string  `json:"participants,omitempty"`
}

type Database interface {
	Connect() error
	Disconnect() error
	SaveMessage(msg *types.Message) error
	GetMessages(roomID string, limit int) ([]*types.Message, error)
	SaveRoom(room *Room) error
	GetRoom(roomID string) (*Room, error)
	GetRooms() ([]*Room, error)
	SaveUser(user *types.User) error
	GetUser(userID string) (*types.User, error)
	AddRoomParticipant(roomID, userID string) error
	RemoveRoomParticipant(roomID, userID string) error
	GetRoomParticipants(roomID string) ([]string, error)
	SaveSettings(key, value string) error
	GetSettings(key string) (string, error)
}

type SQLiteDatabase struct {
	db     *sql.DB
	dbPath string
}

func NewSQLiteDatabase(dbPath string) *SQLiteDatabase {
	return &SQLiteDatabase{
		dbPath: dbPath,
	}
}

func (sdb *SQLiteDatabase) Connect() error {
	var err error
	sdb.db, err = sql.Open("sqlite", sdb.dbPath)
	if err != nil {
		return err
	}
	
	return sdb.createTables()
}

func (sdb *SQLiteDatabase) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS rooms (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			invite_code TEXT UNIQUE,
			is_private BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			public_key TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_blocked BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			room_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			type TEXT DEFAULT 'text',
			encrypted BOOLEAN DEFAULT FALSE,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			signature TEXT,
			FOREIGN KEY (room_id) REFERENCES rooms(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS room_participants (
			room_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (room_id, user_id),
			FOREIGN KEY (room_id) REFERENCES rooms(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
	}
	
	for _, query := range queries {
		if _, err := sdb.db.Exec(query); err != nil {
			return err
		}
	}
	
	// Create indexes for better performance
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_rooms_created_at ON rooms(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_users_public_key ON users(public_key)`,
		`CREATE INDEX IF NOT EXISTS idx_room_participants_room_id ON room_participants(room_id)`,
		`CREATE INDEX IF NOT EXISTS idx_room_participants_user_id ON room_participants(user_id)`,
	}
	
	for _, query := range indexQueries {
		if _, err := sdb.db.Exec(query); err != nil {
			return err
		}
	}
	
	return nil
}

func (sdb *SQLiteDatabase) Disconnect() error {
	if sdb.db != nil {
		return sdb.db.Close()
	}
	return nil
}

func (sdb *SQLiteDatabase) SaveMessage(msg *types.Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	
	if msg.ID == "" || msg.RoomID == "" || msg.UserID == "" {
		return errors.New("message missing required fields")
	}
	
	query := `INSERT OR REPLACE INTO messages (id, room_id, user_id, username, content, type, encrypted, timestamp, signature)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := sdb.db.Exec(query, msg.ID, msg.RoomID, msg.UserID, msg.Username, 
		msg.Content, msg.Type, msg.Encrypted, msg.Timestamp, msg.Signature)
	if err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}
	return nil
}

func (sdb *SQLiteDatabase) GetMessages(roomID string, limit int) ([]*types.Message, error) {
	if roomID == "" {
		return nil, errors.New("room ID is required")
	}
	
	if limit <= 0 {
		limit = 50 // Default limit
	}
	
	query := `SELECT id, room_id, user_id, username, content, type, encrypted, timestamp, signature
			  FROM messages WHERE room_id = ? ORDER BY timestamp DESC LIMIT ?`
	
	rows, err := sdb.db.Query(query, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %v", err)
	}
	defer rows.Close()
	
	var messages []*types.Message
	for rows.Next() {
		msg := &types.Message{}
		err := rows.Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.Username,
			&msg.Content, &msg.Type, &msg.Encrypted, &msg.Timestamp, &msg.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}
		messages = append(messages, msg)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %v", err)
	}
	
	return messages, nil
}

func (sdb *SQLiteDatabase) SaveRoom(room *Room) error {
	query := `INSERT OR REPLACE INTO rooms (id, name, description, invite_code, is_private, created_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := sdb.db.Exec(query, room.ID, room.Name, room.Description, 
		room.InviteCode, room.IsPrivate, room.CreatedAt)
	return err
}

func (sdb *SQLiteDatabase) GetRoom(roomID string) (*Room, error) {
	query := `SELECT id, name, description, invite_code, is_private, created_at
			  FROM rooms WHERE id = ?`
	
	room := &Room{}
	err := sdb.db.QueryRow(query, roomID).Scan(&room.ID, &room.Name, 
		&room.Description, &room.InviteCode, &room.IsPrivate, &room.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("room not found")
	}
	
	return room, err
}

func (sdb *SQLiteDatabase) GetRooms() ([]*Room, error) {
	query := `SELECT id, name, description, invite_code, is_private, created_at
			  FROM rooms ORDER BY created_at DESC`
	
	rows, err := sdb.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rooms []*Room
	for rows.Next() {
		room := &Room{}
		err := rows.Scan(&room.ID, &room.Name, &room.Description, 
			&room.InviteCode, &room.IsPrivate, &room.CreatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	
	return rooms, nil
}

func (sdb *SQLiteDatabase) SaveUser(user *types.User) error {
	query := `INSERT OR REPLACE INTO users (id, username, public_key, created_at, last_seen, is_blocked)
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := sdb.db.Exec(query, user.ID, user.Username, user.PublicKey,
		user.CreatedAt, user.LastSeen, user.IsBlocked)
	return err
}

func (sdb *SQLiteDatabase) GetUser(userID string) (*types.User, error) {
	query := `SELECT id, username, public_key, created_at, last_seen, is_blocked
			  FROM users WHERE id = ?`
	
	user := &types.User{}
	err := sdb.db.QueryRow(query, userID).Scan(&user.ID, &user.Username,
		&user.PublicKey, &user.CreatedAt, &user.LastSeen, &user.IsBlocked)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	
	return user, err
}

func (sdb *SQLiteDatabase) AddRoomParticipant(roomID, userID string) error {
	query := `INSERT OR IGNORE INTO room_participants (room_id, user_id) VALUES (?, ?)`
	_, err := sdb.db.Exec(query, roomID, userID)
	return err
}

func (sdb *SQLiteDatabase) RemoveRoomParticipant(roomID, userID string) error {
	query := `DELETE FROM room_participants WHERE room_id = ? AND user_id = ?`
	_, err := sdb.db.Exec(query, roomID, userID)
	return err
}

func (sdb *SQLiteDatabase) GetRoomParticipants(roomID string) ([]string, error) {
	query := `SELECT user_id FROM room_participants WHERE room_id = ?`
	
	rows, err := sdb.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var participants []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}
	
	return participants, nil
}

func (sdb *SQLiteDatabase) SaveSettings(key, value string) error {
	query := `INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)`
	_, err := sdb.db.Exec(query, key, value)
	return err
}

func (sdb *SQLiteDatabase) GetSettings(key string) (string, error) {
	query := `SELECT value FROM settings WHERE key = ?`
	
	var value string
	err := sdb.db.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", errors.New("setting not found")
	}
	
	return value, err
} 