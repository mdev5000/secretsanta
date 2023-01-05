package appcontext

import (
	"embed"
	"github.com/mdev5000/secretsanta/internal/setup"
)

type AppContext struct {
	SetupService *setup.Service
	SPAContent   embed.FS
}
