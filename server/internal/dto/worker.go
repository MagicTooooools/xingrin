package dto

// WorkerTargetNameResponse is the response for GetTargetName
type WorkerTargetNameResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// WorkerProviderConfigResponse is the response for GetProviderConfig
type WorkerProviderConfigResponse struct {
	Content string `json:"content"`
}
