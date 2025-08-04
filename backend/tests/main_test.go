package main

import (
	"testing"
	"time"
)

// TODO: Implement comprehensive test suite
// TODO: Implement integration tests
// TODO: Implement performance tests
// TODO: Implement security tests

func TestNewNode(t *testing.T) {
	node := NewNode("test-node", "test-address")
	
	if node.ID != "test-node" {
		t.Errorf("Expected node ID to be 'test-node', got '%s'", node.ID)
	}
	
	if node.Address != "test-address" {
		t.Errorf("Expected node address to be 'test-address', got '%s'", node.Address)
	}
	
	if len(node.Peers) != 0 {
		t.Errorf("Expected empty peers map, got %d peers", len(node.Peers))
	}
	
	if len(node.Rooms) != 0 {
		t.Errorf("Expected empty rooms map, got %d rooms", len(node.Rooms))
	}
}

func TestNewRoom(t *testing.T) {
	room := NewRoom("test-room", "Test Room", "A test room")
	
	if room.ID != "test-room" {
		t.Errorf("Expected room ID to be 'test-room', got '%s'", room.ID)
	}
	
	if room.Name != "Test Room" {
		t.Errorf("Expected room name to be 'Test Room', got '%s'", room.Name)
	}
	
	if room.Description != "A test room" {
		t.Errorf("Expected room description to be 'A test room', got '%s'", room.Description)
	}
	
	if len(room.Members) != 0 {
		t.Errorf("Expected empty members map, got %d members", len(room.Members))
	}
	
	if len(room.Messages) != 0 {
		t.Errorf("Expected empty messages slice, got %d messages", len(room.Messages))
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage("test-room", "test-user", "TestUser", "Hello, world!", "text")
	
	if msg.RoomID != "test-room" {
		t.Errorf("Expected room ID to be 'test-room', got '%s'", msg.RoomID)
	}
	
	if msg.UserID != "test-user" {
		t.Errorf("Expected user ID to be 'test-user', got '%s'", msg.UserID)
	}
	
	if msg.Username != "TestUser" {
		t.Errorf("Expected username to be 'TestUser', got '%s'", msg.Username)
	}
	
	if msg.Content != "Hello, world!" {
		t.Errorf("Expected content to be 'Hello, world!', got '%s'", msg.Content)
	}
	
	if msg.Type != "text" {
		t.Errorf("Expected type to be 'text', got '%s'", msg.Type)
	}
	
	if msg.Encrypted {
		t.Error("Expected message to not be encrypted initially")
	}
}

func TestRoomAddMember(t *testing.T) {
	room := NewRoom("test-room", "Test Room", "A test room")
	
	err := room.AddMember("user1", "User1")
	if err != nil {
		t.Errorf("Failed to add member: %v", err)
	}
	
	if len(room.Members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(room.Members))
	}
	
	member, exists := room.Members["user1"]
	if !exists {
		t.Error("Member not found in room")
	}
	
	if member.UserID != "user1" {
		t.Errorf("Expected member user ID to be 'user1', got '%s'", member.UserID)
	}
	
	if member.Username != "User1" {
		t.Errorf("Expected member username to be 'User1', got '%s'", member.Username)
	}
	
	if member.Role != "member" {
		t.Errorf("Expected member role to be 'member', got '%s'", member.Role)
	}
}

func TestRoomRemoveMember(t *testing.T) {
	room := NewRoom("test-room", "Test Room", "A test room")
	room.AddMember("user1", "User1")
	
	err := room.RemoveMember("user1")
	if err != nil {
		t.Errorf("Failed to remove member: %v", err)
	}
	
	if len(room.Members) != 0 {
		t.Errorf("Expected 0 members after removal, got %d", len(room.Members))
	}
}

func TestRoomAddMessage(t *testing.T) {
	room := NewRoom("test-room", "Test Room", "A test room")
	msg := NewMessage("test-room", "user1", "User1", "Hello!", "text")
	
	err := room.AddMessage(msg)
	if err != nil {
		t.Errorf("Failed to add message: %v", err)
	}
	
	if len(room.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(room.Messages))
	}
	
	if room.Messages[0].Content != "Hello!" {
		t.Errorf("Expected message content to be 'Hello!', got '%s'", room.Messages[0].Content)
	}
} 