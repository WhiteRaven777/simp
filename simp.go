package simp

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type DriverName string

const (
	MySQL      = DriverName("mysql")
	PostgreSQL = DriverName("postgres")

	defaultLocalAddress = "127.0.0.1"
	mysqlDefaultPort    = 3306
	postgresDefaultPort = 5432
)

// String returns the converted string from of DriverName.
func (dn DriverName) String() string {
	return string(dn)
}

type Dsn string
type DsnConf struct {
	UserName string
	Password string
	Protocol string
	Address  string
	Port     int
	DbName   string
	Params   map[string]string
}

// String returns the converted string from of Dsn.
func (dsn Dsn) String() string {
	return string(dsn)
}

// DSN returns Dsn based on DsnConf and DriverName.
func (dc DsnConf) DSN(dn DriverName) (dsn Dsn, err error) {
	switch {
	case len(dc.UserName) == 0:
		return "", errors.New("UserName is Empty")
	case (len(dc.Address) > 0 && dc.Address != "localhost" && dc.Address != defaultLocalAddress) &&
		len(dc.Password) == 0:
		return "", errors.New("Password is Empty")
	case len(dc.DbName) == 0:
		return "", errors.New("DbName is Empty")
	}
	switch dn {
	case MySQL:
		return Dsn(fmt.Sprintf(
			"%s%s@%s(%s:%d)/%s%s",
			dc.UserName,
			func() string {
				if len(dc.Password) > 0 {
					return ":" + dc.Password
				}
				return ""
			}(),
			func() string {
				if len(dc.Protocol) == 0 {
					return "tcp"
				} else {
					return dc.Protocol
				}
			}(),
			func() string {
				if len(dc.Address) == 0 {
					return defaultLocalAddress
				} else {
					return dc.Address
				}
			}(),
			func() int {
				switch {
				case dc.Port == 0:
					return mysqlDefaultPort
				default:
					return dc.Port
				}
			}(),
			dc.DbName,
			func() string {
				var buf string
				if len(dc.Params) > 0 {
					for k, v := range dc.Params {
						buf += "&" + k + "=" + v
					}
					buf = strings.Replace(buf, "&", "?", 1)
				}
				return buf
			}(),
		)), nil
	case PostgreSQL:
		return Dsn(fmt.Sprintf(
			"postgres://%s%s@%s:%d/%s%s",
			dc.UserName,
			func() string {
				if len(dc.Password) > 0 {
					return ":" + dc.Password
				}
				return ""
			}(),
			func() string {
				if len(dc.Address) == 0 {
					return defaultLocalAddress
				} else {
					return dc.Address
				}
			}(),
			func() int {
				switch {
				case dc.Port == 0:
					return postgresDefaultPort
				default:
					return dc.Port
				}
			}(),
			dc.DbName,
			func() string {
				var buf string
				if len(dc.Params) > 0 {
					for k, v := range dc.Params {
						buf += "&" + k + "=" + v
					}
					buf = strings.Replace(buf, "&", "?", 1)
				}
				return buf
			}(),
		)), nil
	default:
		return "", errors.New("undefined driver name was used")
	}
}

// DB is a database object.
type DB struct {
	db  *sql.DB
	tx  *sql.Tx
	mx  sync.Mutex
	err error
}

// New returns a new *DB.
// If error, error is set in *DB.err.
// You can get *DB.err with *DB.Error ().
func New(dn DriverName, dsn Dsn) *DB {
	if db, err := sql.Open(dn.String(), dsn.String()); err != nil {
		return &DB{err: err}
	} else {
		return &DB{db: db}
	}
}

// Ping will check if it can connect to the specified database.
func (db *DB) Ping() error {
	db.mx.Lock()
	defer db.mx.Unlock()
	return db.db.Ping()
}

// Begin starts transaction.
func (db *DB) Begin() error {
	var tx *sql.Tx
	tx, db.err = db.db.Begin()
	if db.err != nil {
		return db.err
	}
	db.mx.Lock()
	defer db.mx.Unlock()
	db.tx = tx
	return nil
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
	if db.tx == nil {
		return errors.New("The transaction has not been started or has already been committed or rolled back.")
	}
	err := db.tx.Commit()
	db.tx = nil
	return err
}

// Rollback rolls back transaction.
// Begin needs to be executed before this method is executed.
func (db *DB) Rollback() error {
	db.mx.Lock()
	defer db.mx.Unlock()
	if db.tx == nil {
		return errors.New("The transaction has not been started or has already been committed or rolled back.")
	}
	err := db.tx.Rollback()
	db.tx = nil
	return err
}

// Exec executes a prepared statement with the specified arguments.
// And this method returns a Result summarizing the effect of the statement.
func (db *DB) Exec(query string, args ...interface{}) (ret sql.Result, err error) {
	db.mx.Lock()
	defer db.mx.Unlock()

	if db.err != nil {
		return ret, db.err
	}

	switch db.tx {
	case nil:
		return db.db.Exec(query, args...)
	default:
		return db.tx.Exec(query, args...)
	}
}

// The query executes a prepared query statement with the specified arguments
// and returns the query result as *sql.Rows.
func (db *DB) Query(query string, args ...interface{}) (ret *sql.Rows, err error) {
	db.mx.Lock()
	defer db.mx.Unlock()

	if db.err != nil {
		return ret, db.err
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
func (db *DB) QueryRow(query string, args ...interface{}) (ret *sql.Row) {
	db.mx.Lock()
	defer db.mx.Unlock()

	if db.err != nil {
		return ret
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

// Prepare creates a prepared statement for later queries or executions.
func (db *DB) Prepare(query string, args ...interface{}) (*sql.Stmt, error) {
	db.mx.Lock()
	defer db.mx.Unlock()
	if db.tx == nil {
		return db.db.Prepare(query)
	} else {
		return db.tx.Prepare(query)
	}
}
