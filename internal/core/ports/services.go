package ports

import (
	"context"
	"robinhood/internal/core/domains"
)

type BlogService interface {
	CreateBlog(context.Context, *domains.CreateBlogRequest) (*domains.PopulatedBlog, error)
	CreateBlogTx(context.Context, *domains.CreateBlogRequest) (*domains.PopulatedBlog, error)
	GetBlogByID(context.Context, string) (*domains.PopulatedBlog, error)
	ListBlog(context.Context, *domains.ListBlogRequest) (*domains.ListBlogResponse, error)
	UpdateBlogStatus(context.Context, *domains.UpdateBlogStatusRequest) error
	ArchiveBlog(context.Context, *domains.ArchiveBlogRequest) error
}

type CommentService interface {
	CreateComment(context.Context, *domains.CreateCommentRequest) (*domains.PopulatedComment, error)
	CreateCommentTx(context.Context, *domains.CreateCommentRequest) (*domains.PopulatedComment, error)
	ListComment(context.Context, *domains.ListCommentRequest) ([]domains.PopulatedComment, error)
}

type UserService interface {
	Register(context.Context, *domains.RegisterRequest) error
	Login(context.Context, *domains.LoginRequest) (*domains.LoginResponse, error)
	Update(context.Context, *domains.UpdateUserRequest) (*domains.User, error)
}
