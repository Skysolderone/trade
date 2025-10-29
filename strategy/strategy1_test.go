package strategy

import (
	"testing"

	"trade/db"
)

func TestStrategy1(t *testing.T) {
	db.InitPostgreSql()
	Strategy1()
}
