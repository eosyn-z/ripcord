package main

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"github.com/google/uuid"
	"ripcord/security"
	"ripcord/types"
)

type MessageHandler struct {
	cryptoManager *security.CryptoManager
}

func NewMessageHandler(cryptoManager *security.CryptoManager) *MessageHandler {
	return &MessageHandler{
		cryptoManager: cryptoManager,
	}
}

func NewMessage(roomID, userID, username, content, msgType string) *types.Message {
	msgID := generateMessageID()
	
	if msgType == "" {
		if strings.HasPrefix(content, "/") {
			msgType = types.MessageTypeCommand
		} else {
			msgType = types.MessageTypeText
		}
	}
	
	return &types.Message{
		ID:        msgID,
		RoomID:    roomID,
		UserID:    userID,
		Username:  username,
		Content:   content,
		Type:      msgType,
		Encrypted: false,
		Timestamp: time.Now(),
	}
}

func generateMessageID() string {
	return uuid.New().String()
}

func (mh *MessageHandler) CreateSignedMessage(roomID, userID, username, content, msgType string) (*types.Message, error) {
	msg := NewMessage(roomID, userID, username, content, msgType)
	
	privateKey := mh.cryptoManager.GetPrivateKey()
	if err := msg.Sign(privateKey); err != nil {
		return nil, err
	}
	
	return msg, nil
}


func (mh *MessageHandler) ProcessSlashCommand(msg *types.Message) (*ProtocolMessage, error) {
	cmd, err := msg.ParseSlashCommand()
	if err != nil {
		return nil, err
	}
	
	switch cmd.Command {
	case "join":
		return mh.handleJoinCommand(msg, cmd.Args)
	case "leave":
		return mh.handleLeaveCommand(msg, cmd.Args)
	case "invite":
		return mh.handleInviteCommand(msg, cmd.Args)
	case "block":
		return mh.handleBlockCommand(msg, cmd.Args)
	case "unblock":
		return mh.handleUnblockCommand(msg, cmd.Args)
	case "dm":
		return mh.handleDMCommand(msg, cmd.Args)
	default:
		return nil, errors.New("unknown command: " + cmd.Command)
	}
}

func (mh *MessageHandler) handleJoinCommand(msg *types.Message, args []string) (*ProtocolMessage, error) {
	if len(args) < 1 {
		return nil, errors.New("join command requires invite code")
	}
	
	inviteCode := args[0]
	protocolMsg := NewProtocolMessage(MessageTypeJoin, msg.UserID, generateMessageID())
	protocolMsg.SetPayload(JoinPayload{
		InviteCode: inviteCode,
		Nickname:   msg.Username,
		PublicKey:  hex.EncodeToString(mh.cryptoManager.GetPublicKey()),
	})
	
	return protocolMsg, nil
}

func (mh *MessageHandler) handleLeaveCommand(msg *types.Message, args []string) (*ProtocolMessage, error) {
	reason := ""
	if len(args) > 0 {
		reason = strings.Join(args, " ")
	}
	
	protocolMsg := NewProtocolMessage(MessageTypeLeave, msg.UserID, generateMessageID())
	protocolMsg.RoomID = msg.RoomID
	protocolMsg.SetPayload(LeavePayload{
		RoomID: msg.RoomID,
		Reason: reason,
	})
	
	return protocolMsg, nil
}

func (mh *MessageHandler) handleInviteCommand(msg *types.Message, args []string) (*ProtocolMessage, error) {
	if len(args) < 1 {
		return nil, errors.New("invite command requires username")
	}
	
	username := args[0]
	protocolMsg := NewProtocolMessage(MessageTypeInvite, msg.UserID, generateMessageID())
	protocolMsg.To = username
	protocolMsg.RoomID = msg.RoomID
	
	return protocolMsg, nil
}

func (mh *MessageHandler) handleBlockCommand(msg *types.Message, args []string) (*ProtocolMessage, error) {
	if len(args) < 1 {
		return nil, errors.New("block command requires username")
	}
	
	username := args[0]
	protocolMsg := NewProtocolMessage(MessageTypeBlock, msg.UserID, generateMessageID())
	protocolMsg.SetPayload(UserInfoPayload{
		Nickname:  username,
		IsBlocked: true,
	})
	
	return protocolMsg, nil
}

func (mh *MessageHandler) handleUnblockCommand(msg *types.Message, args []string) (*ProtocolMessage, error) {
	if len(args) < 1 {
		return nil, errors.New("unblock command requires username")
	}
	
	username := args[0]
	protocolMsg := NewProtocolMessage(MessageTypeUnblock, msg.UserID, generateMessageID())
	protocolMsg.SetPayload(UserInfoPayload{
		Nickname:  username,
		IsBlocked: false,
	})
	
	return protocolMsg, nil
}

func (mh *MessageHandler) handleDMCommand(msg *types.Message, args []string) (*ProtocolMessage, error) {
	if len(args) < 2 {
		return nil, errors.New("dm command requires username and message")
	}
	
	username := args[0]
	content := strings.Join(args[1:], " ")
	
	protocolMsg := NewProtocolMessage(MessageTypeDM, msg.UserID, generateMessageID())
	protocolMsg.To = username
	protocolMsg.SetPayload(DMPayload{
		Content:     content,
		IsEncrypted: false,
	})
	
	return protocolMsg, nil
}

 