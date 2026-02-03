package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	coll *mongo.Collection
}

// Initialize collection
func NewRepo(db *mongo.Database)*Repo{
	return &Repo{
		coll:db.Collection("users"),
	}
}

// Helper funcs()
func (r *Repo)FindByEmail(ctx context.Context,email string)(User,error){
	email=strings.ToLower(strings.TrimSpace(email))

	filter:=bson.M{"email":email}

	var u User

	err:=r.coll.FindOne(ctx,filter).Decode(&u)

	// check different types of errors
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return User{},mongo.ErrNoDocuments
		}
		return User{},fmt.Errorf("⚠️ User not found with the email: %w",err)
	}

	// if all, ok..
	return u,nil
}

func (r *Repo)Create(ctx context.Context, u User)(User,error){
	resp,err:=r.coll.InsertOne(ctx, u)

	if err != nil {
		return User{},fmt.Errorf("⚠️ User insertion failed: %w",err)
	}

	id,ok:=resp.InsertedID.(primitive.ObjectID) // id: _id

	if !ok{
		return User{},fmt.Errorf("⚠️ User insertion failed and inserted id is not objectID: %w",err)
	}

	// If all ok, then..
	u.ID=id;
	return u,nil
}