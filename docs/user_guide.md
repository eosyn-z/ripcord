# User Guide - Ripcord

Welcome to Ripcord, a decentralized secure chat platform. This guide will help you get started and make the most of the application.

## Getting Started

### Prerequisites

Before using Ripcord, you need to have:

1. **I2P Router**: Download and install I2P from [geti2p.net](https://geti2p.net/)
2. **Modern Web Browser**: Chrome, Firefox, Safari, or Edge with WebSocket support
3. **Ripcord Application**: The Ripcord backend and frontend files

### Installation

1. **Install I2P Router**
   - Download I2P from [geti2p.net](https://geti2p.net/)
   - Follow the installation instructions for your operating system
   - Start the I2P router and wait for it to connect to the network

2. **Start Ripcord Backend**
   ```bash
   cd backend
   ./ripcord
   ```

3. **Open the Frontend**
   - Open `frontend/index.html` in your web browser
   - Or serve the frontend directory with a web server

### First Time Setup

1. **Configure Your Profile**
   - Click the "Settings" button in the top-right corner
   - Enter your desired username
   - Set your display name (optional)
   - Configure your privacy preferences

2. **Generate Encryption Keys**
   - In the Settings panel, click "Generate New Keys"
   - This creates your personal encryption keys for secure messaging
   - Keep these keys safe - they're stored locally on your device

3. **Configure I2P Connection**
   - Ensure I2P router is running
   - Verify SAM address and port settings
   - Test the connection

## Using the Application

### Interface Overview

The Ripcord interface consists of several main areas:

```
┌─────────────────────────────────────────────────────────────┐
│                    Header Bar                              │
│  Ripcord                    [● Connected]                 │
├─────────────┬─────────────────────────────────────────────┤
│             │                                             │
│   Sidebar   │              Chat Area                     │
│             │                                             │
│ • Rooms     │  ┌─────────────────────────────────────┐    │
│   - Room 1  │  │ Room Name                    [⚙]  │    │
│   - Room 2  │  └─────────────────────────────────────┘    │
│             │                                             │
│ • Users     │  ┌─────────────────────────────────────┐    │
│   - User 1  │  │                                     │    │
│   - User 2  │  │         Message Area               │    │
│             │  │                                     │    │
│ [Create]    │  │                                     │    │
│             │  └─────────────────────────────────────┘    │
│             │                                             │
│             │  ┌─────────────────────────────────────┐    │
│             │  │ [Type your message...]        [Send]│    │
│             │  └─────────────────────────────────────┘    │
└─────────────┴─────────────────────────────────────────────┘
```

### Creating and Joining Rooms

#### Creating a Room
1. Click the "Create Room" button in the sidebar
2. Enter a room name (required)
3. Add a description (optional)
4. Click "Create Room"
5. The room will appear in your room list

#### Joining a Room
1. Click on any room in the room list
2. You'll automatically join the room
3. The room name will appear in the chat header
4. Start typing messages in the input area

### Sending Messages

#### Basic Messaging
1. Click in the message input area at the bottom
2. Type your message
3. Press Enter or click "Send"
4. Your message will appear in the chat

#### Message Formatting
Ripcord supports basic text formatting:
- **Bold**: `**text**` or Ctrl+B
- *Italic*: `*text*` or Ctrl+I
- `Code`: `` `code` `` or Ctrl+K
- Indent: Tab key

#### File Sharing
- **Drag and Drop**: Drag files directly into the message area
- **Paste Images**: Copy an image and paste it (Ctrl+V)
- **File Upload**: Click the upload button (when implemented)

### Managing Your Profile

#### Accessing Settings
1. Click the "Settings" button in the chat header
2. The settings panel will slide in from the right

#### Profile Settings
- **Username**: Your unique identifier
- **Display Name**: How others see your name
- **Avatar**: Upload a profile picture (when implemented)

#### Privacy Settings
- **Show Online Status**: Control who can see if you're online
- **Allow Direct Messages**: Control who can send you private messages
- **Encrypt Messages**: Enable/disable message encryption

#### Notification Settings
- **Enable Notifications**: Browser notifications for new messages
- **Sound Notifications**: Audio alerts for messages
- **Desktop Notifications**: System-level notifications

### User Management

#### Viewing Online Users
- The user list shows all currently online users
- Users are sorted by status (online, away, busy, offline)
- Click on a user to see their profile (when implemented)

#### User Status
- **Online**: User is actively using the application
- **Away**: User hasn't been active for a while
- **Busy**: User has set their status to busy
- **Offline**: User is not connected

#### User Actions
- **Send Direct Message**: Click on a user to start a private conversation
- **View Profile**: See user details and public information
- **Block User**: Prevent a user from contacting you (when implemented)

### Security Features

#### Encryption
- All messages are encrypted end-to-end
- Your private keys are stored locally only
- Messages cannot be read by anyone except the intended recipients

#### Anonymity
- Communication is routed through I2P network
- Your real IP address is hidden
- User identities are pseudonymous

#### Privacy
- No message history is stored on central servers
- You control your own data
- Messages are only stored locally if you choose

### Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Enter | Send message |
| Shift+Enter | New line in message |
| Ctrl+B | Bold text |
| Ctrl+I | Italic text |
| Ctrl+K | Code formatting |
| Tab | Indent text |
| Escape | Close modals/panels |

### Troubleshooting

#### Connection Issues

**Problem**: "Connection Error" status
**Solutions**:
1. Check if I2P router is running
2. Verify SAM address and port in settings
3. Restart the Ripcord backend
4. Check firewall settings

**Problem**: "Disconnected" status
**Solutions**:
1. Check your internet connection
2. Restart the application
3. Check if the backend server is running

#### Message Issues

**Problem**: Messages not sending
**Solutions**:
1. Check your connection status
2. Ensure you're in a room
3. Try refreshing the page
4. Check browser console for errors

**Problem**: Messages not appearing
**Solutions**:
1. Check if you're in the correct room
2. Refresh the page
3. Check if other users are online
4. Verify encryption keys are generated

#### Performance Issues

**Problem**: Slow message loading
**Solutions**:
1. Check your internet connection
2. Reduce the number of open rooms
3. Clear browser cache
4. Restart the application

**Problem**: High memory usage
**Solutions**:
1. Close unused browser tabs
2. Restart the application
3. Clear browser cache
4. Update to the latest version

### Advanced Features

#### Room Management

**Creating Private Rooms**
1. Create a room as usual
2. In room settings, set visibility to "Private"
3. Share the room ID with specific users

**Room Permissions**
- **Admin**: Can manage room settings and members
- **Moderator**: Can moderate messages and users
- **Member**: Can send messages and view content

**Room Settings**
- Change room name and description
- Manage member permissions
- Set room visibility (public/private)
- Configure message retention

#### Message Features

**Message Search**
1. Use Ctrl+F to open search
2. Type your search term
3. Navigate through results
4. Click on a result to jump to that message

**Message History**
- Messages are stored locally
- Scroll up to view older messages
- Use search to find specific content

**Message Reactions**
- React to messages with emojis (when implemented)
- View reaction counts
- Add your own reactions

### Best Practices

#### Security
1. **Keep your keys safe**: Don't share your private keys
2. **Use strong usernames**: Avoid easily identifiable names
3. **Verify connections**: Check the connection status regularly
4. **Update regularly**: Keep the application updated

#### Privacy
1. **Be mindful of what you share**: Even encrypted messages can be screenshotted
2. **Use pseudonyms**: Don't use your real name
3. **Control your visibility**: Adjust privacy settings as needed
4. **Be aware of metadata**: Message timing and frequency can reveal information

#### Communication
1. **Be respectful**: Follow community guidelines
2. **Use clear language**: Avoid ambiguous messages
3. **Be patient**: I2P routing can cause delays
4. **Report issues**: Help improve the platform

### Getting Help

#### Documentation
- Check this user guide for detailed instructions
- Review the developer guide for technical details
- Read the README for installation help

#### Community Support
- Join the Ripcord community room
- Ask questions in the help channel
- Report bugs and suggest features

#### Technical Support
- Check the application logs for error messages
- Verify your I2P router configuration
- Ensure all prerequisites are installed
- Try restarting the application

### Updates and Maintenance

#### Updating the Application
1. Download the latest version
2. Stop the current application
3. Replace the files with new versions
4. Restart the application
5. Check for any configuration changes

#### Backup and Recovery
1. **Backup your keys**: Export your encryption keys
2. **Backup settings**: Save your configuration file
3. **Backup messages**: Export important conversations
4. **Test recovery**: Verify your backups work

#### Regular Maintenance
1. **Update I2P**: Keep your I2P router updated
2. **Clear cache**: Periodically clear browser cache
3. **Check logs**: Review application logs for issues
4. **Test connections**: Verify network connectivity

This user guide covers the essential features of Ripcord. As the application evolves, new features will be added and this guide will be updated accordingly. 