package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf16"
)

// decodeUCS2CharBE 解码 UCS2 大端2字节字符
func decodeUCS2CharBE(charBytes []byte) (string, error) {
	if len(charBytes) != 2 {
		return "", fmt.Errorf("invalid length for UCS2 char: %d", len(charBytes))
	}
	codeUnit := uint16(charBytes[0])<<8 | uint16(charBytes[1])
	runes := utf16.Decode([]uint16{codeUnit})
	if len(runes) == 0 {
		return "", fmt.Errorf("failed to decode UCS2 char")
	}
	return string(runes[0]), nil
}

// decodeUCS2WithAsciiEscapeBE 解码 UCS2 大端字节数组（包含 0x0a + ASCII 形式）
func decodeUCS2WithAsciiEscapeBE(data []byte) (string, error) {
	var output strings.Builder
	i := 0
	for i < len(data) {
		if data[i] == 0x0a && i+1 < len(data) {
			// 如果当前是转义标志 0x0a，且下一个字节存在
			output.WriteByte(data[i+1]) // 直接写入 ASCII 字符
			i += 2                      // 跳过两个字节
		} else if i+1 < len(data) {
			// 否则按两个字节解码为 UCS2 字符
			decodedChar, err := decodeUCS2CharBE(data[i : i+2])
			if err != nil {
				return "", err
			}
			output.WriteString(decodedChar)
			i += 2
		} else {
			break // 剩余 1 字节，不足以解码，退出
		}
	}
	return output.String(), nil
}

func cleanTail(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool {
		return !unicode.IsPrint(r) || r == '\u0000'
	})
}

func main() {
	var hexStr string
	if len(os.Args) > 1 {
		hexStr = os.Args[1]
	} else {
		// 没参数时，从标准输入读取一行
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			hexStr = scanner.Text()
		} else {
			fmt.Fprintln(os.Stderr, "no input data")
			os.Exit(1)
		}
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	result, err := decodeUCS2WithAsciiEscapeBE(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(cleanTail(result))
}
