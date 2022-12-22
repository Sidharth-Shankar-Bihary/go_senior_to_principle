package apierrors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	ErrBadRequest = iota
	ErrInvalidInput
	ErrInvalidQueryParameterValue
	ErrInvalidHeaderValue
	ErrMissingRequiredHeader
	ErrOutOfRangeInput
	ErrInternalDataNotSync
	ErrNotMatchSetting
	ErrUnauthorized
	ErrInvalidAuthenticationInfo
	ErrUsernameOrPasswordIncorrect
	ErrForbidden
	ErrAccountIsDisabled
	ErrAuthenticationFailed
	ErrNotAllowed
	ErrOtpExpired
	ErrInsufficientAccountPermissionsWithOperation
	ErrOtpRequired
	ErrOtpAuthorizationRequired
	ErrOtpIncorrect
	ErrResetPasswordRequired
	ErrResetOTPAccountNameIncorrect
	ErrSignIncorrect
	ErrResetOTPAccountEmailIncorrect
	ErrNotFound
	ErrResourceNotFound
	ErrAccountNotFound
	ErrPageNotFound
	ErrOrderNotFound
	ErrMethodNotAllowed
	ErrRequestTime
	ErrConflict
	ErrAccountAlreadyExists
	ErrAccountBeingCreated
	ErrResourceAlreadyExists
	ErrPhoneVerifiedTimeout
	ErrAwardHasBeenClaimed
	ErrPreconditionRequired
	ErrInternalServerError
	ErrInternalError
	ErrServerBusy
)

type StrErrIdx int

func defaultStrErrMap() map[StrErrIdx]ErrorInfo {
	return map[StrErrIdx]ErrorInfo{
		ErrBadRequest:                 {Code: 400000, Msg: http.StatusText(http.StatusBadRequest), Status: http.StatusBadRequest},
		ErrInvalidInput:               {Code: 400001, Msg: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument},
		ErrInvalidQueryParameterValue: {Code: 400009, Msg: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument},
		ErrInvalidHeaderValue:         {Code: 400004, Msg: "The value provided for one of the HTTP headers was not in the correct format.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument},
		ErrMissingRequiredHeader:      {Code: 400017, Msg: "A required HTTP header was not specified.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument},
		ErrOutOfRangeInput:            {Code: 400020, Msg: "something out of range", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument},
		ErrInternalDataNotSync:        {Code: 400041, Msg: "The specified data not sync in system, please wait a moment.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument},
		ErrNotMatchSetting:            {Code: 400087, Msg: "The specified data not match setting, please adjust your inputs.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument},

		ErrUnauthorized:                {Code: 401001, Msg: http.StatusText(http.StatusUnauthorized), Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated},
		ErrInvalidAuthenticationInfo:   {Code: 401001, Msg: "The authentication information was not provided in the correct format. Verify the value of Authorization header.", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated},
		ErrUsernameOrPasswordIncorrect: {Code: 401002, Msg: "Username or Password is incorrect.", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated},

		ErrForbidden:            {Code: 403000, Msg: http.StatusText(http.StatusForbidden), Status: http.StatusForbidden},
		ErrAccountIsDisabled:    {Code: 403001, Msg: "The specified account is disabled.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrAuthenticationFailed: {Code: 403002, Msg: "Server failed to authenticate the request. Make sure the value of the Authorization header is formed correctly including the signature.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrNotAllowed:           {Code: 403003, Msg: "The request is understood, but it has been refused or access is not allowed.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrOtpExpired:           {Code: 403004, Msg: "OTP is expired.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrInsufficientAccountPermissionsWithOperation: {Code: 403005, Msg: "The account being accessed does not have sufficient permissions to execute this operation.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrOtpRequired:                   {Code: 403007, Msg: "OTP Binding is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrOtpAuthorizationRequired:      {Code: 403008, Msg: "Two-factor authorization is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrOtpIncorrect:                  {Code: 403009, Msg: "OTP is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrResetPasswordRequired:         {Code: 403010, Msg: "Reset password is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrResetOTPAccountNameIncorrect:  {Code: 403011, Msg: "Reset otp failed,requested otp name is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrSignIncorrect:                 {Code: 403012, Msg: "verify sign failed", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},
		ErrResetOTPAccountEmailIncorrect: {Code: 403013, Msg: "Reset otp failed,requested otp email is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied},

		ErrNotFound:         {Code: 404000, Msg: http.StatusText(http.StatusNotFound), Status: http.StatusNotFound},
		ErrResourceNotFound: {Code: 404001, Msg: "The specified resource does not exist.", Status: http.StatusNotFound, GRPCCode: codes.NotFound},
		ErrAccountNotFound:  {Code: 404002, Msg: "cant find any account.", Status: http.StatusNotFound, GRPCCode: codes.NotFound},
		ErrPageNotFound:     {Code: 404003, Msg: "Page Not Found.", Status: http.StatusNotFound, GRPCCode: codes.NotFound},
		ErrOrderNotFound:    {Code: 404004, Msg: "The specified order not found", Status: http.StatusNotFound, GRPCCode: codes.NotFound},

		ErrMethodNotAllowed: {Code: 405001, Msg: "Server has received and recognized the request, but has rejected the specific HTTP method it’s using.", Status: http.StatusMethodNotAllowed, GRPCCode: codes.Unavailable},

		ErrRequestTime: {Code: 408001, Msg: "request time out", Status: http.StatusRequestTimeout, GRPCCode: codes.DeadlineExceeded},

		ErrConflict:              {Code: 409000, Msg: http.StatusText(http.StatusConflict), Status: http.StatusConflict, GRPCCode: codes.AlreadyExists},
		ErrAccountAlreadyExists:  {Code: 409001, Msg: "The specified account already exists.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists},
		ErrAccountBeingCreated:   {Code: 409002, Msg: "The specified account is in the process of being created.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists},
		ErrResourceAlreadyExists: {Code: 409004, Msg: "The specified resource already exists.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists},
		ErrPhoneVerifiedTimeout:  {Code: 409007, Msg: "sms verify time out", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists},
		ErrAwardHasBeenClaimed:   {Code: 409010, Msg: "award has been claimed", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists},

		ErrPreconditionRequired: {Code: 428001, Msg: "Where a client GETs a resource's state, modifies it, and PUTs it back to the server, when meanwhile a third party has modified the state on the server, leading to a conflict.", Status: http.StatusPreconditionRequired, GRPCCode: codes.FailedPrecondition},

		ErrInternalServerError: {Code: 500000, Msg: http.StatusText(http.StatusInternalServerError), Status: http.StatusInternalServerError, GRPCCode: codes.Internal},
		ErrInternalError:       {Code: 500001, Msg: "The server encountered an internal error. Please retry the request.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal},
		ErrServerBusy:          {Code: 500002, Msg: "Server is busy, please retry the request.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal},
	}
}
