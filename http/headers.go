package http

const (
	HeaderHagallIDKey                         = "Hagall-Id"
	HeaderHagallJWTSecretHeaderKey            = "Hagall-Jwt-Secret"
	HeaderHagallJWTChallengeHeaderKey         = "Hagall-Jwt-Challenge"
	HeaderHagallJWTChallengeSolutionHeaderKey = "Hagall-Jwt-Challenge-Solution"
	HeaderHagallRegistrationStateKey          = "Hagall-Registration-State"
	HeaderPosemeshClientID                    = "posemesh-client-id"

	CloudFrontTimezoneNameHeaderKey  = "CloudFront-Viewer-Time-Zone"
	CloudFrontCountryNameHeaderKey   = "CloudFront-Viewer-Country"
	CloudFrontViewerAddressHeaderKey = "CloudFront-Viewer-Address"

	XForwardedForHeaderKey = "X-Forwarded-For"
)

type ClientIDContextValue string
