// Copyright (c) ZStack.io, Inc.

package errors

const (
	ErrServer       = Error("ServerError")
	ErrClient       = Error("ClientError")
	ErrUnclassified = Error("UnclassifiedError")

	ErrDNS            = Error("DNSError")
	ErrEOF            = Error("EOFError")
	ErrNetwork        = Error("NetworkError")
	ErrConnectRefused = Error("ConnectRefusedError")
	ErrConnectReset   = Error("ConnectResetError")
	ErrTimeout        = Error("TimeoutError")

	ErrNotFound         = Error("NotFoundError")
	ErrDuplicateId      = Error("DuplicateIdError")
	ErrNotSupported     = Error("NotSupportedError")
	ErrAccountReadOnly  = Error("AccountReadOnlyError")
	ErrInvalidAccessKey = Error("InvalidAccessKey")

	ErrAggregate = Error("AggregateError")

	ErrParameter = Error("ParameterError")
)
