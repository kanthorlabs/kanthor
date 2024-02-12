package idx

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func New(ns string) string {
	return fmt.Sprintf("%s_%s", ns, ksuid.New().String())
}

func Build(ns, id string) string {
	return fmt.Sprintf("%s_%s", ns, id)
}
