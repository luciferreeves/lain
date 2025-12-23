package email

import (
	"slices"
)

func hasAttribute(attrs []string, attr string) bool {
	return slices.Contains(attrs, attr)
}
