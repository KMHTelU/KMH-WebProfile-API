package importer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Row adalah satu baris data beserta nomor barisnya pada berkas asli.
// Nomor baris dipakai pada pesan error supaya admin bisa langsung membuka
// baris yang bermasalah di Excel.
type Row struct {
	Number int
	values map[string]string
}

func (r Row) String(key string) string {
	return strings.TrimSpace(r.values[key])
}

// IsEmpty menandai baris yang seluruh selnya kosong, biasanya sisa baris
// kosong di bawah data yang ikut terbaca.
func (r Row) IsEmpty() bool {
	for _, value := range r.values {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

// Int32 membaca sel sebagai bilangan bulat. Sel kosong menghasilkan 0.
func (r Row) Int32(key string) (int32, error) {
	raw := r.String(key)
	if raw == "" {
		return 0, nil
	}
	// Excel kerap menyimpan angka bulat sebagai "2025.0".
	if strings.HasSuffix(raw, ".0") {
		raw = strings.TrimSuffix(raw, ".0")
	}
	parsed, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("kolom %s harus berupa angka, ditemukan %q", key, r.String(key))
	}
	return int32(parsed), nil
}

var (
	trueValues  = map[string]bool{"true": true, "1": true, "ya": true, "yes": true, "y": true, "aktif": true, "published": true}
	falseValues = map[string]bool{"false": true, "0": true, "tidak": true, "no": true, "n": true, "nonaktif": true, "draft": true}
)

// Bool membaca sel sebagai nilai ya/tidak. Sel kosong menghasilkan defaultValue.
func (r Row) Bool(key string, defaultValue bool) (bool, error) {
	raw := strings.ToLower(r.String(key))
	if raw == "" {
		return defaultValue, nil
	}
	if trueValues[raw] {
		return true, nil
	}
	if falseValues[raw] {
		return false, nil
	}
	return false, fmt.Errorf("kolom %s harus diisi ya atau tidak, ditemukan %q", key, r.String(key))
}

// Layout tanggal yang diterima, diurutkan dari yang paling lengkap.
var timeLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"02/01/2006 15:04:05",
	"02/01/2006 15:04",
	"02/01/2006",
	"01/02/06",
}

// Time membaca sel sebagai waktu. Sel kosong menghasilkan waktu nol.
func (r Row) Time(key string) (time.Time, error) {
	raw := r.String(key)
	if raw == "" {
		return time.Time{}, nil
	}

	for _, layout := range timeLayouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed, nil
		}
	}

	// Sel bertipe tanggal asli Excel terbaca sebagai angka seri hari.
	if serial, err := strconv.ParseFloat(raw, 64); err == nil {
		if parsed, err := excelize.ExcelDateToTime(serial, false); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"kolom %s bukan tanggal yang dikenali (%q). Gunakan format YYYY-MM-DD atau YYYY-MM-DD HH:MM", key, raw)
}
