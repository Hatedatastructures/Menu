package initialize

import (
	"context"


	sysModel "shack/internal/model/system"
	rbacModel "shack/internal/model/systemRbac"
	"shack/internal/plugin/announcement/model"
	"shack/internal/service/system"

	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const initOrderEnsureTables = system.InitOrderExternal - 1

type ensureTables struct{}

// auto run
func init() {
	system.RegisterInit(initOrderEnsureTables, &ensureTables{})
}

func (e *ensureTables) InitializerName() string {
	return "ensure_tables_created"
}
func (e *ensureTables) InitializeData(ctx context.Context) (next context.Context, err error) {
	return ctx, nil
}

func (e *ensureTables) DataInserted(ctx context.Context) bool {
	return true
}

func (e *ensureTables) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	tables := []interface{}{
		rbacModel.Api{},
		rbacModel.User{},
		rbacModel.SysUserAuthority{},
		rbacModel.SysBaseMenu{},
		rbacModel.SysAuthority{},
		sysModel.JwtBlacklist{},
		rbacModel.SysDictionary{},

		rbacModel.SysOperationRecord{},
		rbacModel.SysDictionaryDetail{},
		rbacModel.SysBaseMenuParameter{},
		rbacModel.SysBaseMenuBtn{},
		rbacModel.SysAuthorityBtn{},
		sysModel.SysExportTemplate{},
		sysModel.Condition{},
		sysModel.JoinTemplate{},
		sysModel.SysParams{},
		sysModel.SysVersion{},
		adapter.CasbinRule{},


		model.Info{},
	}
	for _, t := range tables {
		_ = db.AutoMigrate(&t)
		// 视图 authority_menu 会被当成表来创建，引发冲突错误（更新版本的gorm似乎不会）
		// 由于 AutoMigrate() 基本无需考虑错误，因此显式忽略
	}
	return ctx, nil
}

func (e *ensureTables) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	tables := []interface{}{
		rbacModel.Api{},
		rbacModel.User{},
		rbacModel.SysUserAuthority{},
		rbacModel.SysBaseMenu{},
		rbacModel.SysAuthority{},
		sysModel.JwtBlacklist{},
		rbacModel.SysDictionary{},

		rbacModel.SysOperationRecord{},
		rbacModel.SysDictionaryDetail{},
		rbacModel.SysBaseMenuParameter{},
		rbacModel.SysBaseMenuBtn{},
		rbacModel.SysAuthorityBtn{},
		sysModel.SysExportTemplate{},
		sysModel.Condition{},
		sysModel.JoinTemplate{},

		adapter.CasbinRule{},


		model.Info{},
	}
	yes := true
	for _, t := range tables {
		yes = yes && db.Migrator().HasTable(t)
	}
	return yes
}
