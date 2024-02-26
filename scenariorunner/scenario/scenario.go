package scenario

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aukilabs/go-tooling/pkg/errors"
	hds "github.com/aukilabs/hagall-common/hdsclient"
	httpcmn "github.com/aukilabs/hagall-common/http"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

// Options provides scenario runner configuration
type Options struct {
	Hagall        string               `env:"SCENARIO_HAGALL_ADDR"     help:"The Hagall server address to attack."`
	HDS           string               `env:"SCENARIO_HDS_ADDR"        help:"The HDS server address to retrieve Hagall servers."`
	AppKey        string               `env:"SCENARIO_APP_KEY"         help:"The app key to auth with HDS."`
	AppSecret     string               `env:"SCENARIO_APP_SECRET"      help:"The app secret to auth with HDS."`
	LogLevel      string               `env:"SCENARIO_LOG_LEVEL"       help:"Log level (debug|info|warning|error)."`
	LogIndent     bool                 `env:"SCENARIO_LOG_INDENT"      help:"Indents logs."`
	Scenario      string               `env:"SCENARIO_NAME"            help:"Scenario to run, see --list-scenario for supported scenarios"`
	ListScenario  bool                 `env:"-"                        help:"List Scenarios"`
	SessionAttack SessionAttackOptions `env:"SCENARIO_SESSION_ATTACK"  help:"Session Attack scenario configuration"`
	Help          bool                 `env:"-"                    help:"Show help."`
	Version       bool                 `env:"-"                    help:"Show version."`
}

// Scenario is a interface to run actual scenario
type Scenario interface {
	Run(context.Context, Options) error
}

var (
	allScenarios map[string]Scenario
)

// Init initialize scenarios
func Init() {
	allScenarios = make(map[string]Scenario)
	allScenarios["session-attack"] = sessionAttack{}
	allScenarios["integration-test"] = integrationTest{}
}

// Scenarios return supported scenarios
func Scenarios() map[string]Scenario {
	return allScenarios
}

// Run start running scenario with configured opts
func Run(ctx context.Context, opts Options) error {
	scenario, ok := allScenarios[opts.Scenario]
	if !ok {
		return errors.New("unsupported scenario").
			WithTag("scenario", opts.Scenario)
	}

	if err := scenario.Run(ctx, opts); err != nil {
		return errors.New("failed to run scenario").
			WithTag("scenario", opts.Scenario).
			Wrap(err)
	}

	return nil
}

// connectToHagall is a supporting function to connect to hagall server using hds client with authentication
func connectToHagall(ctx context.Context, opts Options) (*websocket.Conn, error) {
	clientID := uuid.NewString()
	hdsClient := hds.NewClient(
		hds.WithHDSEndpoint(opts.HDS),
		hds.WithEncoder(json.Marshal),
		hds.WithDecoder(json.Unmarshal),
		hds.WithTransport(http.DefaultTransport),
		hds.WithClientID(clientID),
	)

	server, err := hdsClient.GetServerByEndpoint(ctx, hds.GetServerByEndpointIn{
		Endpoint:  opts.Hagall,
		AppKey:    opts.AppKey,
		AppSecret: opts.AppSecret,
	})
	if err != nil {
		return nil, errors.New("getting server info failed").
			WithTag("hds", opts.HDS).
			WithTag("hagall", opts.Hagall).
			Wrap(err)
	}

	wsEndpoint := server.Endpoint
	wsEndpoint = strings.ReplaceAll(wsEndpoint, "host.docker.internal", "localhost")
	wsEndpoint = strings.ReplaceAll(wsEndpoint, "https://", "wss://")
	wsEndpoint = strings.ReplaceAll(wsEndpoint, "http://", "ws://")

	cfg, err := websocket.NewConfig(wsEndpoint, "https://localhost")
	if err != nil {
		return nil, errors.New("creating websocket config failed").
			WithTag("endpoint", wsEndpoint).
			Wrap(err)
	}
	cfg.Header.Set("Authorization", "Bearer "+server.AccessToken)
	cfg.Header.Set("User-Agent", "HDS (Go WebSocket Client golang.org/x/net/websocket)")
	cfg.Header.Set(httpcmn.HeaderPosemeshClientID, clientID)

	conn, err := websocket.DialConfig(cfg)
	if err != nil {
		return nil, errors.New("dialing to websocket failed").
			WithTag("endpoint", wsEndpoint).
			Wrap(err)
	}
	return conn, nil
}
