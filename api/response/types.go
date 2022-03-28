package response

const (
	ContextTagRequestId = "request-id"
	ContextTagStartTime = "start-time"
)

type Result struct {
	Code    string
	Message string
}

var (
	INFO_500 = Result{"PI_7W_500", "Internal Server Error"}
	INFO_400 = Result{"PI_7W_400", "Bad Request"}
	INFO_404 = Result{"PI_7W_404", "Resource not found"}
	INFO_429 = Result{"PI_7W_429", "Too Many Requests"}
	INFO_202 = Result{"PI_7W_202", "Accepted"}
	INFO_200 = Result{"PI_7W_200", "OK"}
)

type Info struct {
	ResultCode    string   `json:"resultCode" example:"PI_7W_200 (200) | PI_7W_202 (202) | PI_7W_400 (400) | PI_7W_404 (404) | PI_7W_429 (429) | PI_7W_500 (500)"`
	ResultMessage string   `json:"resultMessage" example:"description of status code"`
	RequestId     string   `json:"requestId" example:"718527dd-9f10-4b42-bb74-3fc4e0ffea46"`
	Details       []Detail `json:"details,omitempty"`
}

type Detail struct {
	Code    string `json:"code" example:"type of error"`
	Message string `json:"message" example:"description of error"`
	Scope   string `json:"scope"  example:"GECT-NG"`
}

type InfoResponse struct {
	Info Info `json:"info"`
}
