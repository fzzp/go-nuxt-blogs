package app

import (
	"go-nuxt-blogs/db"
	"go-nuxt-blogs/models"
	"go-nuxt-blogs/pkg/errs"
	"go-nuxt-blogs/pkg/luckid"
	"net/http"
	"strconv"

	"github.com/fzzp/gotk/token"
	"github.com/go-chi/chi/v5"
)

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=16"`
	}
	if ok := app.ShouldBindJSON(w, r, &input); !ok {
		return
	}
	user, err := app.Repo.Users.GetByUnique(envelope{"email": input.Email})
	if err != nil {
		panic(db.ConvertApiError(err))
	}
	if err := user.Matches(input.Password); err != nil {
		panic(errs.ErrBadRequest.AsException(err, "密码或账号不匹配"))
	}
	idTx := strconv.Itoa(int(user.ID))
	aPayload := token.NewPayload(idTx, app.Conf.TokenExpire.Duration)
	rPayload := token.NewPayload(idTx, app.Conf.RefreshTokenExpire.Duration)
	aToken, err := app.JWT.GenToken(aPayload)
	rToken, err2 := app.JWT.GenToken(rPayload)
	if err != nil {
		panic(errs.ErrServerError.AsException(err))
	}
	if err2 != nil {
		panic(errs.ErrServerError.AsException(err2))
	}
	data := envelope{"accessToken": aToken, "refreshToken": rToken, "userinfo": user}
	app.SUCC(w, r, data)
}

func (app *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string  `json:"title" validate:"required,min=2"`
		Content string  `json:"content" validate:"required,min=2"`
		Tags    []int64 `json:"tags"`
	}

	if ok := app.ShouldBindJSON(w, r, &input); !ok {
		return
	}

	payload, ok := r.Context().Value(tokenPayloadKey).(*token.Payload)
	if !ok {
		panic(errs.ErrUnauthorized)
	}

	tags := make([]models.Tag, 0)
	for i := range input.Tags {
		tag := models.Tag{
			ID: input.Tags[i],
		}
		tags = append(tags, tag)
	}
	userId, err := strconv.Atoi(payload.UserText)
	if err != nil || userId <= 0 {
		panic(errs.ErrUnauthorized.AsException(err).AsMessage("获取用户id失败"))
	}
	p := models.Posts{
		ID:       luckid.Next(),
		AuthorID: int64(userId),
		Title:    input.Title,
		Content:  input.Content,
		Tags:     tags,
	}
	newID, err := app.Repo.Posts.Create(&p)
	if err != nil {
		panic(db.ConvertApiError(err))
	}
	app.SUCC(w, r, newID)
}

func (app *application) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      int64   `json:"id" validate:"required,min=1"`
		Title   string  `json:"title" validate:"required,min=2"`
		Content string  `json:"content" validate:"required,min=2"`
		AttrID  int64   `json:"attrId" validate:"required,min=1"`
		Tags    []int64 `json:"tags"`
	}

	if ok := app.ShouldBindJSON(w, r, &input); !ok {
		return
	}

	payload, ok := r.Context().Value(tokenPayloadKey).(*token.Payload)
	if !ok {
		panic(errs.ErrUnauthorized)
	}

	tags := make([]models.Tag, 0)
	for i := range input.Tags {
		tag := models.Tag{
			ID: input.Tags[i],
		}
		tags = append(tags, tag)
	}
	userId, err := strconv.Atoi(payload.UserText)
	if err != nil || userId <= 0 {
		panic(errs.ErrUnauthorized.AsException(err).AsMessage("获取用户id失败"))
	}
	p := models.Posts{
		ID:       input.ID,
		AuthorID: int64(userId),
		AttrID:   input.AttrID,
		Title:    input.Title,
		Content:  input.Content,
		Tags:     tags,
	}
	err = app.Repo.Posts.Update(&p)
	if err != nil {
		panic(db.ConvertApiError(err))
	}
	app.SUCC(w, r, "SUCCESS")
}

func (app *application) GetListPostsHandler(w http.ResponseWriter, r *http.Request) {
	pageInt, pageSize := app.GetPagination(r)
	f := db.Filter{
		PageInt:  pageInt,
		PageSize: pageSize,
	}
	list, err := app.Repo.Posts.List(f)
	if err != nil {
		panic(db.ConvertApiError(err))
	}
	app.SUCC(w, r, list)
}

func (app *application) GetPostDetailHandler(w http.ResponseWriter, r *http.Request) {
	idTx := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idTx)
	if id <= 0 {
		panic(errs.ErrNotFound.AsMessage("文章id不存在"))
	}
	post, err := app.Repo.Posts.Get(int64(id))
	if err != nil {
		panic(db.ConvertApiError(err))
	}

	if post == nil {
		panic(errs.ErrNotFound.AsMessage("文章不存在"))
	}

	app.SUCC(w, r, post)
}

func (app *application) GetListTagsHandler(w http.ResponseWriter, r *http.Request) {
	pageInt, pageSize := app.GetPagination(r)
	f := db.Filter{
		PageInt:  pageInt,
		PageSize: pageSize,
	}
	list, err := app.Repo.Posts.ListTags(f)
	if err != nil {
		panic(db.ConvertApiError(err))
	}
	app.SUCC(w, r, list)
}

func (app *application) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TagName string `json:"tagName" validate:"required,min=2"`
	}

	if ok := app.ShouldBindJSON(w, r, &input); !ok {
		return
	}

	newID, err := app.Repo.Posts.CreateTag(models.Tag{TagName: input.TagName, ID: luckid.Next()})
	if err != nil {
		panic(db.ConvertApiError(err))
	}

	app.SUCC(w, r, newID)
}

