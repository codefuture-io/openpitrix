// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"github.com/codefuture-io/openpitrix/pkg/manager"
	"github.com/codefuture-io/openpitrix/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/plugins"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateRuntimeRequest:
		return manager.NewChecker(ctx, r).
			Required("name", "provider", "zone", "runtime_credential_id").
			StringChosen("provider", plugins.GetAvailablePlugins()).
			Exec()
	case *pb.ModifyRuntimeRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_id").
			Exec()
	case *pb.DeleteRuntimesRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_id").
			Exec()
	case *pb.CreateRuntimeCredentialRequest:
		return manager.NewChecker(ctx, r).
			Required("name", "provider", "runtime_credential_content").
			StringChosen("provider", plugins.GetAvailablePlugins()).
			Exec()
	case *pb.ModifyRuntimeCredentialRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_credential_id").
			Exec()
	case *pb.DeleteRuntimeCredentialsRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_credential_id").
			Exec()
	case *pb.DescribeRuntimeProviderZonesRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_credential_id").
			Exec()
	}
	return nil
}
