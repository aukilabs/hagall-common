package models

// PostServerIn represents the input to register a server to HDS.
type PostServerIn struct {
	// The Hagall server endpoint. Optional.
	Endpoint string `json:"endpoint"`

	// The endpoint url signature signed by the wallet
	EndpointSignature string `json:"endpoint_signature"`

	// The Hagall server version.
	Version string `json:"version"`

	// A string that represents a registration state that is used during Hagall
	// server verification to ensure that the registration request originates
	// from the server that started the registration process.
	State string `json:"state"`

	// The modules that the server supports.
	Modules []string `json:"modules"`

	// Feature flags supported by server
	FeatureFlags []string `json:"feature_flags"`

	// The timestamp when endpoint signature was signed
	Timestamp string `json:"timestamp"`
}
