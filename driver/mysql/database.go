package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	driver = "mysql"
)

type Database struct {
	dataSourceName string
	db             *sql.DB
	transaction    *sql.Tx
}

// dataSourceName ("user:password@tcp(ip:port)/database")
func New(dataSourceName string) *Database {
	database := &Database{}

	if dataSourceName != "" {
		database.dataSourceName = dataSourceName
	}

	return database
}

func (database *Database) SetDB(db *sql.DB) {
	database.db = db
}

func (database *Database) GetDB() *sql.DB {
	return database.db
}

func (database *Database) Open() (*Database, error) {

	db, err := sql.Open(driver, database.dataSourceName)
	database.db = db

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	return database, err
}

func (database *Database) Ping() error {
	return database.db.Ping()
}

func (database *Database) Close() error {
	return database.db.Close()
}

func (database *Database) Prepare(sql string) (*sql.Stmt, error) {
	stmt, err := database.db.Prepare(sql)

	return stmt, err
}

func (database *Database) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	if args != nil {
		return database.db.Query(sql, args...)
	}
	return database.db.Query(sql)
}

func (database *Database) Fetch(sql string, args ...interface{}) *sql.Row {
	return database.db.QueryRow(sql, args)
}

func (database *Database) Begin() (*Database, error) {
	transaction, err := database.db.Begin()

	database.transaction = transaction

	return database, err
}

func (database *Database) Rollback() {
	err := database.transaction.Rollback()
	if err != sql.ErrTxDone && err != nil {
		panic(err)
	}
}

func (database *Database) Commit() error {
	return database.transaction.Commit()
}
