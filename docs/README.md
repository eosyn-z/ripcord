# Ripcord - Decentralized Secure Chat Platform

Ripcord is a decentralized, secure chat platform built with Go backend and modern web frontend. It leverages I2P (Invisible Internet Project) for anonymous communication and implements end-to-end encryption for message security.

## Features

- **Decentralized Architecture**: No central server required, peer-to-peer communication
- **I2P Integration**: Anonymous communication over the I2P network
- **End-to-End Encryption**: AES-256-GCM encryption with RSA key exchange
- **Real-time Messaging**: WebSocket-based real-time chat functionality
- **Room-based Chat**: Create and join chat rooms with multiple users
- **Modern UI**: Responsive web interface with dark/light theme support
- **Privacy-focused**: User anonymity and message encryption by default

## Architecture

### Backend (Go)
- **Node Management**: Decentralized node discovery and peer management
- **Room System**: Chat room creation, management, and synchronization
- **Message Handling**: Encrypted message routing and delivery
- **I2P Integration**: SAM (Simple Anonymous Messaging) protocol implementation
- **Database**: SQLite for local message persistence
- **Security**: Cryptographic operations and key management

### Frontend (HTML/CSS/JavaScript)
- **Real-time UI**: WebSocket-based real-time updates
- **Component-based**: Modular JavaScript components for maintainability
- **Responsive Design**: Mobile-friendly interface
- **Theme Support**: Dark/light mode with customizable settings
- **User Management**: User profiles, status indicators, and preferences

## Project Structure

```
ripcord/
├── backend/                 # Go backend application
│   ├── main.go             # Application entry point
│   ├── node.go             # Decentralized node management
│   ├── room.go             # Chat room functionality
│   ├── message.go          # Message handling and encryption
│   ├── protocol.go         # Communication protocol
│   ├── config.go           # Configuration management
│   ├── database/           # Database layer
│   │   └── db.go          # Database interface and SQLite implementation
│   ├── security/           # Cryptographic operations
│   │   └── crypto.go      # Encryption, signing, key management
│   ├── i2p/               # I2P network integration
│   │   └── i2p.go        # SAM protocol implementation
│   └── tests/             # Backend tests
│       └── main_test.go   # Unit tests
├── frontend/              # Web frontend
│   ├── index.html         # Main application page
│   ├── app.js            # Main application logic
│   ├── styles.css        # Application styling
│   ├── components/       # UI components
│   │   ├── ChatPane.js   # Message display component
│   │   ├── RoomList.js   # Room management component
│   │   ├── InputBar.js   # Message input component
│   │   ├── UserList.js   # User list component
│   │   └── SettingsPanel.js # Settings management
│   ├── assets/           # Static assets
│   │   ├── images/       # Image files
│   │   └── fonts/        # Font files
│   └── static/           # Static file serving
│       ├── index.html    # Static version
│       └── favicon.ico   # Application icon
├── docs/                 # Documentation
│   ├── README.md         # This file
│   ├── developer_guide.md # Development guide
│   └── user_guide.md     # User manual
└── tests/               # Test files
    ├── backend_tests/    # Backend test suites
    │   ├── node_test.go  # Node tests
    │   └── message_test.go # Message tests
    └── frontend_tests/   # Frontend test suites
        ├── chat_pane_test.js # Chat component tests
        └── user_list_test.js # User list tests
```

## Prerequisites

### Backend Requirements
- Go 1.21 or later
- I2P router running locally (for I2P network access)
- SQLite3 (included with Go)

### Frontend Requirements
- Modern web browser with WebSocket support
- No additional dependencies (vanilla JavaScript)

## Installation

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/ripcord.git
cd ripcord
```

### 2. Install I2P Router
Download and install I2P from [geti2p.net](https://geti2p.net/)

### 3. Build Backend
```bash
cd backend
go mod init ripcord
go mod tidy
go build -o ripcord .
```

### 4. Run the Application
```bash
# Start the backend server
./ripcord

# Open frontend/index.html in your browser
# Or serve the frontend directory with a web server
```

## Configuration

The application can be configured through the `config.json` file or environment variables:

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
    "algorithm": "AES-256-GCM"
  }
}
```

## Usage

### Starting the Application
1. Ensure I2P router is running
2. Start the backend server: `./ripcord`
3. Open the frontend in your browser
4. Configure your username and settings
5. Create or join a chat room

### Creating a Room
1. Click "Create Room" in the sidebar
2. Enter room name and description
3. Click "Create Room"

### Joining a Room
1. Click on any room in the room list
2. Start chatting immediately

### Security Features
- All messages are encrypted end-to-end
- User identities are pseudonymous
- Communication is routed through I2P for anonymity
- No message history is stored on central servers

## Development

### Backend Development
```bash
cd backend
go run main.go
```

### Frontend Development
```bash
# Serve frontend with a local server
cd frontend
python -m http.server 8000
# Or use any other static file server
```

### Running Tests
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests (when implemented)
cd frontend
npm test
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Add tests for new functionality
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## Security Considerations

- **Key Management**: Private keys are stored locally and never transmitted
- **Message Encryption**: All messages are encrypted with AES-256-GCM
- **Network Anonymity**: I2P provides network-level anonymity
- **No Central Authority**: Decentralized architecture prevents single points of failure
- **Open Source**: Full transparency for security auditing

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [I2P Project](https://geti2p.net/) for anonymous networking
- [Go Programming Language](https://golang.org/) for backend development
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) for real-time communication

## Support

For support and questions:
- Create an issue on GitHub
- Check the documentation in the `docs/` directory
- Review the developer guide for technical details 