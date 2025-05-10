package state

import (
	"github.com/yashranjan1/gator/internal/config"
	"github.com/yashranjan1/gator/internal/database"
)

type State struct {
	DataBase *database.Queries
	Config   *config.Config
}
