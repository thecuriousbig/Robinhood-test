package dto

type Blog struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorId  string `json:"authorId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type PopulatedBlog struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    User   `json:"author"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type CreateBlogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetBlogByIDRequest struct {
	BlogId string `param:"blogId" valid:"required"`
}

type ListBlogRequest struct {
	Page  uint32 `query:"page"`
	Limit uint32 `query:"limit"`
}

type ListBlogResponse struct {
	Blogs   []PopulatedBlog `json:"blogs"`
	HasNext bool            `json:"hasNext"`
}

type UpdateBlogRequest struct {
	BlogId string `param:"blogId" valid:"required"`
	Status string `json:"status" valid:"required"`
}

type ArchiveBlogRequest struct {
	BlogId string `param:"blogId" valid:"required"`
}
