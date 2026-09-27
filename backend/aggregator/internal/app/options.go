package app

import "time"

type Options struct {
	Port int

	PostgresDSN             string
	PostgresMaxConns        int32
	PostgresMinConns        int32
	PostgresConnectTimeout  time.Duration
	PostgresHealthcheckFreq time.Duration
}
