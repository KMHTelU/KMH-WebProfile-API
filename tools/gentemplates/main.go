// Command gentemplates menulis berkas template import ke docs/templates.
//
// Berkas hasil perintah ini identik dengan yang dihasilkan endpoint
// /api/protected/import/{entity}/template, dan disimpan di repositori agar bisa
// dibagikan tanpa harus login lebih dulu. Jalankan ulang setiap kali definisi
// kolom di internal/importer berubah:
//
//	go run ./tools/gentemplates
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/importer"
)

func main() {
	outputDir := filepath.Join("docs", "templates")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "gagal membuat folder %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	for _, spec := range importer.AllSpecs() {
		xlsx, err := importer.BuildXLSX(spec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gagal membuat template xlsx %s: %v\n", spec.Key, err)
			os.Exit(1)
		}
		xlsxPath := filepath.Join(outputDir, fmt.Sprintf("template_import_%s.xlsx", spec.Key))
		if err := os.WriteFile(xlsxPath, xlsx, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "gagal menulis %s: %v\n", xlsxPath, err)
			os.Exit(1)
		}

		csv, err := importer.BuildCSV(spec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gagal membuat template csv %s: %v\n", spec.Key, err)
			os.Exit(1)
		}
		csvPath := filepath.Join(outputDir, fmt.Sprintf("template_import_%s.csv", spec.Key))
		if err := os.WriteFile(csvPath, csv, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "gagal menulis %s: %v\n", csvPath, err)
			os.Exit(1)
		}

		fmt.Printf("dibuat: %s dan %s\n", xlsxPath, csvPath)
	}
}
