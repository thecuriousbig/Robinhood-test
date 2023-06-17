package repositories

import (
	"context"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type commentRepository struct {
	mc  *mongo.Client
	db  string
	cn  string
	col *mongo.Collection
}

func NewCommentRepository(mc *mongo.Client, db string) ports.CommentRepository {
	cn := "comment"
	col := mc.Database(db).Collection(cn)
	// create index
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"blogId": 1},
	})
	return &commentRepository{
		mc:  mc,
		db:  db,
		cn:  "comment",
		col: col,
	}
}

func (r *commentRepository) Create(ctx context.Context, req *domains.CreateCommentRequest) (*domains.Comment, error) {
	bid, _ := primitive.ObjectIDFromHex(req.BlogId)
	aid, _ := primitive.ObjectIDFromHex(req.AuthorId)
	return r.insertOne(ctx, domains.Comment{
		BlogId:   bid,
		AuthorId: aid,
		Content:  req.Content,
	})
}

func (r *commentRepository) CreateTx(ctx context.Context, req *domains.CreateCommentRequest, fn domains.CreateCommentFn) (*domains.PopulatedComment, error) {
	session, err := r.mc.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	res, err := session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		cart, err := fn(sc, req)
		if err != nil {
			return nil, err
		}
		return cart, err
	})
	if err != nil {
		return nil, err
	}

	return res.(*domains.PopulatedComment), nil
}

func (r *commentRepository) List(ctx context.Context, blogId string) ([]domains.PopulatedComment, error) {
	oid, _ := primitive.ObjectIDFromHex(blogId)
	result := []domains.PopulatedComment{}
	pipeline := []bson.M{
		{"$match": bson.M{"blogId": oid}},
		{"$sort": bson.M{"createdAt": -1}},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "authorId",
				"foreignField": "_id",
				"as":           "author",
			},
		},
		{"$unwind": "$author"},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *commentRepository) insertOne(ctx context.Context, in domains.Comment) (*domains.Comment, error) {
	in.CreatedAt = time.Now().UTC()
	result, err := r.col.InsertOne(ctx, in)
	oid, _ := result.InsertedID.(primitive.ObjectID)
	in.ID = oid
	return &in, err
}
