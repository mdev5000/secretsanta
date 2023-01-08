package appcontext

import (
	"embed"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/user"
)

type AppContext struct {
	SetupService *setup.Service
	SPAContent   embed.FS
	UserService  *user.Service
}
