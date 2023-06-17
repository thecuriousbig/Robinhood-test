package bloghdl

import (
	"net/http"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/dto"
	"robinhood/pkg/auth"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s ports.BlogService
	c ports.CommentService
}

func New(s ports.BlogService, c ports.CommentService) *Handler {
	return &Handler{s: s, c: c}
}

// @Summary      Create Blog
// @Tags         Blog
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Router       /blog [post]
// @Param title body string true "blog title"
// @Param content body string true "blog content"
// @Response 200 {object} dto.BaseResponseWithData[dto.PopulatedBlog]
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) CreateBlog(c echo.Context) error {
	ctx := c.Request().Context()
	user := c.Get("user").(*jwt.Token)
	claims, ok := user.Claims.(*auth.JWTCustomClaims)
	if !ok {
		return echo.ErrUnauthorized
	}
	userId := claims.UserId

	var req dto.CreateBlogRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// create blog
	blog, err := h.s.CreateBlog(ctx, &domains.CreateBlogRequest{
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: userId,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponseWithData[dto.PopulatedBlog]{
		BaseResponse: dto.BaseResponse{
			Code: 0,
		},
		Data: dto.PopulatedBlog{
			ID:      blog.ID.Hex(),
			Title:   blog.Title,
			Content: blog.Content,
			Author: dto.User{
				ID:           blog.Author.ID.Hex(),
				Username:     blog.Author.Username,
				Email:        blog.Author.Email,
				ProfileImage: blog.Author.ProfileImage,
			},
			Status:    blog.Status,
			CreatedAt: blog.CreatedAt.String(),
		},
	})
}

// @Summary      Get blog by id
// @Tags         Blog
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Router       /blog/{blogId} [post]
// @Param blogId path string true "blog id"
// @Response 200 {object} dto.BaseResponseWithData[dto.PopulatedBlog]
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) GetBlogByID(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.GetBlogByIDRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// get blog
	blog, err := h.s.GetBlogByID(ctx, req.BlogId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponseWithData[dto.PopulatedBlog]{
		BaseResponse: dto.BaseResponse{
			Code: 0,
		},
		Data: dto.PopulatedBlog{
			ID:      blog.ID.Hex(),
			Title:   blog.Title,
			Content: blog.Content,
			Author: dto.User{
				ID:           blog.Author.ID.Hex(),
				Username:     blog.Author.Username,
				Email:        blog.Author.Email,
				ProfileImage: blog.Author.ProfileImage,
			},
			Status:    blog.Status,
			CreatedAt: blog.CreatedAt.String(),
		},
	})
}

// @Summary      List blog
// @Tags         Blog
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Router       /blog [get]
// @Param page query uint32 true "page number"
// @Param limit query uint32 true "limit per page"
// @Response 200 {object} dto.BaseResponseWithData[dto.ListBlogResponse]
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) ListBlog(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.ListBlogRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// list blog
	blogs, err := h.s.ListBlog(ctx, &domains.ListBlogRequest{
		Page:  req.Page,
		Limit: req.Limit,
	})
	if err != nil {
		return err
	}

	data := make([]dto.PopulatedBlog, len(blogs.Data))
	for i, blog := range blogs.Data {
		data[i] = dto.PopulatedBlog{
			ID:      blog.ID.Hex(),
			Title:   blog.Title,
			Content: blog.Content,
			Author: dto.User{
				ID:           blog.Author.ID.Hex(),
				Username:     blog.Author.Username,
				Email:        blog.Author.Email,
				ProfileImage: blog.Author.ProfileImage,
			},
			Status:    blog.Status,
			CreatedAt: blog.CreatedAt.String(),
		}
	}

	return c.JSON(http.StatusOK, dto.BaseResponseWithData[dto.ListBlogResponse]{
		BaseResponse: dto.BaseResponse{
			Code: 0,
		},
		Data: dto.ListBlogResponse{
			Blogs:   data,
			HasNext: blogs.HasNext,
		},
	})
}

// @Summary      Create comment
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Router       /comment/{blogId} [post]
// @Param blogId path string true "blog id"
// @Param content body string true "comment content"
// @Response 200 {object} dto.BaseResponseWithData[dto.PopulatedComment]
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) CreateComment(c echo.Context) error {
	ctx := c.Request().Context()
	user := c.Get("user").(*jwt.Token)
	claims, ok := user.Claims.(*auth.JWTCustomClaims)
	if !ok {
		return echo.ErrUnauthorized
	}
	userId := claims.UserId

	var req dto.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// comment blog
	comment, err := h.c.CreateComment(ctx, &domains.CreateCommentRequest{
		BlogId:   req.BlogId,
		AuthorId: userId,
		Content:  req.Content,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponseWithData[dto.PopulatedComment]{
		BaseResponse: dto.BaseResponse{
			Code: 0,
		},
		Data: dto.PopulatedComment{
			ID:     comment.ID.Hex(),
			BlogId: comment.BlogId.Hex(),
			Author: dto.User{
				ID:           comment.Author.ID.Hex(),
				Username:     comment.Author.Username,
				Email:        comment.Author.Email,
				ProfileImage: comment.Author.ProfileImage,
			},
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.String(),
		},
	})
}

// @Summary      List comment
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Router       /comment/{blogId} [get]
// @Param blogId path string true "blog id"
// @Response 200 {object} dto.BaseResponseWithData[[]dto.PopulatedComment]
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) ListComment(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.ListCommentRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// list comment
	comments, err := h.c.ListComment(ctx, &domains.ListCommentRequest{
		BlogId: req.BlogId,
	})
	if err != nil {
		return err
	}

	data := make([]dto.PopulatedComment, len(comments))
	for i, cm := range comments {
		data[i] = dto.PopulatedComment{
			ID:     cm.ID.Hex(),
			BlogId: cm.BlogId.Hex(),
			Author: dto.User{
				ID:           cm.Author.ID.Hex(),
				Username:     cm.Author.Username,
				Email:        cm.Author.Email,
				ProfileImage: cm.Author.ProfileImage,
			},
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt.String(),
		}
	}

	return c.JSON(http.StatusOK, dto.BaseResponseWithData[[]dto.PopulatedComment]{
		BaseResponse: dto.BaseResponse{
			Code: 0,
		},
		Data: data,
	})
}

// @Summary      Update blog status
// @Tags         Blog
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Router       /blog/{blogId} [put]
// @Param blogId path string true "blog id"
// @Param status body string true "blog status"
// @Response 200 {object} dto.BaseResponse
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) UpdateBlogStatus(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.UpdateBlogRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// update blog status
	err := h.s.UpdateBlogStatus(ctx, &domains.UpdateBlogStatusRequest{
		BlogId: req.BlogId,
		Status: req.Status,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponse{
		Code: 0,
	})
}

// @Summary      Archive blog
// @Tags         Blog
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Router       /blog/{blogId} [delete]
// @Param blogId path string true "blog id"
// @Response 200 {object} dto.BaseResponse
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) ArchiveBlog(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.ArchiveBlogRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// archive blog
	err := h.s.ArchiveBlog(ctx, &domains.ArchiveBlogRequest{
		BlogId: req.BlogId,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponse{
		Code: 0,
	})
}
