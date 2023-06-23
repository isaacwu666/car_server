// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlayerDao is the data access object for table player.
type PlayerDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns PlayerColumns // columns contains all the column names of Table for convenient usage.
}

// PlayerColumns defines and stores column names for table player.
type PlayerColumns struct {
	Id       string // id
	Phone    string // 电话
	Pwd      string // 密码
	Salt     string // pwd盐
	NickName string // 昵称
}

// playerColumns holds the columns for table player.
var playerColumns = PlayerColumns{
	Id:       "id",
	Phone:    "phone",
	Pwd:      "pwd",
	Salt:     "salt",
	NickName: "nick_name",
}

// NewPlayerDao creates and returns a new DAO object for table data access.
func NewPlayerDao() *PlayerDao {
	return &PlayerDao{
		group:   "default",
		table:   "player",
		columns: playerColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PlayerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PlayerDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PlayerDao) Columns() PlayerColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PlayerDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PlayerDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PlayerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
