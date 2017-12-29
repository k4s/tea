package tea

import (
	"errors"
)

// error codes.
var (
	ErrBadAddr     = errors.New("invalid address")
	ErrBadHeader   = errors.New("invalid header received")
	ErrBadVersion  = errors.New("invalid protocol version")
	ErrTooShort    = errors.New("message is too short")
	ErrTooLong     = errors.New("message is too long")
	ErrClosed      = errors.New("connection closed")
	ErrConnRefused = errors.New("connection refused")
	ErrSendTimeout = errors.New("send time out")
	ErrRecvTimeout = errors.New("receive time out")
	ErrProtoState  = errors.New("incorrect protocol state")
	ErrProtoOp     = errors.New("invalid operation for protocol")
	ErrBadTran     = errors.New("invalid or unsupported transport")
	ErrBadProto    = errors.New("invalid or unsupported protocol")
	ErrPipeFull    = errors.New("pipe full")
	ErrPipeEmpty   = errors.New("pipe empty")
	ErrBadOption   = errors.New("invalid or unsupported option")
	ErrBadValue    = errors.New("invalid option value")
	ErrGarbled     = errors.New("message garbled")
	ErrAddrInUse   = errors.New("address in use")
	ErrBadProperty = errors.New("invalid property name")
	ErrTLSNoConfig = errors.New("missing TLS configuration")
	ErrTLSNoCert   = errors.New("missing TLS certificates")
)
