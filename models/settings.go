package models

type Settings struct {
	ID            int64 `json:"id"`
	TotalPosts    int   `json:"totalPosts"`
	TotalViews    int   `json:"totalViews"`
	TotalComments int   `json:"totalComments"`
}
