package web

type startRequest struct {
	PublicHost   string `json:"publicHost"`
	PublicPort   string `json:"publicPort"`
	UseGoProxyCN bool   `json:"useGoProxyCN"`
}

type startResponse struct {
	JobID string `json:"jobId"`
	State string `json:"state"`
}

type apiError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type indexTemplateData struct {
	InstallMode       string
	DefaultPublicHost string
	DefaultPublicPort string
}
