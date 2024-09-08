package app

import (
	"go-nuxt-blogs/pkg/errs"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fzzp/gotk"
)

type envelope map[string]interface{}

func (app *application) ShouldBindJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	err := gotk.ReadJSON(w, r, dst)
	if err != nil {
		app.FAIL(w, r, errs.ErrBadRequest.AsException(err))
		return false
	}

	slog.InfoContext(r.Context(), "arg", slog.Any("arg", dst))

	err = gotk.CheckStruct(dst)
	if err != nil {
		app.FAIL(w, r, errs.ErrBadRequest.AsException(err, err.Error()))
		return false
	}

	return true
}

func (app *application) GetPagination(r *http.Request) (pageInt, pageSize int) {
	pageInt, _ = strconv.Atoi(r.URL.Query().Get("pageInt"))
	pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageInt <= 0 {
		pageInt = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return
}

// FAIL 请求失败
func (app *application) FAIL(w http.ResponseWriter, r *http.Request, err *gotk.ApiError) {
	slog.ErrorContext(
		r.Context(),
		r.Method+"-"+r.RequestURI,
		slog.String("err", err.Error()),
	)
	gotk.WriteJSON(w, r, err, nil)
}

// SUCC 请求成功
func (app *application) SUCC(w http.ResponseWriter, r *http.Request, data interface{}) {
	gotk.WriteJSON(w, r, errs.ErrOK, data)
}
