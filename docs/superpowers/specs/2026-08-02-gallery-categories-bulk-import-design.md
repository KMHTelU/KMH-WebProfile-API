# Gallery Categories, Bulk Operations, Import Excel, Rename NIM, dan Reset Password

Tanggal: 2026-08-02
Status: disetujui

## Latar Belakang

Frontend membutuhkan daftar kategori gallery untuk memfilter gallery berdasarkan Event.
Selain itu tim admin membutuhkan cara memasukkan dan memperbarui data dalam jumlah banyak,
baik lewat JSON (bulk) maupun lewat berkas Excel/CSV (import).

Basis data tidak memiliki tabel kategori gallery. Satu-satunya pengelompokan yang tersedia
adalah relasi `galleries.event_id -> events.id`, sehingga "kategori" pada konteks ini berarti Event.

## Ruang Lingkup

Termasuk:

- Endpoint daftar kategori gallery dan filter gallery per event.
- Bulk create dan bulk update untuk 8 entitas.
- Import Excel/CSV untuk 5 entitas, lengkap dengan template.
- Perbaikan `INNER JOIN` menjadi `LEFT JOIN` pada members, divisions, events, galleries.
- Melengkapi `member_divisions` yang belum memiliki kode sama sekali.
- Rename field `npm` menjadi `nim` di seluruh lapisan.
- Fitur change password dan forgot password dengan pengiriman email lewat SMTP.

Tidak termasuk:

- Segala hal terkait blog (`blog_posts`, `blog_categories`, `blog_tags`, `blog_post_tags`).
- `contact_messages`.
- Bulk maupun import untuk `users` (dikeluarkan karena menyangkut kredensial).
- Import Excel untuk `homepage_banners` (butuh berkas gambar; cukup bulk saja).
- Import Excel untuk `galleries` dan `gallery_items`.

## Keputusan Desain

### Kategori gallery adalah Event

`GET /api/galleries/categories` mengembalikan Event yang memiliki minimal satu gallery,
beserta jumlah gallery-nya. Frontend memakai daftar ini sebagai opsi filter, lalu memanggil
`GET /api/galleries?event_id=<uuid>`.

Alternatif membuat tabel kategori baru ditolak karena menambah entitas yang tidak dibutuhkan;
Event sudah menjadi pengelompokan alami gallery.

### INNER JOIN diganti LEFT JOIN

`GetAllMembers`, `GetMemberByID`, `GetAllDivisions`, `GetDivisionByID`, `ListEvents`,
`GetEventByID`, `SelectAllGalleries`, dan `SelectGalleryByID` memakai `INNER JOIN` ke `media`
(dan ke `members` untuk divisions, ke `events` untuk galleries).

Akibatnya baris tanpa relasi tersebut tidak pernah muncul. Karena import Excel tidak mengunggah
gambar, seluruh data hasil import akan langsung tak terlihat. Semua join diubah ke `LEFT JOIN`.

Konsekuensi terhadap kontrak API: field relasi kini dapat bernilai `null`. Ini perlu
dikomunikasikan ke tim frontend.

### Partial success, bukan transaksional

Bulk dan import memproses baris satu per satu. Baris yang gagal tidak membatalkan baris lain.
Setiap operasi mengembalikan laporan per-item.

Alasan: berkas Excel hasil kerja manual hampir selalu punya beberapa baris bermasalah.
Menolak seluruh berkas karena satu kesalahan ketik membuat alur kerja admin menyakitkan.

### Kode status HTTP

- Semua item berhasil: `201` untuk create dan import, `200` untuk update.
- Sebagian berhasil: `207 Multi-Status`.
- Semua item gagal, atau payload tidak valid: `400`.

Preview import selalu `200` selama berkas terbaca, karena tidak ada yang ditulis.

### Kolom relasi di Excel memakai nilai human-readable

