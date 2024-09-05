package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.CleanPath)
	mux.Use(middleware.Logger)
	mux.Use(app.RecoverPanic)

	// 跨域中间件
	mux.Use(app.EnableCORS())

	v1 := chi.NewRouter()
	v1.Post("/login", app.LoginHandler)
	v1.Get("/posts", app.GetListPostsHandler) // &pageInt=1&pageSize=10
	v1.Get("/posts/{id:[0-9]+}", app.GetPostDetailHandler)
	v1.Get("/tags", app.GetListTagsHandler) // &pageInt=1&pageSize=10
	v1.Get("/tags/{id:[0-9]+}", app.GetTagDetailHandler)
	v1.Get("/attributes", app.GetListAttributeHandler)

	v1.Get("/img/{slug:[0-9]+}", app.ShowImageHandler)
	v1.Get("/file", app.GetFileHandler)

	// 管理端理由
	v1.Route("/auth", func(r chi.Router) {
		r.Use(app.RequiredAuth)

		r.Post("/createTag", app.CreateTagHandler)
		r.Post("/createPost", app.CreatePostHandler)
		r.Post("/updatePost", app.UpdatePostHandler)
		r.Post("/updateTag", app.UpdateTagHandler)
		r.Post("/savefile", app.SaveFileHandler)
		r.Get("/files", app.GetFilesHandler) // files?pageInt=1&pageSize=10&fileType=1 获取图片
	})

	// 附加/挂载到mux上,方便版本维护
	mux.Mount("/api/v1", app.setApiVersion(v1, "v1"))

	return mux
}
