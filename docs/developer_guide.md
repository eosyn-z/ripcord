# Developer Guide - Ripcord

This guide provides technical details for developers contributing to the Ripcord project.

## Architecture Overview

### Backend Architecture

The backend is built in Go and follows a modular architecture:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Server   │    │  WebSocket Hub  │    │   I2P Manager   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Node Manager  │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │  Room Manager   │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │ Message Handler │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Database      │
                    └─────────────────┘
```

### Frontend Architecture

The frontend uses a component-based architecture:

```
┌─────────────────┐
│   RipcordApp    │
└─────────────────┘
         │
    ┌────┴────┐
    │         │
┌─────────┐ ┌─────────┐
│ChatPane │ │RoomList │
└─────────┘ └─────────┘
    │         │
┌─────────┐ ┌─────────┐
│InputBar │ │UserList │
└─────────┘ └─────────┘
    │
┌─────────┐
│Settings │
└─────────┘
```

## Backend Development

### Project Structure

```
backend/
├── main.go              # Application entry point
├── node.go              # Decentralized node management
├── room.go              # Chat room functionality
├── message.go           # Message handling and encryption
├── protocol.go          # Communication protocol
├── config.go            # Configuration management
├── database/            # Database layer
│   └── db.go           # Database interface
├── security/            # Cryptographic operations
│   └── crypto.go       # Encryption and key management
├── i2p/                 # I2P network integration
│   └── i2p.go          # SAM protocol implementation
└── tests/               # Backend tests
    └── main_test.go     # Unit tests
```

### Key Components

#### Node Management (`node.go`)
- Manages decentralized node discovery
- Handles peer connections and synchronization
- Implements node health monitoring

#### Room System (`room.go`)
- Chat room creation and management
- Member management and permissions
- Room synchronization across nodes

#### Message Handling (`message.go`)
- Message encryption/decryption
- Message signing and verification
- Message routing and delivery

#### I2P Integration (`i2p/i2p.go`)
- SAM (Simple Anonymous Messaging) protocol
- I2P destination management
- Tunnel creation and management

### Database Schema

```sql
-- Users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    public_key BLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Rooms table
CREATE TABLE rooms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,
    content TEXT NOT NULL,
    type TEXT DEFAULT 'text',
    encrypted BOOLEAN DEFAULT FALSE,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    signature BLOB,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Room members table
CREATE TABLE room_members (
    room_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    role TEXT DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id),
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Security Implementation

#### Encryption
- **AES-256-GCM**: Used for message encryption
- **RSA-2048**: Used for key exchange and signing
- **SHA-256**: Used for message hashing

#### Key Management
- Private keys stored locally only
- Public keys shared for encryption
- Key rotation supported
- Secure key generation

### Protocol Specification

#### WebSocket Message Format
```json
{
  "type": "message_type",
  "from": "user_id",
  "to": "user_id", // optional
  "room_id": "room_id", // optional
  "data": "encrypted_data",
  "timestamp": "2023-01-01T00:00:00Z",
  "signature": "base64_signature"
}
```

#### Message Types
- `auth`: User authentication
- `join`: Join room
- `leave`: Leave room
- `message`: Chat message
- `sync`: Room synchronization
- `ping/pong`: Connection health

## Frontend Development

### Project Structure

```
frontend/
├── index.html           # Main application page
├── app.js              # Main application logic
├── styles.css          # Application styling
├── components/         # UI components
│   ├── ChatPane.js    # Message display
│   ├── RoomList.js    # Room management
│   ├── InputBar.js    # Message input
│   ├── UserList.js    # User list
│   └── SettingsPanel.js # Settings
├── assets/            # Static assets
│   ├── images/        # Image files
│   └── fonts/         # Font files
└── static/            # Static file serving
    ├── index.html     # Static version
    └── favicon.ico    # Application icon
```

### Component Architecture

#### RipcordApp (Main Application)
- Manages WebSocket connection
- Coordinates between components
- Handles application state
- Routes messages to components

#### ChatPane Component
- Displays chat messages
- Handles message rendering
- Manages scroll behavior
- Supports message search

#### RoomList Component
- Manages room list display
- Handles room selection
- Supports room filtering
- Manages room state

#### InputBar Component
- Handles message input
- Supports keyboard shortcuts
- Manages message drafts
- Handles file uploads

#### UserList Component
- Displays online users
- Manages user status
- Supports user actions
- Handles user filtering

#### SettingsPanel Component
- Manages user preferences
- Handles theme switching
- Controls privacy settings
- Manages account settings

### State Management

The application uses a simple state management pattern:

