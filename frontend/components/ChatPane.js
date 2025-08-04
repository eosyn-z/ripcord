// TODO: Implement message rendering with proper formatting
// TODO: Implement message timestamps and user avatars
// TODO: Implement message search functionality
// TODO: Implement message reactions and emoji support

class ChatPane {
    constructor() {
        this.messages = [];
        this.currentRoomId = null;
        this.messagesContainer = document.getElementById('chat-messages');
        this.init();
    }
    
    init() {
        this.bindEvents();
    }
    
    bindEvents() {
        // TODO: Bind scroll events for infinite loading
        this.messagesContainer.addEventListener('scroll', (e) => {
            this.handleScroll(e);
        });
    }
    
    addMessage(message) {
        this.messages.push(message);
        this.renderMessage(message);
        this.scrollToBottom();
    }
    
    renderMessage(message) {
        const messageElement = this.createMessageElement(message);
        this.messagesContainer.appendChild(messageElement);
    }
    
    createMessageElement(message) {
        const messageDiv = document.createElement('div');
        messageDiv.className = 'message';
        
        // Check if this is the current user's message
        const isOwnMessage = message.user_id === window.ripcordApp?.currentUser?.id;
        if (isOwnMessage) {
            messageDiv.classList.add('own-message');
        }
        
        const avatar = this.createAvatar(message.username);
        const content = this.createMessageContent(message);
        
        messageDiv.appendChild(avatar);
        messageDiv.appendChild(content);
        
        return messageDiv;
    }
    
    createAvatar(username) {
        const avatarDiv = document.createElement('div');
        avatarDiv.className = 'message-avatar';
        avatarDiv.textContent = username.charAt(0).toUpperCase();
        return avatarDiv;
    }
    
    createMessageContent(message) {
        const contentDiv = document.createElement('div');
        contentDiv.className = 'message-content';
        
        const header = document.createElement('div');
        header.className = 'message-header';
        
        const username = document.createElement('span');
        username.className = 'message-username';
        username.textContent = message.username;
        
        const time = document.createElement('span');
        time.className = 'message-time';
        time.textContent = this.formatTimestamp(message.timestamp);
        
        header.appendChild(username);
        header.appendChild(time);
        
        const text = document.createElement('div');
        text.className = 'message-text';
        text.textContent = message.content;
        
        contentDiv.appendChild(header);
        contentDiv.appendChild(text);
        
        return contentDiv;
    }
    
    formatTimestamp(timestamp) {
        const date = new Date(timestamp);
        const now = new Date();
        const diffInHours = (now - date) / (1000 * 60 * 60);
        
        if (diffInHours < 24) {
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        } else {
            return date.toLocaleDateString();
        }
    }
    
    clearMessages() {
        this.messages = [];
        this.messagesContainer.innerHTML = '';
    }
    
    loadMessageHistory(roomId) {
        // TODO: Load message history from backend
        this.currentRoomId = roomId;
        
        // Request message history from backend
        if (window.ripcordApp) {
            const message = {
                type: 'get_messages',
                room_id: roomId,
                limit: 50
            };
            window.ripcordApp.sendWebSocketMessage(message);
        }
    }
    
    handleScroll(event) {
        const { scrollTop } = event.target;
        
        // Load more messages when scrolling to top
        if (scrollTop === 0 && this.messages.length > 0) {
            this.loadMoreMessages();
        }
    }
    
    loadMoreMessages() {
        // TODO: Implement pagination for message history
        if (window.ripcordApp && this.currentRoomId) {
            const oldestMessageId = this.messages[0]?.id;
            const message = {
                type: 'get_messages',
                room_id: this.currentRoomId,
                before_id: oldestMessageId,
                limit: 20
            };
            window.ripcordApp.sendWebSocketMessage(message);
        }
    }
    
    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }
    
    scrollToMessage(messageId) {
        const messageElement = this.messagesContainer.querySelector(`[data-message-id="${messageId}"]`);
        if (messageElement) {
            messageElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    }
    
    searchMessages(query) {
        // TODO: Implement message search functionality
        const results = this.messages.filter(message => 
            message.content.toLowerCase().includes(query.toLowerCase()) ||
            message.username.toLowerCase().includes(query.toLowerCase())
        );
        
        return results;
    }
    
    highlightMessage(messageId) {
        // Remove previous highlights
        this.messagesContainer.querySelectorAll('.message.highlighted').forEach(el => {
            el.classList.remove('highlighted');
        });
        
        // Add highlight to target message
        const messageElement = this.messagesContainer.querySelector(`[data-message-id="${messageId}"]`);
        if (messageElement) {
            messageElement.classList.add('highlighted');
            this.scrollToMessage(messageId);
        }
    }
    
    // Handle incoming message history
    handleMessageHistory(messages) {
        // Clear existing messages if this is a fresh load
        if (this.messages.length === 0) {
            this.clearMessages();
        }
        
        // Add messages to the beginning for pagination
        messages.reverse().forEach(message => {
            this.messages.unshift(message);
            const messageElement = this.createMessageElement(message);
            this.messagesContainer.insertBefore(messageElement, this.messagesContainer.firstChild);
        });
    }
    
    // Handle new message from WebSocket
    handleNewMessage(message) {
        this.addMessage(message);
    }
    
    // Utility methods
    getMessageCount() {
        return this.messages.length;
    }
    
    getCurrentRoomId() {
        return this.currentRoomId;
    }
    
    // Export for testing
    static createInstance() {
        return new ChatPane();
    }
} 