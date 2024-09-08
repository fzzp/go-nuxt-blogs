package db

import (
	"context"
	"go-nuxt-blogs/models"
)

// UserRepo 定义users表操作方法
type UserRepo interface {
	Create(u *models.User) (int64, error)
	Update(u *models.User) error
	GetByUnique(data map[string]interface{}) (*models.User, error)
}

// 接口检查
var _ UserRepo = (*userRepo)(nil)

type userRepo struct {
	DB Queryable
}

func NewUserRepo(qb Queryable) *userRepo {
	return &userRepo{
		DB: qb,
	}
}

func (store *userRepo) Create(u *models.User) (int64, error) {
	sql := `insert into users (
		id,
		email,
		password,
		username,
		avatar,
		bio
	) values (
		$1, $2, $3, $4, $5, $6
	);`

	return create(store.DB, sql, u.ID, u.Email, u.Password, u.Username, u.Avatar, u.Bio)
}

func (store *userRepo) Update(u *models.User) error {
	sql := `
		update users 
		set 
			username=$1, avatar=$2, bio=$3, updated_at=datetime('now', 'localtime')
		where 
			id = $3 and deleted_at is null;
	`
	return update(store.DB, sql, u.Username, u.Avatar, u.Bio, u.ID)
}

func (store *userRepo) GetByUnique(data map[string]interface{}) (u *models.User, err error) {
	sql := `
		select 
			id, email, password, username, avatar, role, bio, created_at, updated_at
		from 
			users 
		where 
	`
	fields := []string{"id", "email"}
	sqlStr, params := unqQuerySQL(fields, data)
	sql = sql + sqlStr + " and deleted_at is null"

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	u = new(models.User)
	row := stmt.QueryRow(params...)
	err = row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.Avatar,
		&u.Role,
		&u.Bio,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}
