// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix runtime manager
package main

import (
	"github.com/codefuture-io/openpitrix/pkg/config"
	"github.com/codefuture-io/openpitrix/pkg/service/runtime"
)

func main() {
	cfg := config.GetConf()
	runtime.Serve(cfg)
}
