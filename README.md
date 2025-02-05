# RevisiTOUKOM

# 1. Konfigurasi Koneksi Database dengan GORM di Golang. 
## a. Kegunaan
- menghubungkan aplikasi golang ke database mysql menggunakan gorm. dengan fungsi ConnectDB(), aplikasi bisa membuka koneksi ke database dan menyimpannya di variabel global DB, sehingga bisa digunakan di berbagai bagian aplikasi untuk operasi database seperti insert, update, delete, dan query data. Penjelasan berbaris

> **_package config_**
- mendeklarasikan bahwa file ini ada di dalam package config. ini menunjukkan bahwa file ini digunakan untuk konfigurasi aplikasi, khususnya database.

> **_import (_**
**_"log"_**
**_"gorm.io/driver/mysql"_**
**_"gorm.io/gorm"_**
**__)__**
- log: digunakan untuk mencetak pesan log atau error ke terminal.
- gorm.io/driver/mysql: package driver mysql untuk gorm agar bisa menghubungkan aplikasi dengan database mysql.
- gorm.io/gorm: package utama dari gorm yang berisi fitur-fitur orm.

> **_var DB *gorm.DB_**
- deklarasi variabel global DB bertipe *gorm.DB, yang nantinya akan menyimpan koneksi database dan bisa diakses dari berbagai bagian aplikasi.

> **_func ConnectDB(){_**
- mendeklarasikan fungsi ConnectDB(), yang bertanggung jawab untuk membuat koneksi ke database mysql.

> **_dsn := "root:inipassword@tcp(127.0.0.1:3306)/tryoutukom?charset=utf8mb4&parseTime=True&loc=Local"_**
- dsn (data source name) adalah string yang berisi informasi koneksi ke database.
  - root → username database.
  - inipassword → password database.
  - tcp(127.0.0.1:3306) → menghubungkan ke database melalui tcp di localhost pada port 3306.
  - tryoutukom → nama database yang digunakan.
  - charset=utf8mb4 → encoding karakter agar bisa support karakter unik (seperti emoji).
  - parseTime=True → mengaktifkan parsing waktu otomatis.
  - loc=Local → menggunakan zona waktu lokal.

> **_var err error_**
- mendeklarasikan variabel err bertipe error untuk menangkap error yang mungkin terjadi saat koneksi ke database.

> **_DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})_**
- mencoba membuka koneksi ke database menggunakan gorm.
- jika koneksi berhasil, objek koneksi database akan disimpan di variabel DB.
- jika gagal, error akan disimpan di variabel err.

> **_if err != nil {_**
**_log.Fatal("Failed to connect to the database: ", err)_**
**_}_**
- mengecek apakah ada error (err != nil).
- jika terjadi error, akan mencetak pesan "Failed to connect to the database" dan menghentikan program menggunakan log.Fatal().

# 2. Inisialisasi Server dan Konfigurasi Database pada Aplikasi Golang
## a. Kegunaan
- menghubungkan aplikasi ke database menggunakan fungsi ConnectDB().
- melakukan migrasi database dengan AutoMigrate() agar tabel sesuai dengan struktur model.
- mengatur routing aplikasi dengan SetupRoutes().
- menjalankan server web pada port 8080.

>**_package main_**
- mendeklarasikan bahwa file ini adalah bagian dari package main, yang berarti ini adalah titik awal eksekusi program.

>**_import (_**
**_"revisitoukom/config"_**
**_"revisitoukom/models"_**
**_"revisitoukom/routes"_**
)
- "revisitoukom/config" → mengimpor package config, yang berisi konfigurasi database.
- "revisitoukom/models" → mengimpor package models, yang berisi struktur tabel dalam database.
- "revisitoukom/routes" → mengimpor package routes, yang berisi konfigurasi rute (endpoint api).

