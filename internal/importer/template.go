package importer

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// GuideSheetName adalah nama sheet berisi penjelasan tiap kolom.
const GuideSheetName = "Petunjuk"

// BuildXLSX menghasilkan berkas template .xlsx berisi sheet Data dan Petunjuk.
func BuildXLSX(spec EntitySpec) ([]byte, error) {
	file := excelize.NewFile()
	defer file.Close()

	// Sheet bawaan diganti nama agar berkas hanya memuat sheet yang relevan.
	if err := file.SetSheetName(file.GetSheetName(0), DataSheetName); err != nil {
		return nil, err
	}
	if _, err := file.NewSheet(GuideSheetName); err != nil {
		return nil, err
	}

	headerStyle, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"B91C1C"}},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}
	requiredStyle, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"7F1D1D"}},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	if err := writeDataSheet(file, spec, headerStyle, requiredStyle); err != nil {
		return nil, err
	}
	if err := writeGuideSheet(file, spec, headerStyle); err != nil {
		return nil, err
	}

	file.SetActiveSheet(0)

	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func writeDataSheet(file *excelize.File, spec EntitySpec, headerStyle, requiredStyle int) error {
	for index, column := range spec.Columns {
		cell, err := excelize.CoordinatesToCellName(index+1, 1)
		if err != nil {
			return err
		}
		if err := file.SetCellValue(DataSheetName, cell, column.Key); err != nil {
			return err
		}

		style := headerStyle
		if column.Required {
			style = requiredStyle
		}
		if err := file.SetCellStyle(DataSheetName, cell, cell, style); err != nil {
			return err
		}

		// Komentar sel menaruh penjelasan tepat di tempat pengguna mengetik.
		note := column.Description
		if column.Required {
			note = "WAJIB DIISI. " + note
		}
		if err := file.AddComment(DataSheetName, excelize.Comment{
			Cell:   cell,
			Author: "KMH",
			Paragraph: []excelize.RichTextRun{
				{Text: column.Key + "\n", Font: &excelize.Font{Bold: true}},
				{Text: note},
			},
		}); err != nil {
			return err
		}

		width := float64(len(column.Key)) + 6
		if width < 16 {
			width = 16
		}
		if width > 40 {
			width = 40
		}
		columnName, err := excelize.ColumnNumberToName(index + 1)
		if err != nil {
			return err
		}
		if err := file.SetColWidth(DataSheetName, columnName, columnName, width); err != nil {
			return err
		}

		// Dua baris contoh membantu admin melihat format yang diharapkan.
		for exampleIndex, example := range column.Examples {
			exampleCell, err := excelize.CoordinatesToCellName(index+1, exampleIndex+2)
			if err != nil {
				return err
			}
			// Ditulis sebagai teks apa adanya supaya Excel tidak mengubah
			// tampilan tanggal dan NIM yang berawalan angka nol.
			if err := file.SetCellStr(DataSheetName, exampleCell, example); err != nil {
				return err
			}
		}
	}

	if err := file.SetPanes(DataSheetName, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return err
	}
	return nil
}

func writeGuideSheet(file *excelize.File, spec EntitySpec, headerStyle int) error {
	wrapStyle, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
	})
	if err != nil {
		return err
	}
	titleStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 13},
	})
	if err != nil {
		return err
	}

	row := 1
	set := func(cell string, value interface{}) error {
		return file.SetCellValue(GuideSheetName, cell, value)
	}

	if err := set(fmt.Sprintf("A%d", row), spec.Title); err != nil {
		return err
	}
	if err := file.SetCellStyle(GuideSheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), titleStyle); err != nil {
		return err
	}
	row += 2

	defaultNotes := []string{
		"Isi data mulai baris ke-2 pada sheet " + DataSheetName + ". Baris contoh yang sudah ada boleh dihapus atau ditimpa.",
		"Jangan mengubah, menghapus, atau menukar urutan baris header.",
		"Baris yang seluruh selnya kosong akan dilewati.",
		"Kolom bertanda WAJIB harus diisi; baris yang mengosongkannya akan gagal.",
		"Gunakan endpoint import/preview terlebih dahulu untuk memeriksa berkas tanpa menyimpan apa pun.",
		"Bila satu baris gagal, baris lain tetap diproses. Laporan hasil menyebut nomor baris yang bermasalah.",
		"Format tanggal yang diterima: YYYY-MM-DD atau YYYY-MM-DD HH:MM.",
		"Kolom ya/tidak menerima: ya, tidak, true, false, 1, 0.",
	}
	for _, note := range append(defaultNotes, spec.Notes...) {
		if err := set(fmt.Sprintf("A%d", row), "- "+note); err != nil {
			return err
		}
		row++
	}
	row++

	headers := []string{"Kolom", "Wajib", "Tipe", "Penjelasan"}
	for index, header := range headers {
		cell, err := excelize.CoordinatesToCellName(index+1, row)
		if err != nil {
			return err
		}
		if err := set(cell, header); err != nil {
			return err
		}
		if err := file.SetCellStyle(GuideSheetName, cell, cell, headerStyle); err != nil {
			return err
		}
	}
	row++

	for _, column := range spec.Columns {
		required := "opsional"
		if column.Required {
			required = "WAJIB"
		}
		values := []string{column.Key, required, string(column.Type), column.Description}
		for index, value := range values {
			cell, err := excelize.CoordinatesToCellName(index+1, row)
			if err != nil {
				return err
			}
			if err := set(cell, value); err != nil {
				return err
			}
		}
		lastCell, err := excelize.CoordinatesToCellName(4, row)
		if err != nil {
			return err
		}
		if err := file.SetCellStyle(GuideSheetName, lastCell, lastCell, wrapStyle); err != nil {
			return err
		}
		row++
	}

	widths := map[string]float64{"A": 24, "B": 12, "C": 16, "D": 80}
	for column, width := range widths {
		if err := file.SetColWidth(GuideSheetName, column, column, width); err != nil {
			return err
		}
	}
	return nil
}

// BuildCSV menghasilkan template .csv berisi header dan baris contoh.
// Penjelasan kolom tidak ikut karena CSV tidak mengenal sheet tambahan; gunakan
// versi .xlsx bila ingin membaca petunjuk lengkap.
func BuildCSV(spec EntitySpec) ([]byte, error) {
	var buffer bytes.Buffer
	// BOM ditulis agar Excel membuka berkas dengan encoding UTF-8.
	buffer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buffer)
	if err := writer.Write(spec.Headers()); err != nil {
		return nil, err
	}
	for exampleIndex := 0; exampleIndex < 2; exampleIndex++ {
		record := make([]string, 0, len(spec.Columns))
		for _, column := range spec.Columns {
			record = append(record, column.Examples[exampleIndex])
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
