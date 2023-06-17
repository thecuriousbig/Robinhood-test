package ports

import (
	"context"
	"robinhood/internal/core/domains"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogRepository interface {
	Create(context.Context, *domains.CreateBlogRequest) (*domains.Blog, error)
	CreateTx(context.Context, *domains.CreateBlogRequest, domains.CreateBlogFn) (*domains.PopulatedBlog, error)
	GetPopulatedBlogByID(context.Context, string) (*domains.PopulatedBlog, error)
	List(context.Context, *domains.PaginationOptions) ([]domains.PopulatedBlog, error)
	Count(context.Context, *domains.PaginationOptions) (int64, error)
	UpdateStatus(context.Context, *domains.UpdateBlogStatusRequest) error
	Archive(context.Context, *domains.ArchiveBlogRequest) error
}

type CommentRepository interface {
	Create(context.Context, *domains.CreateCommentRequest) (*domains.Comment, error)
	CreateTx(context.Context, *domains.CreateCommentRequest, domains.CreateCommentFn) (*domains.PopulatedComment, error)
	List(context.Context, string) ([]domains.PopulatedComment, error)
}

type UserRepository interface {
	GetByID(context.Context, primitive.ObjectID) (*domains.User, error)
	GetByUsername(context.Context, string) (*domains.User, error)
	Create(context.Context, *domains.CreateUserRequest) (*domains.User, error)
	Update(context.Context, *domains.UpdateUserRequest) (*domains.User, error)
}
