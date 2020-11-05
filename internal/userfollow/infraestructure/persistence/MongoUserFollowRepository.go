package persistence

import (
	"context"
	"log"
	"something/internal/userfollow/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	con *mongo.Collection
}

// NewMongoUserFollowRepository ...
func NewMongoUserFollowRepository(m *mongo.Database) domain.UserFollowRepository {
	return &mongoRepository{
		con: m.Collection("user_follows"),
	}
}

func (r *mongoRepository) FindFollowing(id string) ([]*domain.UserFollow, error) {
	var following []*domain.UserFollow
	cur, err := r.con.Find(context.TODO(), bson.D{primitive.E{Key: "from", Value: id}}, nil)
	if err != nil {
		log.Println(err)
		return following, err
	}

	if err = cur.All(context.TODO(), &following); err != nil {
		log.Println(err)
		return following, err
	}
	return following, nil
}

func (r *mongoRepository) FindFollowers(id string) ([]*domain.UserFollow, error) {
	var followers []*domain.UserFollow
	cur, err := r.con.Find(context.TODO(), bson.D{primitive.E{Key: "to", Value: id}}, nil)
	if err != nil {
		log.Println(err)
		return followers, err
	}

	if err = cur.All(context.TODO(), &followers); err != nil {
		log.Println(err)
		return followers, err
	}
	return followers, nil
}

func (r *mongoRepository) Follow(u *domain.UserFollow) error {
	_, err := r.con.InsertOne(context.TODO(), u)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *mongoRepository) Unfollow(u *domain.UserFollow) error {
	_, err := r.con.DeleteOne(
		context.TODO(),
		bson.D{
			primitive.E{Key: "$and", Value: []interface{}{
				bson.D{primitive.E{Key: "from", Value: u.From}},
				bson.D{primitive.E{Key: "to", Value: u.To}},
			}},
		})
	if err != nil {
		return err
	}
	return nil
}
