package database

import (
	"database/sql"
	"time"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	Connect() error
	Disconnect() error
	SaveMessage(msg *Message) error
	GetMessages(roomID string, limit int) ([]*Message, error)
	SaveRoom(room *Room) error
	GetRoom(roomID string) (*Room, error)
	GetRooms() ([]*Room, error)
	SaveUser(user *User) error
	GetUser(userID string) (*User, error)
	AddRoomParticipant(roomID, userID string) error
	RemoveRoomParticipant(roomID, userID string) error
	GetRoomParticipants(roomID string) ([]string, error)
	SaveSettings(key, value string) error
	GetSettings(key string) (string, error)
}

type Message struct {
	ID        string    `db:"id"`
	RoomID    string    `db:"room_id"`
	UserID    string    `db:"user_id"`
	Username  string    `db:"username"`
	Content   string    `db:"content"`
	Type      string    `db:"type"`
	Encrypted bool      `db:"encrypted"`
	Timestamp time.Time `db:"timestamp"`
	Signature string    `db:"signature"`
}

type Room struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	InviteCode  string    `db:"invite_code"`
	IsPrivate   bool      `db:"is_private"`
	CreatedAt   time.Time `db:"created_at"`
}

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	PublicKey string    `db:"public_key"`
	CreatedAt time.Time `db:"created_at"`
	LastSeen  time.Time `db:"last_seen"`
	IsBlocked bool      `db:"is_blocked"`
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
	sdb.db, err = sql.Open("sqlite3", sdb.dbPath)
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
	
	return nil
}

func (sdb *SQLiteDatabase) Disconnect() error {
	if sdb.db != nil {
		return sdb.db.Close()
	}
	return nil
}

func (sdb *SQLiteDatabase) SaveMessage(msg *Message) error {
	query := `INSERT OR REPLACE INTO messages (id, room_id, user_id, username, content, type, encrypted, timestamp, signature)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := sdb.db.Exec(query, msg.ID, msg.RoomID, msg.UserID, msg.Username, 
		msg.Content, msg.Type, msg.Encrypted, msg.Timestamp, msg.Signature)
	return err
}

func (sdb *SQLiteDatabase) GetMessages(roomID string, limit int) ([]*Message, error) {
	query := `SELECT id, room_id, user_id, username, content, type, encrypted, timestamp, signature
			  FROM messages WHERE room_id = ? ORDER BY timestamp DESC LIMIT ?`
	
	rows, err := sdb.db.Query(query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.Username,
			&msg.Content, &msg.Type, &msg.Encrypted, &msg.Timestamp, &msg.Signature)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
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

func (sdb *SQLiteDatabase) SaveUser(user *User) error {
	query := `INSERT OR REPLACE INTO users (id, username, public_key, created_at, last_seen, is_blocked)
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := sdb.db.Exec(query, user.ID, user.Username, user.PublicKey,
		user.CreatedAt, user.LastSeen, user.IsBlocked)
	return err
}

func (sdb *SQLiteDatabase) GetUser(userID string) (*User, error) {
	query := `SELECT id, username, public_key, created_at, last_seen, is_blocked
			  FROM users WHERE id = ?`
	
	user := &User{}
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