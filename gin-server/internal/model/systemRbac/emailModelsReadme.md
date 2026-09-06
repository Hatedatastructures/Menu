# 邮件系统 Model 说明文档

本目录包含邮件系统的所有Model定义，共9个文件。

## 📁 Model文件列表

### 1. emailTemplate.go - 邮件模板
**核心Model:** `EmailTemplate`

**字段说明:**
- `Code` - 模板编码(唯一标识)
- `Name` - 模板名称
- `Subject` - 邮件主题
- `Content` - HTML内容
- `Type` - 模板类型(system/marketing/notification)
- `Status` - 状态(active/inactive)
- `Variables` - 变量定义JSON

**相关结构:**
- `EmailTemplateVariable` - 变量结构
- `EmailTemplatePreviewRequest/Response` - 预览
- `AiGenerateTemplateRequest/Response` - AI生成

### 2. emailTask.go - 邮件定时任务
**核心Model:** `EmailTask`

**字段说明:**
- `Name` - 任务名称
- `CronExpr` - Cron表达式
- `TemplateCode` - 使用的模板
- `TargetType` - 目标类型
- `TargetConfig` - 目标配置JSON
- `Status` - 任务状态
- `NextRunTime/LastRunTime` - 执行时间
- `TotalRuns/SuccessRuns/FailedRuns` - 执行统计

**相关结构:**
- `EmailTaskTargetConfig` - 目标配置
- `EmailTaskRunResponse` - 执行响应

### 3. emailLog.go - 邮件发送日志
**核心Model:** `EmailLog`

**字段说明:**
- `ToEmail` - 收件人
- `Subject` - 主题
- `TemplateCode` - 模板编码
- `Status` - 发送状态
- `BizType` - 业务类型
- `SentAt/DeliveredAt` - 发送/送达时间
- `ErrorMsg` - 错误信息
- `RetryCount` - 重试次数

**相关结构:**
- `EmailLogListResponse` - 列表响应
- `EmailLogDetailResponse` - 详情响应
- `EmailLogResendResponse` - 重发响应

### 4. emailBlacklist.go - 邮件黑名单
**核心Model:** `EmailBlacklist`

**字段说明:**
- `Email` - 被拉黑邮箱
- `Reason` - 拉黑原因
- `Type` - 拉黑类型(manual/auto)
- `Status` - 状态

**相关结构:**
- `EmailBlacklistListResponse` - 列表响应

### 5. emailLimitLog.go - 限流日志
**核心Model:** `EmailLimitLog`

**字段说明:**
- `Email` - 触犯限流的邮箱
- `IP` - 触犯限流的IP
- `LimitType` - 限流类型(email/ip/global)
- `Action` - 操作
- `Blocked` - 是否被阻止
- `Reason` - 限流原因

**相关结构:**
- `EmailLimitLogListResponse` - 列表响应

### 6. emailConfig.go - 邮件服务商配置
**核心Model:** `EmailConfig`

**字段说明:**
- `Name` - 配置名称
- `Provider` - 服务商(aliyun/tencent/smtp)
- `Config` - 配置JSON
- `IsDefault` - 是否默认
- `Priority` - 优先级
- `Status` - 状态

**相关结构:**
- `EmailConfigDetail` - 配置详情

### 7. emailSend.go - 邮件发送
**请求/响应结构:**
- `SendEmailRequest/Response` - 发送普通邮件
- `SendTemplateEmailRequest/Response` - 发送模板邮件
- `SendBatchEmailRequest/Response` - 批量发送
- `EmailBatchRecipient` - 批量收件人
- `EmailAttachment` - 附件

### 8. verifyCode.go - 验证码
**请求/响应结构:**
- `SendVerifyCodeRequest/Response` - 发送验证码
- `VerifyCodeRequest/Response` - 验证验证码
- `GetVerifyCodeStatusResponse` - 获取状态

### 9. emailStatistics.go - 统计监控
**响应结构:**
- `EmailStatisticsResponse` - 邮件统计
- `EmailStatisticsByDate` - 按日期统计
- `EmailQueueStatusResponse` - 队列状态

