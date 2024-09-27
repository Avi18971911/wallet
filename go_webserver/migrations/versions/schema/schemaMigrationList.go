package schema

import "webserver/migrations/versions"

var SchemaMigrations = []versions.Migration{
	MigrationSchema1,
	MigrationSchema2,
}
