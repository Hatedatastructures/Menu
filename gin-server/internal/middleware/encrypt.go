package middleware

import (
	"bytes"
	"encoding/json"
	"shack/internal/service/systemRbac"
	"shack/internal/utils/aes"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ResponseEncrypted 表示响应是否被加密
	ResponseEncrypted = "X-Response-Encrypted"
	// ResponseEncryptAlgorithm 加密算法标识
	ResponseEncryptAlgorithm = "X-Response-Algorithm"
)

// EncryptResponse 响应加密中间件
// 根据配置决定是否对响应进行AES加密
func EncryptResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 继续处理请求
		c.Next()

		// 获取响应状态码，非成功不加密
		status := c.Writer.Status()
		if status != 200 && status != 201 {
			return
		}

		// 检查是否启用加密
		encryptService := &systemRbac.ConfigEncryptService{}
		enabled := encryptService.IsEncryptEnabled(c)
		if !enabled {
			return
		}

		// 检查加密范围
		scope := encryptService.GetEncryptScope(c)

		if scope == "none" {
			return
		}

		// 检查是否需要加密(根据scope)
		should := shouldEncrypt(c, scope)

		if !should {
			return
		}

		// 标记需要加密（让 CaptureResponseBody 负责加密）
		c.Set("shouldEncryptResponse", true)

	}
}

// shouldEncrypt 根据scope判断是否需要加密
func shouldEncrypt(c *gin.Context, scope string) bool {
	path := c.Request.URL.Path

	switch scope {
	case "all", "global":
		// 全部加密 (global = all)
		return true
	case "partial":
		// 部分加密，排除公开接口和健康检查
		publicPaths := []string{
			"/health",
			"/api/v1/systemRbac/login",
			"/api/v1/systemRbac/register",
			"/api/v1/systemRbac/captcha",
			"/api/v1/systemRbac/public",
		}
		for _, p := range publicPaths {
			if strings.HasPrefix(path, p) {
				return false
			}
		}
		return true
	case "sensitive":
		// 只加密敏感接口(需要登录的)
		return path != "/health" && !strings.HasPrefix(path, "/api/v1/systemRbac/login")
	default:
		return false
	}
}

// getAESKey 获取AES密钥
func getAESKey(c *gin.Context, encryptService *systemRbac.ConfigEncryptService) string {
	// 方案2: 统一使用固定密钥（与前端保持一致）
	// 注意：生产环境应该从安全配置中读取专门的AES密钥
	// 这里使用固定的32字节密钥，与前端保持一致
	return "0123456789abcdef0123456789abcdef"
}

// cleanRSAKey 清理RSA密钥字符串，提取base64内容
func cleanRSAKey(key string) string {
	// 移除头尾标记和换行
	key = strings.ReplaceAll(key, "-----BEGIN RSA PRIVATE KEY-----", "")
	key = strings.ReplaceAll(key, "-----END RSA PRIVATE KEY-----", "")
	key = strings.ReplaceAll(key, "-----BEGIN PUBLIC KEY-----", "")
	key = strings.ReplaceAll(key, "-----END PUBLIC KEY-----", "")
	key = strings.ReplaceAll(key, "\n", "")
	key = strings.ReplaceAll(key, "\r", "")
	key = strings.ReplaceAll(key, " ", "")
	return key
}

// CaptureResponseBody 捕获响应体
// 需在 EncryptResponse 之前使用
func CaptureResponseBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用buffer捕获响应
		blw := &bodyLogWriter{
			body:           bytes.Buffer{},
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		c.Next()

		// 保存响应体供后续中间件使用
		body := blw.body.String()
		c.Set("responseBody", body)

		// 检查是否需要加密（EncryptResponse 设置了标志）
		shouldEncryptResponse := c.GetBool("shouldEncryptResponse")

		// 获取原始的 ResponseWriter
		originalWriter := blw.ResponseWriter

		if shouldEncryptResponse && body != "" {
			// 需要加密

			// 获取AES密钥
			encryptService := &systemRbac.ConfigEncryptService{}
			aesKey := getAESKey(c, encryptService)

			// 验证密钥
			if err := aes.ValidateKey([]byte(aesKey)); err != nil {

				// 写入原始响应
				originalWriter.WriteHeader(200)
				originalWriter.WriteString(body)
				return
			}

			// 加密
			encrypted, err := aes.Encrypt([]byte(body), []byte(aesKey))
			if err != nil {

				// 写入原始响应
				originalWriter.WriteHeader(200)
				originalWriter.WriteString(body)
				return
			}

			// 设置响应头
			originalWriter.Header().Set("Content-Type", "application/json")
			originalWriter.Header().Set("X-Response-Encrypted", "true")
			originalWriter.Header().Set("X-Response-Algorithm", "AES-256-CBC")

			// 写入加密响应
			originalWriter.WriteHeader(200)
			originalWriter.WriteString(`"` + encrypted + `"`)
			return
		}

		// 未加密，需要手动写入原始响应
		if body != "" {
			// 获取状态码（优先使用bodyLogWriter保存的状态码）
			var statusCode int
			if blw.wroteHeader {
				statusCode = blw.statusCode
			} else {
				statusCode = c.Writer.Status()
				if statusCode == 0 {
					statusCode = 200
				}
			}

			// 写入响应
			originalWriter.WriteHeader(statusCode)
			originalWriter.WriteString(body)
		}
	}
}

// bodyLogWriter 自定义响应Writer用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body        bytes.Buffer
	statusCode  int
	wroteHeader bool
}

// WriteHeader 捕获状态码，但不立即发送
func (w *bodyLogWriter) WriteHeader(statusCode int) {
	// 只保存状态码，不调用父类方法
	w.statusCode = statusCode
	w.wroteHeader = true
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	// 只写入buffer，不立即发送响应
	return w.body.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	// 只写入buffer，不立即发送响应
	return w.body.WriteString(s)
}

// Status 获取状态码
func (w *bodyLogWriter) Status() int {
	if w.wroteHeader {
		return w.statusCode
	}
	return w.ResponseWriter.Status()
}

// DecryptRequest 请求解密中间件
// 解密客户端传递的加密请求体
func DecryptRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否有加密的请求体
		encryptedBody := c.GetHeader("X-Encrypted-Data")
		if encryptedBody == "" {
			c.Next()
			return
		}

		// 检查是否启用加密
		encryptService := &systemRbac.ConfigEncryptService{}
		if !encryptService.IsEncryptEnabled(c) {
			c.Next()
			return
		}

		// 获取AES密钥
		aesKey := getAESKey(c, encryptService)
		if aesKey == "" {
			c.JSON(400, gin.H{"error": "AES key not available"})
			c.Abort()
			return
		}

		// 解密请求体
		decrypted, err := aes.Decrypt(encryptedBody, []byte(aesKey))
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to decrypt request body"})
			c.Abort()
			return
		}

		// 解析JSON并设置到请求体
		var jsonBody map[string]interface{}
		if err := json.Unmarshal(decrypted, &jsonBody); err != nil {
			c.JSON(400, gin.H{"error": "Invalid encrypted data format"})
			c.Abort()
			return
		}

		// 将解密后的数据存储到context，供后续处理
		c.Set("decryptedBody", jsonBody)
		c.Next()
	}
}
