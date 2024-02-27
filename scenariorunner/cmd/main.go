package main

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/aukilabs/go-tooling/pkg/cli"
	"github.com/aukilabs/go-tooling/pkg/logs"
	"github.com/aukilabs/hagall-common/scenariorunner/scenario"
	"github.com/segmentio/encoding/json"
)

var (
	// The version number. Set at build.
	version = "v0.0.0"
)

func main() {
	opts := scenario.Options{
		Hagall:    "http://host.docker.internal:4000",
		HDS:       "http://localhost:4002",
		LogLevel:  "info",
		AppKey:    "0x0",
		AppSecret: "0x0",
		SessionAttack: scenario.SessionAttackOptions{
			AttackCount:    1000,
			AttackDuration: time.Second * 5,
			AttackDefer:    time.Millisecond * 1,
		},
		Scenario: "session-attack",
	}
	cli.Register().
		Help("Launches a test scenario on a Hagall server.").
		Options(&opts)
	cli.Load()

	ctx, cancel := cli.ContextWithSignals(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	run(ctx, opts)
}

func run(ctx context.Context, opts scenario.Options) {
	logs.SetLevel(logs.ParseLevel(opts.LogLevel))
	logs.Encoder = json.Marshal
	if opts.LogIndent {
		logs.Encoder = func(v any) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		}
	}

	logs.WithTag("hds", opts.HDS).
		WithTag("hagall", opts.Hagall).
		WithTag("scenario", opts.Scenario).
		Info("start running a scenario on hagall server")

	scenario.Init()

	if opts.ListScenario {
		listScenarios()
		return
	}

	if err := scenario.Run(ctx, opts); err != nil {
		logs.WithTag("scenario", opts.Scenario).
			Fatal(err)
	}
}

func listScenarios() {
	fmt.Println("Supported Scenarios:")
	for name := range scenario.Scenarios() {
		fmt.Printf("  %s\n", name)
	}
}
