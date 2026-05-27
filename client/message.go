package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/chima/CLI_Chat/protocol"
)

func isTerminalAction(action protocol.ActionType) bool {
	return action == protocol.DONE || action == protocol.ERR
}

func isDisplayOnlyAction(action protocol.ActionType) bool {
	switch action {
	case protocol.MESSAGE_ROOM, protocol.SEND_MSG, protocol.MESSAGE_FRIEND:
		return true
	default:
		return false
	}
}

func needsUserReply(action protocol.ActionType) bool {
	return !isTerminalAction(action) && !isDisplayOnlyAction(action)
}

type fileSendCommand struct {
	filename string
	friend   string
}

func (e *fileSendCommand) Error() string {
	return "file send"
}

func parseFileSendCommand(line string) (*fileSendCommand, error) {
	const usage = "/file <friend> <filename>"
	fields := strings.Fields(line)
	if len(fields) == 0 || fields[0] != "/file" {
		return nil, nil
	}
	if len(fields) != 3 {
		return nil, fmt.Errorf("usage: %s (names cannot contain spaces)", usage)
	}
	return &fileSendCommand{friend: fields[1], filename: fields[2]}, nil
}

func buildMessage(line string) (*protocol.Message, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty input")
	}

	if fsc, err := parseFileSendCommand(line); fsc != nil || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fsc
	}

	if strings.HasPrefix(line, "/") {
		msg, err := buildCommandMessage(line)
		if errors.Is(err, errClientOnly) {
			return nil, errClientOnly
		}
		return msg, err
	}

	payload, _ := json.Marshal(map[string]string{"message": line})
	return &protocol.Message{
		Action:  protocol.SEND_MSG,
		Payload: payload,
	}, nil
}

func marshalStringPayload(fields map[string]string) json.RawMessage {
	payload, _ := json.Marshal(fields)
	return payload
}

func marshalAnyPayload(fields map[string]any) json.RawMessage {
	payload, _ := json.Marshal(fields)
	return payload
}

func requireArg(arg, usage string) error {
	if strings.TrimSpace(arg) == "" {
		return fmt.Errorf("usage: %s", usage)
	}
	return nil
}

func requireName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if strings.ContainsAny(name, " \t\r\n") {
		return fmt.Errorf("names cannot contain spaces")
	}
	return nil
}

func requireNameArg(arg, usage string) error {
	if err := requireArg(arg, usage); err != nil {
		return err
	}
	if err := requireName(arg); err != nil {
		return fmt.Errorf("usage: %s: %w", usage, err)
	}
	return nil
}

func splitFields(arg string) []string {
	return strings.Fields(strings.TrimSpace(arg))
}

