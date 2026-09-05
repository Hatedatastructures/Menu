package biz_err

// 成功
const SUCCESS = 0x00000000

// 通用错误码
const (
	// 参数错误 (0x00010500 - 0x000105FF)
	PARAM_ERROR   = 0x00010501 // 参数错误
	PARAM_MISSING = 0x00010502 // 必填参数为空
	PARAM_FORMAT  = 0x00010503 // 参数格式不正确
	JSON_PARSE    = 0x00010505 // JSON格式解析失败

	// 认证授权错误 (0x00010100 - 0x000101FF)
	AUTH_ERROR    = 0x00010101 // 认证失败
	TOKEN_INVALID = 0x00010102 // 令牌无效
	TOKEN_EXPIRED = 0x00010103 // 令牌已过期
	UNAUTHORIZED  = 0x00010104 // 未授权访问
	FORBIDDEN     = 0x00010105 // 权限不足

	// 系统错误 (0x00010200 - 0x000102FF)
	SYSTEM_ERROR = 0x00010201 // 系统内部错误
	SERVICE_BUSY = 0x00010202 // 服务繁忙
	TIMEOUT      = 0x00010203 // 服务超时

	// 数据中间件错误 (0x00010300 - 0x000103FF)
	DB_ERROR    = 0x00010301 // 数据库错误
	CACHE_ERROR = 0x00010302 // 缓存错误
)

// 模块业务错误码
const (
	// 用户模块 (0x00010700 - 0x000107FF)
	USER_ALREADY_EXISTS   = 0x00010701 // 用户已存在
	USER_NOT_FOUND        = 0x00010702 // 用户不存在
	INVALID_CREDENTIALS   = 0x00010703 // 账号或密码错误
	REGISTRATION_FAILED   = 0x00010704 // 注册失败
	LOGIN_FAILED          = 0x00010705 // 登录失败
	PASSWORD_RESET_FAILED = 0x00010706 // 密码重置失败
	USER_UPDATE_FAILED    = 0x00010707 // 更新用户信息失败

	// 商品模块 (0x00010900 - 0x000109FF)
	GOODS_NOT_FOUND  = 0x00010901 // 商品不存在
	GOODS_UNDER_SALE = 0x00010902 // 商品已下架
	GOODS_SOLD_OUT   = 0x00010903 // 商品已售罄
	STOCK_NOT_ENOUGH = 0x00010904 // 库存不足

	// 订单模块 (0x00010800 - 0x000108FF)
	ORDER_NOT_FOUND      = 0x00010801 // 订单不存在
	ORDER_CREATE_FAILED  = 0x00010802 // 订单创建失败
	ORDER_STATUS_ERROR   = 0x00010803 // 订单状态异常
	ORDER_PRICE_ERROR    = 0x00010804 // 订单价格异常
	ORDER_SAVE_FAILED    = 0x00010805 // 订单保存失败
	ORDER_ITEM_NOT_FOUND = 0x00010806 // 订单项不存在
	FORBIDDEN_OPERATION  = 0x00010807 // 禁止该操作

	// 地址模块 (新开一个范围，防止与商品模块冲突)
	ADDRESS_NOT_FOUND   = 0x00010B01 // 未找到地址
	ADDRESS_SAVE_FAILED = 0x00010B02 // 地址保存失败

	// 购物车模块 (0x00010A00 - 0x00010AFF)
	SHOPPING_CART_DATA_ERROR = 0x00010A01 // 购物车数据异常

	// 验证码错误 (0x00010600 - 0x000106FF)
	CAPTCHA_ERROR = 0x00010601 // 验证码错误

	// 未知错误
	UNKNOWN_ERROR = 0x00010FFF

	// 邮件模块 (0x00010C00 - 0x00010CFF)
	EMAIL_NOT_FOUND      = 0x00010C01 // 邮件模板不存在
	EMAIL_SEND_FAILED    = 0x00010C02 // 邮件发送失败
	VERIFY_CODE_ERROR    = 0x00010C03 // 验证码错误
	VERIFY_CODE_EXPIRED  = 0x00010C04 // 验证码已过期
	EMAIL_LIMIT_EXCEEDED = 0x00010C05 // 邮件发送频率超限

	// 字典模块 (0x00010D00 - 0x00010DFF)
	DICT_VALUE_DUPLICATE = 0x00010D01 // 字典值已存在
	DICT_NOT_FOUND       = 0x00010D02 // 字典不存在
	DICT_TYPE_DUPLICATE  = 0x00010D03 // 字典类型已存在

	// Menu 菜谱模块 (0x00010E00 - 0x00010EFF)
	RECIPE_NOT_FOUND       = 0x00010E01 // 菜谱不存在
	INGREDIENT_NOT_FOUND   = 0x00010E02 // 食材不存在
	RECIPE_CREATE_FAILED   = 0x00010E03 // 菜谱创建失败
	RECIPE_UPDATE_FAILED   = 0x00010E04 // 菜谱更新失败
	RECIPE_DELETE_FAILED   = 0x00010E05 // 菜谱删除失败
	INGREDIENT_CREATE_FAILED = 0x00010E06 // 食材创建失败
	INGREDIENT_UPDATE_FAILED = 0x00010E07 // 食材更新失败
	INGREDIENT_DELETE_FAILED = 0x00010E08 // 食材删除失败
	PLAN_NOT_FOUND         = 0x00010E09 // 计划不存在
	PLAN_SAVE_FAILED       = 0x00010E0A // 计划保存失败
	SESSION_NOT_FOUND      = 0x00010E0B // 做饭会话不存在
	SESSION_CREATE_FAILED  = 0x00010E0C // 做饭会话创建失败
	SESSION_UPDATE_FAILED  = 0x00010E0D // 做饭会话更新失败
	FEEDBACK_CREATE_FAILED = 0x00010E0E // 反馈创建失败
	REGISTER_FAILED_MENU   = 0x00010E0F // 注册失败
	LOGIN_FAILED_MENU      = 0x00010E10 // 登录失败
	REFRESH_FAILED         = 0x00010E11 // 刷新令牌失败
)

