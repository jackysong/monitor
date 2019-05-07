package supervisor

import (
	"strings"

	"github.com/kisekivul/utils"
)

//eg:ver:3.0 server:supervisor serial:21 pool:listener poolserial:10 eventname:PROCESS_COMMUNICATION_STDOUT len:54
type Head struct {
	msg        string
	Ver        string `json:"ver"`
	Server     string `json:"server"`
	Serial     int    `json:"serial"`
	Pool       string `json:"pool"`
	PoolSerial int    `json:"poolSerial"`
	State      State  `json:"state"`
	Len        int    `json:"len"`
}

//eg:processname:cat groupname:cat from_state:RUNNING expected:0 pid:2766
type Body struct {
	msg      string
	Process  string `json:"process"`
	Group    string `json:"group"`
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Expected int    `json:"expected"`
	Pid      int    `json:"pid"`
}

func parse(msg string) (data map[string]string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return
	}

	list := strings.Split(msg, " ")
	if len(list) == 0 {
		return
	}

	data = make(map[string]string)
	for _, item := range list {
		if group := strings.Split(item, ":"); len(group) < 2 {
			continue
		} else {
			key := strings.TrimSpace(group[0])
			value := strings.TrimSpace(group[1])
			data[key] = value
		}
	}
	return
}

func (h *Head) Parse() error {
	list := parse(h.msg)
	if list == nil {
		return errorParseHead
	}

	h.Ver = list["ver"]
	h.Server = list["server"]
	h.Serial = utils.Str2Int(list["serial"])
	h.Pool = list["pool"]
	h.PoolSerial = utils.Str2Int(list["poolserial"])
	h.State = State(list["eventname"])
	h.Len = utils.Str2Int(list["len"])

	return nil
}

func (b *Body) Parse() error {
	list := parse(b.msg)
	if list == nil {
		return errorParseBody
	}

	b.Process = list["processname"]
	b.Group = list["groupname"]
	b.Previous = list["from_state"]
	b.Expected = utils.Str2Int(list["expected"])
	b.Pid = utils.Str2Int(list["pid"])

	return nil
}
