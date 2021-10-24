package main

type LoginResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Role string `json:"role"`
	Data interface{} `json:"data"`
}

type IndexResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type CommonInsertionResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type CommonUpdateResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type CommonDeletionResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type CommonDetailResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}