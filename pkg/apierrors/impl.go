package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StrErrFace struct {
	m map[StrErrIdx]ErrorInfo
}

func NewStrErrFace() *StrErrFace {
	return &StrErrFace{m: defaultStrErrMap()}
}

func (f *StrErrFace) RegisterNewErr(idx StrErrIdx, code int, msg string, stat int, grpcCode codes.Code) {
	f.m[idx] = ErrorInfo{Code: code, Msg: msg, Status: stat, GRPCCode: grpcCode}
}

func (f *StrErrFace) RegisterNewErrByInfo(idx StrErrIdx, info ErrorInfo) {
	f.m[idx] = info
}

// CompareErrorCode compare error code is the same or not.
func (f *StrErrFace) CompareErrorCode(targetErrIndex StrErrIdx, errB error) bool {
	if bErr, exists := errors.Cause(errB).(*ErrorInfo); exists {
		target, ok := f.m[targetErrIndex]
		if ok {
			if target.Code == bErr.Code {
				return true
			}
		}
	}
	return false
}

// GetIndexFromError get StrErrIdx from error, index is -1 if not found in map
func (f *StrErrFace) GetIndexFromError(err error) (index StrErrIdx, causeOK bool) {
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		return 0, false
	}
	for k, v := range f.m {
		if v.Code == _err.Code {
			return k, true
		}
	}
	return -1, true
}

// WithErrors is used the customise errors code, if the code is not defined, show the http status info.
func (f *StrErrFace) WithErrors(err error) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		return errors.WithStack(&ErrorInfo{
			Status: f.m[ErrInternalError].Status,
			Code:   f.m[ErrInternalError].Code,
			Msg:    http.StatusText(f.m[ErrInternalError].Status),
		})
	}
	return errors.WithStack(&ErrorInfo{
		Status: _err.Status,
		Code:   _err.Code,
		Msg:    _err.Msg,
	})
}

// ConvertProtoErr Convert ErrorInfo to grpc error
func (f *StrErrFace) ConvertProtoErr(err error) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		return status.Error(f.m[ErrInternalError].GRPCCode, err.Error())
	}
	b, _ := json.Marshal(_err)
	return status.Error(_err.GRPCCode, string(b))
}

// ConvertHTTPErr Convert  grpc error to ErrorInfo
func (f *StrErrFace) ConvertHTTPErr(err error) error {
	if err == nil {
		return nil
	}

	// check err is http error
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if ok {
		return _err
	}

	s := status.Convert(err)
	if s == nil {
		return f.GetErrorsByIndex(ErrInternalError)
	}
	return f.switchCode(s)
}

func (f *StrErrFace) switchCode(s *status.Status) error {
	httpErr := f.GetErrorsByIndex(ErrInternalError)
	switch s.Code() {
	case codes.Unknown:
		httpErr = f.GetErrorsByIndex(ErrInternalError)
	case codes.InvalidArgument:
		httpErr = f.GetErrorsByIndex(ErrInvalidInput)
	case codes.NotFound:
		httpErr = f.GetErrorsByIndex(ErrResourceNotFound)
	case codes.AlreadyExists:
		httpErr = f.GetErrorsByIndex(ErrResourceAlreadyExists)
	case codes.PermissionDenied:
		httpErr = f.GetErrorsByIndex(ErrNotAllowed)
	case codes.Unauthenticated:
		httpErr = f.GetErrorsByIndex(ErrUnauthorized)
	case codes.OutOfRange:
		httpErr = f.GetErrorsByIndex(ErrInvalidInput)
	case codes.Internal:
		httpErr = f.GetErrorsByIndex(ErrInternalError)
	case codes.DataLoss:
		httpErr = f.GetErrorsByIndex(ErrInternalError)
	case codes.Unimplemented:
		httpErr = f.GetErrorsByIndex(ErrResourceNotFound)
	case codes.Unavailable:
		httpErr = f.GetErrorsByIndex(ErrServerBusy)
	}
	return httpErr
}

func (f *StrErrFace) ConvertToJSON(err error) []byte {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		_err = f.GetErrorsByIndex(ErrInternalError).(*ErrorInfo)
	}
	b, err := json.Marshal(_err)
	if err != nil {
		return nil
	}
	return b
}

