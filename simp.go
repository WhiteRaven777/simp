package simp

import (
	"database/sql"
	"errors"
	"sync"
	"time"
)

// DB is a database object.
type DB struct {
	db  *sql.DB
	tx  *sql.Tx
	mx  sync.Mutex
	err error
}

// New returns a new * DB.
// If error, error is set in * DB.err.
// You can get * DB.err with * DB.Error ().
func New(dns string) *DB {
	if db, err := sql.Open("mysql", dns); err != nil {
		return &DB{err: err}
	} else {
		return &DB{db: db}
	}
}

// Begin starts transaction.
func (db *DB) Begin() error {
	if db.err == nil {
		if tx, err := db.db.Begin(); err != nil {
			db.err = err
		} else {
			db.mx.Lock()
			defer db.mx.Unlock()
			db.tx = tx
		}
	}
	return db.err
}

// Close closes database.
func (db *DB) Close() {
	db.err = db.db.Close()
}

// Commit commits the transacrion.
// Begin needs to be executed before this method is executed.
func (db *DB) Commit() error {
	db.mx.Lock()
	defer db.mx.Unlock()
	switch {
	case db.tx == nil:
		db.err = errors.New("The transaction has not been started or has already been committed or rolled back.")
	case db.err == nil:
		if err := db.tx.Commit(); err != nil {
			db.err = err
		} else {
			db.tx = nil
		}
	}
	return db.err
}

// Rollback rolls back transaction.
// Begin needs to be executed before this method is executed.
func (db *DB) Rollback() error {
	db.mx.Lock()
	defer db.mx.Unlock()
	switch {
	case db.tx == nil:
		db.err = errors.New("The transaction has not been started or has already been committed or rolled back.")
	case db.err == nil:
		if err := db.tx.Rollback(); err != nil {
			db.err = err
		} else {
			db.tx = nil
		}
	}
	return db.err
}

// Exec executes a prepared statement with the specified arguments.
// And this method returns a Result summarizing the effect of the statement.
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db.err != nil {
		return nil, db.err
	}
	switch db.tx {
	case nil:
		return db.db.Exec(query, args...)
	default:
		return db.db.Exec(query, args...)
	}
}

// The query executes a prepared query statement with the specified arguments
// and returns the query result as *sql.Rows.
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db.err != nil {
		return nil, db.err
	}
	switch db.tx {
	case nil:
		return db.db.Query(query, args...)
	default:
		return db.tx.Query(query, args...)
	}
}

// QueryRow executes a prepared query statement with the specified arguments.
// Scans the first selected line and returns it as *sql.Row.
// It will be destroyed after that.
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	if db.err != nil {
		return nil
	}
	switch db.tx {
	case nil:
		return db.db.QueryRow(query, args...)
	default:
		return db.tx.QueryRow(query, args...)
	}
}

// The error returns *DB.err.
func (db *DB) Error() error {
	return db.err
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are reused forever.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	db.db.SetConnMaxLifetime(d)
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
// If n <= 0, no idle connections are retained.
func (db *DB) SetMaxIdleConns(n int) {
	db.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (db *DB) SetMaxOpenConns(n int) {
	db.db.SetMaxOpenConns(n)
}