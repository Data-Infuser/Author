package enum

import grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"

type AuthCode int32

const (
	VALID AuthCode = 0
	INTERNAL_EXCEPTION = -1
	PARAMETER_EXCEPTION = -2
	UNREGISTERED_SERVICE = -3
	TERMINATED_SERVICE = -9
	LIMIT_EXCEEDED = -10
	UNAUTHORIZED = -401
)

func (c AuthCode) GetgRPCCode() grpc_author.ApiAuthRes_Code {
	switch c {
	case VALID:
		return grpc_author.ApiAuthRes_VALID

	case INTERNAL_EXCEPTION:
		return grpc_author.ApiAuthRes_INTERNAL_EXCEPTION

	case PARAMETER_EXCEPTION:
		return grpc_author.ApiAuthRes_PARAMETER_EXCEPTION

	case UNREGISTERED_SERVICE:
		return grpc_author.ApiAuthRes_UNREGISTERED_SERVICE

	case TERMINATED_SERVICE:
		return grpc_author.ApiAuthRes_TERMINATED_SERVICE

	case LIMIT_EXCEEDED:
		return grpc_author.ApiAuthRes_LIMIT_EXCEEDED

	case UNAUTHORIZED:
		return grpc_author.ApiAuthRes_UNAUTHORIZED

	default:
		return grpc_author.ApiAuthRes_UNKNOWN
	}
}