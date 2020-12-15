package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/cmd"
	"github.com/avian-digital-forensics/timeline-investigator/cmd/main/server"
	"github.com/avian-digital-forensics/timeline-investigator/configs"
)

func main() {
	ctx := cmd.ContextWithSignal(context.Background())
	if err := run(ctx, os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v - %s\n", time.Now(), err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout io.Writer) error {
	// Get cli arguments
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	cfgPath := flags.String("cfg", "/configs/config.yml", "filepath for the config")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	cfg, err := configs.Get(*cfgPath)
	if err != nil {
		return err
	}

	// init the api-server
	apiServer := server.New(ctx)

	var wg sync.WaitGroup

	apiErrCh := make(chan error, 1)

	// Start the API server.
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := apiServer.Initialize(cfg.MainAPI); err != nil {
			apiErrCh <- fmt.Errorf("API server failed to initialize: %v", err)
		}

		if err := apiServer.Run(cfg.MainAPI); err != nil {
			apiErrCh <- fmt.Errorf("API server has stopped unexpectedly: %v", err)
		}
	}()

	// Wait for a stop event and shutdown server.
	select {
	case <-ctx.Done():
		if err := apiServer.Stop(); err != nil {
			return fmt.Errorf("context done - unable to stop the API server: %v", err)
		}

	case err := <-apiErrCh:
		if err := apiServer.Stop(); err != nil {
			log.Printf("unable to stop the API server: %v", err)
		}
		return err
	}

	wg.Wait()

	return nil
}
