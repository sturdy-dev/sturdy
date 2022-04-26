package sturdytest

import (
	"os"
)

func PsqlDbSourceForTesting() string {
	var host = "127.0.0.1:5432"
	if overrideHost := os.Getenv("E2E_PSQL_HOST"); overrideHost != "" {
		host = overrideHost
	}
	return "postgres://mash:mash@" + host + "/mash?sslmode=disable"
}
