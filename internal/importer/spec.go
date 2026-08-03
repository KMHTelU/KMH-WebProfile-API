// Package importer mendefinisikan format berkas import, membaca berkas
// .xlsx/.csv menjadi baris bernomor, dan menghasilkan berkas template.
//
// Package ini tidak bergantung pada Fiber maupun basis data supaya definisi
// kolom, parser, dan generator template bisa diuji secara terpisah.
package importer

import "strings"

// ColumnType menentukan cara sebuah sel ditafsirkan dan divalidasi.
type ColumnType string

const (
	TypeString   ColumnType = "teks"
	TypeInt      ColumnType = "angka"
	TypeBool     ColumnType = "ya/tidak"
	TypeDate     ColumnType = "tanggal"
	TypeDateTime ColumnType = "tanggal & jam"
)

// ColumnSpec mendeskripsikan satu kolom pada berkas import.
type ColumnSpec struct {
	// Key adalah teks header pada baris pertama sheet Data.
	Key      string
	Type     ColumnType
	Required bool
	// Description dipakai pada sheet Petunjuk.
	Description string
	// Examples mengisi baris contoh pada template.
	Examples [2]string
}

// EntitySpec mendeskripsikan satu berkas template beserta seluruh kolomnya.
type EntitySpec struct {
	// Key dipakai pada URL dan nama berkas, misalnya "members".
	Key     string
	Title   string
	Notes   []string
	Columns []ColumnSpec
}

// Headers mengembalikan urutan header sesuai definisi kolom.
func (s EntitySpec) Headers() []string {
	headers := make([]string, 0, len(s.Columns))
	for _, column := range s.Columns {
		headers = append(headers, column.Key)
	}
	return headers
}

// RequiredHeaders mengembalikan kolom yang wajib ada pada berkas unggahan.
func (s EntitySpec) RequiredHeaders() []string {
	headers := make([]string, 0)
	for _, column := range s.Columns {
		if column.Required {
			headers = append(headers, column.Key)
		}
	}
	return headers
}

// normalizeHeader menyamakan penulisan header supaya berkas yang memakai huruf
// kapital berbeda atau mengandung spasi berlebih tetap dikenali.
func normalizeHeader(header string) string {
	return strings.ToLower(strings.TrimSpace(header))
}
