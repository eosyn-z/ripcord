// TODO: Implement WebSocket connection to backend
// TODO: Implement real-time message handling
// TODO: Implement room management
// TODO: Implement user authentication
// TODO: Implement message encryption/decryption

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
    
    connectToBackend() {
        // TODO: Implement WebSocket connection
        const wsUrl = `ws://${window.location.host}/ws`;
        
        try {
            this.websocket = new WebSocket(wsUrl);
            
            this.websocket.onopen = () => {
                this.updateConnectionStatus('connected', 'Connected');
                this.authenticateUser();
            };
            
            this.websocket.onmessage = (event) => {
                this.handleWebSocketMessage(JSON.parse(event.data));
            };
            
            this.websocket.onclose = () => {
                this.updateConnectionStatus('disconnected', 'Disconnected');
                // Attempt to reconnect after 5 seconds
                setTimeout(() => this.connectToBackend(), 5000);
            };
            
            this.websocket.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.updateConnectionStatus('error', 'Connection Error');
            };
            
        } catch (error) {
            console.error('Failed to connect to backend:', error);
            this.updateConnectionStatus('error', 'Connection Failed');
        }
    }
    
    handleWebSocketMessage(data) {
        switch (data.type) {
            case 'auth_response':
                this.handleAuthResponse(data);
                break;
            case 'room_list':
                this.handleRoomList(data);
                break;
            case 'user_list':
                this.handleUserList(data);
                break;
            case 'message':
                this.handleNewMessage(data);
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
    
    sendMessage() {
        const input = document.getElementById('message-input');
        const content = input.value.trim();
        
        if (!content || !this.currentRoom) {
            return;
        }
        
        const message = {
            type: 'message',
            room_id: this.currentRoom.id,
            content: content,
            timestamp: new Date().toISOString()
        };
        
        this.sendWebSocketMessage(message);
        input.value = '';
    }
    
    createRoom() {
        const nameInput = document.getElementById('room-name');
        const descInput = document.getElementById('room-description');
        
        const roomData = {
            type: 'create_room',
            name: nameInput.value.trim(),
            description: descInput.value.trim()
        };
        
        this.sendWebSocketMessage(roomData);
        this.hideCreateRoomModal();
        
        // Clear form
        nameInput.value = '';
        descInput.value = '';
    }
    
    joinRoom(roomId) {
        const message = {
            type: 'join_room',
            room_id: roomId
        };
        
        this.sendWebSocketMessage(message);
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
        // TODO: Load user preferences from localStorage
    }
    
    loadRooms() {
        const message = {
            type: 'get_rooms'
        };
        this.sendWebSocketMessage(message);
    }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.ripcordApp = new RipcordApp();
}); 