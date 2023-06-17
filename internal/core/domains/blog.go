package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	AuthorId   primitive.ObjectID `bson:"authorId"`
	Status     string             `bson:"status"`
	IsArchived bool               `bson:"isArchived"`
	CreatedAt  time.Time          `bson:"createdAt"`
}

type PopulatedBlog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	Author     User               `bson:"author"`
	Status     string             `bson:"status"`
	IsArchived bool               `bson:"isArchived"`
	CreatedAt  time.Time          `bson:"createdAt"`
}

type CreateBlogRequest struct {
	Title    string
	Content  string
	AuthorId string
}

type ListBlogRequest struct {
	Page  uint32
	Limit uint32
}

type PaginationOptions struct {
	Offset int64
	Limit  int64
}

type ListBlogResponse struct {
	Data    []PopulatedBlog
	HasNext bool
}

type ListBlog struct {
	BLogs []PopulatedBlog `bson:"blogs"`
	Count int64           `bson:"count"`
}

type UpdateBlogStatusRequest struct {
	BlogId string
	Status string
}

type ArchiveBlogRequest struct {
	BlogId string
}