func buildCommandMessage(line string) (*protocol.Message, error) {
	parts := strings.SplitN(line, " ", 2)
	cmd := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(parts[1])
	}

	switch cmd {
	case "/help":
		printClientHelp()
		return nil, errClientOnly

	// auth
	case "/auth/signup":
		if err := requireNameArg(arg, "/auth/signup <username>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.SIGN_UP,
			Payload: marshalStringPayload(map[string]string{"username": arg}),
		}, nil

	case "/auth/login":
		if err := requireNameArg(arg, "/auth/login <name>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.LOGIN,
			Payload: marshalStringPayload(map[string]string{"username": arg}),
		}, nil

	case "/auth/logout":
		return &protocol.Message{Action: protocol.LOGOUT}, nil

	// room
	case "/room/create":
		if err := requireNameArg(arg, "/room/create <name>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.CREATE_ROOM,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/room/join":
		if err := requireNameArg(arg, "/room/join <name>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.JOIN_ROOM,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/room/leave":
		payload := map[string]string{}
		if arg != "" {
			if err := requireNameArg(arg, "/room/leave [<name>]"); err != nil {
				return nil, err
			}
			payload["room"] = arg
		}
		return &protocol.Message{
			Action:  protocol.LEAVE_ROOM,
			Payload: marshalStringPayload(payload),
		}, nil

	case "/room/delete":
		if err := requireNameArg(arg, "/room/delete <name>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.DELETE_ROOM,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/room/edit":
		fields := splitFields(arg)
		if len(fields) < 2 || len(fields) > 3 {
			return nil, fmt.Errorf("usage: /room/edit <room> <new_name> [max_size] (names cannot contain spaces)")
		}
		if err := requireName(fields[0]); err != nil {
			return nil, fmt.Errorf("usage: /room/edit <room> <new_name> [max_size]: %w", err)
		}
		if err := requireName(fields[1]); err != nil {
			return nil, fmt.Errorf("usage: /room/edit <room> <new_name> [max_size]: %w", err)
		}
		payload := map[string]any{
			"room":     fields[0],
			"new_room": fields[1],
		}
		if len(fields) >= 3 {
			maxSize, err := strconv.Atoi(fields[2])
			if err != nil {
				return nil, fmt.Errorf("max_size must be a number")
			}
			payload["max_size"] = maxSize
		}
		return &protocol.Message{
			Action:  protocol.EDIT_ROOM,
			Payload: marshalAnyPayload(payload),
		}, nil

	case "/room/members":
		if err := requireNameArg(arg, "/room/members <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.GET_ROOM_MEMBERS,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/room/kick":
		fields := splitFields(arg)
		if len(fields) != 2 {
			return nil, fmt.Errorf("usage: /room/kick <room> <user> (names cannot contain spaces)")
		}
		return &protocol.Message{
			Action: protocol.DELETE_MEMBER,
			Payload: marshalStringPayload(map[string]string{
				"room": fields[0],
				"nick": fields[1],
			}),
		}, nil

	case "/room/invite":
		fields := splitFields(arg)
		if len(fields) != 2 {
			return nil, fmt.Errorf("usage: /room/invite <room> <user> (names cannot contain spaces)")
		}
		return &protocol.Message{
			Action: protocol.SEND_INVITE_REQUEST,
			Payload: marshalStringPayload(map[string]string{
				"room": fields[0],
				"nick": fields[1],
			}),
		}, nil

	case "/room/invites":
		payload := map[string]string{}
		if arg != "" {
			if err := requireNameArg(arg, "/room/invites [<room>]"); err != nil {
				return nil, err
			}
			payload["room"] = arg
		}
		return &protocol.Message{
			Action:  protocol.SEE_GROUP_INVITE_REQUEST,
			Payload: marshalStringPayload(payload),
		}, nil

	case "/room/list":
		return &protocol.Message{Action: protocol.LIST_ROOMS}, nil

	case "/room/mine":
		return &protocol.Message{Action: protocol.LIST_MY_ROOMS}, nil

	case "/invite/mine":
		return &protocol.Message{Action: protocol.LIST_MY_ROOM_INVITES}, nil

	case "/invite/accept":
		if err := requireNameArg(arg, "/invite/accept <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.ACCEPT_GROUP_INVITE_REQUEST,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/invite/decline":
		if err := requireNameArg(arg, "/invite/decline <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.DELETE_GROUP_INVITE_REQUEST,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	// chat
	case "/switch":
		fields := splitFields(arg)
		if len(fields) != 2 {
			return nil, fmt.Errorf("usage: /switch <room|friend> <name> (names cannot contain spaces)")
		}
		if err := requireName(fields[1]); err != nil {
			return nil, fmt.Errorf("usage: /switch <room|friend> <name>: %w", err)
		}
		payload := map[string]string{"context": strings.ToLower(fields[0])}
		switch strings.ToLower(fields[0]) {
		case "room":
			payload["room"] = fields[1]
		case "friend":
			payload["friend_name"] = fields[1]
		default:
			return nil, fmt.Errorf("context must be room or friend")
		}
		return &protocol.Message{
			Action:  protocol.SWITCH_CONTEXT,
			Payload: marshalStringPayload(payload),
		}, nil

	case "/chat/room":
		fields := splitFields(arg)
		if len(fields) < 2 {
			return nil, fmt.Errorf("usage: /chat/room <room> <message>")
		}
		if err := requireName(fields[0]); err != nil {
			return nil, fmt.Errorf("usage: /chat/room <room> <message>: %w", err)
		}
		return &protocol.Message{
			Action: protocol.MESSAGE_ROOM,
			Payload: marshalStringPayload(map[string]string{
				"room":    fields[0],
				"message": strings.Join(fields[1:], " "),
			}),
		}, nil

	case "/chat/file":
		if err := requireNameArg(arg, "/chat/file <path>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.SEND_FILE,
			Payload: marshalStringPayload(map[string]string{"path": arg}),
		}, nil

	case "/chat/dm":
		fields := splitFields(arg)
		if len(fields) < 2 {
			return nil, fmt.Errorf("usage: /chat/dm <friend> <message>")
		}
		if err := requireName(fields[0]); err != nil {
			return nil, fmt.Errorf("usage: /chat/dm <friend> <message>: %w", err)
		}
		return &protocol.Message{
			Action: protocol.MESSAGE_FRIEND,
			Payload: marshalStringPayload(map[string]string{
				"nick":    fields[0],
				"message": strings.Join(fields[1:], " "),
			}),
		}, nil

	// friend
	case "/friend/add":
		if err := requireNameArg(arg, "/friend/add <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.SEND_FRIEND_REQUEST,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/friend/accept":
		if err := requireNameArg(arg, "/friend/accept <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.ACCEPT_FRIEND_REQUEST,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/friend/remove":
		if err := requireNameArg(arg, "/friend/remove <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.DELETE_FRIEND,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/friend/requests":
		return &protocol.Message{Action: protocol.SEE_FRIEND_REQUEST}, nil

	case "/friend/list":
		return &protocol.Message{Action: protocol.GET_FRIENDS}, nil

	// user
	case "/user/block":
		if err := requireNameArg(arg, "/user/block <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.BLOCK_USER,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/user/unblock":
		if err := requireNameArg(arg, "/user/unblock <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.UNBLOCK_USER,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/user/status":
		if err := requireNameArg(arg, "/user/status <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.GET_USER_STATUS,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	default:
		return nil, fmt.Errorf("unknown command: %s (type /help)", cmd)
	}
}

func buildPromptReply(action protocol.ActionType, serverPayload json.RawMessage, input string) (*protocol.Message, error) {
	input = strings.TrimSpace(input)
	if input == "" && action != protocol.SET_ROOM_PASSWORD {
		return nil, fmt.Errorf("empty input")
	}

	fields := map[string]string{}
	if len(serverPayload) > 0 {
		_ = json.Unmarshal(serverPayload, &fields)
	}

	switch action {
	case protocol.LOGIN:
		if err := requireName(input); err != nil {
			return nil, err
		}
		fields["username"] = input
	case protocol.SET_ROOM_PASSWORD, protocol.GET_ROOM_PASSWORD:
		fields["password"] = input
	case protocol.MESSAGE_ROOM, protocol.SEND_MSG, protocol.MESSAGE_FRIEND:
		fields["message"] = input
	default:
		if _, ok := fields["message"]; !ok {
			fields["message"] = input
		}
	}

	payload, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	return &protocol.Message{
		Action:  action,
		Payload: payload,
	}, nil
}
