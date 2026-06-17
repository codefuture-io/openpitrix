// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package internals3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	opconfig "github.com/codefuture-io/openpitrix/pkg/config"
)

var aConf = opconfig.GetConf().Attachment

var creds = credentials.NewStaticCredentialsProvider(
	aConf.AccessKey,
	aConf.SecretKey,
	"",
)

var cfg = aws.Config{
	Region:      "us-east-1",
	Credentials: creds,
}

var Bucket = aws.String(aConf.BucketName)

var S3 *s3.Client

func init() {
	S3 = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(aConf.Endpoint)
		o.UsePathStyle = true
	})
}
