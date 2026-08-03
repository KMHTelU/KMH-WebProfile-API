package mailer

import (
	"bytes"
	"embed"
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"
)

//go:embed templates/*.html templates/*.txt
var templateFiles embed.FS

var (
	htmlTemplates = htmltemplate.Must(htmltemplate.ParseFS(templateFiles, "templates/*.html"))
	textTemplates = texttemplate.Must(texttemplate.ParseFS(templateFiles, "templates/*.txt"))
)

// ResetPasswordData mengisi email permintaan reset password.
type ResetPasswordData struct {
	AppName          string
	UserName         string
	ResetURL         string
	ExpiresInMinutes int
}

// PasswordChangedData mengisi email notifikasi password berhasil diubah.
type PasswordChangedData struct {
	AppName   string
	UserName  string
	ChangedAt string
}

// BuildResetPassword menyusun email berisi tautan reset password.
func (m *Mailer) BuildResetPassword(to string, data ResetPasswordData) (Message, error) {
	data.AppName = m.config.AppName
	return m.build(to, fmt.Sprintf("Atur Ulang Password %s", m.config.AppName), "reset_password", data)
}

// BuildPasswordChanged menyusun email pemberitahuan bahwa password telah diubah.
func (m *Mailer) BuildPasswordChanged(to string, data PasswordChangedData) (Message, error) {
	data.AppName = m.config.AppName
	return m.build(to, fmt.Sprintf("Password %s Anda Telah Diubah", m.config.AppName), "password_changed", data)
}

func (m *Mailer) build(to, subject, templateName string, data interface{}) (Message, error) {
	var htmlBody bytes.Buffer
	if err := htmlTemplates.ExecuteTemplate(&htmlBody, templateName+".html", data); err != nil {
		return Message{}, fmt.Errorf("gagal menyusun email HTML %s: %w", templateName, err)
	}

	var textBody bytes.Buffer
	if err := textTemplates.ExecuteTemplate(&textBody, templateName+".txt", data); err != nil {
		return Message{}, fmt.Errorf("gagal menyusun email teks %s: %w", templateName, err)
	}

	return Message{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody.String(),
		TextBody: textBody.String(),
	}, nil
}
