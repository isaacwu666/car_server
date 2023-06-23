// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"demo/dao/dao/internal"
)

// internalPlayerDao is internal type for wrapping internal DAO implements.
type internalPlayerDao = *internal.PlayerDao

// playerDao is the data access object for table player.
// You can define custom methods on it to extend its functionality as you wish.
type playerDao struct {
	internalPlayerDao
}

var (
	// Player is globally public accessible object for table player operations.
	Player = playerDao{
		internal.NewPlayerDao(),
	}
)

// Fill with you ideas below.
