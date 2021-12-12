package dto

// CreateRepositoryRequest ...
type CreateRepositoryRequest struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"url"`
}

// CreateRepositoryResponse ...
type CreateRepositoryResponse struct {
	Meta Meta `json:"meta"`
	Data struct {
		ID uint `json:"id"`
	} `json:"data"`
}

// UpdateRepositoryRequest ...
type UpdateRepositoryRequest struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"url"`
}

// UpdateRepositoryResponse ...
type UpdateRepositoryResponse struct {
	Meta Meta `json:"meta"`
}

// DeleteRepositoryResponse ...
type DeleteRepositoryResponse struct {
	Meta Meta `json:"meta"`
}

// GetOneRepositoryResponse ...
type GetOneRepositoryResponse struct {
	Meta Meta       `json:"meta"`
	Data Repository `json:"data"`
}

// Repository ...
type Repository struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
