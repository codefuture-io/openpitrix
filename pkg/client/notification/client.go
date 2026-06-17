// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"

	nfpb "github.com/codefuture-io/notification/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/logger"
	"github.com/codefuture-io/openpitrix/pkg/manager"
	"github.com/codefuture-io/openpitrix/pkg/models"
	"github.com/codefuture-io/openpitrix/pkg/util/jsonutil"
	"github.com/codefuture-io/openpitrix/pkg/util/pbutil"
)

type Client struct {
	nfpb.NotificationClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.NotificationHost, constants.NotificationPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		NotificationClient: nfpb.NewNotificationClient(conn),
	}, nil
}

func SendEmailNotification(ctx context.Context, emailNotifications []*models.EmailNotification) error {
	client, err := NewClient()
	if err != nil {
		logger.Error(ctx, "Failed to create notification client: %+v", err)
		return err
	}
	emailNotifications = models.UniqueEmailNotifications(emailNotifications)
	for _, notification := range emailNotifications {
		_, err := client.CreateNotification(ctx, &nfpb.CreateNotificationRequest{
			ContentType: pbutil.ToProtoString(notification.ContentType),
			Title:       pbutil.ToProtoString(notification.Title),
			Content:     pbutil.ToProtoString(notification.Content),
			Owner:       pbutil.ToProtoString(notification.Owner),
			AddressInfo: pbutil.ToProtoString(
				jsonutil.ToString(map[string][]string{
					constants.NfTypeEmail: notification.Addresses,
				}),
			),
		})

		if err != nil {
			logger.Error(ctx, "Failed to send email, %+v", err)
			return err
		}
	}

	return nil
}
