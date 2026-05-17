package client

import (
	"encoding/json"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/chima/CLI_Chat/protocol"
	"github.com/chzyer/readline"
)

type session struct {
	conn     net.Conn
	rl       *readline.Instance
	incoming chan protocol.Message

	promptMu sync.Mutex
	prompt   *protocol.Message
}

func runSession(conn net.Conn, rl *readline.Instance) {
	s := &session{
		conn:     conn,
		rl:       rl,
		incoming: make(chan protocol.Message, 16),
	}

	go s.readFromServer()
	go s.processServerMessages()
	s.inputLoop()
}

func (s *session) readFromServer() {
	defer close(s.incoming)

	decoder := json.NewDecoder(s.conn)
	for {
		var msg protocol.Message
		if err := decoder.Decode(&msg); err != nil {
			if err != io.EOF {
				printTerminal("\nConnection closed.\n")
			}
			return
		}
		s.incoming <- msg
	}
}

func (s *session) processServerMessages() {
	for msg := range s.incoming {
		s.handleServerMessage(msg)
	}
}

func (s *session) handleServerMessage(msg protocol.Message) {
	if needsUserReply(msg.Action) {
		s.clearCurrentInput()
	}

	if msg.ResponseMsg != "" && (isTerminalAction(msg.Action) || isDisplayOnlyAction(msg.Action) || needsUserReply(msg.Action)) {
		printTerminal("%s\n", msg.ResponseMsg)
	}

	switch {
	case msg.Action == protocol.LOGOUT:
		s.clearPrompt()
	case isTerminalAction(msg.Action):
		s.clearPrompt()
	case isDisplayOnlyAction(msg.Action):
		// response already printed
	case needsUserReply(msg.Action):
		s.setPrompt(&msg)
	default:
		if msg.ResponseMsg == "" {
			printTerminal("[server action %d]\n", msg.Action)
		}
	}
}

func (s *session) setPrompt(msg *protocol.Message) {
	s.promptMu.Lock()
	s.prompt = msg
	s.promptMu.Unlock()
}

func (s *session) clearPrompt() {
	s.promptMu.Lock()
	s.prompt = nil
	s.promptMu.Unlock()
}

func (s *session) pendingPrompt() *protocol.Message {
	s.promptMu.Lock()
	defer s.promptMu.Unlock()
	return s.prompt
}

// clearCurrentInput discards whatever the user is typing (readline buffer).
func (s *session) clearCurrentInput() {
	_, _ = s.rl.WriteStdin([]byte{
		byte(readline.CharLineEnd),
		byte(readline.CharCtrlU),
		byte(readline.CharKill),
	})
}

func (s *session) clearSubmittedInput() {
	printTerminal("\033[1A\033[2K")
	s.rl.Refresh()
}

func (s *session) inputLoop() {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if prompt := s.pendingPrompt(); prompt != nil {
			if err := s.sendPromptReply(prompt, line); err != nil {
				printTerminal("%v\n", err)
			} else {
				s.clearSubmittedInput()
			}
			continue
		}

		msg, err := buildMessage(line)
		if err != nil {
			printTerminal("%v\n", err)
			continue
		}

		if err := writeMessage(s.conn, msg); err != nil {
			printTerminal("Failed to send: %v\n", err)
			return
		}
		s.clearSubmittedInput()
		if msg.Action == protocol.LOGOUT {
			return
		}
	}
}

func (s *session) sendPromptReply(prompt *protocol.Message, input string) error {
	msg, err := buildPromptReply(prompt.Action, prompt.Payload, input)
	if err != nil {
		return err
	}
	s.clearPrompt()
	return writeMessage(s.conn, msg)
}

func writeMessage(conn net.Conn, msg *protocol.Message) error {
	return json.NewEncoder(conn).Encode(msg)
}
