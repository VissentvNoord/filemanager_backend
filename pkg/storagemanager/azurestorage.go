package storagemanager

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureConnection struct {
	Client *azblob.Client
}

func Connect(name string, key string) (*AzureConnection, error) {
	cred, err := azblob.NewSharedKeyCredential(name, key)
	handleError(err)

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", name), cred, nil)
	handleError(err)

	ac := &AzureConnection{
		Client: client,
	}

	return ac, nil
}

func (ac *AzureConnection) UploadFile(blobData []byte, blobName string) error {
	client := ac.Client
	containerName := "testcontainer"

	_, err := client.UploadStream(context.TODO(),
		containerName,
		blobName,
		bytes.NewReader(blobData),
		&azblob.UploadStreamOptions{
			Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
			Tags:     map[string]string{"Year": "2022"},
		})

	if err != nil {
		return err
	}

	return nil
}

func (ac *AzureConnection) DeleteFile(blobName string, containerName string) error {
	client := ac.Client
	_, err := client.DeleteBlob(context.TODO(), containerName, blobName, nil)
	if err != nil {
		return err
	}

	return nil
}

func (ac *AzureConnection) DownloadFile(blobName string, containerName string) ([]byte, error) {
	client := ac.Client
	// Download the blob's contents and ensure that the download worked properly
	blobDownloadResponse, err := client.DownloadStream(context.TODO(), containerName, blobName, nil)
	if err != nil {
		return nil, err
	}

	// Use the bytes.Buffer object to read the downloaded data.
	// RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	reader := blobDownloadResponse.Body
	downloadData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = reader.Close()
	if err != nil {
		return nil, err
	}

	return downloadData, nil
}

func handleError(err error) bool {
	if err != nil {
		panic(err)
	}

	return false
}
