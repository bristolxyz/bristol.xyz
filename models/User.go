package models

import (
	"context"
	"encoding/base64"
	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/env"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// DefaultPFPURL is the PFP for the default user.
var DefaultPFPURL = "/static/png/default_pfp.png"

// Initialises the e-mail index for users.
func init() {
	databaseInitTasks = append(databaseInitTasks, func() {
		cur, err := clients.MongoDatabase.Collection("users").Indexes().List(context.TODO())
		if err != nil {
			sentry.CaptureException(err)
			panic(err)
		}
		var indexes []bson.M
		if err = cur.All(context.TODO(), &indexes); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}
		if len(indexes) == 0 {
			u := true
			_, err = clients.MongoDatabase.Collection("users").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
				Keys:    bson.M{"email": 1},
				Options: &options.IndexOptions{Unique: &u},
			})
			if err != nil {
				sentry.CaptureException(err)
				panic(err)
			}
			InitialUser := env.RequiredEnvs("INITIAL_USER")["INITIAL_USER"]
			s := strings.Split(InitialUser, ":")
			Email := s[0]
			Password := s[1]
			NewUser(Email, Password, true, true, nil, nil)
		}
		cur, err = clients.MongoDatabase.Collection("tokens").Indexes().List(context.TODO())
		if err != nil {
			sentry.CaptureException(err)
			panic(err)
		}
		indexes = []bson.M{}
		if err = cur.All(context.TODO(), &indexes); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}
		if len(indexes) == 0 {
			u := true
			_, err = clients.MongoDatabase.Collection("tokens").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
				Keys:    bson.M{"userId": 1},
				Options: &options.IndexOptions{Unique: &u},
			})
			if err != nil {
				sentry.CaptureException(err)
				panic(err)
			}
		}
	})
}

// User is the model which is used to define a user.
type User struct {
	ID               string  `json:"id" bson:"_id"`
	Email            string  `json:"email" bson:"email"`
	Verified         bool    `json:"verified" bson:"verified"`
	Admin            bool    `json:"admin" bson:"admin"`
	Banned           bool    `json:"banned" bson:"banned"`
	PasswordHashSalt string  `json:"passwordHashSalt" bson:"passwordHashSalt"`
	FirstName        *string `json:"firstName" bson:"firstName"`
	LastName         *string `json:"lastName" bson:"lastName"`
	PFPUrl           string  `json:"pfpUrl" bson:"pfpUrl"`
}

type token struct {
	ID     string `json:"id" bson:"_id"`
	UserID string `json:"userId" bson:"userId"`
}

func (u *User) createPassword(Password string) {
	b, err := bcrypt.GenerateFromPassword([]byte(Password), 12)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	u.PasswordHashSalt = base64.StdEncoding.EncodeToString(b)
}

func set(x bson.M) bson.M {
	return bson.M{"$set": x}
}

// SetPassword is used to set a password in the database.
func (u *User) SetPassword(Password string) {
	u.createPassword(Password)
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"passwordHashSalt": u.PasswordHashSalt}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	_, _ = clients.MongoDatabase.Collection("tokens").DeleteMany(context.TODO(), bson.M{"userId": u.ID})
}

// SetFirstName is used to set a first name in the database.
func (u *User) SetFirstName(FirstName string) {
	u.FirstName = &FirstName
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"firstName": &FirstName}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
}

// SetLastName is used to set a last name in the database.
func (u *User) SetLastName(LastName string) {
	u.LastName = &LastName
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"lastName": &LastName}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
}

func (u *User) setVerifiedEmail(state bool) {
	u.Verified = state
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"verified": state}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
}

// Verify is used to verify an e-mail address.
func (u *User) Verify() {
	u.setVerifiedEmail(true)
}

// ChangeEmail is used to change an email within the database.
func (u *User) ChangeEmail(Email string) {
	u.Email = Email
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"email": Email}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	u.setVerifiedEmail(false)
}

