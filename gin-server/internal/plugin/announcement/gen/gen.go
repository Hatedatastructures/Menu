package main

//go:generate go run gen.go

import (
	"path/filepath"

	"gorm.io/gen"
	"shack/internal/plugin/announcement/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{OutPath: filepath.Join("..", "..", "..", "announcement", "blender", "model", "dao"), Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface})
	g.ApplyBasic(
		new(model.Info),
	)
	g.Execute()
}
