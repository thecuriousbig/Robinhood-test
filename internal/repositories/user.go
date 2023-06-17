package repositories

import (
	"context"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	mc  *mongo.Client
	db  string
	cn  string
	col *mongo.Collection
}

func NewUserRepository(mc *mongo.Client, db string) ports.UserRepository {
	cn := "user"
	col := mc.Database(db).Collection(cn)
	// create index
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"username": 1},
	})
	return &userRepository{
		mc:  mc,
		db:  db,
		cn:  cn,
		col: col,
	}
}

func (r *userRepository) Create(ctx context.Context, req *domains.CreateUserRequest) (*domains.User, error) {
	return r.insertOne(ctx, domains.User{
		Username:     req.Username,
		Password:     req.Password,
		Email:        req.Email,
		ProfileImage: "",
	})
}

func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domains.User, error) {
	return r.findOne(ctx, bson.M{"_id": id}, nil)
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domains.User, error) {
	return r.findOne(ctx, bson.M{"username": username}, nil)
}

func (r *userRepository) Update(ctx context.Context, req *domains.UpdateUserRequest) (*domains.User, error) {
	oid, _ := primitive.ObjectIDFromHex(req.UserId)
	return r.updateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"profileImage": req.ProfileImage}})
}

func (r *userRepository) insertOne(ctx context.Context, in domains.User) (*domains.User, error) {
	in.CreatedAt = time.Now().UTC()
	result, err := r.col.InsertOne(ctx, in)
	oid, _ := result.InsertedID.(primitive.ObjectID)
	in.ID = oid
	return &in, err
}

func (r *userRepository) findOne(ctx context.Context, filter bson.M, opts *options.FindOneOptions) (*domains.User, error) {
	var result domains.User
	if opts == nil {
		opts = &options.FindOneOptions{}
	}
	if err := r.col.FindOne(ctx, filter, opts).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *userRepository) updateOne(ctx context.Context, filter bson.M, update bson.M) (*domains.User, error) {
	var result domains.User
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	return &result, err
}
