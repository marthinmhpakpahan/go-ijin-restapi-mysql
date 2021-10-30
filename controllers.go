package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
)

var db, err = getDB()

// ================================ COMMON ================================ //
func isDataExists(fieldName string, value string, tableName string) bool {
	rows, err := db.Query("SELECT * FROM " +tableName+ " WHERE " +fieldName+ " = ? ", value)
	if err != nil {
		return false
	}

	for rows.Next() {
		return true
	}
	return false
}
// ================================ COMMON ================================ //


// ================================ ADMIN ================================ //
func checkAdmin(_username string, _password string) (Admin, error) {
	fmt.Println("checkAdmin Hit!")

	var admin Admin
	row := db.QueryRow("SELECT id, username from admin WHERE username = ? AND password = ? AND status = 'active' ", _username, _password)
	err = row.Scan(&admin.Id, &admin.Username)
	if err != nil {
		return admin, err
	}
	return admin, nil
}

func adminLogin(w http.ResponseWriter, r *http.Request) {
	var response LoginResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Role = "admin"
	response.Data = Admin{}

	username := r.FormValue("username")
	password := r.FormValue("password")

	admin, err := checkAdmin(username, password)
	if err == nil {
		response.Error = false
		response.Message = "Akun ditemukan"
		response.Data = admin
	} else {
		response.Message = "Akun tidak ditemukan"
	}

	respondWithSuccess(response, w)
}
// ================================ ADMIN ================================ //


// ================================ DOSEN ================================ //
func checkDosen(_username string, _password string) (Dosen, error) {
	fmt.Println("checkDosen Hit!")
	fmt.Println(_username, _password)

	var dosen Dosen
	row := db.QueryRow("SELECT * from dosen WHERE username = ? AND password = ? AND status = 'active' ", _username, _password)
	err = row.Scan(&dosen.Id, &dosen.Username, &dosen.Password, &dosen.Foto, &dosen.NIP, &dosen.NamaLengkap, &dosen.JenisKelamin, &dosen.Status, &dosen.CreatedAt, &dosen.ModifiedAt)
	if err != nil {
		fmt.Println(err)
		return dosen, err
	}
	return dosen, nil
}

func dosenIndex(w http.ResponseWriter, r *http.Request) {
	param_limit := r.URL.Query()["limit"]
	limit := ""
	if len(param_limit) > 0 {
		limit = " LIMIT " + param_limit[0]
	}
	fmt.Println("Limit: ", limit)

	param_status := r.URL.Query()["status"]
	status := ""
	if len(param_status) > 0 {
		status = " AND status = '" + param_status[0] +"'"
	}
	fmt.Println("status: ", status)

	var response IndexResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = make([]string, 0)
	var dosens = []Dosen{}

	rows, err := db.Query("SELECT id, username, foto, nip, nama_lengkap, jenis_kelamin, status, created_at FROM dosen WHERE 1=1 "+status+" ORDER BY id DESC " + limit)
	if err != nil {
		response.Message = "(0) Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
		return
	}

	for rows.Next() {
		var dosen Dosen
		err = rows.Scan(&dosen.Id, &dosen.Username, &dosen.Foto, &dosen.NIP, &dosen.NamaLengkap, &dosen.JenisKelamin, &dosen.Status, &dosen.CreatedAt)
		if err != nil {
			response.Message = "(1) Terjadi kesalahan pada database"
			fmt.Println(err)
			respondWithSuccess(response, w)
			return
		}
		dosens = append(dosens, dosen)
	}

	response.Error = true
	response.Message = "Data ditemukan"
	response.Data = dosens
	respondWithSuccess(response, w)
}

func dosenLogin(w http.ResponseWriter, r *http.Request) {
	var response LoginResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Role = "dosen"
	response.Data = Dosen{}

	username := r.FormValue("username")
	password := r.FormValue("password")

	dosen, err := checkDosen(username, generateMD5(password))
	if err == nil {
		response.Error = false
		response.Message = "Akun ditemukan"
		response.Data = dosen
	} else {
		response.Message = "Akun tidak ditemukan"
	}

	respondWithSuccess(response, w)
}

