package biz

import (
	"math/rand"
	"strings"
	"time"
)

var (
	globalRand *rand.Rand
)

var characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func init() {
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	characters = shuffleCharacters(characters)
}

// 生成短链接核心函数
func generateShortUrl(decimal int64) string {
	return base62Encode(decimal)
}

// base62 编码, 将数字编码为 base62 字符串
func base62Encode(decimal int64) string {
	var result strings.Builder
	for decimal > 0 {
		remainder := decimal % 62
		result.WriteByte(characters[remainder])
		decimal /= 62
	}
	return result.String()
}

func shuffleCharacters(input string) string {
	chars := strings.Split(input, "")
	for i := len(chars) - 1; i > 0; i-- {
		j := globalRand.Intn(i + 1)
		chars[i], chars[j] = chars[j], chars[i] // 随机交换
	}
	return strings.Join(chars, "")
}
