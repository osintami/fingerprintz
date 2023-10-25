// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

func DistanceInMiles(lat1, lon1, lat2, lon2 float64, units rune) float64 {
	p1 := orb.Point{lat1, lon1}
	p2 := orb.Point{lat2, lon2}
	dist := geo.DistanceHaversine(p1, p2)
	switch units {
	case 'M':
		return (dist / 1000) * .621371
	case 'K':
		return (dist / 1000)
	default:
		return dist
	}
}

func EncodeGeoHash(latitude, longitude float64) string {
	return geohash.Encode(latitude, longitude)
}

func DecodeGeoHash(hash string) (float64, float64) {
	return geohash.Decode(hash)
}
