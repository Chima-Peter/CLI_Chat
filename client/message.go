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
	case protocol.SEND_MSG, protocol.MESSAGE_FRIEND:
		return true
	default:
		return false
	}
}

func needsUserReply(action protocol.ActionType) bool {
	return !isTerminalAction(action) && !isDisplayOnlyAction(action)
}

func buildMessage(line string) (*protocol.Message, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty input")
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

	case "/signup":
		if err := requireArg(arg, "/signup <username>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.SIGN_UP,
			Payload: marshalStringPayload(map[string]string{"username": arg}),
		}, nil

	case "/nick":
		if err := requireArg(arg, "/nick <name>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.LOGIN,
			Payload: marshalStringPayload(map[string]string{"username": arg}),
		}, nil

	case "/quit":
		return &protocol.Message{Action: protocol.LOGOUT}, nil

	case "/create":
		if err := requireArg(arg, "/create <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.CREATE_ROOM,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/join":
		if err := requireArg(arg, "/join <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.JOIN_ROOM,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/leave":
		payload := map[string]string{}
		if arg != "" {
			payload["room"] = arg
		}
		return &protocol.Message{
			Action:  protocol.LEAVE_ROOM,
			Payload: marshalStringPayload(payload),
		}, nil

	case "/deleteroom":
		if err := requireArg(arg, "/deleteroom <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.DELETE_ROOM,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/editroom":
		fields := splitFields(arg)
		if len(fields) < 2 {
			return nil, fmt.Errorf("usage: /editroom <room> <new_name> [max_size]")
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

	case "/members":
		if err := requireArg(arg, "/members <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.GET_ROOM_MEMBERS,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/kick":
		fields := splitFields(arg)
		if len(fields) < 2 {
			return nil, fmt.Errorf("usage: /kick <room> <user>")
		}
		return &protocol.Message{
			Action: protocol.DELETE_MEMBER,
			Payload: marshalStringPayload(map[string]string{
				"room": fields[0],
				"nick": fields[1],
			}),
		}, nil

	case "/invite":
		fields := splitFields(arg)
		if len(fields) < 2 {
			return nil, fmt.Errorf("usage: /invite <room> <user>")
		}
		return &protocol.Message{
			Action: protocol.SEND_INVITE_REQUEST,
			Payload: marshalStringPayload(map[string]string{
				"room": fields[0],
				"nick": fields[1],
			}),
		}, nil

	case "/invites":
		payload := map[string]string{}
		if arg != "" {
			payload["room"] = arg
		}
		return &protocol.Message{
			Action:  protocol.SEE_GROUP_INVITE_REQUEST,
			Payload: marshalStringPayload(payload),
		}, nil

	case "/acceptinvite":
		if err := requireArg(arg, "/acceptinvite <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.ACCEPT_GROUP_INVITE_REQUEST,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/declineinvite":
		if err := requireArg(arg, "/declineinvite <room>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.DELETE_GROUP_INVITE_REQUEST,
			Payload: marshalStringPayload(map[string]string{"room": arg}),
		}, nil

	case "/rooms":
		return &protocol.Message{Action: protocol.LIST_ROOMS}, nil

	case "/myrooms":
		return &protocol.Message{Action: protocol.LIST_MY_ROOMS}, nil

	case "/msg":
		if err := requireArg(arg, "/msg <message>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.SEND_MSG,
			Payload: marshalStringPayload(map[string]string{"message": arg}),
		}, nil

	case "/sendfile":
		if err := requireArg(arg, "/sendfile <path>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.SEND_FILE,
			Payload: marshalStringPayload(map[string]string{"path": arg}),
		}, nil

	case "/friendadd":
		if err := requireArg(arg, "/friendadd <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.SEND_FRIEND_REQUEST,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/friendaccept":
		if err := requireArg(arg, "/friendaccept <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.ACCEPT_FRIEND_REQUEST,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/friendrequests":
		return &protocol.Message{Action: protocol.SEE_FRIEND_REQUEST}, nil

	case "/friends":
		return &protocol.Message{Action: protocol.GET_FRIENDS}, nil

	case "/friendremove":
		if err := requireArg(arg, "/friendremove <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.DELETE_FRIEND,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/dm":
		fields := splitFields(arg)
		if len(fields) < 2 {
			return nil, fmt.Errorf("usage: /dm <user> <message>")
		}
		return &protocol.Message{
			Action: protocol.MESSAGE_FRIEND,
			Payload: marshalStringPayload(map[string]string{
				"nick":    fields[0],
				"message": strings.Join(fields[1:], " "),
			}),
		}, nil

	case "/block":
		if err := requireArg(arg, "/block <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.BLOCK_USER,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/unblock":
		if err := requireArg(arg, "/unblock <user>"); err != nil {
			return nil, err
		}
		return &protocol.Message{
			Action:  protocol.UNBLOCK_USER,
			Payload: marshalStringPayload(map[string]string{"nick": arg}),
		}, nil

	case "/status":
		if err := requireArg(arg, "/status <user>"); err != nil {
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
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	fields := map[string]string{}
	if len(serverPayload) > 0 {
		_ = json.Unmarshal(serverPayload, &fields)
	}

	switch action {
	case protocol.LOGIN:
		fields["username"] = input
	case protocol.SET_ROOM_PASSWORD, protocol.GET_ROOM_PASSWORD:
		fields["password"] = input
	case protocol.SEND_MSG, protocol.MESSAGE_FRIEND:
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

func printTerminal(format string, args ...any) {
	fmt.Printf(format, args...)
}
