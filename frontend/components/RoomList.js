// TODO: Implement room filtering and search
// TODO: Implement room categories and favorites
// TODO: Implement room notifications and unread counts
// TODO: Implement room sorting options

class RoomList {
    constructor() {
        this.rooms = [];
        this.currentRoomId = null;
        this.roomListContainer = document.getElementById('room-list');
        this.init();
    }
    
    init() {
        this.bindEvents();
    }
    
    bindEvents() {
        // TODO: Add search functionality
        // TODO: Add room filtering options
    }
    
    updateRooms(rooms) {
        this.rooms = rooms;
        this.renderRooms();
    }
    
    renderRooms() {
        this.roomListContainer.innerHTML = '';
        
        if (this.rooms.length === 0) {
            this.renderEmptyState();
            return;
        }
        
        this.rooms.forEach(room => {
            const roomElement = this.createRoomElement(room);
            this.roomListContainer.appendChild(roomElement);
        });
    }
    
    createRoomElement(room) {
        const roomDiv = document.createElement('div');
        roomDiv.className = 'room-item';
        roomDiv.dataset.roomId = room.id;
        
        if (this.currentRoomId === room.id) {
            roomDiv.classList.add('active');
        }
        
        const nameDiv = document.createElement('div');
        nameDiv.className = 'room-name';
        nameDiv.textContent = room.name;
        
        const descDiv = document.createElement('div');
        descDiv.className = 'room-description';
        descDiv.textContent = room.description || 'No description';
        
        const metaDiv = document.createElement('div');
        metaDiv.className = 'room-meta';
        metaDiv.innerHTML = `
            <span class="member-count">${room.member_count || 0} members</span>
            ${room.unread_count ? `<span class="unread-count">${room.unread_count}</span>` : ''}
        `;
        
        roomDiv.appendChild(nameDiv);
        roomDiv.appendChild(descDiv);
        roomDiv.appendChild(metaDiv);
        
        // Add click event
        roomDiv.addEventListener('click', () => {
            this.selectRoom(room.id);
        });
        
        return roomDiv;
    }
    
    renderEmptyState() {
        const emptyDiv = document.createElement('div');
        emptyDiv.className = 'empty-state';
        emptyDiv.innerHTML = `
            <div class="empty-icon">💬</div>
            <p>No rooms available</p>
            <p class="empty-subtitle">Create a room to start chatting</p>
        `;
        this.roomListContainer.appendChild(emptyDiv);
    }
    
    selectRoom(roomId) {
        // Remove active class from all rooms
        this.roomListContainer.querySelectorAll('.room-item').forEach(item => {
            item.classList.remove('active');
        });
        
        // Add active class to selected room
        const selectedRoom = this.roomListContainer.querySelector(`[data-room-id="${roomId}"]`);
        if (selectedRoom) {
            selectedRoom.classList.add('active');
        }
        
        this.currentRoomId = roomId;
        
        // Join the room
        if (window.ripcordApp) {
            window.ripcordApp.joinRoom(roomId);
        }
    }
    
    addRoom(room) {
        this.rooms.push(room);
        const roomElement = this.createRoomElement(room);
        this.roomListContainer.appendChild(roomElement);
    }
    
    removeRoom(roomId) {
        // Remove from array
        this.rooms = this.rooms.filter(room => room.id !== roomId);
        
        // Remove from DOM
        const roomElement = this.roomListContainer.querySelector(`[data-room-id="${roomId}"]`);
        if (roomElement) {
            roomElement.remove();
        }
        
        // If this was the current room, clear selection
        if (this.currentRoomId === roomId) {
            this.currentRoomId = null;
        }
    }
    
    updateRoom(room) {
        // Update in array
        const index = this.rooms.findIndex(r => r.id === room.id);
        if (index !== -1) {
            this.rooms[index] = room;
        }
        
        // Update in DOM
        const roomElement = this.roomListContainer.querySelector(`[data-room-id="${room.id}"]`);
        if (roomElement) {
            const newRoomElement = this.createRoomElement(room);
            roomElement.replaceWith(newRoomElement);
        }
    }
    
    setUnreadCount(roomId, count) {
        const room = this.rooms.find(r => r.id === roomId);
        if (room) {
            room.unread_count = count;
            this.updateRoom(room);
        }
    }
    
    clearUnreadCount(roomId) {
        this.setUnreadCount(roomId, 0);
    }
    
    searchRooms(query) {
        const filteredRooms = this.rooms.filter(room => 
            room.name.toLowerCase().includes(query.toLowerCase()) ||
            (room.description && room.description.toLowerCase().includes(query.toLowerCase()))
        );
        
        this.renderFilteredRooms(filteredRooms);
    }
    
    renderFilteredRooms(rooms) {
        this.roomListContainer.innerHTML = '';
        
        if (rooms.length === 0) {
            this.renderEmptyState();
            return;
        }
        
        rooms.forEach(room => {
            const roomElement = this.createRoomElement(room);
            this.roomListContainer.appendChild(roomElement);
        });
    }
    
    sortRooms(sortBy = 'name') {
        switch (sortBy) {
            case 'name':
                this.rooms.sort((a, b) => a.name.localeCompare(b.name));
                break;
            case 'activity':
                this.rooms.sort((a, b) => (b.last_activity || 0) - (a.last_activity || 0));
                break;
            case 'members':
                this.rooms.sort((a, b) => (b.member_count || 0) - (a.member_count || 0));
                break;
            case 'unread':
                this.rooms.sort((a, b) => (b.unread_count || 0) - (a.unread_count || 0));
                break;
        }
        
        this.renderRooms();
    }
    
    getCurrentRoom() {
        return this.rooms.find(room => room.id === this.currentRoomId);
    }
    
    getRoomById(roomId) {
        return this.rooms.find(room => room.id === roomId);
    }
    
    getRoomsByCategory(category) {
        return this.rooms.filter(room => room.category === category);
    }
    
    // Handle room updates from WebSocket
    handleRoomUpdate(room) {
        const existingRoom = this.rooms.find(r => r.id === room.id);
        if (existingRoom) {
            this.updateRoom(room);
        } else {
            this.addRoom(room);
        }
    }
    
    handleRoomRemoved(roomId) {
        this.removeRoom(roomId);
    }
    
    // Export for testing
    static createInstance() {
        return new RoomList();
    }
} 