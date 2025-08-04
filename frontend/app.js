// Ripcord Frontend Application
// Handles real-time chat functionality with WebSocket connections

class RipcordApp {
    constructor() {
        this.currentRoom = null;
        this.currentUser = null;
        this.rooms = new Map();
        this.users = new Map();
        this.websocket = null;
        this.components = {};
        
        this.init();
    }
    
    init() {
        this.initializeComponents();
        this.bindEvents();
        this.connectToBackend();
        this.loadUserPreferences();
    }
    
    initializeComponents() {
        // Initialize all UI components
        this.components.chatPane = new ChatPane();
        this.components.roomList = new RoomList();
        this.components.userList = new UserList();
        this.components.inputBar = new InputBar();
        this.components.settingsPanel = new SettingsPanel();
        
        // Set up component interconnections
        this.components.roomList.onRoomSelect = (roomId) => {
            this.selectRoom(roomId);
        };
    }
    
    bindEvents() {
        // Bind DOM events
        document.getElementById('create-room-btn').addEventListener('click', () => {
            this.showCreateRoomModal();
        });
        
        document.getElementById('send-btn').addEventListener('click', () => {
            this.sendMessage();
        });
        
        document.getElementById('message-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });
        
        document.getElementById('create-room-form').addEventListener('submit', (e) => {
            e.preventDefault();
            this.createRoom();
        });
        
        document.getElementById('cancel-create-room').addEventListener('click', () => {
            this.hideCreateRoomModal();
        });
        
        document.getElementById('room-settings-btn').addEventListener('click', () => {
            this.showSettingsPanel();
        });
    }
    
    async connectToBackend() {
        // First get identity via HTTP API
        try {
            const response = await fetch('/api/identity');
            if (response.ok) {
                const identity = await response.json();
                this.currentUser = identity;
                this.loadRooms();
                
                // Then establish WebSocket connection for real-time features
                this.connectWebSocket();
            } else {
                throw new Error('Failed to get identity');
            }
        } catch (error) {
            console.error('Failed to connect to backend:', error);
            this.updateConnectionStatus('error', 'Connection Failed');
            // Retry after 5 seconds
            setTimeout(() => this.connectToBackend(), 5000);
        }
    }
    
    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        this.websocket = new WebSocket(wsUrl);
        
        this.websocket.onopen = () => {
            console.log('WebSocket connected');
            this.updateConnectionStatus('connected', 'Connected');
            
            // Authenticate with WebSocket
            this.sendWebSocketMessage({
                type: 'auth',
                username: this.currentUser.username || 'Anonymous'
            });
        };
        
        this.websocket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleWebSocketMessage(data);
            } catch (error) {
                console.error('Failed to parse WebSocket message:', error);
            }
        };
        
        this.websocket.onclose = () => {
            console.log('WebSocket disconnected');
            this.updateConnectionStatus('error', 'Disconnected');
            
            // Attempt to reconnect after 3 seconds
            setTimeout(() => {
                if (this.websocket.readyState === WebSocket.CLOSED) {
                    this.connectWebSocket();
                }
            }, 3000);
        };
        
        this.websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.updateConnectionStatus('error', 'Connection Error');
        };
    }
    
    handleWebSocketMessage(data) {
        switch (data.type) {
            case 'auth_response':
                this.handleAuthResponse(data);
                break;
            case 'message':
                this.handleNewMessage(data);
                break;
            case 'message_history':
                this.handleMessageHistory(data);
                break;
            case 'room_joined':
                this.handleRoomJoined(data);
                break;
            case 'room_left':
                this.handleRoomLeft(data);
                break;
            case 'user_joined':
                this.handleUserJoined(data);
                break;
            case 'user_left':
                this.handleUserLeft(data);
                break;
            default:
                console.warn('Unknown message type:', data.type);
        }
    }
    
    authenticateUser() {
        // TODO: Implement user authentication
        const authMessage = {
            type: 'auth',
            username: this.getStoredUsername() || 'Anonymous',
            publicKey: this.getStoredPublicKey()
        };
        
        this.sendWebSocketMessage(authMessage);
    }
    
    sendWebSocketMessage(message) {
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify(message));
        } else {
            console.error('WebSocket not connected');
        }
    }
    
    async sendMessage() {
        const input = document.getElementById('message-input');
        const content = input.value.trim();
        
        if (!content || !this.currentRoom) {
            return;
        }
        
        // Sanitize content
        const sanitizedContent = this.sanitizeMessageContent(content);
        if (!sanitizedContent) {
            return;
        }
        
        // Try WebSocket first, fall back to HTTP API
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.sendWebSocketMessage({
                type: 'send_message',
                content: sanitizedContent
            });
            
            // Clear input immediately for better UX
            input.value = '';
            if (this.components.inputBar) {
                this.components.inputBar.clearInput();
            }
        } else {
            // Fallback to HTTP API
            try {
                const response = await fetch('/api/messages/send', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        room_id: this.currentRoom.id,
                        content: content
                    })
                });
                
                if (response.ok) {
                    const message = await response.json();
                    this.components.chatPane.addMessage(message);
                    input.value = '';
                    if (this.components.inputBar) {
                        this.components.inputBar.clearInput();
                    }
                } else {
                    console.error('Failed to send message');
                }
            } catch (error) {
                console.error('Error sending message:', error);
            }
        }
    }
    
    async createRoom() {
        const nameInput = document.getElementById('room-name');
        const descInput = document.getElementById('room-description');
        
        try {
            const response = await fetch('/api/rooms/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name: nameInput.value.trim(),
                    description: descInput.value.trim(),
                    is_private: false
                })
            });
            
            if (response.ok) {
                const room = await response.json();
                this.rooms.set(room.id, room);
                this.components.roomList.addRoom(room);
                this.hideCreateRoomModal();
                
                // Clear form
                nameInput.value = '';
                descInput.value = '';
                
                // Auto-join the newly created room
                this.selectRoom(room.id);
            } else {
                console.error('Failed to create room');
            }
        } catch (error) {
            console.error('Error creating room:', error);
        }
    }
    
    selectRoom(roomId) {
        const room = this.rooms.get(roomId);
        if (room) {
            this.currentRoom = room;
            this.updateCurrentRoomDisplay();
            this.components.chatPane.clearMessages();
            this.components.roomList.setActiveRoom(roomId);
            
            // Join room via WebSocket for real-time updates
            if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
                this.sendWebSocketMessage({
                    type: 'join_room',
                    room_id: roomId
                });
                
                // Request message history
                this.sendWebSocketMessage({
                    type: 'get_messages',
                    room_id: roomId,
                    limit: 50
                });
            } else {
                // Fallback to HTTP API
                this.loadMessages(roomId);
            }
        }
    }
    
    async loadMessages(roomId) {
        try {
            const response = await fetch(`/api/messages?room_id=${roomId}`);
            if (response.ok) {
                const messages = await response.json();
                messages.forEach(message => {
                    this.components.chatPane.addMessage(message);
                });
            }
        } catch (error) {
            console.error('Error loading messages:', error);
        }
    }
    
    leaveRoom() {
        if (!this.currentRoom) return;
        
        const message = {
            type: 'leave_room',
            room_id: this.currentRoom.id
        };
        
        this.sendWebSocketMessage(message);
        this.currentRoom = null;
        this.updateCurrentRoomDisplay();
    }
    
    // Event handlers
    handleAuthResponse(data) {
        if (data.success) {
            this.currentUser = data.user;
            this.storeUserData(data.user);
            this.loadRooms();
        } else {
            console.error('Authentication failed:', data.error);
        }
    }
    
    handleRoomList(data) {
        this.rooms.clear();
        data.rooms.forEach(room => {
            this.rooms.set(room.id, room);
        });
        this.components.roomList.updateRooms(Array.from(this.rooms.values()));
    }
    
    handleUserList(data) {
        this.users.clear();
        data.users.forEach(user => {
            this.users.set(user.id, user);
        });
        this.components.userList.updateUsers(Array.from(this.users.values()));
    }
    
    handleNewMessage(data) {
        this.components.chatPane.addMessage(data.message);
    }
    
    handleMessageHistory(data) {
        if (data.messages && Array.isArray(data.messages)) {
            data.messages.forEach(message => {
                this.components.chatPane.addMessage(message);
            });
        }
    }
    
    handleRoomJoined(data) {
        this.currentRoom = data.room;
        this.updateCurrentRoomDisplay();
        this.components.chatPane.clearMessages();
        this.components.chatPane.loadMessageHistory(data.room.id);
    }
    
    handleRoomLeft(data) {
        if (this.currentRoom && this.currentRoom.id === data.room_id) {
            this.currentRoom = null;
            this.updateCurrentRoomDisplay();
        }
    }
    
    handleUserJoined(data) {
        this.users.set(data.user.id, data.user);
        this.components.userList.addUser(data.user);
    }
    
    handleUserLeft(data) {
        this.users.delete(data.user_id);
        this.components.userList.removeUser(data.user_id);
    }
    
    // UI helpers
    updateConnectionStatus(status, text) {
        const indicator = document.getElementById('status-indicator');
        const statusText = document.getElementById('status-text');
        
        indicator.className = `status-indicator ${status}`;
        statusText.textContent = text;
    }
    
    updateCurrentRoomDisplay() {
        const roomNameElement = document.getElementById('current-room-name');
        if (this.currentRoom) {
            roomNameElement.textContent = this.currentRoom.name;
        } else {
            roomNameElement.textContent = 'Select a room to start chatting';
        }
    }
    
    showCreateRoomModal() {
        document.getElementById('create-room-modal').classList.remove('hidden');
    }
    
    hideCreateRoomModal() {
        document.getElementById('create-room-modal').classList.add('hidden');
    }
    
    showSettingsPanel() {
        this.components.settingsPanel.show();
    }
    
    // Storage helpers
    storeUserData(user) {
        localStorage.setItem('ripcord_username', user.username);
        if (user.publicKey) {
            localStorage.setItem('ripcord_public_key', user.publicKey);
        }
    }
    
    getStoredUsername() {
        return localStorage.getItem('ripcord_username');
    }
    
    getStoredPublicKey() {
        return localStorage.getItem('ripcord_public_key');
    }
    
    loadUserPreferences() {
        // Load user preferences from localStorage
        const username = this.getStoredUsername();
        if (username) {
            this.currentUser = { username };
        }
    }
    
    // Input sanitization
    sanitizeMessageContent(content) {
        // Remove HTML tags and potentially dangerous content
        const div = document.createElement('div');
        div.textContent = content;
        let sanitized = div.textContent || div.innerText || '';
        
        // Remove null bytes and control characters
        sanitized = sanitized.replace(/[\x00-\x1F\x7F]/g, '');
        
        // Trim whitespace
        sanitized = sanitized.trim();
        
        // Limit consecutive whitespace
        sanitized = sanitized.replace(/\s+/g, ' ');
        
        return sanitized;
    }
    
    // Error handling utilities
    showError(message) {
        console.error(message);
        // You could implement a toast notification system here
        alert(message);
    }
    
    showSuccess(message) {
        console.log(message);
        // You could implement a toast notification system here
    }
    
    async loadRooms() {
        try {
            const response = await fetch('/api/rooms');
            if (response.ok) {
                const rooms = await response.json();
                this.rooms.clear();
                rooms.forEach(room => {
                    this.rooms.set(room.id, room);
                });
                this.components.roomList.updateRooms(Array.from(this.rooms.values()));
            }
        } catch (error) {
            console.error('Error loading rooms:', error);
        }
    }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.ripcordApp = new RipcordApp();
}); 