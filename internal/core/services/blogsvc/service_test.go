package blogsvc_test

import (
	"context"
	"errors"
	"robinhood/internal/core/constants"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/core/ports/mocks"
	"robinhood/internal/core/services/blogsvc"
	"robinhood/internal/errmsg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testModule struct {
	br  *mocks.BlogRepository
	ur  *mocks.UserRepository
	svc ports.BlogService
}

type test struct {
	name     string
	args     []interface{}
	mockFn   func(*testModule)
	assertFn func()
}

var (
	ctx  = context.TODO()
	date = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
)

func new(t *testing.T) *testModule {
	br := mocks.NewBlogRepository(t)
	ur := mocks.NewUserRepository(t)
	return &testModule{
		br:  br,
		ur:  ur,
		svc: blogsvc.New(br, ur),
	}
}

func TestCreateBlog(t *testing.T) {
	var err error
	mockReq := &domains.CreateBlogRequest{
		Title:    "title",
		Content:  "content",
		AuthorId: "author_id",
	}

	tests := []test{
		{
			name: "success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.br.On("CreateTx", ctx, mockReq, mock.AnythingOfType("domains.CreateBlogFn")).Return(&domains.PopulatedBlog{}, nil)
			},
			assertFn: func() {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			_, err = tm.svc.CreateBlog(tt.args[0].(context.Context), tt.args[1].(*domains.CreateBlogRequest))
			tt.assertFn()
		})
	}
}

func TestCreateBlogTx(t *testing.T) {
	var result *domains.PopulatedBlog
	var err error
	mockReq := &domains.CreateBlogRequest{
		Title:    "title",
		Content:  "content",
		AuthorId: "author_id",
	}

	tests := []test{
		{
			name: "should return error when create blog failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.br.On("Create", ctx, mockReq).Return(nil, errors.New("error"))
			},
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, errmsg.BlogCreateFailed.Error())
			},
		},
		{
			name: "should return error when get author information failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				oid, _ := primitive.ObjectIDFromHex("testid")
				createdBlog := &domains.Blog{
					AuthorId: oid,
				}
				tm.br.On("Create", ctx, mockReq).Return(createdBlog, nil)
				tm.ur.On("GetByID", ctx, createdBlog.AuthorId).Return(nil, errors.New("error"))
			},
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, errmsg.BlogCreateFailed.Error())
			},
		},
		{
			name: "should return populated blog when success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				oid, _ := primitive.ObjectIDFromHex("testid")
				createdBlog := &domains.Blog{
					ID:         oid,
					Title:      "title",
					Content:    "content",
					AuthorId:   oid,
					Status:     constants.TO_DO,
					IsArchived: false,
					CreatedAt:  date,
				}
				author := &domains.User{
					ID:           oid,
					Username:     "username",
					Email:        "email",
					ProfileImage: "profile_image",
				}
				tm.br.On("Create", ctx, mockReq).Return(createdBlog, nil)
				tm.ur.On("GetByID", ctx, createdBlog.AuthorId).Return(author, nil)
			},
			assertFn: func() {
				assert.NoError(t, err)
				assert.Nil(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			result, err = tm.svc.CreateBlogTx(tt.args[0].(context.Context), tt.args[1].(*domains.CreateBlogRequest))
			tt.assertFn()
		})
	}
}

func TestGetBlogByID(t *testing.T) {
	var result *domains.PopulatedBlog
	var err error
	mockReq := "blog_id"

	tests := []test{
		{
			name: "should return error when get blog failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.br.On("GetPopulatedBlogByID", ctx, mockReq).Return(nil, errors.New("error"))
			},
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, errmsg.BlogGetFailed.Error())
			},
		},
		{
			name: "should return populated blog when success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				oid, _ := primitive.ObjectIDFromHex("testid")
				blog := &domains.PopulatedBlog{
					ID: oid,
				}
				tm.br.On("GetPopulatedBlogByID", ctx, mockReq).Return(blog, nil)
			},
			assertFn: func() {
				assert.NoError(t, err)
				assert.Nil(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			result, err = tm.svc.GetBlogByID(tt.args[0].(context.Context), tt.args[1].(string))
			tt.assertFn()
		})
	}
}

