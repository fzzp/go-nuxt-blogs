package db

import "time"

const defaultTimeout = 5 * time.Second

type Repository struct {
	Users  UserRepo
	Posts  PostsRepo
	Assets AssetsRepo
}

// NewRepository 创建一个Repository仓库，使用 Queryable 接口，同时兼容 sql.DB 和 sql.Tx 接口
func NewRepository(db Queryable) Repository {
	return Repository{
		Users:  NewUserRepo(db),
		Posts:  NewPostsRepo(db),
		Assets: NewAssetsRepo(db),
	}
}
