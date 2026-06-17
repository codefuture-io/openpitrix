// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package attachment

import (
	"google.golang.org/grpc"

	"github.com/codefuture-io/openpitrix/pkg/config"
	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/manager"
	"github.com/codefuture-io/openpitrix/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/pi"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s := Server{}
	manager.NewGrpcServer("attachment-manager", constants.AttachmentManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithMysqlConfig(cfg.Mysql).
		Serve(func(server *grpc.Server) {
			pb.RegisterAttachmentManagerServer(server, &s)
			pb.RegisterAttachmentServiceServer(server, &s)
		})
}
