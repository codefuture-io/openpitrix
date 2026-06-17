// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package attachment

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	"github.com/codefuture-io/openpitrix/pkg/client/internals3"
	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/db"
	"github.com/codefuture-io/openpitrix/pkg/models"
	"github.com/codefuture-io/openpitrix/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/pi"
)

func getAttachments(ctx context.Context, attachmentIds []string) ([]*models.Attachment, error) {
	var as []*models.Attachment
	_, err := pi.Global().DB(ctx).
		Select(models.AttachmentColumns...).
		From(constants.TableAttachment).
		Where(db.Eq(constants.ColumnAttachmentId, attachmentIds)).
		Load(&as)
	return as, err
}

func getAttachment(ctx context.Context, attachmentId string) (*models.Attachment, error) {
	var a models.Attachment
	_, err := pi.Global().DB(ctx).
		Select(models.AttachmentColumns...).
		From(constants.TableAttachment).
		Where(db.Eq(constants.ColumnAttachmentId, attachmentId)).
		Load(&a)
	return &a, err
}

func removeAttachments(ctx context.Context, attachmentIds []string) error {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableAttachment).
		Where(db.Eq(constants.ColumnAttachmentId, attachmentIds)).
		Exec()
	return err
}

func isNoSuchKey(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode() == "NoSuchKey"
	}
	return false
}

func listAttachmentFilenames(ctx context.Context, attachment *models.Attachment) ([]string, error) {
	var filenames []string
	output, err := internals3.S3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: internals3.Bucket,
		Prefix: aws.String(attachment.GetObjectPrefix()),
	})
	if err != nil {
		return nil, err
	}
	for _, o := range output.Contents {
		if o.Key != nil {
			filenames = append(filenames, attachment.RemoveObjectName(*o.Key))
		}
	}
	return filenames, nil
}

func deleteAttachmentFiles(ctx context.Context, attachment *models.Attachment, filename ...string) error {
	var filenames []string
	var err error
	if len(filename) == 0 {
		filenames, err = listAttachmentFilenames(ctx, attachment)
		if err != nil {
			return err
		}
	} else {
		for _, f := range filename {
			filenames = append(filenames, f)
		}
	}

	for _, filename := range filenames {
		_, err := internals3.S3.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: internals3.Bucket,
			Key:    aws.String(attachment.GetObjectName(filename)),
		})
		if err != nil {
			if isNoSuchKey(err) {
				continue
			}
			return err
		}
	}
	return nil
}

type contents interface {
	GetAttachmentContent() map[string][]byte
}

func putAttachmentFiles(ctx context.Context, attachment *models.Attachment, contents contents) error {
	for filename, content := range contents.GetAttachmentContent() {
		_, err := internals3.S3.PutObject(ctx, &s3.PutObjectInput{
			Bucket: internals3.Bucket,
			Key:    aws.String(attachment.GetObjectName(filename)),
			Body:   bytes.NewReader(content),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type getAttachmentReq interface {
	GetFilename() []string
	GetIgnoreContent() bool
}

func getFile(ctx context.Context, attachment *models.Attachment, filename string) (*s3.GetObjectOutput, error) {
	return internals3.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: internals3.Bucket,
		Key:    aws.String(attachment.GetObjectName(filename)),
	})
}

func getAttachmentFiles(ctx context.Context, attachments []*models.Attachment, req getAttachmentReq) ([]*pb.Attachment, error) {
	var err error
	var pbAtts []*pb.Attachment
	for _, a := range attachments {
		var attachmentContent = make(map[string][]byte)
		var filenames []string

		if len(req.GetFilename()) == 0 {
			filenames, err = listAttachmentFilenames(ctx, a)
			if err != nil {
				return nil, err
			}
		} else {
			for _, filename := range req.GetFilename() {
				filenames = append(filenames, filename)
			}
		}

		for _, filename := range filenames {
			var content []byte
			if req.GetIgnoreContent() {
				attachmentContent[filename] = content
				continue
			}
			output, err := getFile(ctx, a, filename)
			if err != nil {
				if isNoSuchKey(err) {
					continue
				}
				return nil, err
			}
			content, err = io.ReadAll(output.Body)
			if err != nil {
				return nil, err
			}
			attachmentContent[filename] = content
		}

		var pbAttachment = models.AttachmentToPb(a)
		pbAttachment.AttachmentContent = attachmentContent
		pbAtts = append(pbAtts, pbAttachment)
	}
	return pbAtts, nil
}
