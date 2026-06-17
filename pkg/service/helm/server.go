package helm

import (
	"github.com/codefuture-io/openpitrix/pkg/config"
	"github.com/codefuture-io/openpitrix/pkg/manager"
	runtimeprovider "github.com/codefuture-io/openpitrix/pkg/service/runtime_provider"

	"google.golang.org/grpc"

	"github.com/codefuture-io/openpitrix/pkg/constants"
	"github.com/codefuture-io/openpitrix/pkg/logger"
	"github.com/codefuture-io/openpitrix/pkg/pb"
	"github.com/codefuture-io/openpitrix/pkg/pi"
)

type Server struct {
	runtimeprovider.Server
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	err := pi.Global().RegisterRuntimeProvider(Provider, ProviderConfig)
	if err != nil {
		logger.Critical(nil, "failed to register provider config: %+v", err)
	}
	s := Server{}
	manager.NewGrpcServer("openpitrix-rp-kubernetes", constants.KubernetesProviderPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterRuntimeProviderManagerServer(server, &s)
		})
}
