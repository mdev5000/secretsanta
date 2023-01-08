// Package mongopool is a utility package for running database tests.
package mongopool

import (
	"context"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/ory/dockertest/v3"
	"log"
	"time"
)

type Client = mongo.Client
type DB = mongo.Database

const testDbUser = "pguser"
const testDbPassword = "secret"

// DbPool manages test access to a database. Conceptually one or more, but currently only one.
type DbPool struct {
	SetupSchema      func(*DB) error
	PurgeDb          func(*DB) error
	sharedDbInstance *DB
	pool             *dockertest.Pool
	resource         *dockertest.Resource
}

func NewDbPool() *DbPool {
	return &DbPool{}
}

// Setup sets up a PostgreSQL database that is loaded via Docker. This function starts up a container instance of the
// database and ensures the database can be reached.
func (d *DbPool) Setup() error {
	var err error

	// Uses a sensible default on windows (tcp/http) and linux/osx (socket).
	d.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not connect to docker: \n%w", err)
	}

	// Pulls an image, creates a container based on it and runs it.
	d.resource, err = d.pool.Run("mongo", "4.2", []string{
		"MONGO_INITDB_ROOT_USERNAME=" + testDbUser,
		"MONGO_INITDB_ROOT_PASSWORD=" + testDbPassword,
	})
	if err != nil {
		log.Printf("failed to start resource: %s", err.Error())
		return fmt.Errorf("could not start resource: \n%w", err)
	}

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet.
	if err := d.pool.Retry(func() error {
		c, err := mongo.Create(fmt.Sprintf(
			"mongodb://%s:%s@localhost:%s",
			testDbUser,
			testDbPassword,
			d.resource.GetPort("27017/tcp"),
		))
		if err != nil {
			log.Printf("failed to create mongo connection: %s", err.Error())
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := c.Connect(ctx); err != nil {
			log.Printf("failed to connect to mongo: %s", err.Error())
			return err
		}
		if err := c.Ping(ctx, nil); err != nil {
			log.Printf("failed to ping to mongo: %s", err.Error())
			return err
		}
		d.sharedDbInstance = c.Database("ss_test_db")
		return nil
		//var err error
		//d.sharedDbInstance, err = postgres.OpenTest(testDbName, testDbUser, testDbPassword, d.resource.GetPort("5432/tcp"))
		//if err != nil {
		//	return err
		//}
		//return d.sharedDbInstance.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to sharedDbInstance via docker: 'n%w", err)
	}

	if d.SetupSchema != nil {
		if err := d.SetupSchema(d.sharedDbInstance); err != nil {
			return err
		}
	}

	return nil
}

// Close removes any docker resources started up for testing.
func (d *DbPool) Close(errIsFatal bool) {
	if err := d.pool.Purge(d.resource); err != nil {
		if errIsFatal {
			log.Fatalf("Could not purge resource: \n%s", err)
		} else {
			log.Printf("Could not purge resource: \n%s", err)
		}
	}
}

// AcquireDb acquires a database instance. You must call close you are finished with the database.  This functions
// currently does 1 thing, but can potentially do 2 at some point.
//
// The first is ensure the database is in a clean state prior to running a test. This means existing database is purged
// from the database.
//
// The second is thing is guarding access to database resources. Currently the database runner only has a single
// database instance, since there is limited testing required. However, at some point it may be required to run multiple
// database instances to improve test performance (and run the db tests in parallel). This function would then act as a
// pool manager, serving database instances as required to test functions.
//
// Ex.
// db, closeDb := acquireDb()
// defer closeDb()
// // do db stuff...
func (d *DbPool) AcquireDb() (*DB, func()) {
	if d.sharedDbInstance == nil {
		panic("dbpool has not been setup, did you run Setup()?")
	}
	return d.sharedDbInstance, func() {
		if d.PurgeDb != nil {
			if err := d.PurgeDb(d.sharedDbInstance); err != nil {
				panic(err)
			}
		}
	}
}
