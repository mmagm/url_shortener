package db

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mmagm/url_shortener/config"
	"github.com/mmagm/url_shortener/logger"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var configuration config.Configuration

func init() {
	log = logger.Logger()
	configuration, _ = config.Load()
	log.Info(configuration)
}

func connectionString() string {
	cnf := configuration.Database

	return strings.Join([]string{"postgres://", cnf.User, ":", cnf.Password, "@", cnf.Host, ":", cnf.Port, "/", cnf.Name, "?sslmode=", cnf.SSLMode, "&binary_parameters=yes"}, "")
}

type Store struct {
	db *sql.DB
}

func NewStore() (store *Store, err error) {
	connStr := connectionString()

	log.Info(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	store = &Store{db: db}

	return store, nil
}

func (store *Store) Close() {
	store.db.Close()
}
