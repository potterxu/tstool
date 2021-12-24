package util

import (
	"fmt"
	"strings"
)

func ArrayInHex(data []byte) string {
	var sb strings.Builder
	sb.WriteString("[")
	for _, d := range data {
		sb.WriteString(" ")
		sb.WriteString(ByteInHex(d))
	}
	sb.WriteString(" ]")
	return sb.String()
}

func ByteInHex(data byte) string {
	return fmt.Sprintf("0x%02x", data)
}