// NewWithMsg is used to replace the error info.
func (f *StrErrFace) NewWithMsg(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		return errors.WithStack(&ErrorInfo{
			Status:   f.m[ErrInternalError].Status,
			Code:     f.m[ErrInternalError].Code,
			Msg:      f.m[ErrInternalError].Msg,
			GRPCCode: f.m[ErrInternalError].GRPCCode,
		})
	}
	err = &ErrorInfo{
		Status:   _err.Status,
		Code:     _err.Code,
		Msg:      message,
		GRPCCode: _err.GRPCCode,
	}
	var msg string
	for i := 0; i < len(args); i++ {
		msg += "%+v"
	}
	return errors.Wrapf(err, msg, args...)
}

func (f *StrErrFace) NewWithMsgFmt(errIndex StrErrIdx, format string, args ...interface{}) error {
	return f.NewWithMsg(f.GetErrorsByIndex(errIndex), fmt.Sprintf(format, args...))
}

// GetCodeWithErrors is used to get code & msg
func (f *StrErrFace) GetCodeWithErrors(err error) (code uint32, msg string, typeError error) {
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		return uint32(f.m[ErrInternalError].Code), f.m[ErrInternalError].Msg, errors.New("unknown error type")
	}
	return uint32(_err.Code), _err.Msg, nil
}

// GetStatus get status from ErrorInfo, if type is not same then status = 500
func (f *StrErrFace) GetStatus(err error) int {
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		return http.StatusInternalServerError
	}
	return _err.Status
}

// GetErrorsByIndex get ErrorInfo by index, if index not found will return nil
func (f *StrErrFace) GetErrorsByIndex(index StrErrIdx) error {
	err, ok := f.m[index]
	if !ok {
		return nil
	}
	return &err
}

func (f *StrErrFace) GetErrorsByCode(code int) error {
	for _, v := range f.m {
		if code == v.Code {
			err := v
			return &err
		}
	}
	return nil
}

// CheckErrorByJSON input json byte convert to ErrorInfo
func (f *StrErrFace) CheckErrorByJSON(b []byte) (bool, error) {
	var causeErr ErrorInfo
	err := json.Unmarshal(b, &causeErr)
	if err != nil {
		return false, nil
	}
	if causeErr.Code != 0 && causeErr.Msg != "" {
		return true, &causeErr
	}
	return false, nil
}

// SetRetInfo set ret_code & ret_message
func (f *StrErrFace) SetRetInfo(retCode *int32, retMsg *string, respErr error) {
	if respErr != nil {
		code, msg, typeError := f.GetCodeWithErrors(respErr)
		if typeError != nil {
			code = uint32(f.m[ErrInternalError].Code)
			msg = respErr.Error()
		}
		*retCode = int32(code)
		*retMsg = msg
	}
}

// Wrap is as the proxy for github.com/pkg/errors.Wrap func.
func (f *StrErrFace) Wrap(errIndex StrErrIdx, message string) error {
	return errors.Wrap(f.GetErrorsByIndex(errIndex), message)
}

// WrapFmt is as the proxy for github.com/pkg/errors.Wrapf func.
func (f *StrErrFace) WrapFmt(errIndex StrErrIdx, format string, args ...interface{}) error {
	return errors.Wrapf(f.GetErrorsByIndex(errIndex), format, args...)
}

// WithMessage is as the proxy for github.com/pkg/errors.WithMessage func.
func (f *StrErrFace) WithMessage(errIndex StrErrIdx, message string) error {
	return errors.WithMessage(f.GetErrorsByIndex(errIndex), message)
}

// WithMsgFmt is as the proxy for github.com/pkg/errors.WithMessagef func.
func (f *StrErrFace) WithMsgFmt(errIndex StrErrIdx, format string, args ...interface{}) error {
	return errors.WithMessagef(f.GetErrorsByIndex(errIndex), format, args...)
}

// WithStack is as the proxy for github.com/pkg/errors.WithStack func.
func (f *StrErrFace) WithStack(errIndex StrErrIdx) error {
	return errors.WithStack(f.GetErrorsByIndex(errIndex))
}

// Is reports whether any error in error's chain matches target.
// The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap.
// An error is considered to match a target if it is equal to that target or if it implements a method Is(error) bool
// such that Is(target) returns true.
func (f *StrErrFace) Is(err, target error) bool {
	return errors.Is(err, target)
}
