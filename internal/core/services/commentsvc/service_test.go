package commentsvc_test

import (
	"context"
	"errors"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/core/ports/mocks"
	"robinhood/internal/core/services/commentsvc"
	"robinhood/internal/errmsg"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testModule struct {
	cr  *mocks.CommentRepository
	ur  *mocks.UserRepository
	svc ports.CommentService
}

type test struct {
	name     string
	args     []interface{}
	mockFn   func(*testModule)
	assertFn func(*testModule)
}

var (
	ctx = context.TODO()
)

func new(t *testing.T) *testModule {
	cr := mocks.NewCommentRepository(t)
	ur := mocks.NewUserRepository(t)
	return &testModule{
		cr:  cr,
		ur:  ur,
		svc: commentsvc.New(cr, ur),
	}
}

func TestCreateComment(t *testing.T) {
	var result *domains.PopulatedComment
	var err error
	mockReq := &domains.CreateCommentRequest{
		BlogId:   "blog-id",
		AuthorId: "author-id",
		Content:  "content",
	}

	tests := []*test{
		{
			name: "return error when create comment failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.cr.On("CreateTx", ctx, mockReq, mock.AnythingOfType("domains.CreateCommentFn")).Return(nil, errors.New("error"))
			},
			assertFn: func(tm *testModule) {
				tm.cr.AssertExpectations(t)
				assert.Error(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "return populated comment when success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				oid, _ := primitive.ObjectIDFromHex("testid")
				tm.cr.On("CreateTx", ctx, mockReq, mock.AnythingOfType("domains.CreateCommentFn")).Return(&domains.PopulatedComment{
					ID:     oid,
					BlogId: oid,
				}, nil)
			},
			assertFn: func(tm *testModule) {
				tm.cr.AssertExpectations(t)
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			result, err = tm.svc.CreateComment(tt.args[0].(context.Context), tt.args[1].(*domains.CreateCommentRequest))
			tt.assertFn(tm)
		})
	}
}

func TestCreateCommentTx(t *testing.T) {
	var result *domains.PopulatedComment
	var err error
	mockReq := &domains.CreateCommentRequest{
		BlogId:   "blog-id",
		AuthorId: "author-id",
		Content:  "content",
	}

	tests := []*test{
		{
			name: "should return error when create comment failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.cr.On("Create", ctx, mockReq).Return(nil, errors.New("error"))
			},
			assertFn: func(tm *testModule) {
				tm.cr.AssertExpectations(t)
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.EqualError(t, err, errmsg.CommentCreateFailed.Error())
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
				tm.cr.On("Create", ctx, mockReq).Return(&domains.Comment{
					AuthorId: oid,
				}, nil)
				tm.ur.On("GetByID", ctx, oid).Return(nil, errors.New("error"))
			},
			assertFn: func(tm *testModule) {
				tm.cr.AssertExpectations(t)
				tm.ur.AssertExpectations(t)
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.EqualError(t, err, errmsg.CommentCreateFailed.Error())
			},
		},
		{
			name: "should return populated comment when success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				oid, _ := primitive.ObjectIDFromHex("testid")
				tm.cr.On("Create", ctx, mockReq).Return(&domains.Comment{
					AuthorId: oid,
				}, nil)
				tm.ur.On("GetByID", ctx, oid).Return(&domains.User{
					ID: oid,
				}, nil)
			},
			assertFn: func(tm *testModule) {
				tm.cr.AssertExpectations(t)
				tm.ur.AssertExpectations(t)
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			result, err = tm.svc.CreateCommentTx(tt.args[0].(context.Context), tt.args[1].(*domains.CreateCommentRequest))
			tt.assertFn(tm)
		})
	}
}

func TestListComment(t *testing.T) {
	var result []domains.PopulatedComment
	var err error
	mockReq := &domains.ListCommentRequest{
		BlogId: "blog-id",
	}

	tests := []*test{
		{
			name: "should return error when list comment failed",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				tm.cr.On("List", ctx, mockReq.BlogId).Return(nil, errors.New("error"))
			},
			assertFn: func(tm *testModule) {
				tm.cr.AssertExpectations(t)
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.EqualError(t, err, errmsg.CommentListFailed.Error())
			},
		},
		{
			name: "should return populated comment when success",
			args: []interface{}{
				ctx,
				mockReq,
			},
			mockFn: func(tm *testModule) {
				oid, _ := primitive.ObjectIDFromHex("testid")
				tm.cr.On("List", ctx, mockReq.BlogId).Return([]domains.PopulatedComment{
					{
						ID:     oid,
						BlogId: oid,
					},
				}, nil)
			},
			assertFn: func(tm *testModule) {
				tm.cr.AssertExpectations(t)
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(t)
			tt.mockFn(tm)
			result, err = tm.svc.ListComment(tt.args[0].(context.Context), tt.args[1].(*domains.ListCommentRequest))
			tt.assertFn(tm)
		})
	}
}
