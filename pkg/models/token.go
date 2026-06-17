// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"github.com/codefuture-io/openpitrix/pkg/constants"

	"github.com/codefuture-io/openpitrix/pkg/db"
	"github.com/codefuture-io/openpitrix/pkg/util/idutil"
)

func NewTokenId() string {
	return idutil.GetUuid("token-")
}

type Token struct {
	TokenId      string
	ClientId     string
	RefreshToken string
	Scope        string
	UserId       string
	Status       string
	CreateTime   time.Time
	StatusTime   time.Time
}

var TokenColumns = db.GetColumnsFromStruct(&Token{})

func NewToken(clientId, userId, scope string) *Token {
	return &Token{
		TokenId:      NewTokenId(),
		ClientId:     clientId,
		RefreshToken: idutil.GetRefreshToken(),
		Scope:        scope,
		UserId:       userId,
		Status:       constants.StatusActive,
		CreateTime:   time.Now(),
		StatusTime:   time.Now(),
	}
}
