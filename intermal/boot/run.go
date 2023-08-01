package boot

import (
	"context"
	systemModel "github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/app"
	"github.com/gogoclouds/gogo/pkg/stream"
)

func Run(config string) {
	app.
		New(context.Background(), config).
		Cache().
		Database().AutoMigrate(tables()).
		HTTP(loadRouter).
		Run()
}

func tables() []any {
	tables := stream.Connect(
		systemModel.Tables,
	)
	return tables.Slice()
}