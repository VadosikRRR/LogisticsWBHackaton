package postgres

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func newRequestID() (string, error) {
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return "", fmt.Errorf("generate request id: %w", err)
	}
	return fmt.Sprintf("tr-%d-%s", time.Now().UTC().UnixNano(), hex.EncodeToString(suffix[:])), nil
}
