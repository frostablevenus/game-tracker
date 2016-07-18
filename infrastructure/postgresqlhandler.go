package infrastructure

import (
	"database/sql"
	//"fmt"
	_ "github.com/lib/pq"

	"game-tracker/interfaces"
)

type PostgresqlHandler struct {
	Conn *sql.DB
}

func (handler *PostgresqlHandler) Execute(statement string, args ...interface{}) (sql.Result, error) {
	res, err := handler.Conn.Exec(statement, args...)
	return res, err
}

func (handler *PostgresqlHandler) Query(statement string, args ...interface{}) (interfaces.Row, error) {
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		return PostgresqlRow{}, err
	}
	r := PostgresqlRow{Rows: rows}
	return r, nil
}

type PostgresqlRow struct {
	Rows *sql.Rows
}

func (r PostgresqlRow) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r PostgresqlRow) Next() bool {
	return r.Rows.Next()
}

func (r PostgresqlRow) Close() error {
	return r.Rows.Close()
}

func NewPostgresqlHandler(dbfileAdr string) (*PostgresqlHandler, error) {
	conn, err := sql.Open("postgres", dbfileAdr)
	if err != nil {
		return new(PostgresqlHandler), err
	}
	postgresqlHandler := new(PostgresqlHandler)
	postgresqlHandler.Conn = conn
	return postgresqlHandler, nil
}
