package storagemanager_test

import (
	"testing"

	"example.com/myproject/pkg/storagemanager"
)

func TestAzureConnect(t *testing.T) {
	t.Run("TestAzureConnect", func(t *testing.T) {
		name := storagemanager.GetCredentials().AzureAccountName
		key := storagemanager.GetCredentials().AzureSecretKey
		connection, err := storagemanager.Connect(name, key)
		if err != nil {
			t.Errorf("Error connecting to client: %v", err)
		}
		if connection.Client == nil {
			t.Errorf("Client is nil: %v", key)
		}
	})
}

func TestAzureUploadFile(t *testing.T) {
	t.Run("TestAzureUploadFile", func(t *testing.T) {
		name := storagemanager.GetCredentials().AzureAccountName
		key := storagemanager.GetCredentials().AzureSecretKey
		connection, err := storagemanager.Connect(name, key)
		if err != nil {
			t.Errorf("Error connecting to client: %v", err)
		}

		bytes := []byte("test")
		err = connection.UploadFile(bytes, "test.txt")
		if err != nil {
			t.Errorf("Error uploading file: %v", err)
		}
	})
}

func TestAzureDownloadFile(t *testing.T) {
	t.Run("TestAzureDownloadFile", func(t *testing.T) {
		name := storagemanager.GetCredentials().AzureAccountName
		key := storagemanager.GetCredentials().AzureSecretKey
		connection, err := storagemanager.Connect(name, key)
		if err != nil {
			t.Errorf("Error connecting to client: %v", err)
		}

		_, err = connection.DownloadFile("test.txt", "testcontainer")
		if err != nil {
			t.Errorf("Error downloading file: %v", err)
		}
	})
}
