package persistence

import (
	"context"
	"log"
	"time"

	ml "github.com/yaroyan/ms/logger/domain/model/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository struct {
	Client *mongo.Client
}

func (l *LogRepository) Insert(e ml.Log) error {
	collection := l.Client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), ml.Log{
		Name:      e.Name,
		Data:      e.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error inserting into logs: ", err)
		return err
	}

	return nil
}

func (l *LogRepository) All() ([]*ml.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := l.Client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error: ", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*ml.Log

	for cursor.Next(ctx) {
		var item ml.Log

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log int slice: ", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *LogRepository) GetOne(id string) (*ml.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := l.Client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry ml.Log
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogRepository) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := l.Client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogRepository) Update(e ml.Log) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := l.Client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(e.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: e.Name},
				{Key: "data", Value: e.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
