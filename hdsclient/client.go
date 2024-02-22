package hdsclient

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/go-tooling/pkg/logs"
	"github.com/aukilabs/hagall-common/crypt"
	httpcmn "github.com/aukilabs/hagall-common/http"
	hsmoketest "github.com/aukilabs/hagall-common/smoketest"
)

type Encoder func(interface{}) ([]byte, error)
type Decoder func([]byte, interface{}) error

// Client represents an HTTP client to communicate with Hagall Discovery Service.
type Client struct {
	// The Hagall endpoint.
	HagallEndpoint string

	// The Hagall Discovery Service endpoint.
	HDSEndpoint string

	// The HTTP transport used to communicate with the Hagall Discovery Service.
	// (optional)
	Transport http.RoundTripper

	// The function to encode request bodies.
	Encode Encoder

	// The function to decode response bodies.
	Decode Decoder

	// Setting RemoteAddr to insert a http header x-real-ip=RemoteAddr to request,
	// allows override of client IP address which read from req.RemoteAddr by default
	RemoteAddr string

	mutex              sync.RWMutex
	serverID           string
	secret             string
	registrationState  string
	lastHealthCheck    time.Time
	clientID           string
	registrationStatus RegistrationStatus

	privateKey *ecdsa.PrivateKey
}

// NewClient creates new client with optional parameters.
func NewClient(opts ...ClientOpts) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}

	if c.Encode == nil {
		c.Encode = json.Marshal
	}

	if c.Decode == nil {
		c.Decode = json.Unmarshal
	}
	return c
}

// Secret returns the attributed JWT secret.
func (c *Client) Secret() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.secret
}

// ServerID returns the attributed server ID.
func (c *Client) ServerID() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.serverID
}

// SetServerData sets internal serverID & server secret state.
func (c *Client) SetServerData(serverID, secret string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.serverID = serverID
	c.secret = secret
}

// LastHealthCheck returns the time when HDS performed the latest healthcheck.
func (c *Client) LastHealthCheck() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.lastHealthCheck
}

// SetLastHealthCheck set the time when HDS performed the latest healthcheck.
func (c *Client) SetLastHealthCheck(t time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastHealthCheck = t
}

func (c *Client) setRegistrationState(v string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.registrationState = v
}

// GetRegistrationState returns current registration state.
func (c *Client) GetRegistrationState() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.registrationState
}

func (c *Client) setRegistrationStatus(v RegistrationStatus) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.registrationStatus = v
}

// GetRegistrationStatus returns current registration status.
func (c *Client) GetRegistrationStatus() RegistrationStatus {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.registrationStatus
}

