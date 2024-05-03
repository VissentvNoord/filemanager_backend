package storagemanager

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	credentials StorageCrendentials
	Conn        *Connections
)

type StorageCrendentials struct {
	MongoURI         string
	AzureAccountName string
	AzureSecretKey   string
}

type Connections struct {
	MongoConnection *MongoConnection
	AzureConnection *AzureConnection
}

func init() {
	// Get the absolute path to the root directory of the project
	rootDir := getRootDir() // Get the current working directory

	// Load environment variables from .env file in the root directory
	envFile := filepath.Join(rootDir, ".env")
	err := godotenv.Load(envFile)
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	cred, err := LoadCredentials()
	if err != nil {
		panic("Error loading credentials: " + err.Error())
	}

	credentials = cred

	_, err = CreateConnections()
	if err != nil {
		panic("Error creating connections: " + err.Error())
	}
}

func LoadCredentials() (StorageCrendentials, error) {
	mongoURI := os.Getenv("MONGO_URI")
	azureAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	azureSecretKey := os.Getenv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	credentials := StorageCrendentials{
		MongoURI:         mongoURI,
		AzureAccountName: azureAccountName,
		AzureSecretKey:   azureSecretKey,
	}

	fmt.Println(credentials)
	return credentials, nil
}

func GetCredentials() *StorageCrendentials {
	return &credentials
}

func getRootDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return currentDir
}

func CreateConnections() (*Connections, error) {
	connections := &Connections{}

	mongoURI := GetCredentials().MongoURI
	mongoConnection, err := ConnectClient(mongoURI)
	if err != nil {
		return nil, err
	}
	connections.MongoConnection = mongoConnection

	azureAccountName := GetCredentials().AzureAccountName
	azureSecretKey := GetCredentials().AzureSecretKey
	azureConnection, err := Connect(azureAccountName, azureSecretKey)
	if err != nil {
		return nil, err
	}

	connections.AzureConnection = azureConnection

	Conn = &Connections{
		MongoConnection: mongoConnection,
		AzureConnection: azureConnection,
	}

	return Conn, nil
}

func UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (File, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return File{}, err
	}
	defer file.Close()

	ID := primitive.NewObjectID()

	// Create a new file record
	newFileRecord := File{
		ID:          ID,
		Name:        fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Size:        fileHeader.Size,
		Date:        primitive.NewDateTimeFromTime(time.Now()),
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return File{}, err
	}

	err = Conn.AzureConnection.UploadFile(bytes, ID.Hex())
	if err != nil {
		return File{}, err
	}

	_, err = Conn.MongoConnection.CreateFileRecord(ctx, newFileRecord)
	if err != nil {
		return File{}, err
	}

	return newFileRecord, nil
}

func DeleteFile(id string) error {
	err := Conn.AzureConnection.DeleteFile(id, "testcontainer")
	if err != nil {
		return err
	}

	err = Conn.MongoConnection.RemoveFileRecord(id)
	if err != nil {
		return err
	}

	return nil
}
