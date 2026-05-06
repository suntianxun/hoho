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

type Client struct {
	Conn    net.Conn
	Encoder *json.Encoder
	Decoder *json.Decoder
}

func Connect(port string) (*Client, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", port))
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn:    conn,
		Encoder: json.NewEncoder(conn),
		Decoder: json.NewDecoder(conn),
	}, nil
}

func (c *Client) ReadState() (StateMessage, error) {
	var msg StateMessage
	err := c.Decoder.Decode(&msg)
	return msg, err
}

func (c *Client) SendCommand(cmd string) error {
	msg := CommandMessage{Type: "command", Cmd: cmd}
	return c.Encoder.Encode(msg)
}
