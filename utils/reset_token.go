package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateResetToken menghasilkan token acak 256 bit dalam bentuk hex.
// Token inilah yang dikirim ke email pengguna.
func GenerateResetToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

// HashResetToken mengubah token menjadi digest SHA-256.
//
// Basis data hanya menyimpan digest ini, bukan token aslinya, supaya bocornya
// isi tabel tidak langsung memberi penyerang token yang bisa dipakai. SHA-256
// dipilih (bukan bcrypt) karena token sudah acak 256 bit sehingga tidak rentan
// ditebak, dan pencarian token perlu dilakukan lewat indeks.
func HashResetToken(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}
