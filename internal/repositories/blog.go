package repositories

import (
	"context"
	"fmt"
	"robinhood/internal/core/constants"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogRepository struct {
	mc  *mongo.Client
	db  string
	cn  string
	col *mongo.Collection
}

func NewBlogRepository(mc *mongo.Client, db string) ports.BlogRepository {
	cn := "blog"
	col := mc.Database(db).Collection(cn)
	return &blogRepository{
		mc:  mc,
		db:  db,
		cn:  cn,
		col: col,
	}
}

func (r *blogRepository) Create(ctx context.Context, req *domains.CreateBlogRequest) (*domains.Blog, error) {
	aid, _ := primitive.ObjectIDFromHex(req.AuthorId)
	return r.insertOne(ctx, domains.Blog{
		Title:      req.Title,
		Content:    req.Content,
		AuthorId:   aid,
		Status:     constants.TO_DO,
		IsArchived: false,
	})
}

func (r *blogRepository) CreateTx(ctx context.Context, req *domains.CreateBlogRequest, fn domains.CreateBlogFn) (*domains.PopulatedBlog, error) {
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

	return res.(*domains.PopulatedBlog), nil
}

func (r *blogRepository) GetPopulatedBlogByID(ctx context.Context, id string) (*domains.PopulatedBlog, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	result := &domains.PopulatedBlog{}
	pipeline := []bson.M{
		{"$match": bson.M{"_id": oid, "isArchived": false}},
		{"$limit": 1},
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
	for cursor.Next(ctx) {
		if err := cursor.Decode(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r *blogRepository) List(ctx context.Context, req *domains.PaginationOptions) ([]domains.PopulatedBlog, error) {
	result := []domains.PopulatedBlog{}
	// aggregate pipeline to get populated blog + count numbers of blogs and map to result object
	pipeline := []bson.M{
		{"$match": bson.M{"isArchived": false}},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": req.Offset},
		{"$limit": req.Limit},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "authorId",
				"foreignField": "_id",
				"as":           "author",
			},
		},
		{"$unwind": "$author"},
		{"$project": bson.M{"comments": 0}},
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

func (r *blogRepository) Count(ctx context.Context, req *domains.PaginationOptions) (int64, error) {
	filter := bson.M{"isArchived": false}
	opts := options.Count().SetSkip(req.Offset).SetLimit(req.Limit)
	return r.col.CountDocuments(ctx, filter, opts)
}

func (r *blogRepository) UpdateStatus(ctx context.Context, req *domains.UpdateBlogStatusRequest) error {
	oid, _ := primitive.ObjectIDFromHex(req.BlogId)
	_, err := r.updateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": req.Status}})
	return err
}

func (r *blogRepository) Archive(ctx context.Context, req *domains.ArchiveBlogRequest) error {
	oid, _ := primitive.ObjectIDFromHex(req.BlogId)
	_, err := r.updateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"isArchived": true}})
	return err
}

func (r *blogRepository) insertOne(ctx context.Context, in domains.Blog) (*domains.Blog, error) {
	in.CreatedAt = time.Now().UTC()
	fmt.Printf("in: %+v\n", in)
	result, err := r.col.InsertOne(ctx, in)
	oid, _ := result.InsertedID.(primitive.ObjectID)
	in.ID = oid
	return &in, err
}

func (r *blogRepository) updateOne(ctx context.Context, filter bson.M, update bson.M) (*domains.Blog, error) {
	var result domains.Blog
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	return &result, err
}
