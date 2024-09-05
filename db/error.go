package db

import (
	"database/sql"
	"errors"
	"go-nuxt-blogs/pkg/errs"
	"strings"

	"github.com/fzzp/gotk"
	"github.com/mattn/go-sqlite3"
)

// ConvertApiError 检查错误，处理常见错误
func ConvertApiError(err error) *gotk.ApiError {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrNotFound.AsException(err)
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			if strings.Contains(err.Error(), ".id") {
				return errs.ErrRecordExists.AsException(err).AsMessage("id 已存在")
			}
			if strings.Contains(err.Error(), "users.email") {
				return errs.ErrRecordExists.AsException(err).AsMessage("email 已存在")
			}
			return errs.ErrRecordExists.AsException(err)
		}
	}

	return errs.ErrServerError.AsException(err)
}
