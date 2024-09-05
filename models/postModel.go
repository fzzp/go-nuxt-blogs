package models

type Posts struct {
	ID        int64   `json:"id"`
	AuthorID  int64   `json:"authorId"`
	AttrID    int64   `json:"attrId"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt *string `json:"-"`

	Author *User `json:"author"`
	Tags   []Tag `json:"tags"`
}

type Tag struct {
	ID        int64  `json:"id"`
	TagName   string `json:"tagName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	PostID int64 `json:"-"` // 方便sql查询，非数据库字段
}

type Attribute struct {
	ID       int64  `json:"id"`
	AttrName string `json:"attrName"`
}