```javascript
class RipcordApp {
    constructor() {
        this.currentRoom = null;
        this.currentUser = null;
        this.rooms = new Map();
        this.users = new Map();
        this.websocket = null;
        this.components = {};
    }
    
    // State update methods
    updateConnectionStatus(status, text) { /* ... */ }
    updateCurrentRoomDisplay() { /* ... */ }
    
    // Component communication
    handleWebSocketMessage(data) { /* ... */ }
}
```

### WebSocket Communication

#### Connection Management
```javascript
connectToBackend() {
    const wsUrl = `ws://${window.location.host}/ws`;
    this.websocket = new WebSocket(wsUrl);
    
    this.websocket.onopen = () => {
        this.updateConnectionStatus('connected', 'Connected');
        this.authenticateUser();
    };
    
    this.websocket.onmessage = (event) => {
        this.handleWebSocketMessage(JSON.parse(event.data));
    };
}
```

#### Message Handling
```javascript
handleWebSocketMessage(data) {
    switch (data.type) {
        case 'auth_response':
            this.handleAuthResponse(data);
            break;
        case 'room_list':
            this.handleRoomList(data);
            break;
        case 'message':
            this.handleNewMessage(data);
            break;
        // ... other message types
    }
}
```

## Development Workflow

### Setting Up Development Environment

1. **Install Go 1.21+**
   ```bash
   # Download from https://golang.org/dl/
   # Or use package manager
   ```

2. **Install I2P Router**
   ```bash
   # Download from https://geti2p.net/
   # Start I2P router
   ```

3. **Clone Repository**
   ```bash
   git clone https://github.com/your-username/ripcord.git
   cd ripcord
   ```

4. **Initialize Go Module**
   ```bash
   cd backend
   go mod init ripcord
   go mod tidy
   ```

### Running in Development Mode

#### Backend
```bash
cd backend
go run main.go
```

#### Frontend
```bash
cd frontend
# Use any static file server
python -m http.server 8000
# Or
npx serve .
```

### Testing

#### Backend Tests
```bash
cd backend
go test ./...
go test -v ./...
go test -cover ./...
```

#### Frontend Tests
```bash
# When test framework is implemented
cd frontend
npm test
```

### Code Style Guidelines

#### Go (Backend)
- Use `gofmt` for formatting
- Follow Go naming conventions
- Write comprehensive tests
- Use meaningful variable names
- Add comments for complex logic

#### JavaScript (Frontend)
- Use ES6+ features
- Follow camelCase naming
- Use meaningful function names
- Add JSDoc comments
- Keep functions small and focused

### Debugging

#### Backend Debugging
```bash
# Run with debug logging
go run main.go -debug

# Use Delve debugger
dlv debug main.go
```

#### Frontend Debugging
- Use browser developer tools
- Enable WebSocket debugging
- Check console for errors
- Use browser network tab

## Deployment

### Backend Deployment
```bash
# Build for production
go build -ldflags="-s -w" -o ripcord .

# Run with production config
./ripcord -config=config.prod.json
```

### Frontend Deployment
```bash
# Copy frontend files to web server
cp -r frontend/* /var/www/ripcord/

# Or use static file hosting service
```

## Contributing Guidelines

### Code Review Process
1. Create feature branch
2. Write tests for new functionality
3. Ensure code follows style guidelines
4. Submit pull request
5. Address review comments
6. Merge after approval

### Testing Requirements
- Unit tests for all new functions
- Integration tests for API endpoints
- Frontend component tests
- End-to-end tests for critical flows

### Documentation
- Update README for new features
- Add inline code comments
- Update API documentation
- Include usage examples

## Security Considerations

### Backend Security
- Validate all input data
- Use prepared statements for database queries
- Implement rate limiting
- Log security events
- Use secure random number generation

### Frontend Security
- Sanitize user input
- Validate data before sending to backend
- Use HTTPS in production
- Implement Content Security Policy
- Secure local storage usage

### Network Security
- Encrypt all communications
- Use I2P for anonymity
- Implement message signing
- Secure key exchange protocols
- Regular security audits

## Performance Optimization

### Backend Optimization
- Use connection pooling
- Implement caching strategies
- Optimize database queries
- Use goroutines for concurrent operations
- Profile memory usage

### Frontend Optimization
- Minimize bundle size
- Use lazy loading for components
- Implement virtual scrolling for large lists
- Optimize image assets
- Use efficient DOM manipulation

## Monitoring and Logging

### Backend Logging
```go
import "log"

log.Printf("User %s joined room %s", userID, roomID)
log.Printf("Message sent: %s", messageID)
```

### Frontend Logging
```javascript
console.log('User joined room:', roomId);
console.error('Connection failed:', error);
```

### Metrics to Monitor
- Active connections
- Message throughput
- Error rates
- Response times
- Memory usage
- CPU usage 