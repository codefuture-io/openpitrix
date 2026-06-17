// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"
	"os"

	"google.golang.org/grpc"

	"github.com/codefuture-io/openpitrix/pkg/config"
	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/db"
	"github.com/codefuture-io/openpitrix/pkg/logger"
	"github.com/codefuture-io/openpitrix/pkg/manager"
	"github.com/codefuture-io/openpitrix/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/pi"
)

type Server struct {
	controller *Controller
}

func Serve(cfg *config.Config) {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Critical(nil, "Failed to get os hostname: %+v", err)
		return
	}
	pi.SetGlobal(cfg)
	jobController := NewController(hostname)
	s := Server{controller: jobController}
	ctx := context.Background()
	ctx = db.NewContext(ctx, cfg.Mysql)
	go jobController.Serve(ctx)

	manager.NewGrpcServer("job-controller", constants.JobManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithMysqlConfig(cfg.Mysql).
		Serve(func(server *grpc.Server) {
			pb.RegisterJobManagerServer(server, &s)
		})
}
