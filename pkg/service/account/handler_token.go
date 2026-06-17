// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"time"

	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/gerr"
	"github.com/codefuture-io/openpitrix/pkg/models"
	"github.com/codefuture-io/openpitrix/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/pi"
	"github.com/codefuture-io/openpitrix/pkg/util/ctxutil"
	"github.com/codefuture-io/openpitrix/pkg/util/jwtutil"
)

var (
	_ pb.TokenManagerServer = (*Server)(nil)
)

func (p *Server) CreateClient(ctx context.Context, req *pb.CreateClientRequest) (*pb.CreateClientResponse, error) {
	sender := ctxutil.GetSender(ctx)
	userId := sender.UserId
	client := models.NewUserClient(userId)
	_, err := pi.Global().DB(ctx).InsertInto(constants.TableUserClient).Record(client).Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	return &pb.CreateClientResponse{
		UserId:       client.UserId,
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
	}, nil
}

func (p *Server) Token(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	if req.GrantType == constants.GrantTypePassword {
		if req.Username == "" {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorParameterShouldNotBeEmpty, "username")
		}
		if req.Password == "" {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorParameterShouldNotBeEmpty, "password")
		}
	}
	if req.GrantType == constants.GrantTypeRefreshToken {
		if req.RefreshToken == "" {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorParameterShouldNotBeEmpty, "refresh_token")
		}
	}
	// validate client credentials
	user, userClient, err := validateClientCredentials(ctx, req)
	if err != nil {
		return nil, err
	}
	// if grant_type is password, switch user
	if req.GrantType == constants.GrantTypePassword {
		var isUserExist bool
		user, isUserExist, _ = validateUserAndGroupExist(ctx, req.Username)
		if !isUserExist {
			return nil, gerr.New(ctx, gerr.NotFound, gerr.ErrorEmailNotExists, req.Username)
		}
		isEmailPasswordMatched := validateUserPassword(ctx, user.UserId, req.Password)
		if !isEmailPasswordMatched {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorEmailPasswordNotMatched)
		}
	}
	userId := user.UserId
	var token *models.Token
	if req.GrantType == constants.GrantTypeRefreshToken {
		token, err = getTokenByRefreshToken(ctx, req.RefreshToken)
		if err != nil {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAuthFailure)
		}
		if token.CreateTime.Add(p.RefreshTokenExpireTime).Unix() <= time.Now().Unix() {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorRefreshTokenExpired)
		}
	} else {
		// reuse exist token
		token, err = getLastToken(ctx, userClient.ClientId, userId, req.Scope)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		// token not exists or expired
		if token == nil || token.CreateTime.Add(p.RefreshTokenExpireTime).Unix() <= time.Now().Unix() {
			// generate access token
			token, err = newToken(ctx, userClient.ClientId, req.Scope, userId)
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
		}
	}

	userId = token.UserId
	accessToken, err := jwtutil.Generate(p.IAMConfig.SecretKey, p.IAMConfig.ExpireTime, userId)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}
	idToken, err := jwtutil.Generate("", p.IAMConfig.ExpireTime, userId)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}

	return &pb.TokenResponse{
		TokenType:    ctxutil.TokenType,
		ExpiresIn:    int32(p.ExpireTime.Seconds()),
		AccessToken:  accessToken,
		IdToken:      idToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
