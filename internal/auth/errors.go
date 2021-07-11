package auth

import "errors"

var (
	errInvalidRequest     = errors.New("invalid request")
	errMissingToken       = errors.New("missing token")
	errInvalidFormatToken = errors.New("token format is not valid")
	errInvalidToken       = errors.New("invalid token")
	errInvalidClaims      = errors.New("invalid claims")
	errNilClaims          = errors.New("nil claims")

	errUnexpectedSigningMethod = errors.New("unexpected signing method")
	errMissingKeyID            = errors.New("missing key id from jwt token")
	errInvalidKeyID            = errors.New("key id is not a string")
	errUnknownKeyID            = errors.New("unknown key id")

	errInternalError = errors.New("internal error")
)
