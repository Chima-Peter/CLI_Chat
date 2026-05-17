package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/chima/CLI_Chat/protocol"
	"github.com/chzyer/readline"
)

func read_from_server(conn net.Conn, rl *readline.Instance) {
	decoder := json.NewDecoder(conn)

	for {
		var response protocol.Message
		if err := decoder.Decode(&response); err != nil {
			conn.Close()
			return
		}

		rl.Clean()
		line := strings.TrimSpace(response.ResponseMsg)
		if line == "" {
			line = fmt.Sprintf("[action %d]", response.Action)
		}
		fmt.Fprintf(rl.Stdout(), "> %s\n", line)
		rl.Refresh()
	}
}
