package infrastructure

import (
	"database/sql"

	"game-tracker/interfaces"
)

type PostgresqlHandler struct {
	Conn *sql.DB
}

func (handler *PostgresqlHandler) Execute(statement string, args ...interface{}) (sql.Result, error) {
	res, err := handler.Conn.Exec(statement, args)
	return res, err
}

func (handler *PostgresqlHandler) Query(statement string, args ...interface{}) (interfaces.Row, error) {
	rows, err := handler.Conn.Query(statement, args)
	r := new(PostgresqlRow)
	if err != nil {
		return r, err
	}
	r.Rows = rows
	return r, nil
}

type PostgresqlRow struct {
	Rows *sql.Rows
}

func (r PostgresqlRow) Scan(dest ...interface{}) error {
	err := r.Rows.Scan(dest...)
	return err
}

func (r PostgresqlRow) Next() bool {
	return r.Rows.Next()
}

func NewPostgresqlHandler(dbfileName string) (*PostgresqlHandler, error) {
	conn, err := sql.Open("postgres", dbfileName)
	if err != nil {
		return new(PostgresqlHandler), err
	}
	postgresqlHandler := new(PostgresqlHandler)
	postgresqlHandler.Conn = conn
	return postgresqlHandler, nil
}