func dosenDisable(w http.ResponseWriter, r *http.Request) {
	var response DeletionResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Dosen tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	_, err = db.Exec("UPDATE dosen SET status = 'deleted' WHERE id = ?", id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Dosen berhasil di Nonaktifkan"
		respondWithSuccess(response, w)
	}
}

func dosenEnable(w http.ResponseWriter, r *http.Request) {
	var response DeletionResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Dosen tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	_, err = db.Exec("UPDATE dosen SET status = 'active' WHERE id = ?", id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Dosen berhasil di Aktifkan"
		respondWithSuccess(response, w)
	}
}

func dosenShow(w http.ResponseWriter, r *http.Request) {
	var dosen Dosen
	var response DetailResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = Dosen{}

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Dosen tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	row := db.QueryRow("SELECT id, username, foto, nip, nama_lengkap, jenis_kelamin, description, status, created_at, modified_at FROM dosen WHERE id = ?", id)
	err = row.Scan(&dosen.Id, &dosen.Username, &dosen.Foto, &dosen.NIP, &dosen.NamaLengkap, &dosen.JenisKelamin, &dosen.Description, &dosen.Status, &dosen.CreatedAt, &dosen.ModifiedAt)
	if err == nil {
		response.Error = false
		response.Message = "Data ditemukan"
		response.Data = dosen
		respondWithSuccess(response, w)
	} else {
		fmt.Println(err)
		response.Message = "(1) ID Dosen tidak ditemukan"
		respondWithSuccess(response, w)
	}
}

func dosenUpdate(w http.ResponseWriter, r *http.Request) {
	var response UpdateResponse
	response.Error = true
	response.Message = "(0) Terjadi kesalahan pada sistem"

	id := r.FormValue("id")
	file, handler, errFile := r.FormFile("foto")
	fullPath := ""
	nip := r.FormValue("nip")
	password := r.FormValue("password")
	namaLengkap := r.FormValue("nama_lengkap")
	jenisKelamin := r.FormValue("jenis_kelamin")
	description := r.FormValue("description")

	if errFile == nil {
		fmt.Println("Uploaded File: %+v\n", handler.Filename)

		photoDir := "images/dosen"
		fileName := ("dosen-*" + ".png")

		// ============ UPLOADING FILES ============ //
		tempFile, err := ioutil.TempFile(photoDir, fileName)
		if err != nil {
			fmt.Println(err)
			response.Message = "(1) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			response.Message = "(2) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}

		tempFile.Write(fileBytes)
		fullPath = tempFile.Name()
		// ============ UPLOADING FILES ============ //
	}

	if fullPath == "" {
		_, err = db.Exec("UPDATE dosen SET nip = ?, nama_lengkap = ?, jenis_kelamin = ?, password = ?, description = ? WHERE id = ?", nip, namaLengkap, jenisKelamin, generateMD5(password), description, id)
	} else {
		_, err = db.Exec("UPDATE dosen SET nip = ?, nama_lengkap = ?, jenis_kelamin = ?, password = ?, description = ?, foto = ? WHERE id = ?", nip, namaLengkap, jenisKelamin, generateMD5(password), description, fullPath, id)
	}

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Dosen berhasil di ubah"
		respondWithSuccess(response, w)
	}
}

