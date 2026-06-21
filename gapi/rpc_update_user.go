package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/tien29CE/Go_simple_bank.git/db/sqlc"
	"github.com/tien29CE/Go_simple_bank.git/pb"
	"github.com/tien29CE/Go_simple_bank.git/util"
	"github.com/tien29CE/Go_simple_bank.git/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateUpdateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hashpassword: %s", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid: true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time: time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "fail to update user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: converUser(&user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); req.Password != nil && err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); req.FullName != nil && err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); req.Email != nil && err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return
}
