package appresult

type AppSuccess struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewAppSuccess(message string, data interface{}) *AppSuccess {
	return &AppSuccess{
		Message: message,
		Data:    data,
	}
}