>**_func main() {_**
- mendeklarasikan fungsi main(), yang merupakan titik awal eksekusi program.

>**_config.ConnectDB()_**
- memanggil fungsi ConnectDB() dari package config untuk menghubungkan aplikasi ke database.

>**_config.DB.AutoMigrate(&models.User{}, &models.Packet{}, &models.Question{}, &models.Order{}, &models.Exam{}, &models.ExamQuestion{})_**
- menjalankan proses migrasi database menggunakan AutoMigrate().
- fungsi ini memastikan bahwa tabel dalam database sesuai dengan struktur yang ada di model User, Packet, Question, Order, Exam, dan ExamQuestion.
- jika tabel belum ada, maka akan dibuat secara otomatis.

>**_router := routes.SetupRoutes()_**
- memanggil fungsi SetupRoutes() dari package routes untuk mengatur routing (endpoint api).
- hasilnya disimpan dalam variabel router, yang akan digunakan untuk menjalankan server.

>**_router.Run("0.0.0.0:8080")_**
- menjalankan server web menggunakan router pada alamat 0.0.0.0:8080.
- 0.0.0.0 berarti server bisa diakses dari jaringan mana saja, dan 8080 adalah port yang digunakan.

# 3. Struct Models
**1. Orders**
| Nama Kolom  | Tipe Data         | Keterangan                            |
|-------------|-------------------|---------------------------------------|
| id          | BIGINT UNSIGNED   | primary key (auto increment)          |
| user_id     | INT               | id pengguna (tidak boleh null)        |
| packet_id   | INT               | id paket (tidak boleh null)           |
| status      | INT               | status order (tidak boleh null)       |
| order_date  | DATE              | tanggal order (tidak boleh null)      |
| amount      | DECIMAL(10,2)     | jumlah pembayaran                     |
| created_at  | DATETIME          | waktu data dibuat (auto create)       |
| updated_at  | DATETIME          | waktu data diperbarui (auto update)   |
### (a). Kegunaan
- mendeklarasikan bahwa file ini ada di dalam package models, yang berarti berisi struktur model data dalam aplikasi.
- mengimpor package time untuk mendukung penggunaan tipe data time.Time, yang digunakan untuk menyimpan informasi tanggal dan waktu.
- mendeklarasikan struct Order, yang merepresentasikan tabel Order dalam database.
- ID sebagai primary key (kunci utama) dalam tabel.
  - tipe datanya uint (unsigned integer), yang berarti nilainya tidak bisa negatif.
- UserID menyimpan id pengguna yang melakukan order.
  - bertipe int, dan tidak boleh null (gorm:"not null"), artinya harus selalu memiliki nilai.
- PacketID menyimpan id paket yang dipesan.
  - bertipe int, dan tidak boleh null.
- Status menyimpan status order (misalnya, 0 = pending, 1 = paid, 2 = canceled).
  - bertipe int, dan tidak boleh null.
- OrderDate menyimpan tanggal order.
  - bertipe time.Time dengan tipe database date (gorm:"type:date", hanya menyimpan tanggal tanpa waktu).
  - tidak boleh null.
- Amount menyimpan jumlah pembayaran dalam bentuk desimal.
  - bertipe float64, dengan format decimal(10,2), artinya maksimal 10 digit dengan 2 angka di belakang koma (contoh: 12345.67).

**2. Exam**
| Nama Kolom   | Tipe Data         | Keterangan                             |
|--------------|-------------------|----------------------------------------|
| id           | BIGINT UNSIGNED   | primary key (auto increment)           |
| user_id      | INT               | id pengguna (tidak boleh null)         |
| packet_id    | INT               | id paket (tidak boleh null)            |
| score        | DECIMAL(10,2)     | nilai ujian                            |
| started_at   | DATETIME          | waktu ujian dimulai (tidak boleh null) |
| ended_at     | DATETIME          | waktu ujian selesai (tidak boleh null) |
| created_at   | DATETIME          | waktu data dibuat (auto create)        |
| updated_at   | DATETIME          | waktu data diperbarui (auto update)    |
### (a). Kegunaan
- ID: kolom ini bertindak sebagai primary key di tabel Exam. Tipe data uint adalah tipe integer yang tidak negatif.
  - gorm:"primaryKey": menandakan bahwa kolom ini adalah primary key.
- OrderID: menyimpan ID pesanan atau order terkait dengan exam. Tipe data int64 digunakan karena ID bisa sangat besar.
  - json:"order_id": digunakan untuk penamaan JSON jika data ini dikirimkan dalam format JSON
- PacketID: menyimpan ID paket terkait dengan exam. Tipe data int64 digunakan untuk menyimpan angka besar.
  - json:"packet_id": penamaan JSON untuk field ini.
- UserID: menyimpan ID pengguna yang mengikuti ujian. Tipe data int64 dipakai untuk ID besar.
  - json:"user_id": penamaan JSON untuk field ini.
- Score: menyimpan skor ujian. Tipe data float64 digunakan karena skor bisa berupa angka desimal.
  - json:"score": penamaan JSON untuk field ini.
- StartedAt: menyimpan waktu mulai ujian. Tipe data time.Time digunakan untuk menangani waktu.
  - json:"started_at": penamaan JSON untuk field ini.
- EndedAt: menyimpan waktu selesai ujian. Tipe data time.Time digunakan untuk menangani waktu.
  - json:"ended_at": penamaan JSON untuk field ini.

**3. ExamQuestions**
| Nama Kolom   | Tipe Data         | Keterangan                             |
|--------------|-------------------|----------------------------------------|
| id           | BIGINT UNSIGNED   | primary key (auto increment)           |
| exam_id      | INT               | id ujian (tidak boleh null)            |
| question_id  | INT               | id soal (tidak boleh null)             |
| user_answer  | TEXT              | jawaban pengguna                       |
| created_at   | DATETIME          | waktu data dibuat (auto create)        |
| updated_at   | DATETIME          | waktu data diperbarui (auto update)    |
### (a). Kegunaan
- ID: kolom ini bertipe uint dan menjadi primary key untuk tabel.
  - gorm:"primaryKey" menandakan bahwa kolom ini adalah primary key di database.
- ExamID: kolom ini bertipe int64 dan berfungsi untuk menyimpan ID dari exam yang terkait.
  - json:"exam_id" menandakan bahwa ketika objek ini di-encode ke dalam format JSON, field ini akan dinamai exam_id.
- QuestionID: kolom ini bertipe int64 dan berfungsi untuk menyimpan ID dari question yang terkait.
  - json:"question_id" menandakan bahwa field ini akan dinamakan question_id dalam format JSON.
- UserAnswer: kolom ini bertipe string dan menyimpan jawaban dari user.
  - json:"user_answer" menandakan bahwa field ini akan dinamakan user_answer dalam format JSON.

