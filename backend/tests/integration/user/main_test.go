package user

import (
	"context"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/testutil/mongopool"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
)

var runDbTests bool
var pool *mongopool.DbPool

func acquireDb(t *testing.T) *mongo.Database {
	if !runDbTests {
		t.SkipNow()
	}
	db, cancel := pool.AcquireDb()
	t.Cleanup(cancel)
	return db
}

func TestMain(m *testing.M) {
	// Do not run any db tests when NODB environment variable is set to 1
	if os.Getenv("NODB") == "1" {
		runDbTests = false
		os.Exit(m.Run())
	}

	runDbTests = true
	pool = mongopool.NewDbPool()
	pool.PurgeDb = func(db *mongopool.DB) error {
		_, err := db.Collection(user.CollectionUsers).DeleteMany(context.Background(), bson.D{})
		return err
	}

	if err := pool.Setup(); err != nil {
		pool.Close(false)
		log.Fatalf("Failed to start pool:\n%s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	pool.Close(true)
	os.Exit(code)
}
