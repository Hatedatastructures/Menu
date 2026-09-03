package initialize

import (
	"fmt"
	"os"

	"shack/internal/global"
	"shack/internal/model/gen"
	"shack/internal/model/system"
	"shack/internal/model/systemRbac"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	case "pgsql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Pgsql.Dbname
		return GormPgSql()
	case "oracle":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Oracle.Dbname
		return GormOracle()
	case "mssql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mssql.Dbname
		return GormMssql()
	case "sqlite":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Sqlite.Dbname
		return GormSqlite()
	default:
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	}
}

// 修复 PostgreSQL sequence（防止 Navicat 生成数据导致自增ID错乱）
func FixSequence(db *gorm.DB, table string) error {
	sql := fmt.Sprintf("SELECT setval(pg_get_serial_sequence('%s','id'), COALESCE(MAX(id),1), true) FROM %s;", table, table)
	return db.Exec(sql).Error
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(

		system.JwtBlacklist{},

		systemRbac.Api{},
		systemRbac.ConfigGroup{},
		systemRbac.SysOperationRecord{},
		systemRbac.SysDictionaryDetail{},
		systemRbac.SysDictionary{},
		systemRbac.SysBaseMenu{},
		systemRbac.SysBaseMenuParameter{},
		systemRbac.SysBaseMenuBtn{},
		systemRbac.User{},
		systemRbac.SysUserAuthority{},
		systemRbac.SysAuthority{},
		systemRbac.SysAuthorityMenu{},
		systemRbac.SysAuthorityBtn{},
		systemRbac.SysMenu{},
		systemRbac.SysFile{},
		systemRbac.Notice{},
		systemRbac.NoticeRecord{},
		systemRbac.ChatBlacklist{},
		systemRbac.ChatConversation{},
		systemRbac.ChatMessage{},
		systemRbac.AiConfig{},
		systemRbac.Template{},

		// 邮件系统表
		systemRbac.EmailTemplate{},
		systemRbac.EmailTask{},
		systemRbac.EmailLog{},
		systemRbac.EmailLimitLog{},
		systemRbac.EmailConfig{},
		systemRbac.EmailBlacklist{},
		systemRbac.ApiTestLog{},

		// AI 题目生成模块
		gen.GenApiKey{},
		gen.GenHistory{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	// 如果是 PostgreSQL，修复 sequence
	if global.GVA_CONFIG.System.DbType == "pgsql" {

		tables := []string{}

		for _, table := range tables {
			if err := FixSequence(db, table); err != nil {
				global.GVA_LOG.Warn("fix sequence failed", zap.String("table", table), zap.Error(err))
			}
		}
	}

	err = bizModel()

	if err != nil {
		global.GVA_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
