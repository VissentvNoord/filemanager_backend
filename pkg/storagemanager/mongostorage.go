package storagemanager

import (
	"context"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database   = "filemanager"
	collection = "files"
)

type MongoConnection struct {
	Client *mongo.Client
	Ctx    context.Context
}

func ConnectClient(connectionString string) (*MongoConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	connection := &MongoConnection{Client: client, Ctx: ctx}

	return connection, err
}

func (conn *MongoConnection) CreateFileRecord(ctx context.Context, file File) (primitive.ObjectID, error) {
	client := conn.Client

	collection := client.Database(database).Collection(collection)
	insertResult, err := collection.InsertOne(ctx, file)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("Error converting InsertedID to ObjectID")
	}

	return insertedID, nil
}

func (conn *MongoConnection) RemoveFileRecord(id string) error {
	client := conn.Client

	objectID, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return err1
	}

	collection := client.Database(database).Collection(collection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}

func (conn *MongoConnection) GetAllFiles() ([]File, error) {
	client := conn.Client
	files := []File{}

	collection := client.Database(database).Collection(collection)
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return files, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var file File
		err := cur.Decode(&file)
		if err != nil {
			return files, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (conn *MongoConnection) GetFile(ID string) (File, error) {
	client := conn.Client
	file := File{}

	collection := client.Database(database).Collection(collection)
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return file, err
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&file)
	if err != nil {
		return file, err
	}

	return file, nil
}