func dosenCreate(w http.ResponseWriter, r *http.Request) {
	var response InsertionResponse
	response.Error = true
	response.Message = "(0) Terjadi kesalahan pada sistem"

	username := r.FormValue("username")
	password := r.FormValue("password")
	file, handler, errFile := r.FormFile("foto")
	fullPath := ""
	nip := r.FormValue("nip")
	namaLengkap := r.FormValue("nama_lengkap")
	jenisKelamin := r.FormValue("jenis_kelamin")
	description := r.FormValue("description")

	if isDataExists("username", username, "dosen") {
		response.Message = "Username ini sudah terdaftar"
		respondWithSuccess(response, w)
		return
	}

	fmt.Println("Uploaded File: %+v\n", handler.Filename)

	if errFile == nil {
		photoDir := "images/dosen"
		fileName := ("dosen-*" + ".png")

		// ============ UPLOADING FILES ============ //
		tempFile, err := ioutil.TempFile(photoDir, fileName)
		if err != nil {
			fmt.Println(err)
			response.Message = "(1) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			response.Message = "(2) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}

		tempFile.Write(fileBytes)
		fullPath = tempFile.Name()
		// ============ UPLOADING FILES ============ //
	}

	_, err = db.Exec("INSERT INTO dosen (username, password, foto, nip, nama_lengkap, jenis_kelamin, description, status, created_at, modified_at) VALUES (?, ?, ?, ?, ?, ?, ?, 'active', NOW(), NOW())", username, generateMD5(password), fullPath, nip, namaLengkap, jenisKelamin, description)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Dosen berhasil di buat"
		respondWithSuccess(response, w)
	}
}

func dosenTotal(w http.ResponseWriter, r *http.Request) {
	var data Counter
	var response DetailResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = data

	row := db.QueryRow("SELECT COUNT(*) as total FROM dosen WHERE status = 'active'")
	err = row.Scan(&data.Total)
	if err == nil {
		response.Error = false
		response.Message = "Data ditemukan"
		response.Data = data
		respondWithSuccess(response, w)
	} else {
		fmt.Println(err)
		response.Message = "Dosen tidak ditemukan"
		respondWithSuccess(response, w)
	}
}
// ================================ DOSEN ================================ //


// ================================ MAHASISWA ================================ //
func checkMahasiswa(_username string, _password string) (Mahasiswa, error) {
	fmt.Println("checkDosen Hit!")
	fmt.Println(_username, _password)

	var mahasiswa Mahasiswa
	row := db.QueryRow("SELECT id, username, foto, nim, nama_lengkap, jenis_kelamin, kelas, tahun_masuk, semester, status, created_at from mahasiswa WHERE username = ? AND password = ? AND status = 'active' ", _username, _password)
	err = row.Scan(&mahasiswa.Id, &mahasiswa.Username, &mahasiswa.Foto, &mahasiswa.NIM, &mahasiswa.NamaLengkap, &mahasiswa.JenisKelamin, &mahasiswa.Kelas, &mahasiswa.TahunMasuk, &mahasiswa.Semester, &mahasiswa.Status, &mahasiswa.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return mahasiswa, err
	}
	return mahasiswa, nil
}

func mahasiswaIndex(w http.ResponseWriter, r *http.Request) {
	param_limit := r.URL.Query()["limit"]
	limit := ""
	if len(param_limit) > 0 {
		limit = " LIMIT " + param_limit[0]
	}
	fmt.Println("Limit: ", limit)

	param_status := r.URL.Query()["status"]
	status := ""
	if len(param_status) > 0 {
		status = " AND status = '" + param_status[0] +"'"
	}
	fmt.Println("status: ", status)

	var response IndexResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = make([]string, 0)
	var mahasiswas = []Mahasiswa{}

	rows, err := db.Query("SELECT id, username, foto, nim, nama_lengkap, jenis_kelamin, kelas, tahun_masuk, semester, status, created_at FROM mahasiswa WHERE 1=1 "+status+" ORDER BY id DESC " + limit)
	if err != nil {
		response.Message = "(0) Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
		return
	}

	for rows.Next() {
		var mahasiswa Mahasiswa
		err = rows.Scan(&mahasiswa.Id, &mahasiswa.Username, &mahasiswa.Foto, &mahasiswa.NIM, &mahasiswa.NamaLengkap, &mahasiswa.JenisKelamin, &mahasiswa.Kelas, &mahasiswa.TahunMasuk, &mahasiswa.Semester, &mahasiswa.Status, &mahasiswa.CreatedAt)
		if err != nil {
			response.Message = "(1) Terjadi kesalahan pada database"
			fmt.Println(err)
			respondWithSuccess(response, w)
			return
		}
		mahasiswas = append(mahasiswas, mahasiswa)
	}

	response.Error = true
	response.Message = "Data ditemukan"
	response.Data = mahasiswas
	respondWithSuccess(response, w)
}

