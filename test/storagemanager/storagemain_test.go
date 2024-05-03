package storagemanager_test

import (
	"testing"

	"example.com/myproject/pkg/storagemanager"
)

func TestCredentials(t *testing.T) {
	// TestConnectClient tests the ConnectClient function
	// It should return a valid client and nil error
	t.Run("TestCredentials", func(t *testing.T) {
		credentials := storagemanager.GetCredentials()
		if credentials.MongoURI == "" {
			t.Errorf("Mongo URI is empty")
		}
	})
}
