package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	marketv1 "github.com/CreatureDev/market/gen/go/market/v1"
	"github.com/CreatureDev/market/pkg/config"
	"github.com/CreatureDev/market/pkg/market"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	errors "golang.org/x/xerrors"
)

type Service struct {
	marketService marketv1.MarketServiceServer
}

func Run(conf config.Config) error {

	marketContext := context.Background()

	mux := runtime.NewServeMux()

	if err := marketv1.RegisterMarketServiceHandlerServer(marketContext, mux, market.NewService(conf)); err != nil {
		return errors.Errorf("failed to register market MarketService handler: %w", err)
	}

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 2)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case e := <-errChan:
		return e
	case s := <-stopChan:
		fmt.Println("Interrupt ", s)
		return nil
	}

}