// UserAuth auhenticates a user to a given server.
func (c *Client) UserAuth(ctx context.Context, in UserAuthIn) (UserAuthResponse, error) {
	body, err := c.Encode(in)
	if err != nil {
		return UserAuthResponse{}, errors.New("encoding body failed").Wrap(err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.HDSEndpoint+"/auth",
		bytes.NewReader(body),
	)

	if c.clientID != "" {
		req.Header.Set(httpcmn.HeaderPosemeshClientID, c.clientID)
	}

	if err != nil {
		return UserAuthResponse{}, errors.New("creating request failed").Wrap(err)
	}
	req.SetBasicAuth(in.AppKey, in.AppSecret)
	req.Header.Set("Content-Type", "application/json")

	var resp UserAuthResponse
	err = c.do(req, &resp)
	return resp, err
}

// VerifyUserAuth verifies a user's identity.
func (c *Client) VerifyUserAuth(token string) error {
	secret := c.Secret()
	if secret == "" {
		return errors.New("hagall server is not registered")
	}

	if err := httpcmn.VerifyHagallUserAccessToken(token, secret); err != nil {
		return errors.New("verifying access token failed").Wrap(err)
	}
	return nil
}

// PostServer registers a server to HDS.
func (c *Client) PostServer(ctx context.Context, in PostServerIn) error {
	if in.State == "" {
		in.State = httpcmn.MakeJWTSecret()
	}
	c.setRegistrationState(in.State)

	if in.Endpoint == "" {
		in.Endpoint = c.HagallEndpoint
	}

	return c.Post(ctx, "/servers", in)
}

// HandleServerRegistration handles Hagall server registration.
//
// When a request succeeds, the Hagall server secret is stored within the client
// and is available by calling the Secret method.
//
// Once the secret is transmitted, further calls to this handler result in a
// FORBIDDEN response.
//
// This handler is meant to be used by a Hagall server under the /registrations
// path.
func (c *Client) HandleServerRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpcmn.MethodNotAllowed(w)
		return
	}

	if c.Secret() != "" {
		httpcmn.Forbidden(w, errors.New("server registration failed").
			WithTag("reason", "server is already registered").
			WithTag("warning", "Did you or someone else just try to register the same address from another Hagall instance?"))
		return
	}

	expectedRegistrationState := c.GetRegistrationState()
	if expectedRegistrationState == "" {
		httpcmn.Forbidden(w, errors.New("non initiated server registration").
			WithTag("warning", "Did you or someone else just try to register the same address from another Hagall instance?"))
		return
	}

	if registrationState := r.Header.Get(httpcmn.HeaderHagallRegistrationStateKey); registrationState != expectedRegistrationState {
		httpcmn.Forbidden(w, errors.New("server registration failed").
			WithTag("current_registration_state", registrationState).
			WithTag("expected_registration_state", expectedRegistrationState).
			WithTag("warning", "Did you or someone else just try to register the same address from another Hagall instance?"))
		return
	}

	id := r.Header.Get(httpcmn.HeaderHagallIDKey)
	if id == "" {
		httpcmn.BadRequest(w, errors.New("server registration failed").
			WithTag("server_id", id))
		return
	}

	secret := r.Header.Get(httpcmn.HeaderHagallJWTSecretHeaderKey)
	if secret == "" {
		httpcmn.BadRequest(w, errors.New("server registration failed").
			WithTag("server_id", id).
			WithTag("server_secret", secret))
		return
	}
	c.SetServerData(id, secret)
	c.SetLastHealthCheck(time.Now())
	c.setRegistrationStatus(RegistrationStatusRegistered)

	httpcmn.OK(w)

	logs.WithTag("server_id", id).
		WithTag("status", c.GetRegistrationStatus()).
		Debug("server registration succeeded")
}

// HandleHealthCheck handles Hagall server health checks.
//
// This handler is meant to be used by a Hagall server under the /health path.
func (c *Client) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpcmn.MethodNotAllowed(w)
		return
	}

	token, err := httpcmn.SignIdentity(c.HagallEndpoint, c.Secret())
	if err != nil {
		httpcmn.InternalServerError(w, errors.New("signing response failed").Wrap(err))
		return
	}

	w.Header().Set("Authorization", httpcmn.MakeAuthorizationHeader(token))
	httpcmn.OK(w)

	c.SetLastHealthCheck(time.Now())
	logs.Debug("health check ok")
}

// GetServers returns a list of the closest servers.
func (c *Client) GetServers(ctx context.Context, in GetServersIn) (GetServersResponse, error) {
	queries := make(map[string]string)
	queries["min_version"] = in.MinVersion

	if len(in.Modules) != 0 {
		queries["modules"] = strings.Join(in.Modules, ",")
	}

	if len(in.FeatureFlags) != 0 {
		queries["feature_flags"] = strings.Join(in.FeatureFlags, ",")
	}

	var servers GetServersResponse
	err := c.GetWithAuth(ctx,
		"/servers",
		in.AppKey,
		in.AppSecret,
		queries,
		&servers,
	)
	return servers, err
}

// GetServerByID returns the server with the given id from HDS.
func (c *Client) GetServerByID(ctx context.Context, in GetServerByIDIn) (ServerResponse, error) {
	var s ServerResponse
	err := c.GetWithAuth(ctx, "/servers/"+in.ServerID, in.AppKey, in.AppSecret, nil, &s)
	return s, err
}

// GetServerByEndpoint returns the server with the given endpoint from HDS.
func (c *Client) GetServerByEndpoint(ctx context.Context, in GetServerByEndpointIn) (ServerResponse, error) {
	var s ServerResponse
	err := c.GetWithAuth(ctx, "/servers?endpoint="+in.Endpoint, in.AppKey, in.AppSecret, nil, &s)
	return s, err
}

