package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/chima/CLI_Chat/protocol"
	"github.com/chzyer/readline"
)

const (
	defaultPrompt = "> "
	clearTerminalSeq = "\033[H\033[2J"
)

type session struct {
	conn     net.Conn
	rl       *readline.Instance
	incoming chan protocol.Message

	displayMu sync.Mutex
	promptMu  sync.Mutex
	prompt    *protocol.Message
}

func runSession(conn net.Conn, rl *readline.Instance) {
	s := &session{
		conn:     conn,
		rl:       rl,
		incoming: make(chan protocol.Message, 16),
	}

	done := make(chan struct{})

	go s.readFromServer(done)
	go s.processServerMessages()
	go func() {
		<-done
		s.rl.Close()
	}()
	s.inputLoop()
}

func (s *session) readFromServer(done chan struct{}) {
	defer close(s.incoming)

	decoder := json.NewDecoder(s.conn)
	for {
		var msg protocol.Message
		if err := decoder.Decode(&msg); err != nil {
			if err != io.EOF {
				s.writeDisplay("Connection closed")
				close(done)
				return
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
	s.displayMu.Lock()
	defer s.displayMu.Unlock()

	switch {
	case msg.Action == protocol.CREATE_FILE_PORT:
		go s.runFileServer()
	case msg.Action == protocol.FILE_PORT_LISTENING:
		var meta struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			Nick     string `json:"nick"`
			Friend   string `json:"friend"`
			Filepath string `json:"filepath"`
		}
		_ = json.Unmarshal(msg.Payload, &meta)
		port, parseErr := strconv.Atoi(meta.Port)
		if parseErr != nil {
			s.writeDisplay("Invalid port received from server: " + meta.Port)
			return
		}
		go func() {
			if err := CreateFileClient(meta.Host, port, meta.Nick, meta.Friend, meta.Filepath); err != nil {
				s.writeDisplay(fmt.Sprintf("Could not send file to %s: %s", meta.Friend, err.Error()))
				return
			}
			s.writeDisplay("File sent to " + meta.Friend + ".")
		}()
	case needsUserReply(msg.Action):
		s.clearCurrentInputLocked()
		s.setPromptLocked(&msg)
	case isTerminalAction(msg.Action):
		if msg.ResponseMsg != "" {
			s.writeDisplayLocked(msg.ResponseMsg)
		}
	case isDisplayOnlyAction(msg.Action):
		if msg.ResponseMsg != "" {
			s.writeDisplayLocked(msg.ResponseMsg)
		}
	default:
		if msg.ResponseMsg == "" {
			s.writeDisplayLocked(fmt.Sprintf("[server action %d]", msg.Action))
		}
	}

	s.rl.Refresh()
}

func (s *session) writeDisplay(msg string) {
	s.displayMu.Lock()
	defer s.displayMu.Unlock()
	s.writeDisplayLocked(msg)
	s.rl.Refresh()
}

func (s *session) writeDisplayLocked(msg string) {
	if msg == "" {
		return
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	pending := s.pendingPrompt()
	s.rl.SetPrompt("")
	_, _ = s.rl.Write([]byte("> " + msg))
	if pending != nil {
		s.setPromptLocked(pending)
	}
}

func (s *session) clearTerminal() {
	s.displayMu.Lock()
	defer s.displayMu.Unlock()
	fmt.Print(clearTerminalSeq)
	s.rl.SetPrompt(defaultPrompt)
	s.rl.Refresh()
}

func (s *session) clearSubmittedInput() {
	s.displayMu.Lock()
	defer s.displayMu.Unlock()
	fmt.Fprint(s.rl.Stdout(), "\033[A\033[2K")
	s.rl.SetPrompt(defaultPrompt)
	s.rl.Refresh()
}

func (s *session) setPromptLocked(msg *protocol.Message) {
	s.promptMu.Lock()
	s.prompt = msg
	s.promptMu.Unlock()

	prompt := defaultPrompt
	if msg != nil && msg.ResponseMsg != "" {
		prompt += msg.ResponseMsg
	}
	s.rl.SetPrompt(prompt)
}

func (s *session) clearPromptLocked() {
	s.promptMu.Lock()
	s.prompt = nil
	s.promptMu.Unlock()
	s.rl.SetPrompt(defaultPrompt)
}

func (s *session) pendingPrompt() *protocol.Message {
	s.promptMu.Lock()
	defer s.promptMu.Unlock()
	return s.prompt
}

func (s *session) clearCurrentInputLocked() {
	_, _ = s.rl.WriteStdin([]byte{
		byte(readline.CharLineEnd),
		byte(readline.CharCtrlU),
		byte(readline.CharKill),
	})
	s.rl.Refresh()
}

func (s *session) inputLoop() {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			msg, _ := buildCommandMessage("/quit")
			if err := writeMessage(s.conn, msg); err != nil {
				return
			}
			return
		}

		if prompt := s.pendingPrompt(); prompt != nil {
			if err := s.sendPromptReply(prompt, line); err != nil {
				s.writeDisplay(err.Error())
			}
			continue
		}

		if line == "clear" {
			s.clearTerminal()
			continue
		}

		msg, err := buildMessage(line)
		if fsc, ok := errors.AsType[*fileSendCommand](err); ok {
			if sendErr := s.sendFileMetadata(fsc.filename, fsc.friend); sendErr != nil {
				s.writeDisplay(fmt.Sprintf("Could not send %q to %s: %s", fsc.filename, fsc.friend, sendErr.Error()))
			} else {
				s.writeDisplay("Sending " + fsc.filename + " to " + fsc.friend + "...")
			}
			s.clearSubmittedInput()
			continue
		}
		if err != nil {
			if errors.Is(err, errClientOnly) {
				continue
			}
			s.writeDisplay(err.Error())
			continue
		}
		if msg == nil {
			continue
		}

		if msg.Action == protocol.LOGOUT {
			if err := writeMessage(s.conn, msg); err != nil {
				s.writeDisplay(err.Error())
				return
			}
			return
		}

		if err := writeMessage(s.conn, msg); err != nil {
			s.writeDisplay(err.Error())
			return
		}
		s.clearSubmittedInput()
	}
}

func (s *session) sendPromptReply(prompt *protocol.Message, input string) error {
	msg, err := buildPromptReply(prompt.Action, prompt.Payload, input)
	if err != nil {
		return err
	}

	s.displayMu.Lock()
	s.clearPromptLocked()
	s.rl.Refresh()
	s.displayMu.Unlock()

	return writeMessage(s.conn, msg)
}

func writeMessage(conn net.Conn, msg *protocol.Message) error {
	return json.NewEncoder(conn).Encode(msg)
}

func (s *session) runFileServer() {
	err := StartFileServer(func(host string, port int) error {
		return writeMessage(s.conn, &protocol.Message{
			Action: protocol.FILE_PORT_LISTENING,
			Payload: marshalStringPayload(map[string]string{
				"host": host,
				"port": strconv.Itoa(port),
			}),
		})
	})
	if err != nil {
		s.writeDisplay("Could not receive files: " + err.Error())
	}
}

func (s *session) sendFileMetadata(filename string, friend string) error {
	zipPath, err := ZipFile(filename)
	if err != nil {
		return err
	}

	msg := &protocol.Message{
		Action: protocol.SEND_FILE,
		Payload: marshalStringPayload(map[string]string{
			"filepath": zipPath,
			"username": friend,
		}),
	}
	if err := writeMessage(s.conn, msg); err != nil {
		return fmt.Errorf("could not reach the server: %w", err)
	}
	return nil
}
