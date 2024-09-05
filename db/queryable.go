package db

import (
	"context"
	"database/sql"
)

// Queryable 将 sqlx.DB 和 sqlx.Tx 常用公共方法提取为一个接口，
// 为了在编写事务代码时，sqlx.Tx 可以重用 “使用sqlx.DB操作数据库” 的方法
type Queryable interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}
