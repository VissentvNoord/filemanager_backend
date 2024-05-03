package storagemanager

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type File struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	UserID      primitive.ObjectID `bson:"userID"`
	Name        string             `bson:"name"`
	ContentType string             `bson:"contentType"`
	Size        int64              `bson:"size"`
	SourceURL   string             `bson:"sourceURL"`
	ParentID    primitive.ObjectID `bson:"parentID"`
	Date        primitive.DateTime `bson:"date"`
}

type Folder struct {
	ParentID primitive.ObjectID `bson:"parentID"`
}
