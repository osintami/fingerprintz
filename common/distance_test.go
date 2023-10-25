// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistance(t *testing.T) {
	lat1 := 29.9055
	lon1 := 96.8766

	lat2 := 20.6534
	lon2 := 105.2253

	miles := DistanceInMiles(lat1, lon1, lat2, lon2, 'M')
	assert.Equal(t, 588.5457424391385, miles)

	kilometers := DistanceInMiles(lat1, lon1, lat2, lon2, 'K')
	assert.Equal(t, 947.1728523525213, kilometers)

	meters := DistanceInMiles(lat1, lon1, lat2, lon2, ' ')
	assert.Equal(t, meters, meters)
}

func TestGeoHash(t *testing.T) {
	lat1 := 29.9055
	lon1 := 96.8766

	hash := EncodeGeoHash(lat1, lon1)
	assert.Equal(t, "wjkfhkpkhmk8", hash)

	lat2, lon2 := DecodeGeoHash(hash)

	// NOTE:  encoding and then decoding won't match exactly, so we need to force rounding
	assert.Equal(t, fmt.Sprintf("%3.6f", lat1), fmt.Sprintf("%3.6f", lat2))
	assert.Equal(t, fmt.Sprintf("%3.6f", lon1), fmt.Sprintf("%3.6f", lon2))
}
