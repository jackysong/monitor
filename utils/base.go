package utils

import (
	"math/rand"
	"time"
)

var Mode string

const (
	RANDOM_NUM   = 0 // 纯数字
	RANDOM_LOWER = 1 // 小写字母
	RANDOM_UPPER = 2 // 大写字母
	RANDOM_ALL   = 3 // 数字、大小写字母
)

//生成随机字符串
func RandomKey(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
