package utils

import (
	"fmt"
	"strconv"
	"time"
	"math"
)

func GetCurrentSwatchTime() string {
	// Biel Mean Time
    // Calculate Biel Mean Time
	now := time.Now()
	bmt := float64(now.Unix()) + time.Hour.Seconds()

    // Calculate the elapsed seconds of the day
    elapsed := math.Mod(bmt, 86400);

    // Get the .beat count for the day
    swatchBeats := elapsed / 86.4;
	swatchTime := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	return swatchTime + "@" + strconv.FormatFloat(swatchBeats, 'f', 1, 64)
}
