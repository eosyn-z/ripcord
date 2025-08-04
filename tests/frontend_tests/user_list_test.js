// TODO: Implement comprehensive UserList tests
// TODO: Test user display and status indicators
// TODO: Test user filtering and search
// TODO: Test user actions and interactions

// Mock DOM environment for testing
const { JSDOM } = require('jsdom');
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<head></head>
<body>
    <div id="user-list"></div>
</body>
</html>
`);

global.window = dom.window;
global.document = dom.window.document;

// Import the component to test
// const UserList = require('../../frontend/components/UserList.js');

describe('UserList', () => {
    let userList;
    let userListContainer;

    beforeEach(() => {
        // Set up DOM elements
        userListContainer = document.getElementById('user-list');
        userList = new UserList();
    });

    afterEach(() => {
        // Clean up
        userListContainer.innerHTML = '';
    });

    describe('Initialization', () => {
        test('should initialize with empty users array', () => {
            expect(userList.users).toEqual([]);
        });

        test('should bind events', () => {
            // Test that events are bound (implementation dependent)
            expect(userList).toBeDefined();
        });
    });

    describe('User Management', () => {
        test('should update users list', () => {
            const users = [
                {
                    id: 'user1',
                    username: 'Alice',
                    status: 'online',
                    display_name: 'Alice Smith'
                },
                {
                    id: 'user2',
                    username: 'Bob',
                    status: 'away',
                    display_name: 'Bob Johnson'
                }
            ];

            userList.updateUsers(users);
            expect(userList.users).toEqual(users);
        });

        test('should render users in DOM', () => {
            const users = [
                {
                    id: 'user1',
                    username: 'Alice',
                    status: 'online',
                    display_name: 'Alice Smith'
                }
            ];

            userList.updateUsers(users);
            const userElement = userListContainer.querySelector('.user-item');
            expect(userElement).toBeTruthy();
            expect(userElement.querySelector('.user-name').textContent).toBe('Alice');
        });

        test('should render empty state when no users', () => {
            userList.updateUsers([]);
            const emptyState = userListContainer.querySelector('.empty-state');
            expect(emptyState).toBeTruthy();
            expect(emptyState.textContent).toContain('No users online');
        });
    });

    describe('User Element Creation', () => {
        test('should create user element with correct structure', () => {
            const user = {
                id: 'user1',
                username: 'Alice',
                status: 'online',
                display_name: 'Alice Smith'
            };

            const userElement = userList.createUserElement(user);
            
            expect(userElement.classList.contains('user-item')).toBe(true);
            expect(userElement.dataset.userId).toBe('user1');
            expect(userElement.querySelector('.user-avatar')).toBeTruthy();
            expect(userElement.querySelector('.user-info')).toBeTruthy();
            expect(userElement.querySelector('.user-name')).toBeTruthy();
            expect(userElement.querySelector('.user-status')).toBeTruthy();
        });

        test('should create avatar with username initial', () => {
            const user = {
                id: 'user1',
                username: 'Alice',
                status: 'online'
            };

            const avatar = userList.createUserAvatar(user);
            expect(avatar.classList.contains('user-avatar')).toBe(true);
            expect(avatar.textContent).toBe('A');
        });

        test('should create avatar with image when available', () => {
            const user = {
                id: 'user1',
                username: 'Alice',
                status: 'online',
                avatar: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=='
            };

            const avatar = userList.createUserAvatar(user);
            const img = avatar.querySelector('img');
            expect(img).toBeTruthy();
            expect(img.src).toBe(user.avatar);
            expect(img.alt).toBe('Alice');
        });

        test('should create user info with correct content', () => {
            const user = {
                id: 'user1',
                username: 'Alice',
                status: 'online',
                display_name: 'Alice Smith'
            };

            const info = userList.createUserInfo(user);
            expect(info.querySelector('.user-name').textContent).toBe('Alice');
            expect(info.querySelector('.user-status').textContent).toBe('Online');
        });

        test('should get correct status text', () => {
            expect(userList.getStatusText('online')).toBe('Online');
            expect(userList.getStatusText('away')).toBe('Away');
            expect(userList.getStatusText('busy')).toBe('Busy');
            expect(userList.getStatusText('offline')).toBe('Offline');
            expect(userList.getStatusText('unknown')).toBe('Unknown');
        });
    });

    describe('User Sorting', () => {
        test('should sort users by status priority', () => {
            const users = [
                { id: 'user1', username: 'Alice', status: 'offline' },
                { id: 'user2', username: 'Bob', status: 'online' },
                { id: 'user3', username: 'Charlie', status: 'away' },
                { id: 'user4', username: 'David', status: 'busy' }
            ];

            const sortedUsers = userList.sortUsers(users);
            expect(sortedUsers[0].status).toBe('online');
            expect(sortedUsers[1].status).toBe('away');
            expect(sortedUsers[2].status).toBe('busy');
            expect(sortedUsers[3].status).toBe('offline');
        });

        test('should sort users by username when status is same', () => {
            const users = [
                { id: 'user1', username: 'Charlie', status: 'online' },
                { id: 'user2', username: 'Alice', status: 'online' },
                { id: 'user3', username: 'Bob', status: 'online' }
            ];

            const sortedUsers = userList.sortUsers(users);
            expect(sortedUsers[0].username).toBe('Alice');
            expect(sortedUsers[1].username).toBe('Bob');
            expect(sortedUsers[2].username).toBe('Charlie');
        });
    });

    describe('User Operations', () => {
        beforeEach(() => {
            const users = [
                { id: 'user1', username: 'Alice', status: 'online' },
                { id: 'user2', username: 'Bob', status: 'away' }
            ];
            userList.updateUsers(users);
        });

        test('should add new user', () => {
            const newUser = { id: 'user3', username: 'Charlie', status: 'online' };
            userList.addUser(newUser);

            expect(userList.users).toHaveLength(3);
            expect(userList.users[2]).toEqual(newUser);
        });

        test('should update existing user', () => {
            const updatedUser = { id: 'user1', username: 'Alice', status: 'busy' };
            userList.updateUser(updatedUser);

            expect(userList.users[0].status).toBe('busy');
        });

        test('should remove user', () => {
            userList.removeUser('user1');

            expect(userList.users).toHaveLength(1);
            expect(userList.users[0].id).toBe('user2');
        });

        test('should set user status', () => {
            userList.setUserStatus('user1', 'busy');

            const user = userList.users.find(u => u.id === 'user1');
            expect(user.status).toBe('busy');
        });
    });

    describe('User Search', () => {
        beforeEach(() => {
            const users = [
                { id: 'user1', username: 'Alice', display_name: 'Alice Smith', status: 'online' },
                { id: 'user2', username: 'Bob', display_name: 'Bob Johnson', status: 'away' },
                { id: 'user3', username: 'Charlie', display_name: 'Charlie Brown', status: 'online' }
            ];
            userList.updateUsers(users);
        });

        test('should search users by username', () => {
            userList.searchUsers('alice');
            const visibleUsers = userListContainer.querySelectorAll('.user-item');
            expect(visibleUsers).toHaveLength(1);
            expect(visibleUsers[0].querySelector('.user-name').textContent).toBe('Alice');
        });

        test('should search users by display name', () => {
            userList.searchUsers('johnson');
            const visibleUsers = userListContainer.querySelectorAll('.user-item');
            expect(visibleUsers).toHaveLength(1);
            expect(visibleUsers[0].querySelector('.user-name').textContent).toBe('Bob');
        });

        test('should be case insensitive', () => {
            userList.searchUsers('CHARLIE');
            const visibleUsers = userListContainer.querySelectorAll('.user-item');
            expect(visibleUsers).toHaveLength(1);
        });

        test('should show empty state for no matches', () => {
            userList.searchUsers('nonexistent');
            const emptyState = userListContainer.querySelector('.empty-state');
            expect(emptyState).toBeTruthy();
        });
    });

    describe('User Filtering', () => {
        beforeEach(() => {
            const users = [
                { id: 'user1', username: 'Alice', status: 'online' },
                { id: 'user2', username: 'Bob', status: 'away' },
                { id: 'user3', username: 'Charlie', status: 'online' },
                { id: 'user4', username: 'David', status: 'offline' }
            ];
            userList.updateUsers(users);
        });

        test('should filter users by online status', () => {
            userList.filterUsersByStatus('online');
            const visibleUsers = userListContainer.querySelectorAll('.user-item');
            expect(visibleUsers).toHaveLength(2);
            visibleUsers.forEach(user => {
                expect(user.querySelector('.user-status').textContent).toBe('Online');
            });
        });

        test('should filter users by away status', () => {
            userList.filterUsersByStatus('away');
            const visibleUsers = userListContainer.querySelectorAll('.user-item');
            expect(visibleUsers).toHaveLength(1);
            expect(visibleUsers[0].querySelector('.user-status').textContent).toBe('Away');
        });

        test('should show empty state for no matches', () => {
            userList.filterUsersByStatus('busy');
            const emptyState = userListContainer.querySelector('.empty-state');
            expect(emptyState).toBeTruthy();
        });
    });

    describe('User Actions', () => {
        test('should handle user click', () => {
            const user = { id: 'user1', username: 'Alice', status: 'online' };
            const userElement = userList.createUserElement(user);
            
            const clickSpy = jest.spyOn(userList, 'showUserActions');
            userElement.click();
            
            expect(clickSpy).toHaveBeenCalledWith(user);
        });

        test('should handle right-click context menu', () => {
            const user = { id: 'user1', username: 'Alice', status: 'online' };
            const userElement = userList.createUserElement(user);
            
            const contextMenuSpy = jest.spyOn(userList, 'showUserContextMenu');
            const event = new Event('contextmenu');
            userElement.dispatchEvent(event);
            
            expect(contextMenuSpy).toHaveBeenCalledWith(event, user);
        });

        test('should show user actions', () => {
            const user = { id: 'user1', username: 'Alice', status: 'online' };
            
            // Mock console.log to avoid output during tests
            const consoleSpy = jest.spyOn(console, 'log').mockImplementation();
            
            userList.showUserActions(user);
            expect(consoleSpy).toHaveBeenCalledWith('User clicked:', 'Alice');
            
            consoleSpy.mockRestore();
        });

        test('should show user context menu', () => {
            const user = { id: 'user1', username: 'Alice', status: 'online' };
            const event = new Event('contextmenu');
            
            // Mock console.log to avoid output during tests
            const consoleSpy = jest.spyOn(console, 'log').mockImplementation();
            
            userList.showUserContextMenu(event, user);
            expect(consoleSpy).toHaveBeenCalledWith('Context menu for user:', 'Alice');
            
            consoleSpy.mockRestore();
        });
    });

    describe('Utility Methods', () => {
        beforeEach(() => {
            const users = [
                { id: 'user1', username: 'Alice', status: 'online' },
                { id: 'user2', username: 'Bob', status: 'away' },
                { id: 'user3', username: 'Charlie', status: 'online' },
                { id: 'user4', username: 'David', status: 'offline' }
            ];
            userList.updateUsers(users);
        });

        test('should get user by ID', () => {
            const user = userList.getUserById('user1');
            expect(user).toBeTruthy();
            expect(user.username).toBe('Alice');
        });

        test('should return null for non-existent user', () => {
            const user = userList.getUserById('nonexistent');
            expect(user).toBeUndefined();
        });

        test('should get online users', () => {
            const onlineUsers = userList.getOnlineUsers();
            expect(onlineUsers).toHaveLength(2);
            onlineUsers.forEach(user => {
                expect(user.status).toBe('online');
            });
        });

        test('should get offline users', () => {
            const offlineUsers = userList.getOfflineUsers();
            expect(offlineUsers).toHaveLength(1);
            offlineUsers.forEach(user => {
                expect(user.status).toBe('offline');
            });
        });

        test('should get user count', () => {
            expect(userList.getUserCount()).toBe(4);
        });

        test('should get online user count', () => {
            expect(userList.getOnlineUserCount()).toBe(2);
        });
    });

    describe('WebSocket Integration', () => {
        test('should handle user update', () => {
            const user = { id: 'user1', username: 'Alice', status: 'online' };
            userList.addUser(user);

            const updatedUser = { id: 'user1', username: 'Alice', status: 'busy' };
            userList.handleUserUpdate(updatedUser);

            expect(userList.users[0].status).toBe('busy');
        });

        test('should handle user removal', () => {
            const user = { id: 'user1', username: 'Alice', status: 'online' };
            userList.addUser(user);

            userList.handleUserRemoved('user1');

            expect(userList.users).toHaveLength(0);
        });
    });

    describe('Error Handling', () => {
        test('should handle invalid user data', () => {
            const invalidUser = {
                id: 'user1',
                // Missing required fields
            };

            expect(() => {
                userList.addUser(invalidUser);
            }).not.toThrow();
        });

        test('should handle DOM manipulation errors', () => {
            // Remove the user list container
            userListContainer.remove();

            expect(() => {
                userList.updateUsers([
                    { id: 'user1', username: 'Alice', status: 'online' }
                ]);
            }).not.toThrow();
        });
    });

    describe('Performance', () => {
        test('should handle large number of users efficiently', () => {
            const startTime = performance.now();

            // Add 1000 users
            const users = [];
            for (let i = 0; i < 1000; i++) {
                users.push({
                    id: `user${i}`,
                    username: `User${i}`,
                    status: i % 2 === 0 ? 'online' : 'away'
                });
            }

            userList.updateUsers(users);

            const endTime = performance.now();
            const duration = endTime - startTime;

            expect(duration).toBeLessThan(1000); // Should complete within 1 second
            expect(userList.getUserCount()).toBe(1000);
        });

        test('should handle frequent user updates efficiently', () => {
            const users = [
                { id: 'user1', username: 'Alice', status: 'online' },
                { id: 'user2', username: 'Bob', status: 'away' }
            ];
            userList.updateUsers(users);

            const startTime = performance.now();

            // Update users 100 times
            for (let i = 0; i < 100; i++) {
                userList.setUserStatus('user1', i % 2 === 0 ? 'online' : 'away');
            }

            const endTime = performance.now();
            const duration = endTime - startTime;

            expect(duration).toBeLessThan(500); // Should complete within 500ms
        });
    });
}); 