func TestListBlog(t *testing.T) {
	var result *domains.ListBlogResponse
	var err error
	mockReq := &domains.ListBlogRequest{
		Page:  1,
		Limit: 10,
	}

	tests := []test{
		{
			name: "should return error when list blog failed",
			args: []interface{}{
				ctx,
				&domains.ListBlogRequest{},
			},
			mockFn: func(tm *testModule) {
				listReq := &domains.PaginationOptions{
					Offset: int64((mockReq.Page - 1) * mockReq.Limit),
					Limit:  int64(mockReq.Limit),
				}
				tm.br.On("List", ctx, listReq).Return(nil, errors.New("error"))
			},
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, errmsg.BlogListFailed.Error())
			},
		},
		{
			name: "should return error when count blog failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				listReq := &domains.PaginationOptions{
					Offset: int64((mockReq.Page - 1) * mockReq.Limit),
					Limit:  int64(mockReq.Limit),
				}
				oid, _ := primitive.ObjectIDFromHex("testid")
				blogs := []domains.PopulatedBlog{
					{ID: oid},
					{ID: oid},
				}
				countReq := &domains.PaginationOptions{
					Offset: int64((mockReq.Page - 1) * mockReq.Limit),
					Limit:  int64(mockReq.Limit + 1),
				}
				tm.br.On("List", ctx, listReq).Return(blogs, nil)
				tm.br.On("Count", ctx, countReq).Return(int64(0), errors.New("error"))
			},
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, errmsg.BlogListFailed.Error())
			},
		},
		{
			name: "should return list blog when success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				listReq := &domains.PaginationOptions{
					Offset: int64((mockReq.Page - 1) * mockReq.Limit),
					Limit:  int64(mockReq.Limit),
				}
				oid, _ := primitive.ObjectIDFromHex("testid")
				blogs := []domains.PopulatedBlog{
					{ID: oid},
					{ID: oid},
				}
				countReq := &domains.PaginationOptions{
					Offset: int64((mockReq.Page - 1) * mockReq.Limit),
					Limit:  int64(mockReq.Limit + 1),
				}
				tm.br.On("List", ctx, listReq).Return(blogs, nil)
				tm.br.On("Count", ctx, countReq).Return(int64(3), nil)
			},
			assertFn: func() {
				assert.NoError(t, err)
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, result.HasNext, true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			result, err = tm.svc.ListBlog(tt.args[0].(context.Context), tt.args[1].(*domains.ListBlogRequest))
			tt.assertFn()
		})
	}
}

func TestUpdateBlogStatus(t *testing.T) {
	var err error

	tests := []test{
		{
			name: "should return error when status is invalid",
			args: []interface{}{
				ctx,
				&domains.UpdateBlogStatusRequest{
					BlogId: "blog_id",
					Status: "some text",
				},
			},
			mockFn: nil,
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, errmsg.BlogInvalidStatus.Error())
			},
		},
		{
			name: "should return error when update blog status failed",
			args: []interface{}{
				ctx,
				&domains.UpdateBlogStatusRequest{
					BlogId: "blog_id",
					Status: constants.IN_PROGRESS,
				},
			},
			mockFn: func(tm *testModule) {
				req := &domains.UpdateBlogStatusRequest{
					BlogId: "blog_id",
					Status: constants.IN_PROGRESS,
				}
				tm.br.On("UpdateStatus", ctx, req).Return(errors.New("error"))
			},
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, "error")
			},
		},
		{
			name: "should update blog status success",
			args: []interface{}{
				ctx,
				&domains.UpdateBlogStatusRequest{
					BlogId: "blog_id",
					Status: constants.DONE,
				},
			},
			mockFn: func(tm *testModule) {
				req := &domains.UpdateBlogStatusRequest{
					BlogId: "blog_id",
					Status: constants.DONE,
				}
				tm.br.On("UpdateStatus", ctx, req).Return(nil)
			},
			assertFn: func() {
				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			if tt.mockFn != nil {
				tt.mockFn(tm)
			}
			err = tm.svc.UpdateBlogStatus(tt.args[0].(context.Context), tt.args[1].(*domains.UpdateBlogStatusRequest))
			tt.assertFn()
		})
	}
}

func TestArchiveBlog(t *testing.T) {
	var err error
	mockReq := &domains.ArchiveBlogRequest{
		BlogId: "blog_id",
	}

	tests := []test{
		{
			name: "should return error when archive blog failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.br.On("Archive", ctx, mockReq).Return(errors.New("error"))
			},
			assertFn: func() {
				assert.Error(t, err)
				assert.EqualError(t, err, errmsg.BlogArchiveFailed.Error())
			},
		},
		{
			name: "should archive blog success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.br.On("Archive", ctx, mockReq).Return(nil)
			},
			assertFn: func() {
				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			err = tm.svc.ArchiveBlog(tt.args[0].(context.Context), tt.args[1].(*domains.ArchiveBlogRequest))
			tt.assertFn()
		})
	}
}
