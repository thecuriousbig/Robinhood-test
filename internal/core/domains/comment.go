package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BlogId    primitive.ObjectID `bson:"blogId"`
	AuthorId  primitive.ObjectID `bson:"authorId"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"createdAt"`
}

type PopulatedComment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BlogId    primitive.ObjectID `bson:"blogId"`
	Author    User               `bson:"author"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"createdAt"`
}

type CreateCommentRequest struct {
	BlogId   string
	AuthorId string
	Content  string
}

type ListCommentRequest struct {
	BlogId string
}
