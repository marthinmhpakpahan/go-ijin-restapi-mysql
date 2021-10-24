package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/user/admin/login", adminLogin).Methods(http.MethodPost)

	// ============================ DOSEN ============================ //
	router.HandleFunc("/user/dosen/login", dosenLogin).Methods(http.MethodPost)
	router.HandleFunc("/user/dosen/create", dosenCreate).Methods(http.MethodPost)
	router.HandleFunc("/user/dosen/update", dosenUpdate).Methods(http.MethodPost)
	router.HandleFunc("/user/dosen/{id}", dosenShow).Methods(http.MethodGet)
	router.HandleFunc("/user/dosen/disable/{id}", dosenDisable).Methods(http.MethodGet)
	router.HandleFunc("/user/dosen/enable/{id}", dosenEnable).Methods(http.MethodGet)
	router.HandleFunc("/user/dosen/", dosenIndex).Methods(http.MethodGet)
	// ============================ DOSEN ============================ //

	// ============================ MAHASISWA ============================ //
	router.HandleFunc("/user/mahasiswa/login", mahasiswaLogin).Methods(http.MethodPost)
	router.HandleFunc("/user/mahasiswa/create", mahasiswaCreate).Methods(http.MethodPost)
	router.HandleFunc("/user/mahasiswa/update", mahasiswaUpdate).Methods(http.MethodPost)
	router.HandleFunc("/user/mahasiswa/{id}", mahasiswaShow).Methods(http.MethodGet)
	router.HandleFunc("/user/mahasiswa/disable/{id}", mahasiswaDisable).Methods(http.MethodGet)
	router.HandleFunc("/user/mahasiswa/enable/{id}", mahasiswaEnable).Methods(http.MethodGet)
	router.HandleFunc("/user/mahasiswa/", mahasiswaIndex).Methods(http.MethodGet)
	// ============================ MAHASISWA ============================ //
}