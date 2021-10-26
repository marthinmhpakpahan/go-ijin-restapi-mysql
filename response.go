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

type InsertionResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type UpdateResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type DeletionResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type DetailResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}