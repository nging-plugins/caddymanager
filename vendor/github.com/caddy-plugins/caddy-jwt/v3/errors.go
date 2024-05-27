package jwt

import "errors"

var (
	errUnknownRuleType                = errors.New(`unknown rule type`)
	errNoTokenFound                   = errors.New(`no token found`)
	errTokenLengthIsZero              = errors.New(`token length is zero`)
	errEmptyPublicKeyFilename         = errors.New(`empty filename for public key provided`)
	errEmptySecretFilename            = errors.New(`empty filename for secret provided`)
	errNoKeybackend                   = errors.New(`there is no keybackend available`)
	errBothHMACRSATokensOnTheSameSite = errors.New(`cannot configure both HMAC and RSA/ECDSA tokens on the same site`)

	// Nested input must be a map or slice
	errNotValidInput = errors.New("Not a valid input: map or slice")
)
