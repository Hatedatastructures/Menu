---
name: service-completer
description: 完成后端service层代码生成。用于生成Go service函数，自动检查待完成文件、获取model结构、管理错误码。触发条件：用户要求"完成service"、"实现service函数"、"写service代码"或类似请求时使用。
---

# Service Completer

## 工作流程



### 1. 获取相关Model
```bash
node .claude/skills/service-completer/scripts/model-check.js
```

### 2. 检查可用错误码
```bash
node .claude/skills/service-completer/scripts/err-code-check.js
```

### 3. 阅读文档
- [utils.md](references/utils.md) - 工具函数文档
- [error_code.md](references/error_code.md) - 错误码规范
- [example.md](references/example.md) - service示例


### 5. 生成Service代码

错误返回格式:
```go
return biz_err.New(biz_err.ERROR_CODE)
return biz_err.New(biz_err.ERROR_CODE, "自定义信息")
```


### 6. 额外注意
#### 1.在往数据库里面写内容的时候,一定要先找到model,通过绑定model,再将数据写入数据库!
```

正确:
func (s *SysOperationRecordService) DeleteSysOperationRecordByIds(c *gin.Context, r req.DeleteSysOperationRecordByIdsReq) (res response.DeleteSysOperationRecordByIdsRes, err error) {
	err = global.GVA_DB.Delete(&[]system.SysOperationRecord{}, "id in (?)", r.Ids).Error
	res = response.DeleteSysOperationRecordByIdsRes{
		Success: err == nil,
	}
	return res, nil
}

func (s *SysOperationRecordService) DeleteSysOperationRecord(c *gin.Context, r req.DeleteSysOperationRecordReq) (res response.DeleteSysOperationRecordRes, err error) {
	var record system.SysOperationRecord
	record.ID = r.ID
	err = global.GVA_DB.Delete(&record).Error
	res = response.DeleteSysOperationRecordRes{
		Success: err == nil,
	}
	return res, nil
}

错误:
	var article struct {
		Id         int64  `json:"id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		Status     int    `json:"status"`
		CategoryId int64  `json:"category_id"`
		UserId     uint   `json:"user_id"`
		CreatedAt  any    `json:"created_at"`
		UpdatedAt  any    `json:"updated_at"`
	}

	err = global.GVA_DB.Table("gra_articles").Where("id = ?", r.Id).First(&article).Error
	if err != nil {
		global.GVA_LOG.Error("获取文章失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var user system.SysUser
	global.GVA_DB.Table("gra_users").Select("id, username").Where("id = ?", article.UserId).First(&user)
```
#### 2.gorm,redis等等实例请从global中引入!



### 7. 创建新错误码(如需)
```bash
node .claude/skills/service-completer/scripts/err-code-create.js <CODE_NAME> <HEX_VALUE> <MESSAGE>
```

## 路径配置
- Service: `server/internal/service/`
- Model: `server/internal/model/`
- Error Code: `server/internal/error/error_code.go`

## 注意事项
- 忽略 `system/` 目录
- 检查mod名称防止引入不存在的模块
- 文件小于500字节视为待完成
