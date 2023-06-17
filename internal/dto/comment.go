package dto

type Comment struct {
	ID        string `json:"id"`
	BlogId    string `json:"blogId"`
	AuthorId  string `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type PopulatedComment struct {
	ID        string `json:"id"`
	BlogId    string `json:"blogId"`
	Author    User   `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type CreateCommentRequest struct {
	BlogId  string `param:"blogId" valid:"required"`
	Content string `json:"content"`
}

type ListCommentRequest struct {
	BlogId string `param:"blogId" valid:"required"`
}