## 🎯 命名规范

### 表名规范
所有表名使用 `sys_email_` 前缀：
- `sys_email_templates` - 邮件模板表
- `sys_email_tasks` - 邮件任务表
- `sys_email_logs` - 邮件日志表
- `sys_email_blacklist` - 黑名单表
- `sys_email_limit_logs` - 限流日志表
- `sys_email_configs` - 邮件配置表

### 字段规范
- JSON字段使用小驼峰: `templateCode`, `createdAt`
- 数据库字段使用下划线: `template_code`, `created_at`
- 时间类型使用 `*time.Time` 指针
- 布尔类型使用 `bool`
- JSON字段使用 `string` 存储

## 📋 GORM标签说明

```go
// 基础标签
gorm:"primarykey"           // 主键
gorm:"not null"             // 不为空
gorm:"uniqueIndex"          // 唯一索引
gorm:"index"                // 普通索引
gorm:"default:value"        // 默认值
gorm:"comment:说明"         // 字段注释

// 类型标签
gorm:"type:varchar(50)"     // 字符串
gorm:"type:text"            // 长文本
gorm:"type:longtext"        // 超长文本
gorm:"type:int"             // 整数
gorm:"type:json"            // JSON类型

// 列名标签
gorm:"column:column_name"   // 指定列名
```

## 🔗 关联关系

```go
// EmailTask -> EmailTemplate (多对一)
TemplateCode string -> EmailTemplate.Code

// EmailLog -> EmailTask (多对一)
TaskID *uint -> EmailTask.ID

// EmailLog -> EmailTemplate (多对一)
TemplateCode string -> EmailTemplate.Code

// EmailLog -> EmailConfig (多对一)
Provider string -> EmailConfig.Provider
```

## ✅ 验证规则

### 邮箱验证
```go
binding:"required,email"  // 必填且格式正确
```

### 枚举验证
```go
// 模板类型
binding:"oneof=system marketing notification"

// 邮件类型
binding:"oneof=html text"

// 验证码类型
binding:"oneof=register login reset_password"

// 任务状态
binding:"oneof=active inactive"
```

### 长度验证
```go
// 验证码固定6位
binding:"len=6"

// 字符串长度
binding:"min=1,max=100"
```

## 📊 使用示例

### 创建邮件模板
```go
template := &systemRbac.EmailTemplate{
    Code:    "register_success",
    Name:    "注册成功邮件",
    Subject: "恭喜{{.name}}注册成功",
    Content: "<html>...</html>",
    Type:    "system",
    Status:  "active",
}
db.Create(template)
```

### 查询邮件日志
```go
var logs []systemRbac.EmailLog
db.Where("to_email = ? AND status = ?", email, "success").
   Order("sent_at DESC").
   Limit(20).
   Find(&logs)
```

### 统计发送量
```go
var total int64
db.Model(&systemRbac.EmailLog{}).
   Where("created_at >= ?", startTime).
   Count(&total)
```

## 🔄 数据迁移

创建表的SQL示例：
```sql
CREATE TABLE `sys_email_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL COMMENT '模板编码',
  `name` varchar(100) NOT NULL COMMENT '模板名称',
  `subject` varchar(200) NOT NULL COMMENT '邮件主题',
  `content` longtext NOT NULL COMMENT '邮件内容HTML',
  `type` varchar(20) NOT NULL DEFAULT 'system' COMMENT '模板类型',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '状态',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email_templates_code` (`code`),
  KEY `idx_email_templates_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件模板表';
```

## 📝 注意事项

1. **软删除**: 所有Model都嵌入了 `global.GVA_MODEL`，包含软删除功能
2. **JSON字段**: JSON类型字段在数据库中存储为TEXT，需要手动序列化/反序列化
3. **时间字段**: 使用指针类型 `*time.Time`，可以为NULL
4. **索引**: 为常用查询字段添加了索引，如 `to_email`, `status`, `created_at`
5. **表前缀**: 统一使用 `sys_email_` 前缀，符合项目规范