func (app *application) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      int64  `json:"id" validate:"required,min=1"`
		TagName string `json:"title" validate:"required,min=2"`
	}

	if ok := app.ShouldBindJSON(w, r, &input); !ok {
		return
	}

	tag := models.Tag{
		ID:      input.ID,
		TagName: input.TagName,
	}

	err := app.Repo.Posts.UpdateTag(tag)
	if err != nil {
		panic(db.ConvertApiError(err))
	}

	app.SUCC(w, r, "SUCCESS")
}

func (app *application) GetTagDetailHandler(w http.ResponseWriter, r *http.Request) {
	idTx := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idTx)
	if id <= 0 {
		panic(errs.ErrNotFound.AsMessage("id不存在"))
	}
	tag, err := app.Repo.Posts.GetOneTag(int64(id))
	if err != nil {
		panic(db.ConvertApiError(err))
	}

	if tag.ID <= 0 {
		panic(errs.ErrNotFound.AsMessage("tag不存在"))
	}

	app.SUCC(w, r, tag)
}

func (app *application) GetListAttributeHandler(w http.ResponseWriter, r *http.Request) {
	list, err := app.Repo.Posts.ListAttributes()
	if err != nil {
		panic(db.ConvertApiError(err))
	}
	app.SUCC(w, r, list)
}

// SaveFileHandler 保存文件到数据库
func (app *application) SaveFileHandler(w http.ResponseWriter, r *http.Request) {
	assets, err := app.FileUp.GetBlob(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.FAIL(w, r, err)
		return
	}
	payload, ok := r.Context().Value(tokenPayloadKey).(*token.Payload)
	if !ok {
		app.FAIL(w, r, errs.ErrServerError)
		return
	}

	fileTypeString := r.FormValue("fileType")
	fileType := models.StringAsFileType(fileTypeString)
	if fileType <= 0 {
		app.FAIL(w, r, errs.ErrBadRequest.AsMessage("文件类型不存在"))
		return
	}

	if len(assets) == 0 {
		app.FAIL(w, r, errs.ErrBadRequest.AsMessage("空文件无法上传"))
		return
	}

	userID, _ := strconv.Atoi(payload.UserText)
	if userID <= 0 {
		app.FAIL(w, r, errs.ErrUnauthorized.AsMessage("用户id不存在"))
		return
	}

	for i := range assets {
		fileID := luckid.Next()
		assets[i].FileType = fileType
		assets[i].ID = fileID
		assets[i].UserId = int64(userID)
	}

	uerr := app.Repo.Assets.SaveFiles(assets)
	if uerr != nil {
		app.FAIL(w, r, errs.ErrServerError.AsMessage(uerr.Error()))
		return
	}

	// 置空data
	for i := range assets {
		assets[i].Data = []byte{}
	}

	app.SUCC(w, r, assets)
}

// GetFilesHandler 获取文件列表
func (app *application) GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	pageInt, pageSize := app.GetPagination(r)
	f := db.Filter{
		PageInt:  pageInt,
		PageSize: pageSize,
	}
	ftInt, _ := strconv.Atoi(r.URL.Query().Get("fileType"))
	fileType := models.FileType(ftInt)
	if fileType.String() == "" {
		fileType = models.IMAGE
	}

	payload, ok := r.Context().Value(tokenPayloadKey).(*token.Payload)
	if !ok {
		app.FAIL(w, r, errs.ErrUnauthorized)
		return
	}

	userID, _ := strconv.Atoi(payload.UserText)
	if userID <= 0 {
		app.FAIL(w, r, errs.ErrUnauthorized.AsMessage("用户id不存在"))
		return
	}

	list, err := app.Repo.Assets.ListFiles(int64(userID), f, fileType)
	if err != nil {
		app.FAIL(w, r, errs.ErrServerError.AsException(err))
		return
	}
	app.SUCC(w, r, list)
}

// GetFileHandler 获取一个文件
func (app *application) GetFileHandler(w http.ResponseWriter, r *http.Request) {
	slugTx := chi.URLParam(r, "slug")
	slug, _ := strconv.Atoi(slugTx)
	if slug <= 0 {
		app.FAIL(w, r, errs.ErrBadRequest)
		return
	}
	a, err := app.Repo.Assets.GetFile(int64(slug))
	if err != nil {
		app.FAIL(w, r, errs.ErrServerError.AsException(err))
		return
	}
	// if a.FileType != models.IMAGE {
	// 	app.FAIL(w, r, errs.ErrBadRequest.AsMessage("不是图片"))
	// 	return
	// }

	app.SUCC(w, r, a)
}

// ShowImageHandler 显示图片，如：<img src="http://xxx/fileID">
func (app *application) ShowImageHandler(w http.ResponseWriter, r *http.Request) {
	slugTx := chi.URLParam(r, "slug")
	slug, _ := strconv.Atoi(slugTx)
	if slug <= 0 {
		app.FAIL(w, r, errs.ErrBadRequest)
		return
	}
	a, err := app.Repo.Assets.GetFile(int64(slug))
	if err != nil {
		app.FAIL(w, r, errs.ErrServerError.AsException(err))
		return
	}
	if a.FileType != models.IMAGE {
		app.FAIL(w, r, errs.ErrBadRequest.AsMessage("不是图片"))
		return
	}

	// NOTE: 这里不再是使用 app.SUCCESS() 作为公共响应
	// 这里必须直接返回文件流
	w.WriteHeader(http.StatusOK)
	w.Write(a.Data)
}
