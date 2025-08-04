// TODO: Implement message formatting and markdown support
// TODO: Implement file upload functionality
// TODO: Implement emoji picker
// TODO: Implement message drafts and auto-save

class InputBar {
    constructor() {
        this.inputElement = document.getElementById('message-input');
        this.sendButton = document.getElementById('send-btn');
        this.currentDraft = '';
        this.autoSaveInterval = null;
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.startAutoSave();
    }
    
    bindEvents() {
        // Handle input changes
        this.inputElement.addEventListener('input', (e) => {
            this.handleInputChange(e);
        });
        
        // Handle key events
        this.inputElement.addEventListener('keydown', (e) => {
            this.handleKeyDown(e);
        });
        
        // Handle send button click
        this.sendButton.addEventListener('click', () => {
            this.sendMessage();
        });
        
        // Handle paste events for file uploads
        this.inputElement.addEventListener('paste', (e) => {
            this.handlePaste(e);
        });
        
        // Handle drag and drop for files
        this.inputElement.addEventListener('dragover', (e) => {
            this.handleDragOver(e);
        });
        
        this.inputElement.addEventListener('drop', (e) => {
            this.handleDrop(e);
        });
    }
    
    handleInputChange(event) {
        const value = event.target.value;
        this.currentDraft = value;
        
        // Enable/disable send button based on content
        this.updateSendButtonState(value.trim().length > 0);
        
        // Auto-resize textarea
        this.autoResize();
    }
    
    handleKeyDown(event) {
        // Send message on Enter (without Shift)
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            this.sendMessage();
        }
        
        // Handle Tab for indentation
        if (event.key === 'Tab') {
            event.preventDefault();
            this.insertAtCursor('\t');
        }
        
        // Handle Ctrl+B for bold
        if (event.ctrlKey && event.key === 'b') {
            event.preventDefault();
            this.wrapSelection('**', '**');
        }
        
        // Handle Ctrl+I for italic
        if (event.ctrlKey && event.key === 'i') {
            event.preventDefault();
            this.wrapSelection('*', '*');
        }
        
        // Handle Ctrl+K for code
        if (event.ctrlKey && event.key === 'k') {
            event.preventDefault();
            this.wrapSelection('`', '`');
        }
    }
    
    handlePaste(event) {
        const items = event.clipboardData.items;
        
        for (let item of items) {
            if (item.type.indexOf('image') !== -1) {
                event.preventDefault();
                this.handleImagePaste(item);
            }
        }
    }
    
    handleDragOver(event) {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'copy';
        this.inputElement.classList.add('drag-over');
    }
    
    handleDrop(event) {
        event.preventDefault();
        this.inputElement.classList.remove('drag-over');
        
        const files = event.dataTransfer.files;
        if (files.length > 0) {
            this.handleFileUpload(files);
        }
    }
    
    sendMessage() {
        const content = this.inputElement.value.trim();
        
        if (!content) {
            return;
        }
        
        // Send the message through the main app
        if (window.ripcordApp) {
            window.ripcordApp.sendMessage();
        }
        // Note: Input will be cleared by the main app after successful sending
    }
    
    clearInput() {
        this.inputElement.value = '';
        this.currentDraft = '';
        this.updateSendButtonState(false);
        this.autoResize();
    }
    
    updateSendButtonState(enabled) {
        this.sendButton.disabled = !enabled;
    }
    
    autoResize() {
        // Reset height to auto to get the correct scrollHeight
        this.inputElement.style.height = 'auto';
        
        // Set the height to scrollHeight
        const newHeight = Math.min(this.inputElement.scrollHeight, 120);
        this.inputElement.style.height = newHeight + 'px';
    }
    
    insertAtCursor(text) {
        const start = this.inputElement.selectionStart;
        const end = this.inputElement.selectionEnd;
        const value = this.inputElement.value;
        
        this.inputElement.value = value.substring(0, start) + text + value.substring(end);
        
        // Set cursor position after inserted text
        this.inputElement.selectionStart = this.inputElement.selectionEnd = start + text.length;
        
        // Trigger input event
        this.inputElement.dispatchEvent(new Event('input'));
    }
    
    wrapSelection(before, after) {
        const start = this.inputElement.selectionStart;
        const end = this.inputElement.selectionEnd;
        const value = this.inputElement.value;
        const selectedText = value.substring(start, end);
        
        if (selectedText) {
            this.inputElement.value = value.substring(0, start) + before + selectedText + after + value.substring(end);
            this.inputElement.selectionStart = start + before.length;
            this.inputElement.selectionEnd = end + before.length;
        } else {
            this.insertAtCursor(before + after);
            this.inputElement.selectionStart = this.inputElement.selectionEnd = start + before.length;
        }
        
        // Trigger input event
        this.inputElement.dispatchEvent(new Event('input'));
    }
    
    handleImagePaste(item) {
        const file = item.getAsFile();
        if (file) {
            this.uploadFile(file);
        }
    }
    
    handleFileUpload(files) {
        for (let file of files) {
            this.uploadFile(file);
        }
    }
    
    uploadFile(file) {
        // TODO: Implement file upload to backend
        console.log('Uploading file:', file.name);
        
        // For now, just insert the filename
        this.insertAtCursor(`[File: ${file.name}]`);
    }
    
    startAutoSave() {
        // Auto-save draft every 5 seconds
        this.autoSaveInterval = setInterval(() => {
            this.saveDraft();
        }, 5000);
    }
    
    stopAutoSave() {
        if (this.autoSaveInterval) {
            clearInterval(this.autoSaveInterval);
            this.autoSaveInterval = null;
        }
    }
    
    saveDraft() {
        if (this.currentDraft.trim()) {
            const roomId = window.ripcordApp?.currentRoom?.id;
            if (roomId) {
                localStorage.setItem(`draft_${roomId}`, this.currentDraft);
            }
        }
    }
    
    loadDraft(roomId) {
        const draft = localStorage.getItem(`draft_${roomId}`);
        if (draft) {
            this.inputElement.value = draft;
            this.currentDraft = draft;
            this.updateSendButtonState(draft.trim().length > 0);
            this.autoResize();
        }
    }
    
    clearDraft(roomId) {
        localStorage.removeItem(`draft_${roomId}`);
    }
    
    focus() {
        this.inputElement.focus();
    }
    
    blur() {
        this.inputElement.blur();
    }
    
    getValue() {
        return this.inputElement.value;
    }
    
    setValue(value) {
        this.inputElement.value = value;
        this.currentDraft = value;
        this.updateSendButtonState(value.trim().length > 0);
        this.autoResize();
    }
    
    isFocused() {
        return document.activeElement === this.inputElement;
    }
    
    // Handle room changes
    onRoomChange(roomId) {
        // Save current draft
        this.saveDraft();
        
        // Load draft for new room
        this.loadDraft(roomId);
    }
    
    // Cleanup
    destroy() {
        this.stopAutoSave();
        this.saveDraft();
    }
    
    // Export for testing
    static createInstance() {
        return new InputBar();
    }
} 