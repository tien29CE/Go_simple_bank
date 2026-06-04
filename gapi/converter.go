package gapi

import (
	db "github.com/tien29CE/Go_simple_bank.git/db/sqlc"
	"github.com/tien29CE/Go_simple_bank.git/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converUser(user *db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
