package db

import (
	"context"
	"go-nuxt-blogs/models"
)

// AssetsRepo 定义users表操作方法
type AssetsRepo interface {
	SaveFiles([]*models.Assets) error
	GetFile(fileId int64) (*models.Assets, error)
	ListFiles(uid int64, f Filter, ft models.FileType) ([]*models.Assets, error)
}

// 接口检查
var _ AssetsRepo = (*assetsRepo)(nil)

// assetsRepo 实现 AssetsRepo 接口
type assetsRepo struct {
	DB Queryable
}

func NewAssetsRepo(db Queryable) *assetsRepo {
	return &assetsRepo{
		DB: db,
	}
}
func (store *assetsRepo) SaveFiles(assets []*models.Assets) error {
	fields := []string{"id", "user_id", "data", "filetype", "filename", "size", "description"}
	values := []interface{}{}
	for _, v := range assets {
		values = append(values, v.ID)
		values = append(values, v.UserId)
		values = append(values, v.Data)
		values = append(values, v.FileType)
		values = append(values, v.Filename)
		values = append(values, v.Size)
		values = append(values, v.Description)

	}
	err := multipleInsert(store.DB, "assets", len(assets), fields, values)
	if err != nil {
		return err
	}
	return nil
}

func (store *assetsRepo) GetFile(fileId int64) (*models.Assets, error) {
	querySQL := `
	SELECT 
		id, user_id, data, filetype, filename, size, description 
	FROM assets 
	WHERE 
		id = $1 AND deleted_at IS NULL;`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, fileId)
	a := new(models.Assets)
	err = row.Scan(
		&a.ID,
		&a.UserId,
		&a.Data,
		&a.FileType,
		&a.Filename,
		&a.Size,
		&a.Description,
	)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (store *assetsRepo) ListFiles(uid int64, f Filter, ft models.FileType) ([]*models.Assets, error) {
	querySQL := `
	SELECT 
		id, user_id, data, filetype, filename, size, description 
	FROM assets 
	WHERE 
		user_id = $1 AND fileType = $2 AND deleted_at IS NULL
	ORDER BY id DESC
	LIMIT $3 OFFSET $4;
	`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, uid, ft, f.limit(), f.offset())
	if err != nil {
		return nil, err
	}
	assets := make([]*models.Assets, 0)
	for rows.Next() {
		a := new(models.Assets)
		err = rows.Scan(
			&a.ID,
			&a.UserId,
			&a.Data,
			&a.FileType,
			&a.Filename,
			&a.Size,
			&a.Description,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	return assets, nil
}
