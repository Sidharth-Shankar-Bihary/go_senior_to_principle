package apierrors

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

type ErrorInfo struct {
	Status   int                    `json:"status"`
	Code     int                    `json:"ret_code"`
	GRPCCode codes.Code             `json:"grpc_code"`
	Msg      string                 `json:"ret_msg"`
	Details  map[string]interface{} `json:"details"`
}

func (e *ErrorInfo) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.Code))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Msg)
	return b.String()
}

// SetDetails set details as you wish =)
func (e *ErrorInfo) SetDetails(details map[string]interface{}) {
	e.Details = details
}

// GetStatus get status from ErrorInfo, if type is not same then status = 500
func GetStatus(err error) int {
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*ErrorInfo)
	if !ok {
		return http.StatusInternalServerError
	}
	return _err.Status
}

// CheckErrorByJSON input json byte convert to ErrorInfo
func CheckErrorByJSON(b []byte) (bool, error) {
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

func SetCustomizeErr(strErr *StrErrFace) (err error) {
	customizeErrMap := map[StrErrIdx]ErrorInfo{}
	for idx, newInfo := range customizeErrMap {
		err = strErr.GetErrorsByIndex(idx)
		if err != nil {
			return strErr.WrapFmt(ErrInternalError, "customize idx(%v) has exist in SErrFac default idx", idx)
		}
		strErr.RegisterNewErrByInfo(idx, newInfo)
	}
	return nil
}
