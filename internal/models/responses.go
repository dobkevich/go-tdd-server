package models

type AddRequest struct {
	A int `query:"a" validate:"required,numeric"`
	B int `query:"b" validate:"required,numeric"`
}

type AddResponse struct {
	Result int `json:"result"`
}

type EchoRequest struct {
	Message string `json:"message" validate:"required,min=1,max=100"`
}

type EchoResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type TimeResponse struct {
	Time string `json:"time"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp,omitempty"`
	Uptime    string `json:"uptime,omitempty"`
}
