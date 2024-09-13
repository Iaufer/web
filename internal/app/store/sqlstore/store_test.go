package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "host=localhost port=5432 user=postgres dbname=postgres_test sslmode=disable password=2505"
	}

	os.Exit(m.Run())
}