func mahasiswaLogin(w http.ResponseWriter, r *http.Request) {
	var response LoginResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Role = "mahasiswa"
	response.Data = Mahasiswa{}

	username := r.FormValue("username")
	password := r.FormValue("password")

	mahasiswa, err := checkMahasiswa(username, generateMD5(password))
	if err == nil {
		response.Error = false
		response.Message = "Akun ditemukan"
		response.Data = mahasiswa
	} else {
		response.Message = "Akun tidak ditemukan"
	}

	respondWithSuccess(response, w)
}

func mahasiswaDisable(w http.ResponseWriter, r *http.Request) {
	var response DeletionResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Mahasiswa tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	_, err = db.Exec("UPDATE mahasiswa SET status = 'deleted' WHERE id = ?", id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Mahasiswa berhasil di Nonaktifkan"
		respondWithSuccess(response, w)
	}
}

func mahasiswaEnable(w http.ResponseWriter, r *http.Request) {
	var response DeletionResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Mahasiswa tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	_, err = db.Exec("UPDATE mahasiswa SET status = 'active' WHERE id = ?", id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Mahasiswa berhasil di Aktifkan"
		respondWithSuccess(response, w)
	}
}

func mahasiswaShow(w http.ResponseWriter, r *http.Request) {
	var mahasiswa Mahasiswa
	var response DetailResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = Mahasiswa{}

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Mahasiswa tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	row := db.QueryRow("SELECT id, username, foto, nim, nama_lengkap, jenis_kelamin, kelas, tahun_masuk, semester, description, status, created_at FROM mahasiswa WHERE id = ?", id)
	err = row.Scan(&mahasiswa.Id, &mahasiswa.Username, &mahasiswa.Foto, &mahasiswa.NIM, &mahasiswa.NamaLengkap, &mahasiswa.JenisKelamin, &mahasiswa.Kelas, &mahasiswa.TahunMasuk, &mahasiswa.Semester, &mahasiswa.Description, &mahasiswa.Status, &mahasiswa.CreatedAt)
	if err == nil {
		response.Error = false
		response.Message = "Data ditemukan"
		response.Data = mahasiswa
		respondWithSuccess(response, w)
	} else {
		fmt.Println(err)
		response.Message = "(1) ID Mahasiswa tidak ditemukan"
		respondWithSuccess(response, w)
	}
}

func mahasiswaUpdate(w http.ResponseWriter, r *http.Request) {
	var response UpdateResponse
	response.Error = true
	response.Message = "(0) Terjadi kesalahan pada sistem"

	id := r.FormValue("id")
	file, handler, errFile := r.FormFile("foto")
	fullPath := ""
	nim := r.FormValue("nim")
	password := r.FormValue("password")
	namaLengkap := r.FormValue("nama_lengkap")
	jenisKelamin := r.FormValue("jenis_kelamin")
	kelas := r.FormValue("kelas")
	tahun_masuk := r.FormValue("tahun_masuk")
	semester := r.FormValue("semester")
	description := r.FormValue("description")

	if errFile == nil {
		fmt.Println("Uploaded File: %+v\n", handler.Filename)

		photoDir := "images/mahasiswa"
		fileName := ("mahasiswa-*" + ".png")

		// ============ UPLOADING FILES ============ //
		tempFile, err := ioutil.TempFile(photoDir, fileName)
		if err != nil {
			fmt.Println(err)
			response.Message = "(1) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			response.Message = "(2) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}

		tempFile.Write(fileBytes)
		fullPath = tempFile.Name()
		// ============ UPLOADING FILES ============ //
	}

	if fullPath == "" {
		_, err = db.Exec("UPDATE mahasiswa SET nim = ?, nama_lengkap = ?, jenis_kelamin = ?, password = ?, kelas = ?, tahun_masuk = ?, semester = ?, description = ? WHERE id = ?", nim, namaLengkap, jenisKelamin, generateMD5(password), kelas, tahun_masuk, semester, description, id)
	} else {
		_, err = db.Exec("UPDATE mahasiswa SET nim = ?, nama_lengkap = ?, jenis_kelamin = ?, password = ?, kelas = ?, tahun_masuk = ?, semester = ?, description = ?, foto = ? WHERE id = ?", nim, namaLengkap, jenisKelamin, generateMD5(password), kelas, tahun_masuk, semester, description, fullPath, id)
	}

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Mahasiswa berhasil di ubah"
		respondWithSuccess(response, w)
	}
}

