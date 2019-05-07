package common

import "encoding/json"

type Action int

const (
	Registe Action = iota
	Offline
	Message
	Warning
)

var actions = [...]string{
	"Registe",
	"Offline",
	"Message",
	"Warning",
}

func (a Action) String() string {
	return actions[a]
}

type Unit struct {
	Code string `json:"code"`
	Data string `json:"data"`
}

func Analyze(data []byte) (*Packet, error) {
	var packet Packet
	err := json.Unmarshal(data, &packet)

	return &packet, err
}
