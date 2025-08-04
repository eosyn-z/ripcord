# Ripcord - Decentralized Secure Chat Platform

A peer-to-peer encrypted chat application built with Go backend and vanilla JavaScript frontend, featuring Ed25519 cryptography, I2P anonymous networking, and real-time WebSocket communication.

## Features

- **End-to-End Encryption**: Ed25519 digital signatures for message authentication
- **Anonymous Networking**: I2P (Invisible Internet Project) integration for privacy
- **Real-time Communication**: WebSocket-based messaging with HTTP API fallback
- **Decentralized Architecture**: Peer-to-peer room management with invite codes
- **Modern UI**: Clean beige-themed interface with responsive design
- **Cross-Platform**: Runs on Windows, macOS, and Linux

## Quick Start

### Prerequisites

- Go 1.19 or later
- Modern web browser with WebSocket support
- I2P router (optional, for anonymous networking)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd ripcord
   ```

2. **Install Go dependencies**
   ```bash
   cd backend
   go mod download
   ```

3. **Build and run the server**
   ```bash
   go run .
   ```

4. **Open your browser**
   ```
   http://localhost:8080
   ```

The application will automatically create a configuration file and database on first run.

## Architecture

### Backend (Go)
- **HTTP API**: RESTful endpoints for room/message management
- **WebSocket Server**: Real-time communication with automatic reconnection
- **SQLite Database**: Local message and room persistence
- **Ed25519 Cryptography**: Message signing and verification
- **I2P Integration**: Anonymous networking via SAM protocol
- **Configuration Management**: JSON-based configuration with defaults

### Frontend (JavaScript)
- **Vanilla JavaScript**: No frameworks, modern ES6+ features
- **WebSocket Client**: Real-time messaging with HTTP fallback
- **Component Architecture**: Modular UI components
- **Responsive Design**: Mobile-friendly beige theme
- **Local Storage**: User preferences and draft messages

## Configuration

The server uses `config.json` for configuration. Default configuration:

```json
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "database": {
    "type": "sqlite",
    "database": "ripcord.db"
  },
  "i2p": {
    "enabled": true,
    "sam_address": "127.0.0.1",
    "sam_port": 7656
  },
  "security": {
    "encryption_enabled": true,
    "key_size": 256,
    "algorithm": "Ed25519"
  }
}
```

## API Documentation

### HTTP API Endpoints

#### Identity
- `GET /api/identity` - Get current user identity and public key

#### Rooms
- `GET /api/rooms` - List all available rooms
- `POST /api/rooms/create` - Create a new room
- `POST /api/rooms/join` - Join a room by invite code
- `POST /api/rooms/leave` - Leave a room

#### Messages
- `GET /api/messages?room_id=<id>` - Get messages for a room
- `POST /api/messages/send` - Send a message to a room

### WebSocket Protocol

#### Authentication
```json
{
  "type": "auth",
  "username": "your-username"
}
```

#### Join Room
```json
{
  "type": "join_room",
  "room_id": "room-uuid"
}
```

#### Send Message
```json
{
  "type": "send_message",
  "content": "Hello, world!"
}
```

#### Message Received
```json
{
  "type": "message",
  "message": {
    "id": "msg-uuid",
    "room_id": "room-uuid",
    "username": "sender",
    "content": "Hello, world!",
    "timestamp": 1234567890
  }
}
```

## Security

### Cryptographic Features
- **Ed25519 Signatures**: All messages are cryptographically signed
- **Key Management**: Automatic key generation and secure storage
- **Message Integrity**: Prevents message tampering and forgery

### Privacy Features
- **I2P Integration**: Optional anonymous networking
- **No Central Server**: Decentralized peer-to-peer architecture
- **Local Storage**: Messages stored locally, not on remote servers

### Input Validation
- **Message Length**: Limited to 2000 characters
- **Room Names**: Limited to 100 characters
- **SQL Injection Protection**: Parameterized queries only
- **XSS Prevention**: Proper input sanitization

## Development

### Project Structure
```
ripcord/
├── backend/
│   ├── main.go              # HTTP server and WebSocket handler
│   ├── config.go            # Configuration management
│   ├── room.go              # Room management
│   ├── message.go           # Message handling and signing
│   ├── node.go              # P2P node management
│   ├── protocol.go          # Protocol definitions
│   ├── database/
│   │   └── db.go            # Database interface and SQLite implementation
│   ├── security/
│   │   └── crypto.go        # Ed25519 cryptography
│   ├── i2p/
│   │   └── i2p.go           # I2P SAM protocol integration
│   └── tests/
│       └── main_test.go     # Unit tests
├── frontend/
│   ├── index.html           # Main HTML page
│   ├── styles.css           # Beige theme styling
│   ├── app.js               # Main application logic
│   ├── components/
│   │   ├── ChatPane.js      # Message display component
│   │   ├── RoomList.js      # Room list component
│   │   ├── UserList.js      # User list component
│   │   ├── InputBar.js      # Message input component
│   │   └── SettingsPanel.js # Settings management
│   └── assets/
│       └── images/
│           └── favicon.ico  # Application icon
├── go.mod                   # Go module dependencies
├── config.json              # Server configuration (auto-generated)
├── FIXES_LOG.md             # Detailed fix documentation
└── README.md                # This file
```

### Running Tests
```bash
cd backend
go test ./...
```

### Building for Production
```bash
cd backend
go build -o ripcord .
```

## I2P Integration

For anonymous networking, install and configure I2P:

1. **Install I2P**: Download from [geti2p.net](https://geti2p.net)
2. **Start I2P Router**: Run I2P and wait for network integration
3. **Enable SAM**: In I2P router console, enable SAM application bridge
4. **Configure Ripcord**: Set `i2p.enabled: true` in config.json

## Troubleshooting

### Common Issues

#### "Failed to connect to backend"
- Ensure the server is running on the correct port
- Check firewall settings
- Verify configuration file is valid JSON

#### "WebSocket connection failed"
- Check browser console for detailed error messages
- Ensure WebSocket support is enabled in browser
- Try the HTTP API fallback functionality

#### "Database connection failed"
- Ensure write permissions in the data directory
- Check disk space availability
- Verify SQLite is properly installed

#### "I2P connection failed"
- Ensure I2P router is running and integrated
- Verify SAM bridge is enabled on port 7656
- Check I2P router console for connectivity issues

### Debug Mode

Set environment variable for verbose logging:
```bash
export RIPCORD_DEBUG=true
go run .
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Create a pull request

### Code Style
- Follow Go standard formatting (`gofmt`)
- Use meaningful variable names
- Add comments for public functions
- Maintain consistent error handling patterns

## Security Considerations

### Production Deployment
- Change CORS settings from `*` to specific domains
- Implement rate limiting
- Use HTTPS/WSS for encrypted transport
- Set up proper authentication and authorization
- Configure firewall rules appropriately
- Regular security updates

### Known Limitations
- SQLite is suitable for development but consider PostgreSQL for production
- Current authentication is basic - implement proper session management
- No built-in DDoS protection - use reverse proxy in production

## License

[Specify your license here]

## Support

For issues and questions:
1. Check the troubleshooting section above
2. Review FIXES_LOG.md for detailed technical information
3. Search existing issues in the repository
4. Create a new issue with detailed reproduction steps

## Changelog

### v1.0.0 - Initial Release
- Complete Ed25519 cryptographic implementation
- Full WebSocket real-time communication
- SQLite database integration
- I2P anonymous networking support
- Responsive beige-themed UI
- Comprehensive security and input validation
- HTTP API with WebSocket fallback
- Modular component architecture

---

**Note**: This is a development version. For production use, implement additional security measures and consider scalability requirements.