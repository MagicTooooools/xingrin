package domain

type UpdateRequiredPayload struct {
	Version string `json:"version"`
	Image   string `json:"image"`
}
