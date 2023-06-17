package commentsvc

import (
	"context"
	"log"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/errmsg"
)

type commentService struct {
	cr ports.CommentRepository
	ur ports.UserRepository
}

func New(cr ports.CommentRepository, ur ports.UserRepository) ports.CommentService {
	return &commentService{cr: cr, ur: ur}
}

func (s *commentService) CreateComment(ctx context.Context, req *domains.CreateCommentRequest) (*domains.PopulatedComment, error) {
	return s.cr.CreateTx(ctx, req, s.CreateCommentTx)
}

func (s *commentService) CreateCommentTx(ctx context.Context, req *domains.CreateCommentRequest) (*domains.PopulatedComment, error) {
	// create comment first
	comment, err := s.cr.Create(ctx, req)
	if err != nil {
		log.Printf("[commentService::CreateCommentTx::Create] error => %+v", err)
		return nil, errmsg.CommentCreateFailed
	}

	// get author information
	author, err := s.ur.GetByID(ctx, comment.AuthorId)
	if err != nil {
		log.Printf("[commentService::CreateCommentTx::GetByID] error => %+v", err)
		return nil, errmsg.CommentCreateFailed
	}

	return &domains.PopulatedComment{
		ID:     comment.ID,
		BlogId: comment.BlogId,
		Author: domains.User{
			ID:           author.ID,
			Username:     author.Username,
			Email:        author.Email,
			ProfileImage: author.ProfileImage,
		},
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (s *commentService) ListComment(ctx context.Context, req *domains.ListCommentRequest) ([]domains.PopulatedComment, error) {
	comments, err := s.cr.List(ctx, req.BlogId)
	if err != nil {
		log.Printf("[commentService::ListComment::List] error => %+v", err)
		return nil, errmsg.CommentListFailed
	}
	return comments, nil
}
