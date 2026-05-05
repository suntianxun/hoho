package client

import (
	"encoding/json"
	"net"
)

type VarState struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	IsPrivate bool   `json:"is_private"`
}

type StateMessage struct {
	Type string     `json:"type"`
	Vars []VarState `json:"vars"`
}

type CommandMessage struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
}

func Connect(port string) (net.Conn, error) {
	return net.Dial("tcp", "127.0.0.1:"+port)
}

func ReadState(conn net.Conn) (StateMessage, error) {
	var msg StateMessage
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&msg)
	return msg, err
}

func SendCommand(conn net.Conn, cmd string) error {
	msg := CommandMessage{Type: "command", Cmd: cmd}
	encoder := json.NewEncoder(conn)
	return encoder.Encode(msg)
}
