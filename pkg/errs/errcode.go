package errs

import (
	"net/http"

	"github.com/fzzp/gotk"
)

/**
 * 定义系统的所有错误码, 包括公共错误码和具体业务错误码
 */

var (
	// 定义公共错误码
	ErrOK                    = gotk.NewApiError(http.StatusOK, "10000", "请求成功")                     // 200
	ErrRecordExists          = gotk.NewApiError(http.StatusConflict, "10001", "记录已存在")              // 409
	ErrBadRequest            = gotk.NewApiError(http.StatusBadRequest, "10004", "入参错误")             // 400
	ErrUnauthorized          = gotk.NewApiError(http.StatusUnauthorized, "10005", "验证失败")           // 401
	ErrForbidden             = gotk.NewApiError(http.StatusForbidden, "10006", "未经授权")              // 403
	ErrNotFound              = gotk.NewApiError(http.StatusNotFound, "10007", "未找到")                // 404
	ErrMethodNotAllowed      = gotk.NewApiError(http.StatusMethodNotAllowed, "10008", "请求方法不支持")    // 405
	ErrNotAcceptable         = gotk.NewApiError(http.StatusNotAcceptable, "10009", "请求头无效")         // 406
	ErrRequestTimeout        = gotk.NewApiError(http.StatusRequestTimeout, "10010", "请求超时")         // 408
	ErrRequestEntityTooLarge = gotk.NewApiError(http.StatusRequestEntityTooLarge, "10011", "请求体过大") // 413
	ErrUnprocessableEntity   = gotk.NewApiError(http.StatusUnprocessableEntity, "10012", "实体错误")    // 422
	ErrTooManyRequests       = gotk.NewApiError(http.StatusTooManyRequests, "10013", "请求繁忙")        // 429
	ErrServerError           = gotk.NewApiError(http.StatusInternalServerError, "10014", "务器内部错误")  // 500
	ErrBadGateway            = gotk.NewApiError(http.StatusBadGateway, "10015", "网关错误")             // 502
	ErrServiceUnavailable    = gotk.NewApiError(http.StatusServiceUnavailable, "10016", "无法处理请求")   // 503
	ErrGatewayTimeout        = gotk.NewApiError(http.StatusGatewayTimeout, "10017", "服务器未响应")       // 504
)
