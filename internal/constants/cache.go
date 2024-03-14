package constants

import "fmt"

func CacheKeyAppTier(id string) string {
	return fmt.Sprintf("sdk/app/%s/tier", id)
}
