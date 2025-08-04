// TODO: Implement user status indicators (online, away, busy)
// TODO: Implement user profiles and avatars
// TODO: Implement user search and filtering
// TODO: Implement user actions (message, block, etc.)

class UserList {
    constructor() {
        this.users = [];
        this.userListContainer = document.getElementById('user-list');
        this.init();
    }
    
    init() {
        this.bindEvents();
    }
    
    bindEvents() {
        // TODO: Add user search functionality
        // TODO: Add user filtering options
    }
    
    updateUsers(users) {
        this.users = users;
        this.renderUsers();
    }
    
    renderUsers() {
        this.userListContainer.innerHTML = '';
        
        if (this.users.length === 0) {
            this.renderEmptyState();
            return;
        }
        
        // Sort users by status and name
        const sortedUsers = this.sortUsers(this.users);
        
        sortedUsers.forEach(user => {
            const userElement = this.createUserElement(user);
            this.userListContainer.appendChild(userElement);
        });
    }
    
    createUserElement(user) {
        const userDiv = document.createElement('div');
        userDiv.className = 'user-item';
        userDiv.dataset.userId = user.id;
        
        const avatar = this.createUserAvatar(user);
        const info = this.createUserInfo(user);
        
        userDiv.appendChild(avatar);
        userDiv.appendChild(info);
        
        // Add click event for user actions
        userDiv.addEventListener('click', () => {
            this.showUserActions(user);
        });
        
        // Add right-click context menu
        userDiv.addEventListener('contextmenu', (e) => {
            e.preventDefault();
            this.showUserContextMenu(e, user);
        });
        
        return userDiv;
    }
    
    createUserAvatar(user) {
        const avatarDiv = document.createElement('div');
        avatarDiv.className = 'user-avatar';
        
        if (user.avatar) {
            const img = document.createElement('img');
            img.src = user.avatar;
            img.alt = user.username;
            avatarDiv.appendChild(img);
        } else {
            avatarDiv.textContent = user.username.charAt(0).toUpperCase();
        }
        
        // Add status indicator
        const statusIndicator = document.createElement('div');
        statusIndicator.className = `status-indicator ${user.status || 'offline'}`;
        avatarDiv.appendChild(statusIndicator);
        
        return avatarDiv;
    }
    
    createUserInfo(user) {
        const infoDiv = document.createElement('div');
        infoDiv.className = 'user-info';
        
        const nameDiv = document.createElement('div');
        nameDiv.className = 'user-name';
        nameDiv.textContent = user.username;
        
        const statusDiv = document.createElement('div');
        statusDiv.className = 'user-status';
        statusDiv.textContent = this.getStatusText(user.status);
        
        infoDiv.appendChild(nameDiv);
        infoDiv.appendChild(statusDiv);
        
        return infoDiv;
    }
    
    getStatusText(status) {
        switch (status) {
            case 'online':
                return 'Online';
            case 'away':
                return 'Away';
            case 'busy':
                return 'Busy';
            case 'offline':
                return 'Offline';
            default:
                return 'Unknown';
        }
    }
    
    sortUsers(users) {
        return users.sort((a, b) => {
            // Sort by status priority first
            const statusOrder = { 'online': 0, 'away': 1, 'busy': 2, 'offline': 3 };
            const statusA = statusOrder[a.status] || 3;
            const statusB = statusOrder[b.status] || 3;
            
            if (statusA !== statusB) {
                return statusA - statusB;
            }
            
            // Then sort by username
            return a.username.localeCompare(b.username);
        });
    }
    
    renderEmptyState() {
        const emptyDiv = document.createElement('div');
        emptyDiv.className = 'empty-state';
        emptyDiv.innerHTML = `
            <div class="empty-icon">👥</div>
            <p>No users online</p>
            <p class="empty-subtitle">Users will appear here when they join</p>
        `;
        this.userListContainer.appendChild(emptyDiv);
    }
    
    addUser(user) {
        // Check if user already exists
        const existingUser = this.users.find(u => u.id === user.id);
        if (existingUser) {
            this.updateUser(user);
            return;
        }
        
        this.users.push(user);
        const userElement = this.createUserElement(user);
        this.userListContainer.appendChild(userElement);
    }
    
    removeUser(userId) {
        // Remove from array
        this.users = this.users.filter(user => user.id !== userId);
        
        // Remove from DOM
        const userElement = this.userListContainer.querySelector(`[data-user-id="${userId}"]`);
        if (userElement) {
            userElement.remove();
        }
    }
    
    updateUser(user) {
        // Update in array
        const index = this.users.findIndex(u => u.id === user.id);
        if (index !== -1) {
            this.users[index] = user;
        }
        
        // Update in DOM
        const userElement = this.userListContainer.querySelector(`[data-user-id="${user.id}"]`);
        if (userElement) {
            const newUserElement = this.createUserElement(user);
            userElement.replaceWith(newUserElement);
        }
    }
    
    setUserStatus(userId, status) {
        const user = this.users.find(u => u.id === userId);
        if (user) {
            user.status = status;
            this.updateUser(user);
        }
    }
    
    showUserActions(user) {
        // TODO: Implement user action menu
        console.log('User clicked:', user.username);
    }
    
    showUserContextMenu(event, user) {
        // TODO: Implement context menu
        console.log('Context menu for user:', user.username);
    }
    
    searchUsers(query) {
        const filteredUsers = this.users.filter(user => 
            user.username.toLowerCase().includes(query.toLowerCase()) ||
            (user.display_name && user.display_name.toLowerCase().includes(query.toLowerCase()))
        );
        
        this.renderFilteredUsers(filteredUsers);
    }
    
    renderFilteredUsers(users) {
        this.userListContainer.innerHTML = '';
        
        if (users.length === 0) {
            this.renderEmptyState();
            return;
        }
        
        const sortedUsers = this.sortUsers(users);
        sortedUsers.forEach(user => {
            const userElement = this.createUserElement(user);
            this.userListContainer.appendChild(userElement);
        });
    }
    
    filterUsersByStatus(status) {
        const filteredUsers = this.users.filter(user => user.status === status);
        this.renderFilteredUsers(filteredUsers);
    }
    
    getUserById(userId) {
        return this.users.find(user => user.id === userId);
    }
    
    getOnlineUsers() {
        return this.users.filter(user => user.status === 'online');
    }
    
    getOfflineUsers() {
        return this.users.filter(user => user.status === 'offline');
    }
    
    getUserCount() {
        return this.users.length;
    }
    
    getOnlineUserCount() {
        return this.getOnlineUsers().length;
    }
    
    // Handle user updates from WebSocket
    handleUserUpdate(user) {
        const existingUser = this.users.find(u => u.id === user.id);
        if (existingUser) {
            this.updateUser(user);
        } else {
            this.addUser(user);
        }
    }
    
    handleUserRemoved(userId) {
        this.removeUser(userId);
    }
    
    // Export for testing
    static createInstance() {
        return new UserList();
    }
} 