// DeleteServer unregisters a server from HDS.
func (c *Client) DeleteServer(ctx context.Context) error {
	return c.Delete(ctx, "/servers")
}

// PostSession registers a session to HDS.
func (c *Client) PostSession(ctx context.Context, in PostSessionIn) error {
	return c.Post(ctx, "/sessions", in)
}

// Get sends a GET request to the given path and stores result in the given
// output.
func (c *Client) Get(ctx context.Context, path string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.HDSEndpoint+path,
		nil,
	)
	if err != nil {
		return errors.New("creating request failed").Wrap(err)
	}

	return c.do(req, out)
}

// GetWithAuth sends a GET request to the given path and stores result in
// the given output. It uses app key and app secret for basic authentication.
func (c *Client) GetWithAuth(
	ctx context.Context, path, appKey, appSecret string, queries map[string]string, out interface{},
) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.HDSEndpoint+path,
		nil,
	)

	if len(queries) > 0 {
		query := req.URL.Query()
		for k, v := range queries {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	if err != nil {
		return errors.New("creating request failed").Wrap(err)
	}
	req.SetBasicAuth(appKey, appSecret)

	if len(c.RemoteAddr) > 0 {
		req.Header.Set(httpcmn.XForwardedForHeaderKey, c.RemoteAddr)
	}

	if len(c.clientID) > 0 {
		req.Header.Set(httpcmn.HeaderPosemeshClientID, c.clientID)
	}

	return c.do(req, out)
}

// Post sends a POST request to the given path with the given input.
func (c *Client) Post(ctx context.Context, path string, in interface{}) error {
	body, err := c.Encode(in)
	if err != nil {
		return errors.New("encoding body failed").Wrap(err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.HDSEndpoint+path,
		bytes.NewReader(body),
	)
	if err != nil {
		return errors.New("creating request failed").Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if len(c.RemoteAddr) > 0 {
		req.Header.Set(httpcmn.XForwardedForHeaderKey, c.RemoteAddr)
	}

	return c.do(req, nil)
}

// Delete sends a DELETE request to the given path.
func (c *Client) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		c.HDSEndpoint+path,
		nil,
	)
	if err != nil {
		return errors.New("creating request failed").Wrap(err)
	}

	return c.do(req, nil)
}

// SendSmokeTestResult sends smoke test results to HDS /smoke-test-results endpoint.
func (c *Client) SendSmokeTestResult(ctx context.Context, results hsmoketest.SmokeTestResults) error {
	return c.Post(ctx, "/smoke-test-results", results)
}

func (c *Client) do(req *http.Request, out interface{}) error {
	secret := c.Secret()
	authorization := req.Header.Get("Authorization")
	if authorization == "" && secret != "" {
		token, err := httpcmn.SignIdentity(c.HagallEndpoint, secret)
		if err != nil {
			return errors.New("signing request failed").Wrap(err)
		}
		req.Header.Set("Authorization", httpcmn.MakeAuthorizationHeader(token))
	}

	res, err := c.Transport.RoundTrip(req)
	if err != nil {
		return errors.New("request failed").Wrap(err)
	}
	defer res.Body.Close()

	var body []byte
	if res.Body != nil {
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return errors.New("reading response failed").Wrap(err)
		}
	}

	if res.StatusCode >= 400 {
		return errors.New("request failed").
			WithTag("status", res.Status).
			WithTag("status_code", res.StatusCode).
			WithTag("message", string(body))
	}

	if out == nil {
		return nil
	}

	if err = c.Decode(body, out); err != nil {
		return errors.New("decoding response failed").Wrap(err)
	}
	return nil
}

