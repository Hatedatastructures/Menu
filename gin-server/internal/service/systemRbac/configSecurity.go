package systemRbac

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ConfigSecurityService struct{}

// GetSecurityConfig 获取安全配置-后台使用
func (s *ConfigSecurityService) GetSecurityConfig(
	ctx *gin.Context,
) (rs res.GetSecurityConfigRes, err error) {
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "security").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rs = res.GetSecurityConfigRes{
				EncryptEnabled:       false,
				EncryptScope:         "partial",
				EncryptPublicKey:     "",
				EncryptPrivateKey:    "",
				DisableDevtool:       false,
				TokenName:            "Authorization",
				TokenTimeout:         86400,
				TokenActiveTimeout:   86400,
				TokenIsConcurrent:    true,
				TokenIsShare:         true,
				TokenStyle:           "uuid",
				TokenIsReadHeader:    true,
				TokenIsLog:           false,
				TokenIsPrint:         true,
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询安全配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	rs = res.GetSecurityConfigRes{
		EncryptEnabled:       getBoolFromMap(config, "encryptEnabled", false),
		EncryptScope:          getStringFromMap(config, "encryptScope", "partial"),
		EncryptPublicKey:      getStringFromMap(config, "encryptPublicKey", ""),
		EncryptPrivateKey:    getStringFromMap(config, "encryptPrivateKey", ""),
		DisableDevtool:        getBoolFromMap(config, "disableDevtool", false),
		TokenName:            getStringFromMap(config, "tokenName", "Authorization"),
		TokenTimeout:         getIntFromMap(config, "tokenTimeout", 86400),
		TokenActiveTimeout:   getIntFromMap(config, "tokenActiveTimeout", 86400),
		TokenIsConcurrent:    getBoolFromMap(config, "tokenIsConcurrent", true),
		TokenIsShare:         getBoolFromMap(config, "tokenIsShare", true),
		TokenStyle:           getStringFromMap(config, "tokenStyle", "uuid"),
		TokenIsReadHeader:    getBoolFromMap(config, "tokenIsReadHeader", true),
		TokenIsLog:           getBoolFromMap(config, "tokenIsLog", false),
		TokenIsPrint:         getBoolFromMap(config, "tokenIsPrint", true),
	}
	return rs, nil
}

// SaveSecurityConfig 保存安全配置-后台使用
func (s *ConfigSecurityService) SaveSecurityConfig(
	ctx *gin.Context,
	r req.SaveSecurityConfigReq,
) (rs res.SaveSecurityConfigRes, err error) {
	configMap := map[string]interface{}{
		"encryptEnabled":      r.EncryptEnabled,
		"encryptScope":        r.EncryptScope,
		"encryptPublicKey":    r.EncryptPublicKey,
		"encryptPrivateKey":  r.EncryptPrivateKey,
		"disableDevtool":     r.DisableDevtool,
		"tokenName":          r.TokenName,
		"tokenTimeout":       r.TokenTimeout,
		"tokenActiveTimeout": r.TokenActiveTimeout,
		"tokenIsConcurrent":  r.TokenIsConcurrent,
		"tokenIsShare":       r.TokenIsShare,
		"tokenStyle":         r.TokenStyle,
		"tokenIsReadHeader":  r.TokenIsReadHeader,
		"tokenIsLog":         r.TokenIsLog,
		"tokenIsPrint":       r.TokenIsPrint,
	}
	configJson, _ := json.Marshal(configMap)

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "security").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   "security",
				Name:   "安全配置",
				Config: datatypes.JSON(configJson),
				Status: 1,
			}
			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return rs, biz_err.New(biz_err.DB_ERROR, "创建安全配置失败")
			}
			rs = res.SaveSecurityConfigRes{Version: group.Version}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询安全配置失败")
	}

	group.Config = datatypes.JSON(configJson)
	group.Version++
	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存安全配置失败")
	}

	// 清除Redis缓存，确保下次读取最新配置
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Del(ctx, "system:security:config")
	}

	rs = res.SaveSecurityConfigRes{Version: group.Version}
	return rs, nil
}

// GenerateRSAKeys 生成RSA密钥对-后台使用
func (s *ConfigSecurityService) GenerateRSAKeys(
	ctx *gin.Context,
) (rs res.GenerateRSAKeysRes, err error) {
	// 生成2048位RSA密钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "生成RSA密钥失败")
	}

	// 编码私钥
	privateKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509EncodePrivateKey(privateKey),
	})

	// 编码公钥
	publicKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509EncodePublicKey(&privateKey.PublicKey),
	})

	rs = res.GenerateRSAKeysRes{
		PublicKey:  string(publicKeyBytes),
		PrivateKey: string(privateKeyBytes),
	}
	return rs, nil
}

// x509EncodePrivateKey RSA私钥编码
func x509EncodePrivateKey(key *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(key)
}

// x509EncodePublicKey RSA公钥编码
func x509EncodePublicKey(key *rsa.PublicKey) []byte {
	pkixBytes, _ := x509.MarshalPKIXPublicKey(key)
	return pkixBytes
}

