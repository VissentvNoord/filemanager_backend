package storagemanager_test

import (
	"context"
	"testing"

	"example.com/myproject/pkg/storagemanager"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestConnectClient(t *testing.T) {
	t.Run("TestConnectClient", func(t *testing.T) {
		mongoURI := storagemanager.GetCredentials().MongoURI
		connection, err := storagemanager.ConnectClient(mongoURI)
		if err != nil {
			t.Errorf("Error connecting to client: %v", err)
		}
		if connection.Client == nil {
			t.Errorf("Client is nil")
		}
	})
}

func TestCreateAndRemoveFileRecord(t *testing.T) {
	t.Run("TestCreateAndRemoveFileRecord", func(t *testing.T) {
		connection, err := storagemanager.ConnectClient(storagemanager.GetCredentials().MongoURI)
		if err != nil {
			t.Errorf("Error connecting to client: %v", err)
		}

		// Start a new session
		session, err := connection.Client.StartSession()
		if err != nil {
			t.Errorf("Error starting session: %v", err)
		}
		defer session.EndSession(context.Background())

		// Start a new transaction
		err = session.StartTransaction()
		if err != nil {
			t.Errorf("Error starting transaction: %v", err)
			return
		}

		file := storagemanager.File{
			ID:          primitive.NewObjectID(),
			Name:        "test.txt",
			ContentType: "TEST",
		}

		// Insert the file record into MongoDB
		createdID, createError := connection.CreateFileRecord(context.Background(), file)
		if createError != nil {
			session.AbortTransaction(context.Background())
			t.Errorf("Error creating file: %v", createError)
			return
		}

		// Remove the file record from MongoDB
		removeError := connection.RemoveFileRecord(createdID.Hex())
		if removeError != nil {
			session.AbortTransaction(context.Background())
			t.Errorf("Error removing file: %v", removeError)
			return
		}

		err = session.CommitTransaction(context.Background())
		if err != nil {
			t.Errorf("Error committing transaction: %v", err)
		}
	})
}
