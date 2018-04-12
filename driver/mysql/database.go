package mysql

import (
	"database/sql"
)

const (
	driver = "mysql"
)

// Database records dataSourceName 、db
type Database struct {
	dataSourceName string
	db             *sql.DB
}

// New returns  a newly initialized Database that implements Database
// dataSourceName ("user:password@tcp(ip:port)/database")
func New(dataSourceName string) *Database {
	database := &Database{}

	if dataSourceName != "" {
		database.dataSourceName = dataSourceName
	}

	return database
}

// SetDB sets db
func (database *Database) SetDB(db *sql.DB) {
	database.db = db
}

// GetDB returns db
func (database *Database) GetDB() *sql.DB {
	return database.db
}

// Open returns mysql driver database
func (database *Database) Open() (*Database, error) {

	db, err := sql.Open(driver, database.dataSourceName)
	database.db = db

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	return database, err
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (database *Database) Ping() error {
	return database.db.Ping()
}

// Close closes the database, releasing any open resources.
func (database *Database) Close() error {
	return database.db.Close()
}

// Prepare creates a prepared statement for later queries or executions.
func (database *Database) Prepare(sql string) (*sql.Stmt, error) {
	stmt, err := database.db.Prepare(sql)

	return stmt, err
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (database *Database) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	if args != nil {
		return database.db.Query(sql, args...)
	}
	return database.db.Query(sql)
}

// Fetch executes a query that is expected to return at most one row.
// Fetch is wrapped for db.QueryRow
func (database *Database) Fetch(sql string, args ...interface{}) *sql.Row {
	return database.db.QueryRow(sql, args)
}

// Transaction records Tx
type Transaction struct {
	Tx *sql.Tx
}

// NewTx returns a newly initialized Transaction that implements Transaction
func NewTx() *Transaction {
	return &Transaction{}
}

// Begin starts a transaction. The default isolation level is dependent on
// the driver.
func (database *Database) Begin() (*Transaction, error) {
	transaction, err := database.db.Begin()
	tx := NewTx()
	tx.Tx = transaction
	return tx, err
}

// Rollback aborts the transaction.
func (transaction *Transaction) Rollback() error {
	return transaction.Tx.Rollback()
}

// Commit commits the transaction.
func (transaction *Transaction) Commit() error {
	return transaction.Tx.Commit()
}

// PrepareAndExecute creates a prepared statement for later queries or executions.
func (transaction *Transaction) PrepareAndExecute(queryBuilder *QueryBuilder) (sql.Result, error) {
	stmt, err := transaction.Tx.Prepare(queryBuilder.GetSQL())

	if err != nil {
		panic(err)
	}
	return stmt.Exec(queryBuilder.GetParams()...)
}