// ChangePFPUrl is used to change the PFP URL within the database.
func (u *User) ChangePFPUrl(URL string) {
	u.PFPUrl = URL
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"pfpUrl": URL}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	u.setVerifiedEmail(false)
}

func (u *User) setBanState(state bool) {
	u.Banned = state
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"banned": state}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
}

// Ban is used to drop the ban hammer.
func (u *User) Ban() {
	u.setBanState(true)
	_, _ = clients.MongoDatabase.Collection("tokens").DeleteMany(context.TODO(), set(bson.M{"userId": u.ID}))
}

// Unban is used to unban a user.
func (u *User) Unban() {
	u.setBanState(false)
}

func (u *User) setAdminState(state bool) {
	u.Admin = state
	_, err := clients.MongoDatabase.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, set(bson.M{"admin": state}))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
}

// PromoteToAdmin is used to promote a user to admin.
func (u *User) PromoteToAdmin() {
	u.setAdminState(true)
}

// DemoteFromAdmin is used to demote a user from admin.
func (u *User) DemoteFromAdmin() {
	u.setAdminState(false)
}

// GetUserByEmail is used to check if a user exists with their e-mail address and return it if so.
func GetUserByEmail(Email string) *User {
	Email = strings.ToLower(Email)
	var usr User
	err := clients.MongoDatabase.Collection("users").FindOne(context.TODO(), bson.M{"email": Email}).Decode(&usr)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return &usr
}

// GetUserByID is used to get the user by their user ID.
func GetUserByID(ID string) *User {
	var usr User
	err := clients.MongoDatabase.Collection("users").FindOne(context.TODO(), bson.M{"_id": ID}).Decode(&usr)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return &usr
}

// GetUserByToken is used to get the user by their token.
func GetUserByToken(Token string) *User {
	var t token
	err := clients.MongoDatabase.Collection("tokens").FindOne(context.TODO(), bson.M{"_id": Token}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return GetUserByID(t.UserID)
}

// LoginUser is used to login a user. If the email/password is correct, a user pointer and a string pointer for a token will be returned.
// If it is not correct, these will be nil.
func LoginUser(Email string, Password string) (*User, *string) {
	u := GetUserByEmail(Email)
	if u == nil {
		// E-mail incorrect.
		return nil, nil
	}
	AlledgedPassword := []byte(Password)
	RealPasswordHashSalt, err := base64.StdEncoding.DecodeString(u.PasswordHashSalt)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	err = bcrypt.CompareHashAndPassword(RealPasswordHashSalt, AlledgedPassword)
	if err != nil {
		// Password incorrect.
		return nil, nil
	}
	ID := uuid.Must(uuid.NewUUID()).String()
	_, err = clients.MongoDatabase.Collection("tokens").InsertOne(context.TODO(), bson.M{"_id": ID, "userId": u.ID})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return u, &ID
}

// LogoutUser is used to logout a user.
func LogoutUser(Token string) {
	_, err := clients.MongoDatabase.Collection("tokens").DeleteOne(context.TODO(), bson.M{"_id": Token})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
}

// NewUser is used to create a new user and new token with it.
func NewUser(Email, Password string, Admin, Verified bool, FirstName, LastName *string) (*User, string) {
	UserID := uuid.Must(uuid.NewUUID()).String()
	u := &User{
		ID:        UserID,
		Email:     Email,
		Verified:  Verified,
		Admin:     Admin,
		Banned:    false,
		FirstName: FirstName,
		LastName:  LastName,
		PFPUrl:    DefaultPFPURL,
	}
	u.createPassword(Password)
	_, err := clients.MongoDatabase.Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	Token := uuid.Must(uuid.NewUUID()).String()
	_, err = clients.MongoDatabase.Collection("tokens").InsertOne(context.TODO(), bson.M{"_id": Token, "userId": u.ID})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return u, Token
}
