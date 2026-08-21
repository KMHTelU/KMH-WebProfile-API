package importer

// Definisi kolom setiap template. Kolom relasi diisi nilai yang dikenal admin
// (NIM, slug) dan diterjemahkan ke UUID oleh layer service.

var membersSpec = EntitySpec{
	Key:   "members",
	Title: "Template Import Data Anggota",
	Notes: []string{
		"NIM dipakai sebagai penanda unik. Baris dengan NIM yang sudah terdaftar akan ditolak.",
		"Foto anggota tidak bisa diimpor lewat berkas ini. Unggah foto menyusul lewat endpoint upload foto anggota.",
	},
	Columns: []ColumnSpec{
		{Key: "name", Type: TypeString, Required: true,
			Description: "Nama lengkap anggota.",
			Examples:    [2]string{"Aulia Rahman", "Bunga Lestari"}},
		{Key: "nim", Type: TypeString, Required: true,
			Description: "Nomor Induk Mahasiswa. Harus unik, belum pernah dipakai anggota lain.",
			Examples:    [2]string{"1301213045", "1301213046"}},
		{Key: "email", Type: TypeString,
			Description: "Alamat email aktif. Kosongkan bila tidak ada.",
			Examples:    [2]string{"aulia@student.telkomuniversity.ac.id", ""}},
		{Key: "phone", Type: TypeString,
			Description: "Nomor telepon atau WhatsApp.",
			Examples:    [2]string{"081234567890", ""}},
		{Key: "bio", Type: TypeString,
			Description: "Deskripsi singkat anggota.",
			Examples:    [2]string{"Anggota divisi Media sejak 2024.", ""}},
		{Key: "instagram_url", Type: TypeString,
			Description: "Tautan lengkap profil Instagram, diawali https://.",
			Examples:    [2]string{"https://instagram.com/aulia", ""}},
		{Key: "faculty", Type: TypeString,
			Description: "Fakultas anggota.",
			Examples:    [2]string{"Fakultas Informatika", "Fakultas Rekayasa Industri"}},
		{Key: "study_program", Type: TypeString,
			Description: "Program studi anggota.",
			Examples:    [2]string{"S1 Informatika", "S1 Sistem Informasi"}},
		{Key: "cohort_year", Type: TypeInt,
			Description: "Tahun angkatan masuk kuliah, empat digit. Kosongkan bila tidak diketahui.",
			Examples:    [2]string{"2021", "2022"}},
		{Key: "period_start", Type: TypeInt, Required: true,
			Description: "Tahun mulai kepengurusan, empat digit.",
			Examples:    [2]string{"2024", "2025"}},
		{Key: "period_end", Type: TypeInt, Required: true,
			Description: "Tahun akhir kepengurusan, empat digit.",
			Examples:    [2]string{"2025", "2026"}},
		{Key: "is_active", Type: TypeBool,
			Description: "Status keanggotaan. Dikosongkan berarti aktif.",
			Examples:    [2]string{"ya", "ya"}},
	},
}

var divisionsSpec = EntitySpec{
	Key:   "divisions",
	Title: "Template Import Data Divisi",
	Notes: []string{
		"Slug dipakai sebagai penanda unik divisi dan tidak boleh sama dengan divisi yang sudah ada.",
		"Kolom coordinator_nim merujuk anggota yang sudah terdaftar. Impor data anggota terlebih dahulu bila belum ada.",
		"Ikon divisi tidak bisa diimpor lewat berkas ini. Unggah menyusul lewat endpoint upload ikon divisi.",
	},
	Columns: []ColumnSpec{
		{Key: "name", Type: TypeString, Required: true,
			Description: "Nama divisi.",
			Examples:    [2]string{"Media dan Informasi", "Hubungan Masyarakat"}},
		{Key: "slug", Type: TypeString, Required: true,
			Description: "Penanda unik untuk URL. Huruf kecil, kata dipisah tanda hubung.",
			Examples:    [2]string{"media-informasi", "humas"}},
		{Key: "subtitle", Type: TypeString,
			Description: "Kalimat singkat yang tampil pada hero halaman detail divisi.",
			Examples:    [2]string{"Wajah KMH di dunia digital.", ""}},
		{Key: "description", Type: TypeString,
			Description: "Penjelasan tugas dan cakupan kerja divisi.",
			Examples:    [2]string{"Mengelola konten dan publikasi KMH.", ""}},
		{Key: "responsibilities", Type: TypeString,
			Description: "Daftar tanggung jawab divisi. Pisahkan antar butir dengan tanda titik koma (;).",
			Examples:    [2]string{"Mengelola media sosial KMH; Membuat konten publikasi kegiatan", ""}},
		{Key: "programs", Type: TypeString,
			Description: "Daftar program kerja. Tulis tiap program dengan format Nama | Deskripsi, lalu pisahkan antar program dengan tanda titik koma (;).",
			Examples:    [2]string{"KMH Podcast | Podcast bulanan seputar kegiatan KMH; Kelas Desain | Pelatihan desain untuk anggota", ""}},
		{Key: "coordinator_nim", Type: TypeString,
			Description: "NIM anggota yang menjadi koordinator. Anggota harus sudah terdaftar. Kosongkan bila belum ditentukan.",
			Examples:    [2]string{"1301213045", ""}},
		{Key: "is_active", Type: TypeBool,
			Description: "Status divisi. Dikosongkan berarti aktif.",
			Examples:    [2]string{"ya", "ya"}},
	},
}

