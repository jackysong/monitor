package supervisor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kisekivul/utils"
)

type Program struct {
	Name  string `json:"name"`
	State State  `json:"state"`
	Last  string `json:"last"`
}

func reload() {
	output, err := utils.ExecShell("supervisorctl")
	if utils.ErrorCheck(err) {
		return
	}

	now := time.Now()
	for _, line := range strings.Split(strings.TrimSuffix(strings.TrimSpace(output), "supervisor>"), "\n") {
		if line != "" {
			temp := strings.Fields(line)
			//cache
			c.set(temp[0], getLast(now, line), getState(temp[1]))
		}
	}
}

func getLast(t time.Time, line string) string {
	if list := strings.Split(line, "uptime"); len(list) > 1 {
		var (
			h, m, s int
			date    string
		)

		if temp := strings.Split(list[1], "days,"); len(temp) > 1 {
			h = 24 * utils.Str2Int(strings.TrimSpace(temp[0]))
			date = temp[1]
		} else {
			date = temp[0]
		}

		ts := strings.Split(strings.TrimSpace(date), ":")
		h += utils.Str2Int(ts[0])
		m = utils.Str2Int(ts[1])
		s = utils.Str2Int(ts[2])

		t = t.Add(-1 * (time.Hour*time.Duration(h) + time.Minute*time.Duration(m) + time.Second*time.Duration(s)))
	}
	return t.Format("2006-01-02 15:04:05")
}

func listen(handler func(*Head, *Body, ...interface{}), params ...interface{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		ready()

		head, err := getHead(reader)
		if err != nil {
			fail(err)
			continue
		}

		body, err := getBody(reader, head.Len)
		if err != nil {
			fail(err)
			continue
		}

		update(head, body)

		handler(head, body, params)

		finish()
	}
}

func getHead(reader *bufio.Reader) (*Head, error) {
	msg, err := reader.ReadString('\n')
	if utils.ErrorCheck(err) {
		return nil, err
	}

	h := &Head{msg: msg}
	err = h.Parse()
	if utils.ErrorCheck(err) {
		return nil, err
	}
	return h, nil
}

func getBody(reader *bufio.Reader, length int) (*Body, error) {
	msg := make([]byte, length)
	l, err := reader.Read(msg)
	if utils.ErrorCheck(err) {
		return nil, err
	}

	if l != length {
		return nil, errorBodyLength
	}

	b := &Body{msg: string(msg)}
	err = b.Parse()
	if utils.ErrorCheck(err) {
		return nil, err
	}
	return b, nil
}

func update(head *Head, body *Body) {
	switch head.State {
	case ADDED:
		c.set(body.Group, time.Now().Format("2006-01-02 15:04:05"), ADDED)
	case REMOVED:
		c.del(body.Group)
	default:
		c.set(body.Process, time.Now().Format("2006-01-02 15:04:05"), head.State)
	}
}

func ready() {
	fmt.Fprint(os.Stdout, "READY\n")
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Fprint(os.Stdout, "Result 2\nFAIL")
}

func finish() {
	fmt.Fprint(os.Stdout, "RESULT 2\nOK")
}
