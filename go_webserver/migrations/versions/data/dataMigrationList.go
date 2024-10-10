package data

import (
	"time"
	"webserver/migrations/versions"
)

const timeout = time.Minute * 1

var Migrations = []versions.Migration{
	MigrationData1,
	MigrationData2,
}
