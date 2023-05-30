package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Version struct {
	Version      int       `bson:"version"`
	Content      string    `bson:"content"`
	CreatedTime  time.Time `bson:"created_time"`
	ModifiedTime time.Time `bson:"modified_time"`
}

type Versioning struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type VFileMetadata struct {
	Filename  string    `bson:"filename"`
	Versions  []Version `bson:"versions"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewVersioning() (*Versioning, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	collection := client.Database("myfilesdb").Collection("files")

	return &Versioning{
		client:     client,
		collection: collection,
	}, nil
}

func (v *Versioning) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	v.client.Disconnect(ctx)
}

func (v *Versioning) GetLatestVersion(filename string) (int, error) {
	filter := bson.M{"filename": filename}
	opts := options.FindOne().SetSort(bson.M{"versions.version": -1})
	result := v.collection.FindOne(context.Background(), filter, opts)
	if result.Err() == mongo.ErrNoDocuments {
		return 0, nil
	} else if result.Err() != nil {
		return 0, result.Err()
	}

	var fileMetadata VFileMetadata
	err := result.Decode(&fileMetadata)
	if err != nil {
		return 0, err
	}

	latestVersion := 0
	if len(fileMetadata.Versions) > 0 {
		latestVersion = fileMetadata.Versions[0].Version
	}

	return latestVersion, nil
}

func (v *Versioning) CreateVersion(filename string, content string) error {
	fileMetadata := &VFileMetadata{
		Filename:  filename,
		Versions:  []Version{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	newVersion := Version{
		Version:      1,
		Content:      content,
		CreatedTime:  time.Now().UTC(),
		ModifiedTime: time.Now().UTC(),
	}

	fileMetadata.Versions = append(fileMetadata.Versions, newVersion)

	_, err := v.collection.InsertOne(context.Background(), fileMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (v *Versioning) AddVersion(filename string, content []byte) error {
	latestVersion, err := v.GetLatestVersion(filename)
	if err != nil {
		return err
	}

	newVersion := latestVersion + 1

	if latestVersion == 0 {
		// If no previous versions exist, use InsertOne
		newDocument := bson.M{
			"filename": filename,
			"versions": []bson.M{
				{
					"version":       newVersion,
					"content":       content,
					"created_time":  time.Now().UTC(),
					"modified_time": time.Now().UTC(),
				},
			},
			"updated_at": time.Now().UTC(),
		}

		_, err = v.collection.InsertOne(context.Background(), newDocument)
		if err != nil {
			return err
		}
	} else {
		// If previous versions exist, use UpdateOne
		filter := bson.M{"filename": filename}
		update := bson.M{
			"$push": bson.M{"versions": bson.M{
				"version":       newVersion,
				"content":       content,
				"created_time":  time.Now().UTC(),
				"modified_time": time.Now().UTC(),
			}},
			"$set": bson.M{
				"updated_at": time.Now().UTC(),
			},
		}

		_, err = v.collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func (v *Versioning) AddVersion(filename string, data []byte) error {
	latestVersion, err := v.GetLatestVersion(filename)
	if err != nil {
		return err
	}

	newVersion := 1
	if latestVersion != 0 {
		newVersion = latestVersion + 1
	}

	newVersionDocument := bson.M{
		"version":       newVersion,
		"content":       data,
		"created_time":  time.Now().UTC(),
		"modified_time": time.Now().UTC(),
	}

	_, err = v.collection.InsertOne(context.Background(), newVersionDocument)
	if err != nil {
		return err
	}

	return nil
}
*/
