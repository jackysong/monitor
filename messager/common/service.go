package common

import (
	"strconv"
	"time"

	"github.com/kisekivul/messager/protocol"
	"github.com/kisekivul/utils"
)

func Now() int {
	return utils.Str2Int(strconv.FormatInt(time.Now().Unix(), 10))
}

func Send(conn Connection, data []byte) error {
	_, err := conn.Conn.Write(protocol.Enpack((data)))
	return err
}
