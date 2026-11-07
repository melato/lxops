package util

import (
	"fmt"
	"math"
)

// FormatZfsSize formats a size given in GB as a ZFS size string.
// ZFS size suffixes are powers of 1024, so GB means GiB and MB means MiB.
// The result uses GB when the size is a whole number of GB, otherwise MB, rounded to the nearest MB.
func FormatZfsSize(g float32) (string, error) {
	if g < 0 {
		return "", fmt.Errorf("negative size: %f", g)
	}
	mb := int64(math.Round(float64(g) * 1024))
	if mb < 1 {
		return "", fmt.Errorf("size must be at least 1 MB: %f GB", g)
	}
	if mb%1024 == 0 {
		return fmt.Sprintf("%dGB", mb/1024), nil
	}
	return fmt.Sprintf("%dMB", mb), nil
}
