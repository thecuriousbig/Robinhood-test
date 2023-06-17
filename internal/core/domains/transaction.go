package domains

import "context"

type CreateBlogFn func(context.Context, *CreateBlogRequest) (*PopulatedBlog, error)
type CreateCommentFn func(context.Context, *CreateCommentRequest) (*PopulatedComment, error)
