package importer

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// DataSheetName adalah nama sheet yang dibaca dari berkas .xlsx. Bila sheet
// dengan nama ini tidak ada, parser memakai sheet pertama.
const DataSheetName = "Data"

// Parse membaca berkas unggahan menjadi kumpulan baris sesuai spec.
//
// Baris yang seluruh selnya kosong diabaikan, sehingga baris sisa di bawah data
// tidak dilaporkan sebagai kesalahan.
func Parse(spec EntitySpec, fileName string, content []byte, maxRows int) ([]Row, error) {
	var table [][]string
	var err error

	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".xlsx":
		table, err = readXLSX(content)
	case ".csv":
		table, err = readCSV(content)
	default:
		return nil, fmt.Errorf("format berkas tidak didukung, gunakan .xlsx atau .csv")
	}
	if err != nil {
		return nil, err
	}

	if len(table) == 0 {
		return nil, fmt.Errorf("berkas kosong, tidak ada baris header")
	}

	columnIndex, err := mapHeaders(spec, table[0])
	if err != nil {
		return nil, err
	}

	rows := make([]Row, 0, len(table)-1)
	for offset, record := range table[1:] {
		values := make(map[string]string, len(columnIndex))
		for key, index := range columnIndex {
			if index < len(record) {
				values[key] = record[index]
			}
		}

		// Nomor baris dihitung dari 1 dan memperhitungkan baris header.
		row := Row{Number: offset + 2, values: values}
		if row.IsEmpty() {
			continue
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("tidak ada baris data yang bisa diproses")
	}
	if len(rows) > maxRows {
		return nil, fmt.Errorf("jumlah baris data %d melebihi batas %d per berkas", len(rows), maxRows)
	}

	return rows, nil
}

// mapHeaders mencocokkan header berkas dengan kolom pada spec dan memastikan
// seluruh kolom wajib tersedia.
func mapHeaders(spec EntitySpec, header []string) (map[string]int, error) {
	position := make(map[string]int, len(header))
	for index, cell := range header {
		position[normalizeHeader(cell)] = index
	}

	columnIndex := make(map[string]int, len(spec.Columns))
	missing := make([]string, 0)
	for _, column := range spec.Columns {
		index, found := position[normalizeHeader(column.Key)]
		if !found {
			if column.Required {
				missing = append(missing, column.Key)
			}
			continue
		}
		columnIndex[column.Key] = index
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("kolom wajib tidak ditemukan pada berkas: %s. Unduh ulang template untuk melihat format yang benar",
			strings.Join(missing, ", "))
	}
	return columnIndex, nil
}

func readXLSX(content []byte) ([][]string, error) {
	file, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("berkas Excel tidak bisa dibaca: %w", err)
	}
	defer file.Close()

	sheet := DataSheetName
	if !hasSheet(file, sheet) {
		sheets := file.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("berkas Excel tidak memiliki sheet")
		}
		sheet = sheets[0]
	}

	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet %s: %w", sheet, err)
	}
	return rows, nil
}

func hasSheet(file *excelize.File, name string) bool {
	for _, sheet := range file.GetSheetList() {
		if sheet == name {
			return true
		}
	}
	return false
}

func readCSV(content []byte) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(stripBOM(content)))
	// Jumlah kolom dibiarkan bebas agar baris yang lebih pendek tidak menggagalkan
	// seluruh berkas; sel yang hilang diperlakukan sebagai kosong.
	reader.FieldsPerRecord = -1

	records := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("berkas CSV tidak bisa dibaca: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}

// stripBOM membuang penanda BOM yang biasa ditambahkan Excel saat menyimpan CSV.
func stripBOM(content []byte) []byte {
	return bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})
}