Kolom relasi diisi nilai yang dikenal admin (NIM anggota, slug divisi, nama role), bukan UUID.
Backend me-resolve ke UUID dan melaporkan error bila tidak ditemukan.

### Duplikat dianggap error

Baris yang bentrok dengan data existing (NIM, slug, nama role yang sudah ada) ditandai gagal
pada baris tersebut; baris lain tetap diproses. Upsert tidak dipakai agar import tidak
diam-diam menimpa data.

Duplikat di dalam satu berkas juga dideteksi. Pemeriksaan ke basis data saja tidak cukup:
dua baris ber-NIM sama akan lolos karena baris pertama belum tersimpan saat baris kedua
diperiksa, dan pada mode preview tidak ada yang tersimpan sama sekali.

### Rename npm menjadi nim

Perubahan menyeluruh: kolom basis data, field JSON request dan response, kolom template Excel,
dan swagger. Ini breaking change bagi frontend dan perlu dikomunikasikan.

Kolom basis data diubah lewat migrasi `ALTER TABLE members RENAME COLUMN npm TO nim`, bukan
dengan menyunting migrasi awal, karena migrasi awal sudah dijalankan di lingkungan yang ada.

### Reset password memakai tautan berisi token acak

`POST /api/forgot-password` mengirim email berisi
`{FRONTEND_URL}/reset-password?token=...`. Token acak 256 bit, berlaku 60 menit, sekali pakai.

Basis data hanya menyimpan digest SHA-256 token, bukan token aslinya, supaya bocornya isi tabel
tidak langsung memberi penyerang token yang bisa dipakai. SHA-256 dipilih alih-alih bcrypt karena
token sudah acak 256 bit sehingga tidak rentan ditebak, sementara pencarian token perlu lewat indeks.

Endpoint selalu membalas sukses dengan pesan generik, baik email terdaftar maupun tidak, agar
tidak bisa dipakai memetakan alamat email mana yang punya akun. Setiap permintaan baru
membatalkan token lama, dan jumlah permintaan dibatasi 3 per akun per jam.

Email dikirim di goroutine terpisah agar respons API tidak menunggu SMTP; kegagalan kirim
dicatat di log. Bila `SMTP_HOST` kosong, email tidak dikirim dan isinya hanya dicetak ke log,
sehingga pengembangan lokal berjalan tanpa kredensial.

### Change password mewajibkan password lama

`POST /api/protected/auth/change-password` meminta password lama dan password baru yang berbeda.
Setelah berhasil, tautan reset yang masih menganggur dibatalkan dan email notifikasi dikirim.

Token akses yang sudah beredar tetap berlaku sampai kedaluwarsa; JWT bersifat stateless dan
pencabutan sesi berada di luar cakupan pekerjaan ini.

## Arsitektur

Mengikuti pola yang sudah ada: `routes -> handlers -> services -> repositories -> generated (sqlc)`.

Komponen baru:

- `utils/bulk.go` — tipe laporan hasil bulk dan responder yang memilih 200/201/207/400.
- `utils/nulls.go` — helper konversi ke tipe nullable sqlc.
- `utils/reset_token.go` — pembuatan token acak dan digest SHA-256.
- `internal/importer/` — package mandiri tanpa ketergantungan ke Fiber maupun database:
  - `spec.go` — `ColumnSpec` dan `EntitySpec` yang mendeskripsikan kolom sebuah template.
  - `row.go` — baris bernomor beserta konversi nilai sel ke int, bool, dan waktu.
  - `parser.go` — membaca `.xlsx`/`.csv` menjadi baris bernomor.
  - `template.go` — menghasilkan berkas template dari `EntitySpec`.
  - `specs.go` — definisi kolom untuk kelima entitas.
- `internal/mailer/` — klien SMTP dan template email HTML beserta versi teksnya.
- `internal/services/entity_core_service.go` — fungsi create/update per entitas tanpa
  pemeriksaan token dan tanpa activity log, dipakai bersama oleh endpoint tunggal, bulk, dan import.
