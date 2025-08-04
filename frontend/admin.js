class RipcordNodeAdmin {
    constructor() {
        this.currentSection = 'overview';
        this.nodeInfo = null;
        this.rooms = new Map();
        this.peers = new Map();
        this.settings = {};
        this.refreshInterval = null;
        
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.loadInitialData();
        this.startPeriodicRefresh();
    }
    
    bindEvents() {
        // Navigation
        document.querySelectorAll('.nav-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const section = e.target.dataset.section;
                this.switchSection(section);
            });
        });
        
        // Overview actions
        document.getElementById('refresh-status').addEventListener('click', () => {
            this.refreshNodeStatus();
        });
        
        document.getElementById('restart-node').addEventListener('click', () => {
            this.restartNode();
        });
        
        // Room management
        document.getElementById('create-admin-room').addEventListener('click', () => {
            this.showCreateRoomModal();
        });
        
        document.getElementById('refresh-rooms').addEventListener('click', () => {
            this.loadRooms();
        });
        
        document.getElementById('create-admin-room-form').addEventListener('submit', (e) => {
            e.preventDefault();
            this.createAdminRoom();
        });
        
        document.getElementById('cancel-create-admin-room').addEventListener('click', () => {
            this.hideCreateRoomModal();
        });
        
        // Peer management
        document.getElementById('refresh-peers').addEventListener('click', () => {
            this.loadPeers();
        });
        
        document.getElementById('clear-blocked-peers').addEventListener('click', () => {
            this.clearBlockedPeers();
        });
        
        // Settings
        document.getElementById('save-settings').addEventListener('click', () => {
            this.saveSettings();
        });
        
        document.getElementById('reset-settings').addEventListener('click', () => {
            this.resetSettings();
        });
        
        document.getElementById('reload-config').addEventListener('click', () => {
            this.reloadConfiguration();
        });
        
        // Modal close events
        document.getElementById('close-room-details').addEventListener('click', () => {
            this.hideRoomDetailsModal();
        });
    }
    
    switchSection(section) {
        // Update navigation
        document.querySelectorAll('.nav-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector(`[data-section="${section}"]`).classList.add('active');
        
        // Update sections
        document.querySelectorAll('.admin-section').forEach(sec => {
            sec.classList.remove('active');
        });
        document.getElementById(`${section}-section`).classList.add('active');
        
        this.currentSection = section;
        
        // Load section-specific data
        switch (section) {
            case 'overview':
                this.loadNodeOverview();
                break;
            case 'rooms':
                this.loadRooms();
                break;
            case 'peers':
                this.loadPeers();
                break;
            case 'settings':
                this.loadSettings();
                break;
        }
    }
    
    async loadInitialData() {
        await this.loadNodeIdentity();
        await this.loadNodeOverview();
    }
    
    async loadNodeIdentity() {
        try {
            const response = await fetch('/api/admin/identity');
            if (response.ok) {
                this.nodeInfo = await response.json();
                this.updateNodeIdentityDisplay();
            } else {
                // Fallback to regular identity endpoint
                const fallbackResponse = await fetch('/api/identity');
                if (fallbackResponse.ok) {
                    this.nodeInfo = await fallbackResponse.json();
                    this.updateNodeIdentityDisplay();
                }
            }
        } catch (error) {
            console.error('Failed to load node identity:', error);
            this.showError('Failed to load node identity');
        }
    }
    
    updateNodeIdentityDisplay() {
        if (!this.nodeInfo) return;
        
        document.getElementById('node-public-key').textContent = 
            this.truncateKey(this.nodeInfo.public_key || 'Not Available');
        document.getElementById('node-fingerprint').textContent = 
            this.nodeInfo.fingerprint || 'Not Available';
        document.getElementById('node-nickname').textContent = 
            this.nodeInfo.nickname || 'Anonymous';
    }
    
    async loadNodeOverview() {
        try {
            // Load I2P status
            const i2pResponse = await fetch('/api/admin/i2p/status');
            if (i2pResponse.ok) {
                const i2pStatus = await i2pResponse.json();
                this.updateI2PStatus(i2pStatus);
            } else {
                this.updateI2PStatus({ status: 'unavailable', address: 'Not Available' });
            }
            
            // Load server statistics
            const statsResponse = await fetch('/api/admin/stats');
            if (statsResponse.ok) {
                const stats = await statsResponse.json();
                this.updateServerStats(stats);
            } else {
                this.updateServerStats({ uptime: 0, rooms: 0, messages: 0, peers: 0 });
            }
            
            this.updateNodeStatus('online', 'Node Online');
        } catch (error) {
            console.error('Failed to load node overview:', error);
            this.updateNodeStatus('error', 'Connection Error');
        }
    }
    
    updateI2PStatus(status) {
        const addressElement = document.getElementById('i2p-address');
        const statusElement = document.getElementById('i2p-status');
        
        addressElement.textContent = status.address || 'Not Available';
        
        statusElement.textContent = status.status || 'Unknown';
        statusElement.className = 'status-badge';
        
        switch (status.status) {
            case 'connected':
                statusElement.classList.add('online');
                break;
            case 'connecting':
                statusElement.classList.add('warning');
                break;
            default:
                statusElement.classList.add('offline');
                break;
        }
    }
    
    updateServerStats(stats) {
        document.getElementById('server-uptime').textContent = 
            this.formatUptime(stats.uptime || 0);
        document.getElementById('active-rooms-count').textContent = 
            stats.rooms || 0;
        document.getElementById('total-messages-count').textContent = 
            stats.messages || 0;
        document.getElementById('active-peers-count').textContent = 
            stats.peers || 0;
    }
    
    updateNodeStatus(status, text) {
        const indicator = document.getElementById('node-status-indicator');
        const statusText = document.getElementById('node-status-text');
        
        indicator.className = `status-indicator ${status}`;
        statusText.textContent = text;
    }
    
    async loadRooms() {
        try {
            const response = await fetch('/api/admin/rooms');
            if (response.ok) {
                const rooms = await response.json();
                this.displayRooms(rooms);
            } else {
                // Fallback to regular rooms endpoint
                const fallbackResponse = await fetch('/api/rooms');
                if (fallbackResponse.ok) {
                    const rooms = await fallbackResponse.json();
                    this.displayRooms(rooms);
                }
            }
        } catch (error) {
            console.error('Failed to load rooms:', error);
            this.showError('Failed to load rooms');
        }
    }
    
    displayRooms(rooms) {
        const container = document.getElementById('admin-rooms-list');
        container.innerHTML = '';
        
        if (rooms.length === 0) {
            container.innerHTML = '<p class="empty-state">No rooms found.</p>';
            return;
        }
        
        rooms.forEach(room => {
            const roomCard = this.createRoomCard(room);
            container.appendChild(roomCard);
        });
    }
    
    createRoomCard(room) {
        const card = document.createElement('div');
        card.className = 'admin-room-card';
        
        card.innerHTML = `
            <div class="room-header">
                <h3 class="room-title">${this.escapeHtml(room.name)}</h3>
                <div class="room-invite-code">${room.invite_code || 'N/A'}</div>
            </div>
            <div class="room-stats">
                <div class="room-stat">
                    <div class="room-stat-value">${room.participant_count || 0}</div>
                    <div class="room-stat-label">Participants</div>
                </div>
                <div class="room-stat">
                    <div class="room-stat-value">${room.message_count || 0}</div>
                    <div class="room-stat-label">Messages</div>
                </div>
                <div class="room-stat">
                    <div class="room-stat-value">${room.is_private ? 'Private' : 'Public'}</div>
                    <div class="room-stat-label">Type</div>
                </div>
            </div>
            <div class="room-actions">
                <button class="btn btn-small btn-primary" onclick="adminApp.viewRoomDetails('${room.id}')">
                    View Details
                </button>
                <button class="btn btn-small btn-secondary" onclick="adminApp.copyInviteCode('${room.invite_code || ''}')">
                    Copy Invite
                </button>
                <button class="btn btn-small btn-warning" onclick="adminApp.kickAllFromRoom('${room.id}')">
                    Kick All
                </button>
                <button class="btn btn-small btn-danger" onclick="adminApp.deleteRoom('${room.id}')">
                    Delete
                </button>
            </div>
        `;
        
        return card;
    }
    
    async loadPeers() {
        try {
            const response = await fetch('/api/admin/peers');
            if (response.ok) {
                const peers = await response.json();
                this.displayPeers(peers);
            } else {
                this.displayPeers([]);
            }
            
            // Load connection logs
            const logsResponse = await fetch('/api/admin/logs/connections');
            if (logsResponse.ok) {
                const logs = await logsResponse.json();
                this.displayConnectionLogs(logs);
            }
        } catch (error) {
            console.error('Failed to load peers:', error);
            this.showError('Failed to load peer information');
        }
    }
    
    displayPeers(peers) {
        const container = document.getElementById('peers-list');
        container.innerHTML = '';
        
        if (peers.length === 0) {
            container.innerHTML = '<p class="empty-state">No peers connected.</p>';
            return;
        }
        
        peers.forEach(peer => {
            const peerCard = this.createPeerCard(peer);
            container.appendChild(peerCard);
        });
    }
    
    createPeerCard(peer) {
        const card = document.createElement('div');
        card.className = 'peer-card';
        
        const statusClass = peer.status === 'connected' ? 'online' : 
                          peer.status === 'connecting' ? 'warning' : 'offline';
        
        card.innerHTML = `
            <div class="peer-header">
                <div class="peer-id">${this.truncateKey(peer.id || 'Unknown')}</div>
                <div class="peer-status">
                    <span class="status-badge ${statusClass}">${peer.status || 'Unknown'}</span>
                </div>
            </div>
            <div class="peer-info">
                <div class="peer-info-item">
                    <div class="peer-info-label">I2P Address</div>
                    <div class="peer-info-value">${this.truncateAddress(peer.i2p_address || 'N/A')}</div>
                </div>
                <div class="peer-info-item">
                    <div class="peer-info-label">Connected Since</div>
                    <div class="peer-info-value">${this.formatTimestamp(peer.connected_at)}</div>
                </div>
                <div class="peer-info-item">
                    <div class="peer-info-label">Messages Sent</div>
                    <div class="peer-info-value">${peer.message_count || 0}</div>
                </div>
                <div class="peer-info-item">
                    <div class="peer-info-label">Trust Level</div>
                    <div class="peer-info-value">${peer.trust_level || 'Unknown'}</div>
                </div>
            </div>
            <div class="peer-actions">
                <button class="btn btn-small btn-secondary" onclick="adminApp.verifyPeer('${peer.id}')">
                    Verify
                </button>
                <button class="btn btn-small btn-warning" onclick="adminApp.blockPeer('${peer.id}')">
                    Block
                </button>
                <button class="btn btn-small btn-danger" onclick="adminApp.disconnectPeer('${peer.id}')">
                    Disconnect
                </button>
            </div>
        `;
        
        return card;
    }
    
    displayConnectionLogs(logs) {
        const container = document.getElementById('connection-logs');
        container.innerHTML = '';
        
        if (logs.length === 0) {
            container.innerHTML = '<p class="empty-state">No connection logs available.</p>';
            return;
        }
        
        logs.forEach(log => {
            const logEntry = document.createElement('div');
            logEntry.className = 'log-entry';
            
            logEntry.innerHTML = `
                <div class="log-timestamp">${this.formatTimestamp(log.timestamp)}</div>
                <div class="log-message">${this.escapeHtml(log.message)}</div>
                <div class="log-level ${log.level || 'info'}">${log.level || 'INFO'}</div>
            `;
            
            container.appendChild(logEntry);
        });
    }
    
    async loadSettings() {
        try {
            const response = await fetch('/api/admin/settings');
            if (response.ok) {
                this.settings = await response.json();
                this.populateSettingsForm();
            } else {
                this.loadDefaultSettings();
            }
        } catch (error) {
            console.error('Failed to load settings:', error);
            this.loadDefaultSettings();
        }
    }
    
    loadDefaultSettings() {
        this.settings = {
            i2p: {
                host: '127.0.0.1',
                port: 7656,
                tunnel_length: 2
            },
            server: {
                port: 8080,
                max_peers: 50,
                message_retention_days: 30
            },
            security: {
                auto_block_malicious: true,
                require_signature_verification: true,
                rate_limit_per_minute: 60
            }
        };
        this.populateSettingsForm();
    }
    
    populateSettingsForm() {
        // I2P settings
        document.getElementById('i2p-host').value = this.settings.i2p?.host || '127.0.0.1';
        document.getElementById('i2p-port').value = this.settings.i2p?.port || 7656;
        document.getElementById('tunnel-length').value = this.settings.i2p?.tunnel_length || 2;
        
        // Server settings
        document.getElementById('server-port').value = this.settings.server?.port || 8080;
        document.getElementById('max-peers').value = this.settings.server?.max_peers || 50;
        document.getElementById('message-retention').value = this.settings.server?.message_retention_days || 30;
        
        // Security settings
        document.getElementById('auto-block-malicious').checked = this.settings.security?.auto_block_malicious || false;
        document.getElementById('require-signature-verification').checked = this.settings.security?.require_signature_verification || true;
        document.getElementById('rate-limit').value = this.settings.security?.rate_limit_per_minute || 60;
    }
    
    // Action methods
    async refreshNodeStatus() {
        await this.loadNodeOverview();
        this.showSuccess('Node status refreshed');
    }
    
    async restartNode() {
        if (!confirm('Are you sure you want to restart the node? This will disconnect all users.')) {
            return;
        }
        
        try {
            const response = await fetch('/api/admin/restart', { method: 'POST' });
            if (response.ok) {
                this.showSuccess('Node restart initiated');
            } else {
                this.showError('Failed to restart node');
            }
        } catch (error) {
            console.error('Error restarting node:', error);
            this.showError('Error restarting node');
        }
    }
    
    async createAdminRoom() {
        const name = document.getElementById('admin-room-name').value.trim();
        const description = document.getElementById('admin-room-description').value.trim();
        const isPrivate = document.getElementById('admin-room-private').checked;
        
        if (!name) {
            this.showError('Room name is required');
            return;
        }
        
        try {
            const response = await fetch('/api/rooms/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: name,
                    description: description,
                    is_private: isPrivate
                })
            });
            
            if (response.ok) {
                this.hideCreateRoomModal();
                this.loadRooms();
                this.showSuccess('Room created successfully');
            } else {
                this.showError('Failed to create room');
            }
        } catch (error) {
            console.error('Error creating room:', error);
            this.showError('Error creating room');
        }
    }
    
    async saveSettings() {
        const settings = {
            i2p: {
                host: document.getElementById('i2p-host').value,
                port: parseInt(document.getElementById('i2p-port').value),
                tunnel_length: parseInt(document.getElementById('tunnel-length').value)
            },
            server: {
                port: parseInt(document.getElementById('server-port').value),
                max_peers: parseInt(document.getElementById('max-peers').value),
                message_retention_days: parseInt(document.getElementById('message-retention').value)
            },
            security: {
                auto_block_malicious: document.getElementById('auto-block-malicious').checked,
                require_signature_verification: document.getElementById('require-signature-verification').checked,
                rate_limit_per_minute: parseInt(document.getElementById('rate-limit').value)
            }
        };
        
        try {
            const response = await fetch('/api/admin/settings', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(settings)
            });
            
            if (response.ok) {
                this.settings = settings;
                this.showSuccess('Settings saved successfully');
            } else {
                this.showError('Failed to save settings');
            }
        } catch (error) {
            console.error('Error saving settings:', error);
            this.showError('Error saving settings');
        }
    }
    
    // Utility methods
    startPeriodicRefresh() {
        this.refreshInterval = setInterval(() => {
            if (this.currentSection === 'overview') {
                this.loadNodeOverview();
            }
        }, 30000); // Refresh every 30 seconds
    }
    
    truncateKey(key) {
        if (!key || key.length <= 16) return key;
        return key.substring(0, 8) + '...' + key.substring(key.length - 8);
    }
    
    truncateAddress(address) {
        if (!address || address.length <= 32) return address;
        return address.substring(0, 16) + '...' + address.substring(address.length - 16);
    }
    
    formatUptime(seconds) {
        const days = Math.floor(seconds / 86400);
        const hours = Math.floor((seconds % 86400) / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        
        if (days > 0) return `${days}d ${hours}h ${minutes}m`;
        if (hours > 0) return `${hours}h ${minutes}m`;
        return `${minutes}m`;
    }
    
    formatTimestamp(timestamp) {
        if (!timestamp) return 'N/A';
        return new Date(timestamp).toLocaleString();
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    
    showCreateRoomModal() {
        document.getElementById('create-admin-room-modal').classList.remove('hidden');
    }
    
    hideCreateRoomModal() {
        document.getElementById('create-admin-room-modal').classList.add('hidden');
        document.getElementById('create-admin-room-form').reset();
    }
    
    showSuccess(message) {
        // TODO: Implement toast notifications
        alert('Success: ' + message);
    }
    
    showError(message) {
        // TODO: Implement toast notifications
        alert('Error: ' + message);
    }
}

// Global functions for onclick handlers
window.copyToClipboard = function(elementId) {
    const element = document.getElementById(elementId);
    const text = element.textContent;
    
    navigator.clipboard.writeText(text).then(() => {
        adminApp.showSuccess('Copied to clipboard');
    }).catch(err => {
        console.error('Failed to copy:', err);
        adminApp.showError('Failed to copy to clipboard');
    });
};

// Initialize the admin app
let adminApp;
document.addEventListener('DOMContentLoaded', () => {
    adminApp = new RipcordNodeAdmin();
});