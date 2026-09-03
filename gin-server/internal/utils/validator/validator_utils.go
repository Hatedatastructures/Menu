// Package validator provides parameter validation utilities
// Author: Done-0
// Created: 2025-08-24
package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidErrRes validation error result struct
type ValidErrRes struct {
	Field string // Error field name
	Tag   string // Error tag
	Value any    // Error value
}

// ValidError 校验错误（实现 error 接口）
type ValidError struct {
	Errors []ValidErrRes
}

// Error 实现 error 接口（把多个错误拼接起来）
func (e *ValidError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, err := range e.Errors {
		sb.WriteString(fmt.Sprintf(
			"%s 校验失败(%s); ",
			err.Field,
			err.Tag,
		))
	}
	return sb.String()
}

// NewValidator global validator instance
var NewValidator = validator.New()

// Validate parameter validator
// 校验通过返回 nil，失败返回 error
func Validate(data any) error {
	err := NewValidator.Struct(data)
	if err == nil {
		return nil
	}

	var errs []ValidErrRes
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, ValidErrRes{
			Field: e.Field(),
			Tag:   e.Tag(),
			Value: e.Value(),
		})
	}

	return &ValidError{
		Errors: errs,
	}
}