var memberDivisionsSpec = EntitySpec{
	Key:   "member-divisions",
	Title: "Template Import Penugasan Anggota ke Divisi",
	Notes: []string{
		"Anggota dan divisi harus sudah terdaftar sebelum berkas ini diimpor.",
		"Satu anggota hanya boleh punya satu baris untuk divisi yang sama.",
	},
	Columns: []ColumnSpec{
		{Key: "member_nim", Type: TypeString, Required: true,
			Description: "NIM anggota yang ditugaskan. Harus sudah terdaftar.",
			Examples:    [2]string{"1301213045", "1301213046"}},
		{Key: "division_slug", Type: TypeString, Required: true,
			Description: "Slug divisi tujuan. Harus sudah terdaftar.",
			Examples:    [2]string{"media-informasi", "humas"}},
		{Key: "role_title", Type: TypeString,
			Description: "Jabatan di divisi, misalnya Ketua, Wakil, atau Staff.",
			Examples:    [2]string{"Ketua", "Staff"}},
	},
}

var eventsSpec = EntitySpec{
	Key:   "events",
	Title: "Template Import Data Kegiatan",
	Notes: []string{
		"Slug dipakai sebagai penanda unik kegiatan dan tidak boleh sama dengan kegiatan yang sudah ada.",
		"Gambar sampul tidak bisa diimpor lewat berkas ini. Unggah menyusul lewat endpoint media, lalu perbarui kegiatan.",
	},
	Columns: []ColumnSpec{
		{Key: "title", Type: TypeString, Required: true,
			Description: "Judul kegiatan.",
			Examples:    [2]string{"Malam Keakraban 2025", "Seminar Karier Teknologi"}},
		{Key: "slug", Type: TypeString, Required: true,
			Description: "Penanda unik untuk URL. Huruf kecil, kata dipisah tanda hubung.",
			Examples:    [2]string{"makrab-2025", "seminar-karier-teknologi"}},
		{Key: "description", Type: TypeString,
			Description: "Penjelasan kegiatan.",
			Examples:    [2]string{"Kegiatan pengakraban anggota baru KMH.", ""}},
		{Key: "event_type", Type: TypeString,
			Description: "Diisi internal atau external.",
			Examples:    [2]string{"internal", "external"}},
		{Key: "start_time", Type: TypeDateTime, Required: true,
			Description: "Waktu mulai. Format YYYY-MM-DD HH:MM.",
			Examples:    [2]string{"2025-11-02 09:00", "2025-12-10 13:30"}},
		{Key: "end_time", Type: TypeDateTime,
			Description: "Waktu selesai. Format YYYY-MM-DD HH:MM.",
			Examples:    [2]string{"2025-11-02 17:00", ""}},
		{Key: "location", Type: TypeString,
			Description: "Nama tempat kegiatan.",
			Examples:    [2]string{"Gedung Serbaguna Telkom University", ""}},
		{Key: "google_maps_url", Type: TypeString,
			Description: "Tautan lokasi Google Maps, diawali https://.",
			Examples:    [2]string{"https://maps.app.goo.gl/contoh", ""}},
		{Key: "registration_url", Type: TypeString,
			Description: "Tautan pendaftaran peserta, diawali https://.",
			Examples:    [2]string{"https://bit.ly/daftar-makrab", ""}},
		{Key: "status", Type: TypeString,
			Description: "Diisi upcoming, ongoing, atau finished.",
			Examples:    [2]string{"upcoming", "upcoming"}},
		{Key: "is_published", Type: TypeBool,
			Description: "Tampilkan di situs publik. Dikosongkan berarti tidak.",
			Examples:    [2]string{"tidak", "ya"}},
	},
}

var rolesSpec = EntitySpec{
	Key:   "roles",
	Title: "Template Import Data Role Pengguna",
	Notes: []string{
		"Nama role harus unik. Baris dengan nama yang sudah ada akan ditolak.",
	},
	Columns: []ColumnSpec{
		{Key: "name", Type: TypeString, Required: true,
			Description: "Nama role, misalnya admin atau editor.",
			Examples:    [2]string{"editor", "viewer"}},
		{Key: "description", Type: TypeString,
			Description: "Penjelasan kewenangan role.",
			Examples:    [2]string{"Dapat mengelola konten tanpa mengubah pengguna.", ""}},
	},
}

// specsByKey memetakan nama entitas pada URL ke definisi templatenya.
var specsByKey = map[string]EntitySpec{
	membersSpec.Key:         membersSpec,
	divisionsSpec.Key:       divisionsSpec,
	memberDivisionsSpec.Key: memberDivisionsSpec,
	eventsSpec.Key:          eventsSpec,
	rolesSpec.Key:           rolesSpec,
}

// SpecFor mengembalikan definisi template untuk sebuah entitas.
func SpecFor(key string) (EntitySpec, bool) {
	spec, found := specsByKey[key]
	return spec, found
}

// AllSpecs mengembalikan seluruh definisi template, dipakai saat membuat berkas
// template statis di docs/templates.
func AllSpecs() []EntitySpec {
	return []EntitySpec{membersSpec, divisionsSpec, memberDivisionsSpec, eventsSpec, rolesSpec}
}
