package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"shack/internal/global"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AutoInitDB 自动初始化数据库: 检测空库则执行SQL文件
func AutoInitDB() {
	if !global.GVA_CONFIG.System.AutoInitDB {
		return
	}
	if global.GVA_DB == nil {
		global.GVA_LOG.Warn("auto-init-db: 数据库连接为空,跳过自动初始化")
		return
	}

	db := global.GVA_DB

	// 检查关键表是否有数据
	if !isDatabaseEmpty(db) {
		global.GVA_LOG.Info("auto-init-db: 数据库已有数据,跳过自动初始化")
		return
	}

	sqlDir := global.GVA_CONFIG.System.InitSQLDir
	if sqlDir == "" {
		sqlDir = "../docs/sql"
	}

	// 支持相对路径: 相对于当前工作目录
	if !filepath.IsAbs(sqlDir) {
		wd, err := os.Getwd()
		if err != nil {
			global.GVA_LOG.Error("auto-init-db: 获取工作目录失败", zap.Error(err))
			return
		}
		sqlDir = filepath.Join(wd, sqlDir)
	}

	if _, err := os.Stat(sqlDir); os.IsNotExist(err) {
		global.GVA_LOG.Warn("auto-init-db: SQL目录不存在,跳过", zap.String("dir", sqlDir))
		return
	}

	global.GVA_LOG.Info("auto-init-db: 检测到空数据库,开始自动初始化...", zap.String("dir", sqlDir))

	if err := executeSQLFiles(db, sqlDir); err != nil {
		global.GVA_LOG.Error("auto-init-db: 执行SQL文件失败", zap.Error(err))
		return
	}

	global.GVA_LOG.Info("auto-init-db: 数据库自动初始化完成!")
}

// isDatabaseEmpty 检查关键表是否为空
func isDatabaseEmpty(db *gorm.DB) bool {
	// 检查几个关键表是否有数据,任一有数据就认为不是空库
	tables := []string{
		"sys_rbac_authorities",
		"sys_rbac_base_menus",
		"sys_rbac_users",
	}

	for _, table := range tables {
		var count int64
		// 先检查表是否存在
		if db.Migrator().HasTable(table) {
			if err := db.Table(table).Count(&count).Error; err != nil {
				// 表存在但查询出错,继续检查其他表
				continue
			}
			if count > 0 {
				return false
			}
		}
	}
	return true
}

// executeSQLFiles 读取并执行目录下所有SQL文件(按文件名排序)
func executeSQLFiles(db *gorm.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取SQL目录失败: %w", err)
	}

	// 收集所有 .sql 文件并排序
	var sqlFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		global.GVA_LOG.Warn("auto-init-db: SQL目录下没有找到.sql文件", zap.String("dir", dir))
		return nil
	}

	for _, name := range sqlFiles {
		path := filepath.Join(dir, name)
		global.GVA_LOG.Info("auto-init-db: 执行SQL文件", zap.String("file", name))

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取SQL文件 %s 失败: %w", name, err)
		}

		sql := string(content)
		if strings.TrimSpace(sql) == "" {
			continue
		}

		// 按分号分割并逐条执行(跳过空语句)
		statements := splitSQL(sql)
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if err := db.Exec(stmt).Error; err != nil {
				// 记录警告但继续执行(部分语句可能因表已存在而失败)
				global.GVA_LOG.Warn("auto-init-db: SQL语句执行失败(继续执行)",
					zap.String("file", name),
					zap.String("sql", truncate(stmt, 200)),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}

// splitSQL 按分号分割SQL语句(忽略注释中的分号)
func splitSQL(sql string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 跳过纯注释行
		if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "#") {
			current.WriteString(line)
			current.WriteString("\n")
			continue
		}

		for i, ch := range line {
			switch ch {
			case '\'':
				if !inDoubleQuote {
					inSingleQuote = !inSingleQuote
				}
			case '"':
				if !inSingleQuote {
					inDoubleQuote = !inDoubleQuote
				}
			case ';':
				if !inSingleQuote && !inDoubleQuote {
					current.WriteString(line[:i+1])
					stmt := strings.TrimSpace(current.String())
					if stmt != "" {
						statements = append(statements, stmt)
					}
					current.Reset()
					// 把这一行剩余部分加入新的current
					if i+1 < len(line) {
						current.WriteString(line[i+1:])
					}
					current.WriteString("\n")
					goto nextLine
				}
			}
		}
		current.WriteString(line)
		current.WriteString("\n")
	nextLine:
	}

	// 最后一个没有分号的语句
	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
