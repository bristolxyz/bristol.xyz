package models

import (
	"context"

	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// VerificationCode defines a verification code model.
type VerificationCode struct {
	ID     string `json:"id" bson:"_id"`
	UserID string `json:"userId" bson:"userId"`
}

// NewVerificationCode is used to create a verification code for a user.
func NewVerificationCode(u *User) (*VerificationCode, error) {
	// Create the verification ID.
	VerificationID := uuid.Must(uuid.NewUUID()).String()

	// Create the DB structure.
	item := &VerificationCode{
		ID:     VerificationID,
		UserID: u.ID,
	}

	// Insert the verification code into the DB.
	_, err := clients.MongoDatabase.Collection("verificationCodes").InsertOne(context.TODO(), item)
	if err != nil {
		return nil, err
	}

	// Return the DB item and no error.
	return item, nil
}

// GetVerificationCode is used to get a verification code if possible.
func GetVerificationCode(ID string) *VerificationCode {
	res := clients.MongoDatabase.Collection("verificationCodes").FindOne(context.TODO(), bson.M{"_id": ID})
	if res.Err() != nil {
		return nil
	}
	var v VerificationCode
	err := res.Decode(&v)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}
	return &v
}

// Verify is used to verify the user and then delete the object from the DB.
func (v *VerificationCode) Verify() {
	u := GetUserByID(v.UserID)
	if u != nil {
		u.Verify()
		_, err := clients.MongoDatabase.Collection("verificationCodes").DeleteOne(context.TODO(), v)
		if err != nil {
			sentry.CaptureException(err)
		}
	}
}
