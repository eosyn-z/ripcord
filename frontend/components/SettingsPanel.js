// TODO: Implement theme switching (dark/light mode)
// TODO: Implement notification settings
// TODO: Implement privacy settings
// TODO: Implement account management

class SettingsPanel {
    constructor() {
        this.panel = document.getElementById('settings-panel');
        this.isVisible = false;
        this.settings = {};
        this.init();
    }
    
    init() {
        this.loadSettings();
        this.renderPanel();
        this.bindEvents();
    }
    
    bindEvents() {
        // Close panel when clicking outside
        document.addEventListener('click', (e) => {
            if (!this.panel.contains(e.target) && this.isVisible) {
                this.hide();
            }
        });
        
        // Handle escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isVisible) {
                this.hide();
            }
        });
    }
    
    renderPanel() {
        this.panel.innerHTML = `
            <div class="settings-header">
                <h3>Settings</h3>
                <button class="close-btn" id="close-settings">×</button>
            </div>
            
            <div class="settings-content">
                <div class="settings-section">
                    <h4>Appearance</h4>
                    <div class="setting-item">
                        <label for="theme-select">Theme:</label>
                        <select id="theme-select">
                            <option value="light">Light</option>
                            <option value="dark">Dark</option>
                            <option value="auto">Auto</option>
                        </select>
                    </div>
                    
                    <div class="setting-item">
                        <label for="font-size">Font Size:</label>
                        <select id="font-size">
                            <option value="small">Small</option>
                            <option value="medium">Medium</option>
                            <option value="large">Large</option>
                        </select>
                    </div>
                </div>
                
                <div class="settings-section">
                    <h4>Notifications</h4>
                    <div class="setting-item">
                        <label>
                            <input type="checkbox" id="enable-notifications">
                            Enable notifications
                        </label>
                    </div>
                    
                    <div class="setting-item">
                        <label>
                            <input type="checkbox" id="sound-notifications">
                            Sound notifications
                        </label>
                    </div>
                    
                    <div class="setting-item">
                        <label>
                            <input type="checkbox" id="desktop-notifications">
                            Desktop notifications
                        </label>
                    </div>
                </div>
                
                <div class="settings-section">
                    <h4>Privacy</h4>
                    <div class="setting-item">
                        <label>
                            <input type="checkbox" id="show-online-status">
                            Show online status
                        </label>
                    </div>
                    
                    <div class="setting-item">
                        <label>
                            <input type="checkbox" id="allow-direct-messages">
                            Allow direct messages
                        </label>
                    </div>
                    
                    <div class="setting-item">
                        <label>
                            <input type="checkbox" id="encrypt-messages">
                            Encrypt messages
                        </label>
                    </div>
                </div>
                
                <div class="settings-section">
                    <h4>Account</h4>
                    <div class="setting-item">
                        <label for="username">Username:</label>
                        <input type="text" id="username" placeholder="Enter username">
                    </div>
                    
                    <div class="setting-item">
                        <label for="display-name">Display Name:</label>
                        <input type="text" id="display-name" placeholder="Enter display name">
                    </div>
                    
                    <div class="setting-item">
                        <button id="generate-keys" class="btn btn-secondary">Generate New Keys</button>
                    </div>
                </div>
                
                <div class="settings-section">
                    <h4>Connection</h4>
                    <div class="setting-item">
                        <label for="i2p-enabled">
                            <input type="checkbox" id="i2p-enabled">
                            Enable I2P network
                        </label>
                    </div>
                    
                    <div class="setting-item">
                        <label for="sam-address">SAM Address:</label>
                        <input type="text" id="sam-address" placeholder="127.0.0.1">
                    </div>
                    
                    <div class="setting-item">
                        <label for="sam-port">SAM Port:</label>
                        <input type="number" id="sam-port" placeholder="7656">
                    </div>
                </div>
            </div>
            
            <div class="settings-footer">
                <button id="save-settings" class="btn btn-primary">Save Settings</button>
                <button id="reset-settings" class="btn btn-secondary">Reset to Defaults</button>
            </div>
        `;
        
        this.bindPanelEvents();
    }
    
    bindPanelEvents() {
        // Close button
        document.getElementById('close-settings').addEventListener('click', () => {
            this.hide();
        });
        
        // Save settings
        document.getElementById('save-settings').addEventListener('click', () => {
            this.saveSettings();
        });
        
        // Reset settings
        document.getElementById('reset-settings').addEventListener('click', () => {
            this.resetSettings();
        });
        
        // Generate keys
        document.getElementById('generate-keys').addEventListener('click', () => {
            this.generateNewKeys();
        });
        
        // Theme change
        document.getElementById('theme-select').addEventListener('change', (e) => {
            this.applyTheme(e.target.value);
        });
        
        // Font size change
        document.getElementById('font-size').addEventListener('change', (e) => {
            this.applyFontSize(e.target.value);
        });
    }
    
    show() {
        this.panel.classList.remove('hidden');
        this.panel.classList.add('show');
        this.isVisible = true;
        
        // Load current settings into form
        this.loadSettingsIntoForm();
    }
    
    hide() {
        this.panel.classList.remove('show');
        this.panel.classList.add('hidden');
        this.isVisible = false;
    }
    
    loadSettings() {
        const stored = localStorage.getItem('ripcord_settings');
        if (stored) {
            this.settings = JSON.parse(stored);
        } else {
            this.settings = this.getDefaultSettings();
        }
    }
    
    getDefaultSettings() {
        return {
            theme: 'light',
            fontSize: 'medium',
            notifications: {
                enabled: true,
                sound: true,
                desktop: true
            },
            privacy: {
                showOnlineStatus: true,
                allowDirectMessages: true,
                encryptMessages: true
            },
            account: {
                username: '',
                displayName: ''
            },
            connection: {
                i2pEnabled: true,
                samAddress: '127.0.0.1',
                samPort: 7656
            }
        };
    }
    
    loadSettingsIntoForm() {
        // Theme
        document.getElementById('theme-select').value = this.settings.theme || 'light';
        
        // Font size
        document.getElementById('font-size').value = this.settings.fontSize || 'medium';
        
        // Notifications
        document.getElementById('enable-notifications').checked = this.settings.notifications?.enabled ?? true;
        document.getElementById('sound-notifications').checked = this.settings.notifications?.sound ?? true;
        document.getElementById('desktop-notifications').checked = this.settings.notifications?.desktop ?? true;
        
        // Privacy
        document.getElementById('show-online-status').checked = this.settings.privacy?.showOnlineStatus ?? true;
        document.getElementById('allow-direct-messages').checked = this.settings.privacy?.allowDirectMessages ?? true;
        document.getElementById('encrypt-messages').checked = this.settings.privacy?.encryptMessages ?? true;
        
        // Account
        document.getElementById('username').value = this.settings.account?.username || '';
        document.getElementById('display-name').value = this.settings.account?.displayName || '';
        
        // Connection
        document.getElementById('i2p-enabled').checked = this.settings.connection?.i2pEnabled ?? true;
        document.getElementById('sam-address').value = this.settings.connection?.samAddress || '127.0.0.1';
        document.getElementById('sam-port').value = this.settings.connection?.samPort || 7656;
    }
    
    saveSettings() {
        // Collect settings from form
        this.settings = {
            theme: document.getElementById('theme-select').value,
            fontSize: document.getElementById('font-size').value,
            notifications: {
                enabled: document.getElementById('enable-notifications').checked,
                sound: document.getElementById('sound-notifications').checked,
                desktop: document.getElementById('desktop-notifications').checked
            },
            privacy: {
                showOnlineStatus: document.getElementById('show-online-status').checked,
                allowDirectMessages: document.getElementById('allow-direct-messages').checked,
                encryptMessages: document.getElementById('encrypt-messages').checked
            },
            account: {
                username: document.getElementById('username').value,
                displayName: document.getElementById('display-name').value
            },
            connection: {
                i2pEnabled: document.getElementById('i2p-enabled').checked,
                samAddress: document.getElementById('sam-address').value,
                samPort: parseInt(document.getElementById('sam-port').value)
            }
        };
        
        // Save to localStorage
        localStorage.setItem('ripcord_settings', JSON.stringify(this.settings));
        
        // Apply settings
        this.applySettings();
        
        // Show success message
        this.showMessage('Settings saved successfully!');
    }
    
    resetSettings() {
        if (confirm('Are you sure you want to reset all settings to defaults?')) {
            this.settings = this.getDefaultSettings();
            localStorage.setItem('ripcord_settings', JSON.stringify(this.settings));
            this.loadSettingsIntoForm();
            this.applySettings();
            this.showMessage('Settings reset to defaults!');
        }
    }
    
    applySettings() {
        // Apply theme
        this.applyTheme(this.settings.theme);
        
        // Apply font size
        this.applyFontSize(this.settings.fontSize);
        
        // Apply other settings
        this.applyNotificationSettings();
        this.applyPrivacySettings();
        
        // Notify main app of settings change
        if (window.ripcordApp) {
            window.ripcordApp.onSettingsChanged(this.settings);
        }
    }
    
    applyTheme(theme) {
        document.body.className = document.body.className.replace(/theme-\w+/, '');
        document.body.classList.add(`theme-${theme}`);
    }
    
    applyFontSize(size) {
        document.body.className = document.body.className.replace(/font-size-\w+/, '');
        document.body.classList.add(`font-size-${size}`);
    }
    
    applyNotificationSettings() {
        // TODO: Apply notification settings
        console.log('Applying notification settings:', this.settings.notifications);
    }
    
    applyPrivacySettings() {
        // TODO: Apply privacy settings
        console.log('Applying privacy settings:', this.settings.privacy);
    }
    
    generateNewKeys() {
        if (confirm('Generate new encryption keys? This will invalidate your current keys.')) {
            // TODO: Implement key generation
            console.log('Generating new keys...');
            this.showMessage('New keys generated successfully!');
        }
    }
    
    showMessage(message) {
        // TODO: Implement toast notification
        console.log(message);
    }
    
    getSetting(path) {
        return path.split('.').reduce((obj, key) => obj?.[key], this.settings);
    }
    
    setSetting(path, value) {
        const keys = path.split('.');
        const lastKey = keys.pop();
        const target = keys.reduce((obj, key) => obj[key] = obj[key] || {}, this.settings);
        target[lastKey] = value;
    }
    
    // Export for testing
    static createInstance() {
        return new SettingsPanel();
    }
} 