- `internal/services/import_service.go` — orkestrasi preview/commit dan resolusi relasi.
- `internal/services/auth_password_service.go` — change, forgot, dan reset password.
- `internal/repositories/lookup_repository.go` — pencarian berdasarkan kunci human-readable.
- `internal/{repositories,services,handlers}/member_division_*.go` — stack `member_divisions`.
- `tools/gentemplates/` — perintah yang menulis berkas template statis ke `docs/templates/`.

`internal/importer` dan `internal/mailer` sengaja dipisah dari layer service supaya definisi
kolom, parser, generator template, dan tata letak email dapat diuji tanpa database maupun HTTP.

Endpoint tunggal yang sudah ada direfaktor agar memanggil fungsi core yang sama dengan bulk dan
import, sehingga pemetaan field dan pemeriksaan duplikat hanya ditulis di satu tempat.

## Kontrak API

### Gallery

| Method | Path | Auth | Keterangan |
| --- | --- | --- | --- |
| GET | `/api/galleries/categories` | tidak | Event yang punya gallery + `gallery_count` |
| GET | `/api/galleries?event_id=&is_public=&limit=&start=` | tidak | Filter opsional |

Route `/galleries/categories` harus didaftarkan sebelum `/galleries/:id`.

### Bulk

`POST /api/protected/<res>/bulk` dan `PUT /api/protected/<res>/bulk` untuk:
`members`, `divisions`, `member-divisions`, `events`, `galleries`, `gallery-items`,
`homepage-banners`, `roles`.

Body: `{ "items": [ ... ] }`, 1 sampai 500 item. Pada bulk update setiap item wajib memuat `id`.

Bentuk respons:

```json
{
  "status": "partial",
  "message": "Bulk create members processed",
  "data": {
    "total": 3,
    "succeeded": 2,
    "failed": 1,
    "results": [
      { "index": 0, "status": "success", "id": "..." },
      { "index": 1, "status": "failed", "error": "NIM 1301213045 sudah terdaftar atas nama Aulia Rahman" },
      { "index": 2, "status": "success", "id": "..." }
    ]
  }
}
```

Route `/bulk` harus didaftarkan sebelum `/:id`.

### Import

Untuk `members`, `divisions`, `member-divisions`, `events`, `roles`:

| Method | Path | Keterangan |
| --- | --- | --- |
| GET | `/api/protected/import/<res>/template?format=xlsx\|csv` | Unduh template |
| POST | `/api/protected/import/<res>/preview` | Validasi saja, tidak menulis DB |
| POST | `/api/protected/import/<res>` | Eksekusi |

Preview dan import menerima multipart dengan field `file`. Format yang diterima `.xlsx` dan `.csv`,
maksimal 5 MB dan 500 baris data.

Laporan import memakai nomor baris berkas yang sebenarnya, misalnya
`baris 7: division_slug 'humas' tidak ditemukan`.

## Kolom Template

Kolom bertanda `*` wajib diisi.

| Entitas | Kolom |
| --- | --- |
| members | `name*`, `nim*`, `email`, `phone`, `bio`, `instagram_url`, `period_start*`, `period_end*`, `is_active` |
| divisions | `name*`, `slug*`, `description`, `coordinator_nim`, `is_active` |
| member_divisions | `member_nim*`, `division_slug*`, `role_title` |
| events | `title*`, `slug*`, `description`, `event_type`, `start_time*`, `end_time`, `location`, `google_maps_url`, `registration_url`, `status`, `is_published` |
| roles | `name*`, `description` |

Setiap template berisi dua sheet: **Data** (baris header dan dua baris contoh) dan
**Petunjuk** (arti tiap kolom, wajib atau opsional, format tanggal, cara mengisi kolom relasi).

Template tersedia lewat endpoint yang menghasilkannya saat diminta, dan sebagai berkas statis
di `docs/templates/` yang ikut di-commit.