func mahasiswaCreate(w http.ResponseWriter, r *http.Request) {
	var response InsertionResponse
	response.Error = true
	response.Message = "(0) Terjadi kesalahan pada sistem"

	username := r.FormValue("username")
	password := r.FormValue("password")
	file, handler, errFile := r.FormFile("foto")
	fullPath := ""
	nim := r.FormValue("nim")
	namaLengkap := r.FormValue("nama_lengkap")
	jenisKelamin := r.FormValue("jenis_kelamin")
	kelas := r.FormValue("kelas")
	tahun_masuk := r.FormValue("tahun_masuk")
	semester := r.FormValue("semester")
	description := r.FormValue("description")

	if isDataExists("username", username, "mahasiswa") {
		response.Message = "Username ini sudah terdaftar"
		respondWithSuccess(response, w)
		return
	}

	fmt.Println("Uploaded File: %+v\n", handler.Filename)

	if errFile == nil {
		photoDir := "images/mahasiswa"
		fileName := ("mahasiswa-*" + ".png")

		// ============ UPLOADING FILES ============ //
		tempFile, err := ioutil.TempFile(photoDir, fileName)
		if err != nil {
			fmt.Println(err)
			response.Message = "(1) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			response.Message = "(2) Terjadi kesalahan pada upload File"
			respondWithSuccess(response, w)
			return
		}

		tempFile.Write(fileBytes)
		fullPath = tempFile.Name()
		// ============ UPLOADING FILES ============ //
	}

	_, err = db.Exec("INSERT INTO mahasiswa (username, password, foto, nim, nama_lengkap, jenis_kelamin, kelas, tahun_masuk, semester, description, status, created_at, modified_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'active', NOW(), NOW())", username, generateMD5(password), fullPath, nim, namaLengkap, jenisKelamin, kelas, tahun_masuk, semester, description)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Akun Mahasiswa berhasil di buat"
		respondWithSuccess(response, w)
	}
}

func mahasiswaTotal(w http.ResponseWriter, r *http.Request) {
	var data Counter
	var response DetailResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = data

	row := db.QueryRow("SELECT COUNT(*) as total FROM mahasiswa WHERE status = 'active'")
	err = row.Scan(&data.Total)
	if err == nil {
		response.Error = false
		response.Message = "Data ditemukan"
		response.Data = data
		respondWithSuccess(response, w)
	} else {
		fmt.Println(err)
		response.Message = "Mahasiswa tidak ditemukan"
		respondWithSuccess(response, w)
	}
}
// ================================ MAHASISWA ================================ //

