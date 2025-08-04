// TODO: Implement comprehensive ChatPane tests
// TODO: Test message rendering and display
// TODO: Test scroll behavior and pagination
// TODO: Test message search functionality

// Mock DOM environment for testing
const { JSDOM } = require('jsdom');
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<head></head>
<body>
    <div id="chat-messages"></div>
</body>
</html>
`);

global.window = dom.window;
global.document = dom.window.document;

// Import the component to test
// const ChatPane = require('../../frontend/components/ChatPane.js');

describe('ChatPane', () => {
    let chatPane;
    let messagesContainer;

    beforeEach(() => {
        // Set up DOM elements
        messagesContainer = document.getElementById('chat-messages');
        chatPane = new ChatPane();
    });

    afterEach(() => {
        // Clean up
        messagesContainer.innerHTML = '';
    });

    describe('Initialization', () => {
        test('should initialize with empty messages array', () => {
            expect(chatPane.messages).toEqual([]);
            expect(chatPane.currentRoomId).toBeNull();
        });

        test('should bind scroll events', () => {
            const spy = jest.spyOn(chatPane, 'handleScroll');
            const event = new Event('scroll');
            messagesContainer.dispatchEvent(event);
            expect(spy).toHaveBeenCalled();
        });
    });

    describe('Message Management', () => {
        test('should add message to messages array', () => {
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            chatPane.addMessage(message);
            expect(chatPane.messages).toHaveLength(1);
            expect(chatPane.messages[0]).toEqual(message);
        });

        test('should render message in DOM', () => {
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            chatPane.addMessage(message);
            const messageElement = messagesContainer.querySelector('.message');
            expect(messageElement).toBeTruthy();
            expect(messageElement.querySelector('.message-username').textContent).toBe('TestUser');
            expect(messageElement.querySelector('.message-text').textContent).toBe('Hello, world!');
        });

        test('should clear messages', () => {
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            chatPane.addMessage(message);
            expect(chatPane.messages).toHaveLength(1);

            chatPane.clearMessages();
            expect(chatPane.messages).toHaveLength(0);
            expect(messagesContainer.children).toHaveLength(0);
        });
    });

    describe('Message Rendering', () => {
        test('should create message element with correct structure', () => {
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            const messageElement = chatPane.createMessageElement(message);
            
            expect(messageElement.classList.contains('message')).toBe(true);
            expect(messageElement.querySelector('.message-avatar')).toBeTruthy();
            expect(messageElement.querySelector('.message-content')).toBeTruthy();
            expect(messageElement.querySelector('.message-username')).toBeTruthy();
            expect(messageElement.querySelector('.message-time')).toBeTruthy();
            expect(messageElement.querySelector('.message-text')).toBeTruthy();
        });

        test('should mark own messages correctly', () => {
            // Mock current user
            window.ripcordApp = {
                currentUser: { id: 'user1' }
            };

            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            const messageElement = chatPane.createMessageElement(message);
            expect(messageElement.classList.contains('own-message')).toBe(true);
        });

        test('should create avatar with username initial', () => {
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            const avatar = chatPane.createAvatar(message.username);
            expect(avatar.classList.contains('message-avatar')).toBe(true);
            expect(avatar.textContent).toBe('T');
        });

        test('should format timestamp correctly', () => {
            const now = new Date();
            const oneHourAgo = new Date(now.getTime() - 60 * 60 * 1000);
            const oneDayAgo = new Date(now.getTime() - 25 * 60 * 60 * 1000);

            const time1 = chatPane.formatTimestamp(oneHourAgo.toISOString());
            const time2 = chatPane.formatTimestamp(oneDayAgo.toISOString());

            expect(time1).toMatch(/^\d{1,2}:\d{2}$/); // HH:MM format
            expect(time2).toMatch(/^\d{1,2}\/\d{1,2}\/\d{4}$/); // MM/DD/YYYY format
        });
    });

    describe('Scroll Behavior', () => {
        test('should scroll to bottom when new message added', () => {
            const scrollSpy = jest.spyOn(chatPane, 'scrollToBottom');
            
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            chatPane.addMessage(message);
            expect(scrollSpy).toHaveBeenCalled();
        });

        test('should load more messages when scrolling to top', () => {
            const loadMoreSpy = jest.spyOn(chatPane, 'loadMoreMessages');
            
            // Simulate scroll to top
            messagesContainer.scrollTop = 0;
            const event = new Event('scroll');
            messagesContainer.dispatchEvent(event);

            expect(loadMoreSpy).toHaveBeenCalled();
        });

        test('should scroll to specific message', () => {
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            chatPane.addMessage(message);
            
            const scrollIntoViewSpy = jest.spyOn(Element.prototype, 'scrollIntoView');
            chatPane.scrollToMessage('msg1');
            
            expect(scrollIntoViewSpy).toHaveBeenCalledWith({
                behavior: 'smooth',
                block: 'center'
            });
        });
    });

    describe('Message Search', () => {
        beforeEach(() => {
            const messages = [
                {
                    id: 'msg1',
                    user_id: 'user1',
                    username: 'Alice',
                    content: 'Hello, how are you?',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                },
                {
                    id: 'msg2',
                    user_id: 'user2',
                    username: 'Bob',
                    content: 'I am doing well, thanks!',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                },
                {
                    id: 'msg3',
                    user_id: 'user1',
                    username: 'Alice',
                    content: 'That is great to hear!',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                }
            ];

            messages.forEach(msg => chatPane.addMessage(msg));
        });

        test('should search messages by content', () => {
            const results = chatPane.searchMessages('hello');
            expect(results).toHaveLength(1);
            expect(results[0].content).toBe('Hello, how are you?');
        });

        test('should search messages by username', () => {
            const results = chatPane.searchMessages('alice');
            expect(results).toHaveLength(2);
            expect(results[0].username).toBe('Alice');
            expect(results[1].username).toBe('Alice');
        });

        test('should return empty array for no matches', () => {
            const results = chatPane.searchMessages('nonexistent');
            expect(results).toHaveLength(0);
        });

        test('should be case insensitive', () => {
            const results = chatPane.searchMessages('HELLO');
            expect(results).toHaveLength(1);
        });
    });

    describe('Message History', () => {
        test('should load message history for room', () => {
            const roomId = 'room1';
            const sendMessageSpy = jest.spyOn(window.ripcordApp, 'sendWebSocketMessage');

            chatPane.loadMessageHistory(roomId);

            expect(chatPane.currentRoomId).toBe(roomId);
            expect(sendMessageSpy).toHaveBeenCalledWith({
                type: 'get_messages',
                room_id: roomId,
                limit: 50
            });
        });

        test('should handle message history response', () => {
            const messages = [
                {
                    id: 'msg1',
                    user_id: 'user1',
                    username: 'Alice',
                    content: 'Hello',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                },
                {
                    id: 'msg2',
                    user_id: 'user2',
                    username: 'Bob',
                    content: 'Hi there',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                }
            ];

            chatPane.handleMessageHistory(messages);

            expect(chatPane.messages).toHaveLength(2);
            expect(messagesContainer.children).toHaveLength(2);
        });
    });

    describe('Message Highlighting', () => {
        test('should highlight specific message', () => {
            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            chatPane.addMessage(message);
            chatPane.highlightMessage('msg1');

            const messageElement = messagesContainer.querySelector('.message');
            expect(messageElement.classList.contains('highlighted')).toBe(true);
        });

        test('should remove previous highlights when highlighting new message', () => {
            const messages = [
                {
                    id: 'msg1',
                    user_id: 'user1',
                    username: 'Alice',
                    content: 'Hello',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                },
                {
                    id: 'msg2',
                    user_id: 'user2',
                    username: 'Bob',
                    content: 'Hi',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                }
            ];

            messages.forEach(msg => chatPane.addMessage(msg));

            chatPane.highlightMessage('msg1');
            chatPane.highlightMessage('msg2');

            const highlightedMessages = messagesContainer.querySelectorAll('.message.highlighted');
            expect(highlightedMessages).toHaveLength(1);
        });
    });

    describe('Utility Methods', () => {
        test('should return correct message count', () => {
            expect(chatPane.getMessageCount()).toBe(0);

            const message = {
                id: 'msg1',
                user_id: 'user1',
                username: 'TestUser',
                content: 'Hello, world!',
                timestamp: new Date().toISOString(),
                type: 'text'
            };

            chatPane.addMessage(message);
            expect(chatPane.getMessageCount()).toBe(1);
        });

        test('should return current room ID', () => {
            expect(chatPane.getCurrentRoomId()).toBeNull();

            chatPane.currentRoomId = 'room1';
            expect(chatPane.getCurrentRoomId()).toBe('room1');
        });
    });

    describe('Error Handling', () => {
        test('should handle invalid message data', () => {
            const invalidMessage = {
                id: 'msg1',
                // Missing required fields
            };

            expect(() => {
                chatPane.addMessage(invalidMessage);
            }).not.toThrow();
        });

        test('should handle DOM manipulation errors', () => {
            // Remove the messages container
            messagesContainer.remove();

            expect(() => {
                chatPane.addMessage({
                    id: 'msg1',
                    user_id: 'user1',
                    username: 'TestUser',
                    content: 'Hello, world!',
                    timestamp: new Date().toISOString(),
                    type: 'text'
                });
            }).not.toThrow();
        });
    });

    describe('Performance', () => {
        test('should handle large number of messages efficiently', () => {
            const startTime = performance.now();

            // Add 1000 messages
            for (let i = 0; i < 1000; i++) {
                chatPane.addMessage({
                    id: `msg${i}`,
                    user_id: 'user1',
                    username: 'TestUser',
                    content: `Message ${i}`,
                    timestamp: new Date().toISOString(),
                    type: 'text'
                });
            }

            const endTime = performance.now();
            const duration = endTime - startTime;

            expect(duration).toBeLessThan(1000); // Should complete within 1 second
            expect(chatPane.getMessageCount()).toBe(1000);
        });
    });
}); 