package test

import (
	"testing"

	"trade/db"
)

func TestDb(t *testing.T) {
	db.InitPostgreSql()
}
