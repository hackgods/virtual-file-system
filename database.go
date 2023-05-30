package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	client     *mongo.Client
	collection *mongo.Collection
}


type FileMetadata struct {
	Filename  string `bson:"filename,omitempty"`
	FileSize  int64  `bson:"filesize,omitempty"`
	Checksum  string `bson:"checksum,omitempty"`
	Timestamp time.Time
}

func NewDatabase() *Database {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	db := client.Database("myfilesdb")
	collection := db.Collection("files")

	return &Database{
		client:     client,
		collection: collection,
	}
}

func (db *Database) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db.client.Disconnect(ctx)
}

func (db *Database) SaveFileMetadata(file *FileMetadata) error {
	_, err := db.collection.InsertOne(context.Background(), file)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateFileMetadata(file *FileMetadata) error {
	filter := bson.M{"filename": file.Filename}
	update := bson.M{"$set": bson.M{
		"filesize": file.FileSize,
		"checksum": file.Checksum,
	}}
	_, err := db.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteFileMetadata(filename string) error {
	filter := bson.M{"filename": filename}
	_, err := db.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetFileMetadata(filename string) (*FileMetadata, error) {
	filter := bson.M{"filename": filename}
	result := db.collection.FindOne(context.Background(), filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	var file FileMetadata
	err := result.Decode(&file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
