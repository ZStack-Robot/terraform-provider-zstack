package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/kataras/golog"
)

func GenRequestId(bytes int) string {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		golog.Errorf("Fail to generate Request Id: %s", err)
		return ""
	}
	return hex.EncodeToString(b)
}
