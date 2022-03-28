package response

import "github.com/gin-gonic/gin"

const GectNgScope = "GECT-NG"

func ProduceDefaultResponse(c *gin.Context, statusCode int, ambit string, err error) {
	switch statusCode {
	case 400:
		_ = c.Error(BadRequestError{Ambit: ambit, InnerErr: []error{err}})
	case 404:
		_ = c.Error(NotFoundError{Ambit: ambit, InnerErr: []error{err}})
	default:
		_ = c.Error(InternalServerError{Ambit: ambit, InnerErr: []error{err}})
	}
}

type InternalServerError struct {
	Ambit    string
	InnerErr []error
}

func (ise InternalServerError) Error() string {

	if len(ise.InnerErr) > 0 {
		return ise.InnerErr[0].Error()
	}
	return ""
}

func (ise InternalServerError) Info(reqId string, _ int64) InfoResponse {
	i := Info{
		ResultCode:    INFO_500.Code,
		ResultMessage: INFO_500.Message,
		RequestId:     reqId,
	}

	a := ise.Ambit
	if a == "" {
		a = "server-error"
	}

	dets := make([]Detail, 0)
	for _, e := range ise.InnerErr {

		dets = append(dets, Detail{Code: a, Message: e.Error(), Scope: GectNgScope})
	}

	i.Details = dets
	return InfoResponse{Info: i}
}

type BadRequestError struct {
	Ambit    string
	InnerErr []error
}

func (bre BadRequestError) Error() string {
	if len(bre.InnerErr) > 0 {
		return bre.InnerErr[0].Error()
	}
	return ""
}

func (bre BadRequestError) Info(reqId string, _ int64) InfoResponse {
	i := Info{
		ResultCode:    INFO_400.Code,
		ResultMessage: INFO_400.Message,
		RequestId:     reqId,
	}

	a := bre.Ambit
	if a == "" {
		a = "bad-request"
	}

	dets := make([]Detail, 0)
	for _, e := range bre.InnerErr {
		dets = append(dets, Detail{Code: a, Message: e.Error(), Scope: GectNgScope})
	}

	i.Details = dets
	return InfoResponse{Info: i}
}

type NotFoundError struct {
	Ambit    string
	InnerErr []error
}

func (nfe NotFoundError) Error() string {
	if len(nfe.InnerErr) > 0 {
		return nfe.InnerErr[0].Error()
	}
	return ""
}

func (nfe NotFoundError) Info(reqId string, _ int64) InfoResponse {
	i := Info{
		ResultCode:    INFO_404.Code,
		ResultMessage: INFO_404.Message,
		RequestId:     reqId,
	}

	a := nfe.Ambit
	if a == "" {
		a = "not-found"
	}

	dets := make([]Detail, 0)
	for _, e := range nfe.InnerErr {
		dets = append(dets, Detail{Code: a, Message: e.Error(), Scope: GectNgScope})
	}

	i.Details = dets
	return InfoResponse{Info: i}
}

type TooManyRequestsError struct {
	Ambit    string
	InnerErr []error
}

func (tre TooManyRequestsError) Error() string {
	if len(tre.InnerErr) > 0 {
		return tre.InnerErr[0].Error()
	}
	return ""
}

func (tre TooManyRequestsError) Info(reqId string, _ int64) InfoResponse {
	i := Info{
		ResultCode:    INFO_429.Code,
		ResultMessage: INFO_429.Message,
		RequestId:     reqId,
	}

	a := tre.Ambit
	if a == "" {
		a = "too-many-requests"
	}

	dets := make([]Detail, 0)
	for _, e := range tre.InnerErr {
		dets = append(dets, Detail{Code: a, Message: e.Error(), Scope: GectNgScope})
	}

	i.Details = dets
	return InfoResponse{Info: i}
}

/*
func InternalServerErrorInfo(reqId string, timeInfoMs int64, det Detail) Info {
	i := Info{
		ResultCode:    INFO_500.Code,
		ResultMessage: INFO_500.Message,
		ExecutionTime: timeInfoMs,
		RequestId:     reqId,
	}

	if det.Code != "" {
		i.Details = []Detail{det}
	}

	return i
}

func BadRequestInfo(reqId string, timeInfoMs int64, det Detail) Info {
	i := Info{
		ResultCode:    INFO_400.Code,
		ResultMessage: INFO_400.Message,
		ExecutionTime: timeInfoMs,
		RequestId:     reqId,
	}

	if det.Code != "" {
		i.Details = []Detail{det}
	}

	return i
}

func NotFoundInfo(reqId string, timeInfoMs int64, det Detail) Info {
	i := Info{
		ResultCode:    INFO_404.Code,
		ResultMessage: INFO_404.Message,
		ExecutionTime: timeInfoMs,
		RequestId:     reqId,
	}

	if det.Code != "" {
		i.Details = []Detail{det}
	}

	return i
}
*/

func SuccessInfo(reqId string, _ int64) InfoResponse {
	i := Info{
		ResultCode:    INFO_200.Code,
		ResultMessage: INFO_200.Message,
		RequestId:     reqId,
		Details:       nil,
	}

	return InfoResponse{Info: i}
}

func AcceptedInfo(reqId string, _ int64) InfoResponse {
	i := Info{
		ResultCode:    INFO_202.Code,
		ResultMessage: INFO_202.Message,
		RequestId:     reqId,
		Details:       nil,
	}

	return InfoResponse{Info: i}
}
