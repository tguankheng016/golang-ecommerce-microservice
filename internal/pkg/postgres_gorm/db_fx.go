package postgresgorm

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"gormdbfx",
		gormDbProviders,
		gormDbInvokes,
	)

	gormDbProviders = fx.Options(
		fx.Provide(
			NewGormDB,
		),
	)

	gormDbInvokes = fx.Options(
		fx.Invoke(func(db *gorm.DB) error {
			return RunGooseMigration(db)
		}),
	)
)
