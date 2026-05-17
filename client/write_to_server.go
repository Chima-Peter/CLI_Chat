package client

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/chima/CLI_Chat/protocol"
	"github.com/chzyer/readline"
)

func write_to_server(conn net.Conn, rl *readline.Instance) {
	encoder := json.NewEncoder(conn)

	for {
		line, err := rl.Readline()
		if err != nil {
			rl.Clean()
			_ = encoder.Encode(&protocol.Message{Action: protocol.LOGOUT})
			conn.Close()
			return
		}

		fmt.Fprint(rl.Stdout(), "\033[A\033[2K")

		msg, err := buildMessage(line)
		if err != nil {
			fmt.Fprintf(rl.Stdout(), "> ERR: %s\n", err.Error())
			rl.Refresh()
			continue
		}

		if err := encoder.Encode(msg); err != nil {
			return
		}
	}
}