// CodeMsg 映射错误码到对应的错误信息
var CodeMsg = map[int]string{
	SUCCESS: "操作成功",

	// 参数错误
	PARAM_ERROR:   "参数错误",
	PARAM_MISSING: "必填参数为空",
	PARAM_FORMAT:  "参数格式不正确",
	JSON_PARSE:    "JSON格式解析失败",

	// 认证授权
	AUTH_ERROR:    "认证失败",
	TOKEN_INVALID: "令牌无效",
	TOKEN_EXPIRED: "令牌已过期",
	UNAUTHORIZED:  "未授权访问",
	FORBIDDEN:     "权限不足",

	// 系统
	SYSTEM_ERROR: "系统内部错误",
	SERVICE_BUSY: "服务繁忙",
	TIMEOUT:      "服务超时",

	// 中间件
	DB_ERROR:    "数据库错误",
	CACHE_ERROR: "缓存错误",

	// 用户模块
	USER_ALREADY_EXISTS:   "用户已存在",
	USER_NOT_FOUND:        "用户不存在",
	INVALID_CREDENTIALS:   "账号或密码错误",
	REGISTRATION_FAILED:   "注册失败",
	LOGIN_FAILED:          "登录失败",
	PASSWORD_RESET_FAILED: "密码重置失败",
	USER_UPDATE_FAILED:    "更新用户信息失败",

	// 商品模块
	GOODS_NOT_FOUND:  "商品不存在",
	GOODS_UNDER_SALE: "商品已下架",
	GOODS_SOLD_OUT:   "商品已售罄",
	STOCK_NOT_ENOUGH: "库存不足",

	// 订单模块
	ORDER_NOT_FOUND:      "订单不存在",
	ORDER_CREATE_FAILED:  "订单创建失败",
	ORDER_STATUS_ERROR:   "订单状态异常",
	ORDER_PRICE_ERROR:    "订单价格异常",
	ORDER_SAVE_FAILED:    "订单保存失败",
	ORDER_ITEM_NOT_FOUND: "订单项不存在",
	FORBIDDEN_OPERATION:  "禁止该操作",

	// 地址模块
	ADDRESS_NOT_FOUND:   "未找到地址",
	ADDRESS_SAVE_FAILED: "地址保存失败",

	// 购物车模块
	SHOPPING_CART_DATA_ERROR: "购物车数据异常",

	// 验证码
	CAPTCHA_ERROR: "验证码错误",

	// 邮件模块
	EMAIL_NOT_FOUND:      "邮件模板不存在",
	EMAIL_SEND_FAILED:    "邮件发送失败",
	VERIFY_CODE_ERROR:    "验证码错误",
	VERIFY_CODE_EXPIRED:  "验证码已过期",
	EMAIL_LIMIT_EXCEEDED: "邮件发送频率超限",

	// 字典模块
	DICT_VALUE_DUPLICATE: "字典值已存在",
	DICT_NOT_FOUND:       "字典不存在",
	DICT_TYPE_DUPLICATE:  "字典类型已存在",

	// Menu 菜谱模块
	RECIPE_NOT_FOUND:         "菜谱不存在",
	INGREDIENT_NOT_FOUND:     "食材不存在",
	RECIPE_CREATE_FAILED:     "菜谱创建失败",
	RECIPE_UPDATE_FAILED:     "菜谱更新失败",
	RECIPE_DELETE_FAILED:     "菜谱删除失败",
	INGREDIENT_CREATE_FAILED: "食材创建失败",
	INGREDIENT_UPDATE_FAILED: "食材更新失败",
	INGREDIENT_DELETE_FAILED: "食材删除失败",
	PLAN_NOT_FOUND:           "计划不存在",
	PLAN_SAVE_FAILED:         "计划保存失败",
	SESSION_NOT_FOUND:        "做饭会话不存在",
	SESSION_CREATE_FAILED:    "做饭会话创建失败",
	SESSION_UPDATE_FAILED:    "做饭会话更新失败",
	FEEDBACK_CREATE_FAILED:   "反馈创建失败",
	REGISTER_FAILED_MENU:     "注册失败",
	LOGIN_FAILED_MENU:        "登录失败",
	REFRESH_FAILED:           "刷新令牌失败",

	// 未知错误
	UNKNOWN_ERROR: "未知错误",
}

// GetMessage 根据错误码获取对应的错误信息
func GetMessage(code int) string {
	if msg, ok := CodeMsg[code]; ok {
		return msg
	}
	return CodeMsg[UNKNOWN_ERROR]
}
