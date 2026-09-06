package initialize

import (
	"shack/internal/global"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
	// 业务表

	)
	if err != nil {
		return err
	}
	return nil
}
