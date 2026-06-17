// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"github.com/codefuture-io/openpitrix/pkg/apigateway"
	"github.com/codefuture-io/openpitrix/pkg/config"
)

func main() {
	cfg := config.GetConf()
	apigateway.Serve(cfg)
}
