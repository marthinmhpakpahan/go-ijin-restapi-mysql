package main

type Admin struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Dosen struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Foto string `json:"foto"`
	NIP string `json:"nip"`
	NamaLengkap string `json:"nama_lengkap"`
	JenisKelamin string `json:"jenis_kelamin"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type Mahasiswa struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Foto string `json:"foto"`
	NIM string `json:"nim"`
	NamaLengkap string `json:"nama_lengkap"`
	JenisKelamin string `json:"jenis_kelamin"`
	Kelas string `json:"kelas"`
	TahunMasuk string `json:"tahun_masuk"`
	Semester string `json:"semester"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type RequestType struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type Request struct {
	Id int64 `json:"id"`
	RequestTypeId string `json:"request_type_id"`
	RequestName string `json:"request_name"`
	Description string `json:"description"`
	DosenId string `json:"dosen_id"`
	DosenName string `json:"dosen_name"`
	MahasiswaId string `json:"mahasiswa_id"`
	MahasiswaName string `json:"mahasiswa_name"`
	StartDatetime string `json:"start_datetime"`
	EndDatetime string `json:"end_datetime"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type Counter struct {
	Total int64 `json:"total"`
}