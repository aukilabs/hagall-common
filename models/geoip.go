package models

import (
	"net"
	"time"
)

// GeoIPData stores GeoIP data
type GeoIPData struct {
	Network   net.IPNet
	ASN       int
	Latitude  float64
	Longitude float64
}

// UpdateASN is used to update ASN record
type UpdateASN struct {
	Network net.IPNet
	ASN     int
	Version time.Time
}

// UpdateCity is used to update City record
type UpdateCity struct {
	Network   net.IPNet
	Latitude  float64
	Longitude float64
	Version   time.Time
}
