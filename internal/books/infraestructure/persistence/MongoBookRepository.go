package persistence

import (
	"context"
	"errors"
	"log"
	"something/internal/books/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	con *mongo.Collection
}

// NewMongoBookRepository ...
func NewMongoBookRepository(m *mongo.Database) domain.BookRepository {
	return &mongoRepository{
		con: m.Collection("books"),
	}
}

func (r *mongoRepository) Find(criteria *domain.BookCriteria) ([]*domain.Book, error) {
	findOptions := options.Find()
	findOptions.SetSkip((criteria.Page - 1) * criteria.PerPage)
	findOptions.SetLimit(criteria.PerPage)

	var books []*domain.Book

	query := generateQueryWithCriteria(criteria)

	cur, err := r.con.Find(context.TODO(), query, findOptions)
	if err != nil {
		log.Println(err)
		return books, err
	}
	if err = cur.All(context.TODO(), &books); err != nil {
		log.Println(err)
		return books, err
	}
	return books, nil
}

func generateQueryWithCriteria(criteria *domain.BookCriteria) bson.D {
	query := bson.D{}
	if criteria.Query != "" {
		regex := primitive.Regex{Pattern: criteria.Query, Options: "i"}
		condition := primitive.E{Key: "title", Value: regex}
		query = append(query, condition)
	}
	if criteria.Author != "" {
		regex := primitive.Regex{Pattern: criteria.Author, Options: "i"}
		condition := primitive.E{Key: "author", Value: regex}
		query = append(query, condition)
	}
	if criteria.Genre != "" {
		regex := primitive.Regex{Pattern: criteria.Genre, Options: "i"}
		condition := primitive.E{Key: "genre", Value: regex}
		query = append(query, condition)
	}
	return query
}

func (r *mongoRepository) FindByID(id string) (*domain.Book, error) {
	var result *domain.Book
	err := r.con.FindOne(
		context.TODO(),
		bson.D{primitive.E{Key: "id", Value: id}},
		options.FindOne()).Decode(&result)
	if result == nil {
		log.Println(err)
		return nil, errors.New("book not found")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func (r *mongoRepository) Update(book *domain.Book) error {
	_, err := r.con.UpdateOne(context.TODO(), bson.M{"id": book.ID}, bson.D{
		{"$set", bson.D{
			primitive.E{Key: "title", Value: book.Title},
			primitive.E{Key: "description", Value: book.Description},
			primitive.E{Key: "author", Value: book.Author},
			primitive.E{Key: "genre", Value: book.Genre},
			primitive.E{Key: "pages", Value: book.Pages},
		},
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *mongoRepository) Save(book *domain.Book) error {
	_, err := r.con.InsertOne(context.TODO(), book)
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
