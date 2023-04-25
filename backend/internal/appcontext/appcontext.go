package appcontext

import (
	"embed"
	"github.com/mdev5000/secretsanta/internal/family"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/transactions"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppContext struct {
	Db             *mongo.Database
	SetupService   *setup.Service
	SPAContent     embed.FS
	UserService    *user.Service
	FamilyService  *family.Service
	TransactionMgr *transactions.NoTransactionMgr
}
