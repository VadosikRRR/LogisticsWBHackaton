package app

import "time"

type Options struct {
	Port               int
	AggregatorURL      string
	MLBackendURL       string
	TransportURL       string
	UseMockML          bool
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	RedisDialTimeout   time.Duration
	RedisReadTimeout   time.Duration
	RedisWriteTimeout  time.Duration
	DefaultWindowLimit int
	VehicleCapacity    float64
	SafetyBuffer       float64
	MaxVehiclesPerSlot int
}