// ================================ REQUEST ================================ //
func requestIndex(w http.ResponseWriter, r *http.Request) {
	var requests = []Request{}
	var response IndexResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = make([]string, 0)

	param_limit := r.URL.Query()["limit"]
	limit := ""
	if len(param_limit) > 0 {
		limit = " LIMIT " + param_limit[0]
	}
	fmt.Println("Limit: ", limit)

	param_status := r.URL.Query()["status"]
	status := ""
	if len(param_status) > 0 {
		status = " AND r.status = '" + param_status[0] +"'"
	}
	fmt.Println("status: ", status)

	query := `SELECT r.*, rt.name as request_name, d.nama_lengkap as dosen_name, m.nama_lengkap as mahasiswa_name FROM requests r
			INNER JOIN request_type rt ON r.request_type_id = rt.id
			INNER JOIN dosen d ON r.dosen_id = d.id
			INNER JOIN mahasiswa m ON r.mahasiswa_id = m.id
			WHERE 1=1 ` + status + ` ORDER BY id DESC ` + limit
	rows, err := db.Query(query)
	if err != nil {
		response.Message = "(0) Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
		return
	}

	for rows.Next() {
		var request Request
		err = rows.Scan(&request.Id, &request.RequestTypeId, &request.Description, &request.DosenId, &request.MahasiswaId, &request.StartDatetime, &request.EndDatetime, &request.Status, &request.CreatedAt, &request.ModifiedAt, &request.RequestName, &request.DosenName, &request.MahasiswaName)
		if err != nil {
			response.Message = "(1) Terjadi kesalahan pada database"
			fmt.Println(err)
			respondWithSuccess(response, w)
			return
		}
		requests = append(requests, request)
	}

	response.Error = true
	response.Message = "Data ditemukan"
	response.Data = requests
	respondWithSuccess(response, w)
}

func requestDosen(w http.ResponseWriter, r *http.Request) {
	var requests = []Request{}
	var response IndexResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = make([]string, 0)

	param_limit := r.URL.Query()["limit"]
	limit := ""
	if len(param_limit) > 0 {
		limit = " LIMIT " + param_limit[0]
	}
	fmt.Println("Limit: ", limit)

	param_status := r.URL.Query()["status"]
	status := ""
	if len(param_status) > 0 {
		status = " AND STATUS = '" + param_status[0] +"'"
	}
	fmt.Println("status: ", status)

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Dosen tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	rows, err := db.Query("SELECT * FROM requests WHERE dosen_id = ? " + status + limit, id)
	if err != nil {
		response.Message = "(0) Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
		return
	}

	for rows.Next() {
		var request Request
		err = rows.Scan(&request.Id, &request.RequestTypeId, &request.Description, &request.DosenId, &request.MahasiswaId, &request.StartDatetime, &request.EndDatetime, &request.Status, &request.CreatedAt, &request.ModifiedAt)
		if err != nil {
			response.Message = "(1) Terjadi kesalahan pada database"
			fmt.Println(err)
			respondWithSuccess(response, w)
			return
		}
		requests = append(requests, request)
	}

	response.Error = true
	response.Message = "Data ditemukan"
	response.Data = requests
	respondWithSuccess(response, w)
}

func requestMahasiswa(w http.ResponseWriter, r *http.Request) {
	var requests = []Request{}
	var response IndexResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = make([]string, 0)

	param_limit := r.URL.Query()["limit"]
	limit := ""
	if len(param_limit) > 0 {
		limit = " LIMIT " + param_limit[0]
	}
	fmt.Println("Limit: ", limit)

	param_status := r.URL.Query()["status"]
	status := ""
	if len(param_status) > 0 {
		status = " AND STATUS = '" + param_status[0] +"'"
	}
	fmt.Println("status: ", status)

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Mahasiswa tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	rows, err := db.Query("SELECT * FROM requests WHERE mahasiswa_id = ? " + status + limit, id)
	if err != nil {
		response.Message = "(0) Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
		return
	}

	for rows.Next() {
		var request Request
		err = rows.Scan(&request.Id, &request.RequestTypeId, &request.Description, &request.DosenId, &request.MahasiswaId, &request.StartDatetime, &request.EndDatetime, &request.Status, &request.CreatedAt, &request.ModifiedAt)
		if err != nil {
			response.Message = "(1) Terjadi kesalahan pada database"
			fmt.Println(err)
			respondWithSuccess(response, w)
			return
		}
		requests = append(requests, request)
	}

	response.Error = true
	response.Message = "Data ditemukan"
	response.Data = requests
	respondWithSuccess(response, w)
}

