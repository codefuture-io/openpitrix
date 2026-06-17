// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service_config

import (
	"google.golang.org/grpc"

	"github.com/codefuture-io/openpitrix/pkg/config"
	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/manager"
	"github.com/codefuture-io/openpitrix/pkg/pb"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	s := Server{}
	manager.NewGrpcServer("service-config", constants.ServiceConfigPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterServiceConfigServer(server, &s)
		})
}
