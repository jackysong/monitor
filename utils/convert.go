package utils

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
	"time"
)

func Float642Str(f float64) string {
	return strconv.FormatFloat(f, 'g', 4, 64)
}

func Str2Float64(s string) float64 {
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return t
}

func Str2Int64(s string) int64 {
	t, err := strconv.ParseInt(s, 10, 64)
	ErrorCheck(err)
	return t
}

func Str2Int(str string) int {
	t, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return t
}

func Int2Str(t int) string {
	str := strconv.Itoa(t)
	return str
}

func Int642Str(t int64) string {
	str := strconv.FormatInt(int64(t), 10)
	return str
}

func Date2Str(date int64) string {
	str := time.Unix(date, 0).Format("2006-01-02 15:04:05")
	return str
}

func Object2Int(object interface{}) int {
	return int(Object2Int64(object))
}

func Object2Int64(object interface{}) int64 {
	if object == nil {
		return 0
	}

	var res int64
	switch reflect.TypeOf(object).String() {
	case "[]uint8":
		res = Str2Int64(Object2Str(object))
	case "int64":
		res = reflect.ValueOf(object).Int()
	}
	return res
}

func Object2Str(object interface{}) string {
	if object == nil {
		return ""
	}
	return string(reflect.ValueOf(object).Bytes())
}

func Str2Unix(date string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
	if ErrorCheck(err) {
		return 0, err
	}
	return t.Unix(), nil
}

func BytesToStr(b []byte) string {
	return bytes.NewBuffer(b).String()
}

func Str2Bool(str string) bool {
	t, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return t
}

func List2Map(list []interface{}) map[interface{}]interface{} {
	dic := make(map[interface{}]interface{})
	for _, value := range list {
		dic[value] = value
	}
	return dic
}

func Int2Bytes(n int) []byte {
	x := int32(n)

	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, x)
	return buff.Bytes()
}

func Bytes2Int(b []byte) int {
	buff := bytes.NewBuffer(b)

	var x int32
	binary.Read(buff, binary.BigEndian, &x)

	return int(x)
}
