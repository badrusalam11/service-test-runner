package domain

// ProjectResponse represents an individual project's information.
type ProjectResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ShowProjectResponse is a slice of ProjectResponse.
type ShowProjectResponse []ProjectResponse
