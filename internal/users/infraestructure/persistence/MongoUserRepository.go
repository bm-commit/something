package persistence

import (
	"context"
	"errors"
	"log"
	"something/internal/users/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	con *mongo.Collection
}

// NewMongoUsersRepository ...
func NewMongoUsersRepository(m *mongo.Database) domain.UserRepository {
	return &mongoRepository{
		con: m.Collection("users"),
	}
}

func (r *mongoRepository) Find(criteria *domain.UserCriteria) ([]*domain.User, error) {
	findOptions := options.Find()
	findOptions.SetSkip((criteria.Page - 1) * criteria.PerPage)
	findOptions.SetLimit(criteria.PerPage)

	var users []*domain.User

	query := generateQueryWithCriteria(criteria)

	cur, err := r.con.Find(context.TODO(), query, findOptions)
	if err != nil {
		log.Println(err)
		return users, err
	}
	if err = cur.All(context.TODO(), &users); err != nil {
		log.Println(err)
		return users, err
	}
	return users, nil
}

func generateQueryWithCriteria(criteria *domain.UserCriteria) bson.D {
	query := bson.D{}
	if criteria.Query != "" {
		regex := primitive.Regex{Pattern: criteria.Query, Options: "i"}
		orCondition := primitive.E{Key: "$or", Value: []interface{}{
			bson.D{primitive.E{Key: "name", Value: regex}},
			bson.D{primitive.E{Key: "username", Value: regex}},
		}}
		query = append(query, orCondition)
	}
	return query
}

func (r *mongoRepository) FindByID(id string) (*domain.User, error) {
	var result *domain.User
	err := r.con.FindOne(
		context.TODO(),
		bson.D{primitive.E{Key: "id", Value: id}},
		options.FindOne()).Decode(&result)
	if result == nil {
		log.Println(err)
		return nil, errors.New("user not found")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func (r *mongoRepository) FindByEmail(email string) (*domain.User, error) {
	var user *domain.User
	err := r.con.FindOne(
		context.TODO(),
		bson.D{primitive.E{Key: "email", Value: email}},
		options.FindOne()).Decode(&user)
	if user == nil {
		return nil, errors.New("email not found")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func (r *mongoRepository) FindByUsername(username string) (*domain.User, error) {
	var user *domain.User
	err := r.con.FindOne(
		context.TODO(),
		bson.D{primitive.E{Key: "username", Value: username}},
		options.FindOne()).Decode(&user)
	if user == nil {
		return nil, errors.New("username not found")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func (r *mongoRepository) Update(user *domain.User) error {
	_, err := r.con.UpdateOne(context.TODO(), bson.M{"id": user.ID}, bson.D{
		{"$set", bson.D{
			primitive.E{Key: "name", Value: user.Name},
			primitive.E{Key: "username", Value: user.Username},
		},
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *mongoRepository) UpdateInterests(userID, bookID, status string) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.con.UpdateOne(context.TODO(), bson.M{"id": userID}, bson.D{
		{"$set", bson.D{
			primitive.E{Key: "interests." + bookID, Value: status},
		},
		},
	}, opts)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *mongoRepository) Save(user *domain.User) error {
	_, err := r.con.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *mongoRepository) Delete(id string) error {
	_, err := r.con.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoRepository) DeleteInterest(userID, bookID string) error {
	_, err := r.con.UpdateOne(context.TODO(), bson.M{"id": userID}, bson.D{
		{"$unset", bson.D{
			primitive.E{Key: "interests." + bookID, Value: ""},
		},
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