Format tanggal yang diterima: `YYYY-MM-DD`, `YYYY-MM-DD HH:mm`, `YYYY-MM-DD HH:mm:ss`, RFC3339,
serta nilai tanggal asli Excel.
Nilai boolean yang diterima: `true`/`false`, `1`/`0`, `ya`/`tidak`, `yes`/`no`, `aktif`/`nonaktif`.

## Kontrak API Autentikasi

| Method | Path | Auth | Keterangan |
| --- | --- | --- | --- |
| POST | `/api/forgot-password` | tidak | Kirim tautan reset ke email |
| POST | `/api/reset-password` | tidak | Tukar token dengan password baru |
| POST | `/api/protected/auth/change-password` | ya | Ganti password sendiri |

Variabel lingkungan baru didokumentasikan di `.env.example`: `APP_NAME`, `FRONTEND_URL`,
`SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_EMAIL`, `SMTP_FROM_NAME`,
`SMTP_ENCRYPTION`, `PASSWORD_RESET_TOKEN_TTL_MINUTES`, dan `PASSWORD_RESET_MAX_PER_HOUR`.

Konfigurasi SMTP dibuat generik lewat env sehingga Gmail, Brevo, Resend, SendGrid, maupun mail
server domain sendiri bisa dipakai tanpa mengubah kode.

## Perubahan Basis Data

Dua migrasi dbmate baru:

- `20260802231500_rename_members_npm_to_nim.sql` — mengganti nama kolom `members.npm` menjadi `nim`.
- `20260802231600_create_password_reset_tokens.sql` — tabel `password_reset_tokens` berisi
  `id`, `user_id`, `token_hash`, `expires_at`, `used_at`, `created_at`, dengan indeks unik pada
  `token_hash` dan indeks pada `user_id`.

Perubahan pada berkas query sqlc:

- Perbaikan `InsertMemberDivision` yang tidak mengisi kolom `id`, padahal kolom itu primary key
  tanpa nilai default sehingga query selalu gagal.
- `UpdateMember` sebelumnya tidak pernah menyimpan kolom NIM meski request memuatnya; kini disimpan.
- Query baru: `UpdateMemberDivision`, `GetMemberDivisionByID`, `GetMemberDivisionByPair`,
  `UpdateHomepageBanner`, `GetHomepageBannerByID`, `UpdateGalleryItem`, `GetGalleryItemByID`,
  `SelectGalleryCategories`, `UpdateUserPassword`, seluruh query `password_reset_tokens`,
  serta query lookup `GetMemberByNIM`, `GetDivisionBySlug`, `GetEventBySlug`, `GetRoleByName`.

## Penanganan Error

Kegagalan tingkat permintaan (token tidak valid, berkas tidak terbaca, jumlah baris melebihi batas)
mengembalikan `*fiber.Error` seperti endpoint lain.

Kegagalan tingkat baris dikumpulkan ke dalam laporan dan tidak menghentikan pemrosesan.
Pesan error ditulis dalam bentuk yang bisa ditindaklanjuti admin, menyebut kolom dan nilai
yang bermasalah, bukan meneruskan pesan mentah dari driver basis data.

Satu activity log ringkasan dicatat per operasi bulk atau import, bukan satu log per baris,
agar tabel `activity_logs` tidak membengkak.

## Verifikasi

- `sqlc generate` berhasil dan `go build ./...` bersih.
- `go vet ./...` bersih.
- Template yang dihasilkan endpoint dapat dibuka, lalu diisi dan diunggah kembali lewat
  preview dan import.
- `GET /api/galleries/categories` mengembalikan jumlah yang sesuai isi basis data.
- Gallery tanpa event, member tanpa foto, dan event tanpa cover kini muncul di endpoint list.

## Dependensi Baru

- `github.com/xuri/excelize/v2` untuk membaca dan menulis `.xlsx`.
- `github.com/wneessen/go-mail` untuk pengiriman SMTP beserta penanganan STARTTLS dan SSL.
