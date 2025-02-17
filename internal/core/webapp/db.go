package webapp

import (
	"errors"

	"github.com/euiko/go-fullstack-boilerplate/internal/core/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultDbName = "default"

var (
	// for holds all the database connections
	dbInstances  map[string]*gorm.DB = make(map[string]*gorm.DB)
	errNotOpened                     = errors.New("database not opened")
)

// DB returns the selected database connection by its name, default to default connection
func DB(names ...string) *gorm.DB {
	name := defaultDbName
	if len(names) > 0 {
		name = names[0]
	}

	// ensure the database is opened
	if _, ok := dbInstances[name]; !ok {
		log.Error("database not opened",
			log.WithField("name", name),
			log.WithError(errNotOpened),
		)
		// exit early
		panic(errNotOpened)
	}

	return dbInstances[name]
}

func OpenDB(settings DatabaseSettings, names ...string) error {
	var (
		config gorm.Config
		name   = defaultDbName
	)

	// use supplied name if any
	if len(names) > 0 {
		name = names[0]
	}

	// ensure the database is not already opened
	if _, ok := dbInstances[name]; ok {
		return errors.New("database already opened")
	}

	// TODO: add gorm configurations
	gormDb, err := gorm.Open(postgres.Open(settings.Uri), &config)
	if err != nil {
		return err
	}

	// configure connection pool
	sqlDb, err := gormDb.DB()
	if err != nil {
		return err
	}
	sqlDb.SetMaxIdleConns(settings.MaxIdleConns)
	sqlDb.SetMaxOpenConns(settings.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(settings.ConnMaxLifetime)

	dbInstances[name] = gormDb
	return nil
}

func initializeDB(settings DatabaseSettings) error {
	// TODO: support multiple databases
	// TODO: support database other than postgres
	return OpenDB(settings)
}

func closeDB() {
	for name, db := range dbInstances {
		sqlDb, err := db.DB()
		if err != nil {
			// skip closing if there is an error
			continue
		}

		// close the underlying sql connection pool
		if err := sqlDb.Close(); err != nil {
			log.Error("failed to close the database connection",
				log.WithField("name", name),
				log.WithError(err),
			)
		}
	}
}
