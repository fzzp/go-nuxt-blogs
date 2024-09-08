package db

import (
	"context"
	"go-nuxt-blogs/models"
	"strconv"
	"strings"
)

// PostsRepo 定义posts表操作方法
type PostsRepo interface {
	Create(p *models.Posts) (int64, error)
	Update(p *models.Posts) error
	Get(id int64) (p *models.Posts, err error)
	List(f Filter) (list []*models.Posts, err error)

	CreateTag(t models.Tag) (int64, error)
	UpdateTag(t models.Tag) error
	ListTags(f Filter) (list []models.Tag, err error)
	GetOneTag(id int64) (t models.Tag, err error)

	ListAttributes() (list []*models.Attribute, err error)
}

// 接口检查
var _ PostsRepo = (*postsRepo)(nil)

// postsRepo 实现 PostsRepo 接口
type postsRepo struct {
	DB Queryable
}

func NewPostsRepo(db Queryable) *postsRepo {
	return &postsRepo{DB: db}
}

func (store *postsRepo) Create(p *models.Posts) (int64, error) {
	sql := `
		insert into posts (id, author_id, attr_id, title, content) values($1, $2, $3, $4, $5);
	`
	newID, err := create(store.DB, sql, p.ID, p.AuthorID, p.AttrID, p.Title, p.Content)
	if err != nil {
		return 0, err
	}

	if len(p.Tags) <= 0 {
		return newID, nil
	}

	fields := []string{"post_id", "tag_id"}
	values := []interface{}{}

	for _, t := range p.Tags {
		values = append(values, newID)
		values = append(values, t.ID)
	}

	err = multipleInsert(store.DB, "posts_tags", len(p.Tags), fields, values)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (store *postsRepo) Update(p *models.Posts) error {
	sql := `
		update posts 
		set
			title=$1,
			content=$2,
			attr_id=$3,
			updated_at=datetime('now', 'localtime')
		where 
			id=$4 and deleted_at is null;
	`
	err := update(
		store.DB,
		sql,
		p.Title,
		p.Content,
		p.AttrID,
		p.ID,
	)

	if err != nil {
		return err
	}

	// 删除关系
	delSQL := `delete from posts_tags where post_id = $1`
	_, err = store.DB.Exec(delSQL, p.ID)
	if err != nil {
		return err
	}

	// 更新 tag
	if len(p.Tags) <= 0 {
		return nil
	}

	fields := []string{"post_id", "tag_id"}
	values := []interface{}{}

	for _, t := range p.Tags {
		values = append(values, p.ID)
		values = append(values, t.ID)
	}

	err = multipleInsert(store.DB, "posts_tags", len(p.Tags), fields, values)
	if err != nil {
		return err
	}

	return err
}

func (store *postsRepo) Get(id int64) (p *models.Posts, err error) {
	sql := `
		select 
			p.id, 
			p.author_id, 
			p.title, 
			p.content,
			p.created_at, 
			p.updated_at,
			p.attr_id,
			p.views,
			u.username, 
			u.avatar
		from 
			posts p
		left join 
			users u
		on 
			p.author_id = u.id
		where 
			p.id = $1 and p.deleted_at is null 
		limit 1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	p = new(models.Posts)
	u := new(models.User)
	err = row.Scan(
		&p.ID,
		&p.AuthorID,
		&p.Title,
		&p.Content,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.AttrID,
		&p.Views,
		&u.Username,
		&u.Avatar,
	)
	if err != nil {
		return nil, err
	}
	u.ID = p.AttrID
	p.Author = u

	list, err := store.GetTagsByPostIDs(id)
	if err != nil {
		return p, err
	}

	p.Tags = list

	return p, nil
}

func (store *postsRepo) List(f Filter) (list []*models.Posts, err error) {
	sql := `
	select 
		p.id, 
		p.author_id, 
		p.attr_id,
		p.views,
		p.title, 
		p.content,
		p.created_at, 
		p.updated_at,
		u.id as userId, 
		u.username, 
		u.avatar
	from 
		posts p
	left join 
		users u
	on 
		p.author_id = u.id
	where 
		p.deleted_at is null
	order by p.id desc
	limit $2 offset $3;
	`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, f.limit(), f.offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list = make([]*models.Posts, 0)
	for rows.Next() {
		p := new(models.Posts)
		u := new(models.User)
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.AttrID,
			&p.Views,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			&u.ID,
			&u.Username,
			&u.Avatar,
		)
		if err != nil {
			return nil, err
		}
		p.Author = u
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 查询对应的tag
	pids := make([]int64, len(list))
	for i, post := range list {
		pids[i] = post.ID
	}

	tags, err := store.GetTagsByPostIDs(pids...)
	if err != nil {
		return list, err
	}

	if len(tags) > 0 {
		for i := 0; i < len(list); i++ {
			for j := 0; j < len(tags); j++ {
				if list[i].ID == tags[j].PostID {
					if list[i].Tags == nil {
						list[i].Tags = make([]models.Tag, 0)
					}
					list[i].Tags = append(list[i].Tags, tags[j])
				}
			}
		}
	}

	return list, nil
}

func (store *postsRepo) CreateTag(t models.Tag) (int64, error) {
	sql := `
		insert into tags (id, tag_name) values($1, $2);
	`
	return create(store.DB, sql, t.ID, t.TagName)
}

func (store *postsRepo) UpdateTag(t models.Tag) error {
	sql := `
		update tags 
		set
			tag_name=$1,
			updated_at=datetime('now', 'localtime')
		where 
			id=$2;
	`
	err := update(
		store.DB,
		sql,
		t.TagName,
		t.ID,
	)

	return err
}

func (store *postsRepo) ListTags(f Filter) (list []models.Tag, err error) {
	sql := `
	select 
		t.id, 
		t.tag_name, 
		t.created_at, 
		t.updated_at
	from 
		tags t
	order by 
		t.id desc
	limit $1 offset $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, f.limit(), f.offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list = make([]models.Tag, 0)
	for rows.Next() {
		t := new(models.Tag)
		err = rows.Scan(
			&t.ID,
			&t.TagName,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, *t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (store *postsRepo) GetOneTag(id int64) (t models.Tag, err error) {
	sql := `
		select 
			t.id, 
			t.tag_name, 
			t.created_at, 
			t.updated_at
		from 
			tags t
		where 
			t.id = $1
		limit 1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, sql)
	if err != nil {
		return models.Tag{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(
		&t.ID,
		&t.TagName,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return models.Tag{}, err
	}
	return t, nil
}

func (store *postsRepo) ListAttributes() (list []*models.Attribute, err error) {
	sql := `select id, attr_name from attributes`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list = make([]*models.Attribute, 0)
	for rows.Next() {
		a := new(models.Attribute)
		err = rows.Scan(
			&a.ID,
			&a.AttrName,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (store *postsRepo) GetTagsByPostIDs(pids ...int64) (list []models.Tag, err error) {
	sql := `
	select 
		t.id, 
		t.tag_name, 
		t.created_at, 
		t.updated_at,
		pt.post_id
	from 
		posts_tags pt
	inner join 
		tags t 
	on 
		pt.tag_id = t.id
	where 
		pt.post_id in
	`

	whereSQL := []string{}
	for _, id := range pids {
		whereSQL = append(whereSQL, strconv.Itoa(int(id)))
	}

	sql += " (" + strings.Join(whereSQL, ",") + ") " + " order by t.id desc "

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := store.DB.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list = make([]models.Tag, 0)
	for rows.Next() {
		t := new(models.Tag)
		err = rows.Scan(
			&t.ID,
			&t.TagName,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.PostID,
		)
		if err != nil {
			return nil, err
		}
		if t.ID > 0 {
			list = append(list, *t)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
