package caching

import (
	"fmt"
	"time"
)

func App(id string) (string, time.Duration) {
	return fmt.Sprintf("sdk/app/%s", id), time.Hour * 24
}