// Pair pairs the client with HDS, registering an endpoint when necessary.
func (c *Client) Pair(ctx context.Context, in PairIn) error {
	// run register right after function starts
	registrationTimer := time.NewTimer(1 * time.Millisecond)
	defer registrationTimer.Stop()

	var registrationCount int
	retried := 1

	retryInterval := in.RegistrationInterval
	register := func() error {
		defer func() {
			// schedule next registration check
			status := c.GetRegistrationStatus()
			switch status {
			case RegistrationStatusPendingVerification:
				fallthrough
			case RegistrationStatusRegistered:
				retried = 0
				// wait for hds registration callback back, in heathcheck TTL
				retryInterval = in.HealthCheckTTL
			default:
				retryInterval = time.Duration(retried * int(in.RegistrationInterval))
				retried++
			}
			logs.WithTag("retry_interval", retryInterval.Seconds()).
				WithTag("retry_count", retried).
				WithTag("status", status).
				Debug("checking registration")
			registrationTimer.Reset(retryInterval)
		}()

		if time.Since(c.lastHealthCheck) <= in.HealthCheckTTL {
			return nil
		}

		c.setRegistrationStatus(RegistrationStatusRegistering)
		c.SetServerData("", "")
		registrationCount++

		var endpointSignature, timestamp string

		if c.privateKey != nil {
			sig, ts, err := crypt.SignWithTimestamp(c.privateKey, in.Endpoint)
			if err != nil {
				return errors.New("error signing endpoint").Wrap(err)
			}
			endpointSignature = sig
			timestamp = ts
		}

		logs.WithTag("registration_count", registrationCount).
			WithTag("endpoint", in.Endpoint).
			WithTag("version", in.Version).
			WithTag("modules", in.Modules).
			WithTag("feature_flags", in.FeatureFlags).
			WithTag("endpoint_signature", endpointSignature).
			WithTag("timestamp", timestamp).
			Debug("registering hagall to hds")

		if err := c.PostServer(ctx, PostServerIn{
			Endpoint:          in.Endpoint,
			Version:           in.Version,
			Modules:           in.Modules,
			FeatureFlags:      in.FeatureFlags,
			EndpointSignature: endpointSignature,
			Timestamp:         timestamp,
		}); err != nil {
			if retried >= in.RegistrationRetries {
				c.setRegistrationStatus(RegistrationStatusFailed)
				return errors.New("registering hagall to hds failed").
					WithTag("registration_count", registrationCount).
					WithTag("endpoint", in.Endpoint).
					WithTag("version", in.Version).
					WithTag("modules", in.Modules).
					WithTag("feature_flags", in.FeatureFlags).
					WithTag("retry_count", retried).
					Wrap(err)
			}

			logs.Error(err)
			return nil
		}

		logs.WithTag("registration_count", registrationCount).
			WithTag("endpoint", in.Endpoint).
			WithTag("version", in.Version).
			WithTag("modules", in.Modules).
			WithTag("feature_flags", in.FeatureFlags).
			WithTag("retry_count", retried).
			Info("hagall is successfully registered to hds")

		if c.GetRegistrationStatus() != RegistrationStatusRegistered {
			c.setRegistrationStatus(RegistrationStatusPendingVerification)
		}

		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-registrationTimer.C:
			if err := register(); err != nil {
				return err
			}
		}
	}
}

// Unpair deregisters Hagall from HDS.
func (c *Client) Unpair() error {
	logs.WithTag("status", c.GetRegistrationStatus()).
		Debug("unpairing server")
	if c.GetRegistrationStatus() != RegistrationStatusRegistered {
		return nil
	}

	dialer := &net.Dialer{
		Timeout: 3 * time.Second,
	}

	// update transport with shorter timeout
	transport := http.DefaultTransport.(*http.Transport)
	transport.DialContext = dialer.DialContext

	c.Transport = transport

	// create a new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.DeleteServer(ctx); err != nil {
		return errors.New("delete server failed").Wrap(err)
	}

	c.setRegistrationStatus(RegistrationStatusInit)
	logs.Info("unpair succeed")

	return nil
}

// isRetryableError returns false for error with status code in between 400 and 499.
func isRetryableError(err error) bool {
	rErr, ok := err.(errors.Error)
	if !ok {
		return true
	}

	if rErr.Tags()["status_code"] == "" {
		return true
	}

	statusCodeStr := rErr.Tags()["status_code"]
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		logs.Error(errors.New("failed to convert status code").WithTag("status_code", statusCodeStr))
		return true
	}

	return !(statusCode >= 400 && statusCode < 500)
}
