package blogsvc

import (
	"context"
	"log"
	"robinhood/internal/core/constants"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/errmsg"
)

type blogService struct {
	br ports.BlogRepository
	ur ports.UserRepository
}

func New(br ports.BlogRepository, ur ports.UserRepository) ports.BlogService {
	return &blogService{br: br, ur: ur}
}

func (s *blogService) CreateBlog(ctx context.Context, req *domains.CreateBlogRequest) (*domains.PopulatedBlog, error) {
	return s.br.CreateTx(ctx, req, s.CreateBlogTx)
}

func (s *blogService) CreateBlogTx(ctx context.Context, req *domains.CreateBlogRequest) (*domains.PopulatedBlog, error) {
	// create blog
	blog, err := s.br.Create(ctx, req)
	if err != nil {
		log.Printf("[blogService::CreateBlogTx::Create] error => %+v", err)
		return nil, errmsg.BlogCreateFailed
	}

	// get author information
	author, err := s.ur.GetByID(ctx, blog.AuthorId)
	if err != nil {
		log.Printf("[commentService::GetByID::Create] error => %+v", err)
		return nil, errmsg.BlogCreateFailed
	}

	return &domains.PopulatedBlog{
		ID:      blog.ID,
		Title:   blog.Title,
		Content: blog.Content,
		Author: domains.User{
			ID:           author.ID,
			Username:     author.Username,
			Email:        author.Email,
			ProfileImage: author.ProfileImage,
		},
		Status:     blog.Status,
		IsArchived: blog.IsArchived,
		CreatedAt:  blog.CreatedAt,
	}, nil
}

func (s *blogService) GetBlogByID(ctx context.Context, id string) (*domains.PopulatedBlog, error) {
	blog, err := s.br.GetPopulatedBlogByID(ctx, id)
	if err != nil {
		log.Printf("[blogService::GetBlogByID::GetPopulatedBlogByID] error => %+v", err)
		return nil, errmsg.BlogGetFailed
	}
	return blog, nil
}

func (s *blogService) ListBlog(ctx context.Context, req *domains.ListBlogRequest) (*domains.ListBlogResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}

	blogs, err := s.br.List(ctx, &domains.PaginationOptions{
		Offset: int64((req.Page - 1) * req.Limit),
		Limit:  int64(req.Limit),
	})
	if err != nil {
		log.Printf("[blogService::ListBlog::List] error => %+v", err)
		return nil, errmsg.BlogListFailed
	}

	// count + 1 to check if there is next page
	count, err := s.br.Count(ctx, &domains.PaginationOptions{
		Offset: int64((req.Page - 1) * req.Limit),
		Limit:  int64(req.Limit + 1),
	})
	if err != nil {
		log.Printf("[blogService::ListBlog::Count] error => %+v", err)
		return nil, errmsg.BlogListFailed
	}

	lenBlogs := len(blogs)
	hasNextPage := false
	if count > int64(lenBlogs) {
		hasNextPage = true
	}
	result := &domains.ListBlogResponse{
		HasNext: hasNextPage,
		Data:    blogs,
	}

	return result, nil
}

func (s *blogService) UpdateBlogStatus(ctx context.Context, req *domains.UpdateBlogStatusRequest) error {
	// check status is valid
	switch req.Status {
	case constants.TO_DO:
	case constants.IN_PROGRESS:
	case constants.DONE:
	default:
		return errmsg.BlogInvalidStatus
	}
	return s.br.UpdateStatus(ctx, req)
}

func (s *blogService) ArchiveBlog(ctx context.Context, req *domains.ArchiveBlogRequest) error {
	if err := s.br.Archive(ctx, req); err != nil {
		log.Printf("[blogService::ArchiveBlog] error => %+v", err)
		return errmsg.BlogArchiveFailed
	}
	return nil
}
