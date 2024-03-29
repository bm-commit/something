package persistence

import (
	"context"
	"errors"
	"log"
	"something/internal/bookreviews/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	con *mongo.Collection
}

// NewMongoBookReviewRepository ...
func NewMongoBookReviewRepository(m *mongo.Database) domain.BookReviewRepository {
	return &mongoRepository{
		con: m.Collection("book_reviews"),
	}
}

func (r *mongoRepository) Find(bookID string) ([]*domain.BookReview, error) {
	var bookReviews []*domain.BookReview

	cur, err := r.con.Find(context.TODO(), bson.D{primitive.E{Key: "bookid", Value: bookID}}, nil)
	if err != nil {
		log.Println(err)
		return bookReviews, err
	}

	if err = cur.All(context.TODO(), &bookReviews); err != nil {
		log.Println(err)
		return bookReviews, err
	}

	return bookReviews, nil
}

func (r *mongoRepository) FindByID(id string) (*domain.BookReview, error) {
	var result *domain.BookReview
	err := r.con.FindOne(
		context.TODO(),
		bson.D{primitive.E{Key: "id", Value: id}},
		options.FindOne()).Decode(&result)
	if result == nil {
		log.Println(err)
		return nil, errors.New("book review not found")
	}
	return result, nil
}

func (r *mongoRepository) FindReviews(criteria *domain.BookReviewCriteria) ([]*domain.BookReviewShort, error) {

	var bookReviews []*domain.BookReviewShort

	sortStage := bson.D{{"$sort", bson.D{{"rating", criteria.Sort}}}}
	groupStage := bson.D{
		{"$group",
			bson.D{
				{"_id", "$bookid"},
				{"total",
					bson.D{{"$sum", 1}}},
				{"rating", bson.D{{"$avg", "$rating"}}},
			},
		}}
	limit := bson.D{{"$limit", 25}}

	cur, err := r.con.Aggregate(
		context.TODO(),
		mongo.Pipeline{
			sortStage,
			groupStage,
			limit,
		})
	if err != nil {
		log.Println(err)
		return bookReviews, err
	}

	if err = cur.All(context.TODO(), &bookReviews); err != nil {
		log.Println(err)
		return bookReviews, err
	}

	return bookReviews, nil
}

func (r *mongoRepository) Update(bookReview *domain.BookReview) error {
	_, err := r.con.UpdateOne(context.TODO(), bson.M{"id": bookReview.ID}, bson.D{
		{"$set", bson.D{
			primitive.E{Key: "text", Value: bookReview.Text},
		},
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *mongoRepository) Save(bookReview *domain.BookReview) error {
	_, err := r.con.InsertOne(context.TODO(), bookReview)
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
