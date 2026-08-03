package mailer

import (
	"strings"
	"testing"
)

func newTestMailer() *Mailer {
	return New(Config{AppName: "KMH Tel-U"})
}

func TestBuildResetPasswordIncludesLinkInBothBodies(t *testing.T) {
	mailer := newTestMailer()
	resetURL := "https://leviathanbolu.my.id/reset-password?token=abc123"

	message, err := mailer.BuildResetPassword("admin@example.com", ResetPasswordData{
		UserName:         "Aulia",
		ResetURL:         resetURL,
		ExpiresInMinutes: 60,
	})
	if err != nil {
		t.Fatalf("BuildResetPassword: %v", err)
	}

	if message.To != "admin@example.com" {
		t.Errorf("To = %q", message.To)
	}
	if !strings.Contains(message.Subject, "KMH Tel-U") {
		t.Errorf("subjek harus memuat nama aplikasi, dapat %q", message.Subject)
	}

	// Penerima yang memakai klien email teks harus tetap bisa melihat tautannya.
	for name, body := range map[string]string{"HTML": message.HTMLBody, "teks": message.TextBody} {
		if !strings.Contains(body, resetURL) {
			t.Errorf("body %s tidak memuat tautan reset", name)
		}
		if !strings.Contains(body, "Aulia") {
			t.Errorf("body %s tidak memuat nama penerima", name)
		}
		if !strings.Contains(body, "60") {
			t.Errorf("body %s tidak menyebut masa berlaku", name)
		}
	}
}

// Template harus tetap terbaca ketika nama pengguna kosong.
func TestBuildResetPasswordWithoutUserName(t *testing.T) {
	message, err := newTestMailer().BuildResetPassword("admin@example.com", ResetPasswordData{
		ResetURL:         "https://leviathanbolu.my.id/reset-password?token=abc",
		ExpiresInMinutes: 60,
	})
	if err != nil {
		t.Fatalf("BuildResetPassword: %v", err)
	}
	if !strings.Contains(message.TextBody, "Halo,") {
		t.Errorf("sapaan tanpa nama harus berbunyi \"Halo,\", dapat:\n%s", message.TextBody)
	}
}

func TestBuildPasswordChanged(t *testing.T) {
	message, err := newTestMailer().BuildPasswordChanged("admin@example.com", PasswordChangedData{
		UserName:  "Aulia",
		ChangedAt: "2 Agustus 2026 23:10 WIB",
	})
	if err != nil {
		t.Fatalf("BuildPasswordChanged: %v", err)
	}
	for name, body := range map[string]string{"HTML": message.HTMLBody, "teks": message.TextBody} {
		if !strings.Contains(body, "2 Agustus 2026 23:10 WIB") {
			t.Errorf("body %s tidak memuat waktu perubahan", name)
		}
	}
}

// Tanpa SMTP_HOST, Send tidak boleh gagal: pengembangan lokal berjalan tanpa
// kredensial dan isi email cukup dicatat di log.
func TestSendWithoutSMTPHostIsNoOp(t *testing.T) {
	if err := newTestMailer().Send(Message{To: "admin@example.com", Subject: "tes"}); err != nil {
		t.Errorf("Send tanpa SMTP_HOST harus berhasil, dapat: %v", err)
	}
}
