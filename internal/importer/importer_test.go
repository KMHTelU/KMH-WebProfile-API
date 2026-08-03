package importer

import (
	"strings"
	"testing"
	"time"
)

// Template yang dihasilkan harus bisa dibaca kembali oleh parser. Tanpa ini,
// perubahan pada generator bisa diam-diam menghasilkan berkas yang ditolak saat
// diunggah kembali.
func TestGeneratedTemplatesRoundTrip(t *testing.T) {
	for _, spec := range AllSpecs() {
		t.Run(spec.Key, func(t *testing.T) {
			xlsx, err := BuildXLSX(spec)
			if err != nil {
				t.Fatalf("BuildXLSX: %v", err)
			}
			rows, err := Parse(spec, "template.xlsx", xlsx, 100)
			if err != nil {
				t.Fatalf("Parse xlsx: %v", err)
			}
			// Template memuat dua baris contoh, tetapi baris contoh kedua boleh
			// kosong seluruhnya pada spec yang semua kolom opsionalnya kosong.
			if len(rows) == 0 {
				t.Fatal("tidak ada baris contoh yang terbaca dari template xlsx")
			}
			if rows[0].Number != 2 {
				t.Errorf("nomor baris pertama = %d, ingin 2", rows[0].Number)
			}

			csv, err := BuildCSV(spec)
			if err != nil {
				t.Fatalf("BuildCSV: %v", err)
			}
			if _, err := Parse(spec, "template.csv", csv, 100); err != nil {
				t.Fatalf("Parse csv: %v", err)
			}
		})
	}
}

func TestParseRejectsMissingRequiredColumn(t *testing.T) {
	content := "name,email\nAulia,aulia@example.com\n"
	_, err := Parse(membersSpec, "data.csv", []byte(content), 100)
	if err == nil {
		t.Fatal("ingin error karena kolom nim tidak ada, tetapi berhasil")
	}
	if !strings.Contains(err.Error(), "nim") {
		t.Errorf("pesan error harus menyebut kolom yang hilang, dapat: %v", err)
	}
}

func TestParseIgnoresBlankRowsAndTrailingSpaces(t *testing.T) {
	content := " NAME ,nim,period_start,period_end\n" +
		"Aulia,1301213045,2024,2025\n" +
		",,,\n" +
		"Bunga,1301213046,2024,2025\n"

	rows, err := Parse(membersSpec, "data.csv", []byte(content), 100)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("jumlah baris = %d, ingin 2", len(rows))
	}
	// Header ditulis " NAME " namun harus tetap dikenali sebagai kolom name.
	if got := rows[0].String("name"); got != "Aulia" {
		t.Errorf("name baris pertama = %q, ingin %q", got, "Aulia")
	}
	// Baris kosong dilewati, sehingga baris kedua berasal dari baris ke-4 berkas.
	if rows[1].Number != 4 {
		t.Errorf("nomor baris kedua = %d, ingin 4", rows[1].Number)
	}
}

func TestParseRejectsRowsBeyondLimit(t *testing.T) {
	var builder strings.Builder
	builder.WriteString("name,nim,period_start,period_end\n")
	for i := 0; i < 5; i++ {
		builder.WriteString("Aulia,130121304,2024,2025\n")
	}

	_, err := Parse(membersSpec, "data.csv", []byte(builder.String()), 3)
	if err == nil {
		t.Fatal("ingin error karena melebihi batas baris, tetapi berhasil")
	}
}

func TestRowInt32(t *testing.T) {
	row := Row{Number: 2, values: map[string]string{
		"kosong":  "",
		"bulat":   "2025",
		"desimal": "2025.0",
		"huruf":   "dua ribu",
	}}

	if got, err := row.Int32("kosong"); err != nil || got != 0 {
		t.Errorf("sel kosong = (%d, %v), ingin (0, nil)", got, err)
	}
	if got, err := row.Int32("bulat"); err != nil || got != 2025 {
		t.Errorf("sel bulat = (%d, %v), ingin (2025, nil)", got, err)
	}
	// Excel kerap menuliskan bilangan bulat dengan akhiran .0
	if got, err := row.Int32("desimal"); err != nil || got != 2025 {
		t.Errorf("sel desimal = (%d, %v), ingin (2025, nil)", got, err)
	}
	if _, err := row.Int32("huruf"); err == nil {
		t.Error("sel berisi huruf harus menghasilkan error")
	}
}

func TestRowBool(t *testing.T) {
	row := Row{Number: 2, values: map[string]string{
		"kosong": "",
		"ya":     "Ya",
		"tidak":  "TIDAK",
		"satu":   "1",
		"aneh":   "mungkin",
	}}

	if got, _ := row.Bool("kosong", true); !got {
		t.Error("sel kosong harus memakai nilai default")
	}
	if got, _ := row.Bool("ya", false); !got {
		t.Error("\"Ya\" harus dibaca sebagai true")
	}
	if got, _ := row.Bool("tidak", true); got {
		t.Error("\"TIDAK\" harus dibaca sebagai false")
	}
	if got, _ := row.Bool("satu", false); !got {
		t.Error("\"1\" harus dibaca sebagai true")
	}
	if _, err := row.Bool("aneh", false); err == nil {
		t.Error("nilai tak dikenal harus menghasilkan error")
	}
}

func TestRowTime(t *testing.T) {
	row := Row{Number: 2, values: map[string]string{
		"kosong": "",
		"tanggal": "2025-11-02",
		"jam":     "2025-11-02 09:30",
		"serial":  "45963",
		"rusak":   "02 Nov",
	}}

	if got, err := row.Time("kosong"); err != nil || !got.IsZero() {
		t.Errorf("sel kosong = (%v, %v), ingin waktu nol", got, err)
	}

	parsed, err := row.Time("tanggal")
	if err != nil {
		t.Fatalf("Time tanggal: %v", err)
	}
	if parsed.Year() != 2025 || parsed.Month() != time.November || parsed.Day() != 2 {
		t.Errorf("tanggal = %v, ingin 2025-11-02", parsed)
	}

	parsed, err = row.Time("jam")
	if err != nil {
		t.Fatalf("Time jam: %v", err)
	}
	if parsed.Hour() != 9 || parsed.Minute() != 30 {
		t.Errorf("jam = %v, ingin 09:30", parsed)
	}

	// Sel bertipe tanggal asli Excel terbaca sebagai angka seri hari.
	if _, err := row.Time("serial"); err != nil {
		t.Errorf("angka seri Excel harus bisa dibaca: %v", err)
	}

	if _, err := row.Time("rusak"); err == nil {
		t.Error("format tanggal tak dikenal harus menghasilkan error")
	}
}
