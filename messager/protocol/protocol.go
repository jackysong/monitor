package protocol

import "github.com/kisekivul/utils"

const (
	HEADER        = "Header"
	HEADER_LENGTH = 6
	ENCODE_LENGTH = 4
)

var (
	CODE        = "999999"
	CODE_LENGTH = 6
)

func Initialize(code string) {
	CODE = code
	CODE_LENGTH = len([]byte(code))
}

func Enpack(message []byte) []byte {
	return append(append(append(append([]byte(HEADER), utils.Int2Bytes(CODE_LENGTH)...), []byte(CODE)...), utils.Int2Bytes(len(message))...), message...)
}

func Depack(buffer []byte) ([][]byte, []byte, []byte) {
	var (
		i      int
		code   []byte
		list   [][]byte
		length = len(buffer)
	)

	for i = 0; i < length; {
		if length < i+HEADER_LENGTH+ENCODE_LENGTH {
			return list, code, buffer[i:]
		}

		if string(buffer[i:i+HEADER_LENGTH]) == HEADER {
			code_length := utils.Bytes2Int(buffer[i+HEADER_LENGTH : i+HEADER_LENGTH+ENCODE_LENGTH])
			if length < i+HEADER_LENGTH+ENCODE_LENGTH+code_length+ENCODE_LENGTH {
				break
			}
			code = buffer[i+HEADER_LENGTH+ENCODE_LENGTH : i+HEADER_LENGTH+ENCODE_LENGTH+code_length]

			data_length := utils.Bytes2Int(buffer[i+HEADER_LENGTH+ENCODE_LENGTH+code_length : i+HEADER_LENGTH+ENCODE_LENGTH+code_length+ENCODE_LENGTH])
			if length < i+HEADER_LENGTH+ENCODE_LENGTH+code_length+ENCODE_LENGTH+data_length {
				break
			}
			data := buffer[i+HEADER_LENGTH+ENCODE_LENGTH+code_length+ENCODE_LENGTH : i+HEADER_LENGTH+ENCODE_LENGTH+code_length+ENCODE_LENGTH+data_length]

			list = append(list, data)
			i = i + HEADER_LENGTH + ENCODE_LENGTH + code_length + ENCODE_LENGTH + data_length
		} else {
			i++
		}
	}
	return list, code, buffer[i:]
}
