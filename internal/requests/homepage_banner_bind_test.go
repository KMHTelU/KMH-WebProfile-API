package requests

// Regresi bug "banner tidak bisa ditambah": decoder form Fiber memperlakukan
// titik pada nama field (mis. "data.title") sebagai path struct bersarang,
// sehingga field tidak pernah terisi. Kontrak yang benar memakai nama datar.

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func bindBanner(t *testing.T, fields map[string]string) HomepageBannerRequest {
	t.Helper()
	app := fiber.New()
	app.Post("/t", func(c fiber.Ctx) error {
		var req HomepageBannerRequest
		if err := c.Bind().Form(&req); err != nil {
			return c.Status(400).SendString("bind error: " + err.Error())
		}
		return c.JSON(req)
	})

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range fields {
		_ = w.WriteField(k, v)
	}
	fw, _ := w.CreateFormFile("media", "banner.png")
	_, _ = fw.Write([]byte("fakeimg"))
	_ = w.Close()

	httpReq := httptest.NewRequest("POST", "/t", body)
	httpReq.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := app.Test(httpReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != 200 {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("binding gagal: status=%d body=%s", resp.StatusCode, raw)
	}
	raw, _ := io.ReadAll(resp.Body)
	var out HomepageBannerRequest
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("gagal unmarshal respons: %v", err)
	}
	return out
}

func TestBannerBindFlatFields(t *testing.T) {
	out := bindBanner(t, map[string]string{
		"title":      "Judul",
		"subtitle":   "Sub",
		"is_active":  "true",
		"start_date": "2026-08-18T00:00:00.000Z",
		"end_date":   "2026-09-18T00:00:00.000Z",
	})

	if out.Title != "Judul" {
		t.Errorf("title tidak ter-bind, dapat %q", out.Title)
	}
	if out.Subtitle != "Sub" {
		t.Errorf("subtitle tidak ter-bind, dapat %q", out.Subtitle)
	}
	if !out.IsActive {
		t.Error("is_active seharusnya true")
	}
	if out.StartDate.IsZero() || out.EndDate.IsZero() {
		t.Errorf("tanggal tidak ter-parse: start=%v end=%v", out.StartDate, out.EndDate)
	}
}

// Banner nonaktif (is_active=false) harus tetap bisa dibuat.
func TestBannerBindInactive(t *testing.T) {
	out := bindBanner(t, map[string]string{
		"title":      "Judul",
		"is_active":  "false",
		"start_date": "2026-08-18T00:00:00.000Z",
		"end_date":   "2026-09-18T00:00:00.000Z",
	})
	if out.IsActive {
		t.Error("is_active seharusnya false")
	}
}
