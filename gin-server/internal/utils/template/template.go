package template

import (
	"bytes"
	biz_err "shack/internal/error"
	"text/template"
)

// 渲染模板
func RenderTemplate(tpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", biz_err.New(biz_err.UNKNOWN_ERROR, "模板解析错误")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", biz_err.New(biz_err.UNKNOWN_ERROR, "模板渲染错误")
	}

	return buf.String(), nil
}

// 校验模板
func ValidateTemplate(tpl string) error {
	_, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return biz_err.New(biz_err.UNKNOWN_ERROR, "模板格式错误")
	}
	return nil
}