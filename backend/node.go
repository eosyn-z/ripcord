package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
	"ripcord/security"
)

type Node struct {
	ID             string
	cryptoManager  *security.CryptoManager
	roomManager    *RoomManager
	messageHandler *MessageHandler
	peers          map[string]*Peer
	isRunning      bool
	mu             sync.RWMutex
	startTime      time.Time
}

type Peer struct {
	ID        string
	PublicKey string
	Nickname  string
	Address   string
	LastSeen  time.Time
	Status    string
	IsBlocked bool
}

const (
	PeerStatusConnected    = "connected"
	PeerStatusDisconnected = "disconnected"
	PeerStatusBlocked      = "blocked"
)

func NewNode(cryptoManager *security.CryptoManager, roomManager *RoomManager, messageHandler *MessageHandler) *Node {
	nodeID := hex.EncodeToString(cryptoManager.GetPublicKey())
	
	return &Node{
		ID:             nodeID,
		cryptoManager:  cryptoManager,
		roomManager:    roomManager,
		messageHandler: messageHandler,
		peers:          make(map[string]*Peer),
		isRunning:      false,
		startTime:      time.Now(),
	}
}

func (n *Node) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	if n.isRunning {
		return fmt.Errorf("node already running")
	}
	
	log.Println("Starting node:", n.ID[:16]+"...")
	
	n.isRunning = true
	
	go n.heartbeatLoop()
	
	return nil
}

func (n *Node) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	if !n.isRunning {
		return nil
	}
	
	log.Println("Stopping node...")
	n.isRunning = false
	
	for _, peer := range n.peers {
		peer.Status = PeerStatusDisconnected
	}
	
	return nil
}

func (n *Node) IsRunning() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.isRunning
}

func (n *Node) AddPeer(publicKey, nickname, address string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	if publicKey == hex.EncodeToString(n.cryptoManager.GetPublicKey()) {
		return fmt.Errorf("cannot add self as peer")
	}
	
	peer := &Peer{
		ID:        publicKey,
		PublicKey: publicKey,
		Nickname:  nickname,
		Address:   address,
		LastSeen:  time.Now(),
		Status:    PeerStatusConnected,
		IsBlocked: false,
	}
	
	n.peers[publicKey] = peer
	log.Printf("Added peer: %s (%s)", nickname, publicKey[:16]+"...")
	
	return nil
}

func (n *Node) RemovePeer(publicKey string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	if peer, exists := n.peers[publicKey]; exists {
		peer.Status = PeerStatusDisconnected
		delete(n.peers, publicKey)
		log.Printf("Removed peer: %s", publicKey[:16]+"...")
		return nil
	}
	
	return fmt.Errorf("peer not found")
}

func (n *Node) BlockPeer(publicKey string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	if peer, exists := n.peers[publicKey]; exists {
		peer.IsBlocked = true
		peer.Status = PeerStatusBlocked
		log.Printf("Blocked peer: %s", publicKey[:16]+"...")
		return nil
	}
	
	return fmt.Errorf("peer not found")
}

func (n *Node) UnblockPeer(publicKey string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	if peer, exists := n.peers[publicKey]; exists {
		peer.IsBlocked = false
		peer.Status = PeerStatusConnected
		log.Printf("Unblocked peer: %s", publicKey[:16]+"...")
		return nil
	}
	
	return fmt.Errorf("peer not found")
}

func (n *Node) GetPeers() []Peer {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	peers := make([]Peer, 0, len(n.peers))
	for _, peer := range n.peers {
		peers = append(peers, *peer)
	}
	
	return peers
}

func (n *Node) BroadcastMessage(msg *ProtocolMessage) error {
	privateKey := n.cryptoManager.GetPrivateKey()
	if err := msg.Sign(privateKey); err != nil {
		return err
	}
	
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	
	n.mu.RLock()
	activePeers := make([]*Peer, 0)
	for _, peer := range n.peers {
		if peer.Status == PeerStatusConnected && !peer.IsBlocked {
			activePeers = append(activePeers, peer)
		}
	}
	n.mu.RUnlock()
	
	for _, peer := range activePeers {
		go func(p *Peer) {
			if err := n.sendToPeer(p, data); err != nil {
				log.Printf("Failed to send message to peer %s: %v", p.ID[:16]+"...", err)
			}
		}(peer)
	}
	
	return nil
}

func (n *Node) sendToPeer(peer *Peer, data []byte) error {
	// TODO: Implement actual network communication via I2P
	log.Printf("Sending message to peer %s (address: %s)", peer.Nickname, peer.Address)
	return nil
}

func (n *Node) ProcessIncomingMessage(data []byte, fromPeer string) error {
	msg, err := ParseProtocolMessage(data)
	if err != nil {
		return err
	}
	
	n.mu.RLock()
	peer, exists := n.peers[fromPeer]
	n.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("message from unknown peer: %s", fromPeer[:16]+"...")
	}
	
	if peer.IsBlocked {
		log.Printf("Ignoring message from blocked peer: %s", fromPeer[:16]+"...")
		return nil
	}
	
	// TODO: Verify message signature
	// TODO: Process different message types
	log.Printf("Processing %s message from %s", msg.Type, peer.Nickname)
	
	return nil
}

func (n *Node) heartbeatLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if !n.IsRunning() {
				return
			}
			n.sendHeartbeat()
		}
	}
}

func (n *Node) sendHeartbeat() {
	msg := NewProtocolMessage(MessageTypeHeartbeat, n.ID, generateMessageID())
	
	payload := HeartbeatPayload{
		Nickname:    n.cryptoManager.GetNickname(),
		PublicKey:   hex.EncodeToString(n.cryptoManager.GetPublicKey()),
		I2PAddress:  "localhost", // TODO: Get actual I2P address
		ActiveRooms: []string{},  // TODO: Get active rooms
	}
	
	msg.SetPayload(payload)
	
	if err := n.BroadcastMessage(msg); err != nil {
		log.Printf("Failed to send heartbeat: %v", err)
	}
}

func (n *Node) GetInfo() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	return map[string]interface{}{
		"id":         n.ID,
		"nickname":   n.cryptoManager.GetNickname(),
		"is_running": n.isRunning,
		"peers":      len(n.peers),
		"uptime":     time.Since(n.startTime).String(),
	}
} 