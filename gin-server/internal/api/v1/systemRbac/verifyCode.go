package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type VerifyCodeApi struct{}

// SendVerifyCodeHandler
// @Tags systemRbacverifyCodeApi
// @Summary SendVerifyCodeHandler 发送验证码-前台使用
// @Description SendVerifyCodeHandler 发送验证码-前台使用
// @Param data body req.SendVerifyCodeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SendVerifyCodeRes}
// @Router /email/send-code [POST]
func (s *VerifyCodeApi) SendVerifyCodeHandler(c *gin.Context) {
	var req req.SendVerifyCodeReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := verifyCodeService.SendVerifyCode(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// VerifyCodeHandler
// @Tags systemRbacverifyCodeApi
// @Summary VerifyCodeHandler 验证验证码-前台使用
// @Description VerifyCodeHandler 验证验证码-前台使用
// @Param data body req.VerifyCodeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.VerifyCodeRes}
// @Router /email/verify-code [POST]
func (s *VerifyCodeApi) VerifyCodeHandler(c *gin.Context) {
	var req req.VerifyCodeReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := verifyCodeService.VerifyCode(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetVerifyCodeStatusHandler
// @Tags systemRbacverifyCodeApi
// @Summary GetVerifyCodeStatusHandler 获取验证码状态-前台使用
// @Description GetVerifyCodeStatusHandler 获取验证码状态-前台使用
// @Param data body req.GetVerifyCodeStatusReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetVerifyCodeStatusRes}
// @Router /email/code/status [GET]
func (s *VerifyCodeApi) GetVerifyCodeStatusHandler(c *gin.Context) {
	var req req.GetVerifyCodeStatusReq
	// query 参数
	{
		val := c.Query("email")
		if val != "" {
			req.Email = val

		}
	}
	{
		val := c.Query("type")
		if val != "" {
			req.Type = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := verifyCodeService.GetVerifyCodeStatus(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
