package mongodb

import (
	"app/be/internal/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type User struct {
	*Base
}

func NewUserRepo() *User {
	return &User{
		Base: &Base{
			collection: GetDB().Collection("users"),
		},
	}
}

func (ins *User) AddUser(ctx context.Context, object interface{}) (userId string, err error) {
	id, err := ins.Base.Add(ctx, object)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (ins *User) GetUserByIdentifier(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := ins.Base.Get(ctx, map[string]interface{}{"identifier": username}, &user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &user, nil
}
func (ins *User) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := ins.Base.Get(ctx, map[string]interface{}{"email": email}, &user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
func (ins *User) CheckUserExisting(ctx context.Context, username string) (bool, error) {
	count, err := ins.Base.Count(ctx, map[string]interface{}{"username": username})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
func (ins *User) GetUserById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	idObject, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = ins.Base.Get(ctx, map[string]interface{}{"_id": idObject}, &user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &user, nil
}
func (ins *User) UpdatePwd(ctx context.Context, id string, newPwd string) error {
	idObject, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = ins.Base.Update(ctx, map[string]interface{}{"_id": idObject}, map[string]interface{}{"$set": map[string]interface{}{"password": newPwd}})
	if err != nil {
		return err
	}
	return nil
}
func (ins *User) UpdateRefreshToken(ctx context.Context, id string, refreshToken string) error {
	idObject, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = ins.Base.Update(ctx, map[string]interface{}{"_id": idObject}, map[string]interface{}{"$set": map[string]interface{}{"refresh_token": refreshToken}})
	if err != nil {
		return err
	}
	return nil
}
func (ins *User) UpdateTOTPSecret(ctx context.Context, id, secKey string) error {
	idObject, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if err := ins.Base.Update(ctx, map[string]interface{}{"_id": idObject}, map[string]interface{}{"$set": map[string]interface{}{"totpSecret": secKey}}); err != nil {
		return err
	}

	return nil
}
