// Package mailer mengirim email transaksional lewat SMTP.
//
// Package ini sengaja tidak bergantung pada Fiber maupun database supaya isi
// dan tata letak email bisa diuji tanpa menjalankan server atau basis data.
package mailer

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/wneessen/go-mail"
)

// Config memuat seluruh pengaturan SMTP beserta identitas pengirim.
type Config struct {
	Host       string
	Port       int
	Username   string
	Password   string
	FromEmail  string
	FromName   string
	Encryption string
	AppName    string
}

// Mailer mengirim email. Bila Host kosong, Mailer berjalan dalam mode nonaktif:
// email tidak dikirim, hanya dicatat di log. Ini membuat pengembangan lokal
// tetap jalan tanpa kredensial SMTP.
type Mailer struct {
	config  Config
	enabled bool
}

func New(config Config) *Mailer {
	return &Mailer{
		config:  config,
		enabled: config.Host != "",
	}
}

// Message adalah satu email siap kirim.
type Message struct {
	To       string
	Subject  string
	HTMLBody string
	TextBody string
}

// Send mengirim email dan menunggu sampai selesai.
func (m *Mailer) Send(msg Message) error {
	if !m.enabled {
		log.Infof("SMTP nonaktif, email ke %s tidak dikirim. Subjek: %s", msg.To, msg.Subject)
		log.Infof("Isi email (plain text):\n%s", msg.TextBody)
		return nil
	}

	email := mail.NewMsg()
	if err := email.FromFormat(m.config.FromName, m.config.FromEmail); err != nil {
		return fmt.Errorf("alamat pengirim tidak valid: %w", err)
	}
	if err := email.To(msg.To); err != nil {
		return fmt.Errorf("alamat penerima tidak valid: %w", err)
	}
	email.Subject(msg.Subject)
	email.SetBodyString(mail.TypeTextPlain, msg.TextBody)
	email.AddAlternativeString(mail.TypeTextHTML, msg.HTMLBody)

	client, err := m.newClient()
	if err != nil {
		return err
	}
	if err := client.DialAndSend(email); err != nil {
		return fmt.Errorf("gagal mengirim email: %w", err)
	}
	return nil
}

// SendAsync mengirim email di goroutine terpisah supaya request HTTP tidak
// menunggu SMTP. Kegagalan hanya dicatat di log karena pemanggil sudah terlanjur
// menerima respons.
func (m *Mailer) SendAsync(msg Message) {
	go func() {
		if err := m.Send(msg); err != nil {
			log.Errorf("Gagal mengirim email ke %s (%s): %v", msg.To, msg.Subject, err)
		}
	}()
}

func (m *Mailer) newClient() (*mail.Client, error) {
	options := []mail.Option{
		mail.WithPort(m.config.Port),
		mail.WithTimeout(15 * time.Second),
	}

	switch m.config.Encryption {
	case "ssl", "tls":
		// WithSSL dipilih daripada WithSSLPort karena yang terakhir memaksa
		// port 465 dan mengabaikan SMTP_PORT yang sudah diatur.
		options = append(options, mail.WithSSL())
	case "none":
		options = append(options, mail.WithTLSPolicy(mail.NoTLS))
	default:
		options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))
	}

	// Relay tanpa autentikasi tetap didukung: kredensial hanya dipasang bila diisi.
	if m.config.Username != "" {
		options = append(options,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(m.config.Username),
			mail.WithPassword(m.config.Password),
		)
	}

	client, err := mail.NewClient(m.config.Host, options...)
	if err != nil {
		return nil, fmt.Errorf("gagal menyiapkan koneksi SMTP: %w", err)
	}
	return client, nil
}
