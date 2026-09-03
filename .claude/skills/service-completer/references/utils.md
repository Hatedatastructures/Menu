# 工具函数文档

## JWT/认证相关

文件位置: `server/internal/utils/handleRequestFunc.go`

```go
utils.ClearToken(c *gin.Context)
utils.SetToken(c *gin.Context, token string, maxAge int)
utils.GetToken(c *gin.Context) string
utils.GetClaims(c *gin.Context) (*systemReq.CustomClaims, error)
utils.GetUserID(c *gin.Context) uint
utils.GetUserInfo(c *gin.Context) *systemReq.CustomClaims
utils.GetUserName(c *gin.Context) string
utils.GetNickname(c *gin.Context) string
```

## 定时器

文件位置: `server/internal/utils/timer`

```go
task := timer.NewTimerTask()
task.AddTaskByFunc(taskName, spec, func() {}, option...)
task.AddTaskByJob(taskName, spec, job, option...)
task.StartTask(taskName)
task.StopTask(taskName)
task.Remove(taskName, id)
task.Clear(taskName)
task.Close()
```

## 邮件

文件位置: `server/internal/utils/email`

```go
email.NewBuilder().
    To(emails...).
    CC(emails...).
    Subject("subject").
    Template("template").
    Data(data).
    Send()
```

## 雪花算法

文件位置: `server/internal/utils/snowflake_utils.go`
