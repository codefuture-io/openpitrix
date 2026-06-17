// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

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
	manager.NewGrpcServer("repo-manager", constants.RepoManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		WithBuilder(s.Builder).
		WithMysqlConfig(cfg.Mysql).
		Serve(func(server *grpc.Server) {
			pb.RegisterRepoManagerServer(server, &s)
		})
}