func requestDosenAccept(w http.ResponseWriter, r *http.Request) {
	var response DeletionResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Request tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	_, err = db.Exec("UPDATE requests SET status = 'accepted', modified_at = NOW() WHERE id = ?", id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Request berhasil di Setujui"
		respondWithSuccess(response, w)
	}
}

func requestDosenReject(w http.ResponseWriter, r *http.Request) {
	var response DeletionResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Request tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	_, err = db.Exec("UPDATE requests SET status = 'rejected', modified_at = NOW() WHERE id = ?", id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Request berhasil di Tolak"
		respondWithSuccess(response, w)
	}
}

func requestMahasiswaDelete(w http.ResponseWriter, r *http.Request) {
	var response DeletionResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"

	_id := mux.Vars(r)["id"]
	id, err := stringToInt64(_id)

	if err != nil {
		response.Message = "(0) ID Request tidak ditemukan"
		respondWithSuccess(response, w)
		return
	}

	_, err = db.Exec("UPDATE requests SET status = 'deleted', modified_at = NOW() WHERE id = ?", id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Request berhasil di hapus"
		respondWithSuccess(response, w)
	}
}

func requestMahasiswaUpdate(w http.ResponseWriter, r *http.Request) {
	var response UpdateResponse
	response.Error = true
	response.Message = "(0) Terjadi kesalahan pada sistem"

	id := r.FormValue("id")
	request_type_id := r.FormValue("request_type_id")
	description := r.FormValue("description")
	dosen_id := r.FormValue("dosen_id")
	mahasiswa_id := r.FormValue("mahasiswa_id")
	start_datetime := r.FormValue("start_datetime")
	end_datetime := r.FormValue("end_datetime")

	_, err = db.Exec("UPDATE requests SET request_type_id = ?, description = ?, dosen_id = ?, mahasiswa_id = ?, start_datetime = ?, end_datetime = ?, modified_at = NOW() WHERE id = ?", request_type_id, description, dosen_id, mahasiswa_id, start_datetime, end_datetime, id)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Request ke Dosen berhasil di ubah"
		respondWithSuccess(response, w)
	}
}

func requestMahasiswaCreate(w http.ResponseWriter, r *http.Request) {
	var response InsertionResponse
	response.Error = true
	response.Message = "(0) Terjadi kesalahan pada sistem"

	request_type_id := r.FormValue("request_type_id")
	description := r.FormValue("description")
	dosen_id := r.FormValue("dosen_id")
	mahasiswa_id := r.FormValue("mahasiswa_id")
	start_datetime := r.FormValue("start_datetime")
	end_datetime := r.FormValue("end_datetime")

	_, err = db.Exec("INSERT INTO requests (request_type_id, description, dosen_id, mahasiswa_id, start_datetime, end_datetime, status, created_at, modified_at) VALUES (?, ?, ?, ?, ?, ?, 'pending', NOW(), NOW())", request_type_id, description, dosen_id, mahasiswa_id, start_datetime, end_datetime)

	if err != nil {
		response.Message = "Terjadi kesalahan pada database"
		fmt.Println(err)
		respondWithSuccess(response, w)
	} else {
		response.Error = false
		response.Message = "Request untuk Dosen berhasil di buat"
		respondWithSuccess(response, w)
	}
}

func requestTotal(w http.ResponseWriter, r *http.Request) {
	var data Counter
	var response DetailResponse
	response.Error = true
	response.Message = "Terjadi kesalahan pada sistem"
	response.Data = data

	param_status := r.URL.Query()["status"]
	status := ""
	if len(param_status) > 0 {
		status = " AND STATUS = '" + param_status[0] +"'"
	}
	fmt.Println("status: ", status)

	row := db.QueryRow("SELECT COUNT(*) as total FROM requests WHERE 1=1 " + status)
	err = row.Scan(&data.Total)
	if err == nil {
		response.Error = false
		response.Message = "Data ditemukan"
		response.Data = data
		respondWithSuccess(response, w)
	} else {
		fmt.Println(err)
		response.Message = "Request tidak ditemukan"
		respondWithSuccess(response, w)
	}
}
// ================================ REQUEST ================================ //