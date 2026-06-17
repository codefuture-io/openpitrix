// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	pbam "github.com/codefuture-io/iam/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/manager"
)

type Client struct {
	pbam.AccessManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.AMServiceHost, constants.AMServicePort)
	if err != nil {
		return nil, err
	}

	return &Client{
		AccessManagerClient: pbam.NewAccessManagerClient(conn),
	}, nil
}
