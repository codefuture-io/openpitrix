// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix runtime provider manager
package main

import (
	"github.com/codefuture-io/openpitrix/pkg/config"
	"github.com/codefuture-io/openpitrix/pkg/service/runtime_provider"
)

func main() {
	cfg := config.GetConf()
	runtime_provider.Serve(cfg)
}
