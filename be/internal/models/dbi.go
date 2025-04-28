package models

import "context"

type UserRepo interface {
	AddUser(ctx context.Context, object interface{}) (userId string, err error)
	GetUserByIdentifier(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CheckUserExisting(ctx context.Context, username string) (bool, error)
	GetUserById(ctx context.Context, id string) (*User, error)
	UpdatePwd(ctx context.Context, id string, newPwd string) error
	UpdateRefreshToken(ctx context.Context, id string, refreshToken string) error
	UpdateTOTPSecret(ctx context.Context, id, secKey string) error
}

type TransRepo interface {
	AddTransaction(ctx context.Context, object interface{}) (
		string, error)
	GetTransactionByUserId(ctx context.Context, identifier string) (
		[]*Transaction, error)
}
type DbHub struct {
	DbcUser  UserRepo
	DbcTrans TransRepo
}
