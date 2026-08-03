package routes

import (
	"testing"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/handlers"
	"github.com/gofiber/fiber/v3"
)

func setupTestApp(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New()
	// Handler cukup kosong: pengujian ini hanya memeriksa pendaftaran rute,
	// bukan menjalankan handler-nya.
	InitializeRoutes(&handlers.Handler{}).SetupRoutes(app)
	return app
}

func registeredPaths(app *fiber.App, method string) []string {
	paths := make([]string, 0)
	for _, route := range app.GetRoutes() {
		if route.Method == method {
			paths = append(paths, route.Path)
		}
	}
	return paths
}

func indexOf(paths []string, target string) int {
	for index, path := range paths {
		if path == target {
			return index
		}
	}
	return -1
}

// Rute literal harus terdaftar sebelum rute berparameter yang bentuknya sama,
// jika tidak "categories" dan "bulk" akan tertangkap sebagai nilai :id.
func TestLiteralRoutesRegisteredBeforeParameterRoutes(t *testing.T) {
	app := setupTestApp(t)

	cases := []struct {
		method  string
		literal string
		param   string
	}{
		{fiber.MethodGet, "/api/galleries/categories", "/api/galleries/:id"},
		{fiber.MethodPost, "/api/protected/members/bulk", "/api/protected/members/:id/upload"},
		{fiber.MethodPut, "/api/protected/members/bulk", "/api/protected/members/:id"},
		{fiber.MethodPut, "/api/protected/divisions/bulk", "/api/protected/divisions/:id"},
		{fiber.MethodPut, "/api/protected/events/bulk", "/api/protected/events/:id"},
		{fiber.MethodPut, "/api/protected/galleries/bulk", "/api/protected/galleries/:id"},
		{fiber.MethodPut, "/api/protected/roles/bulk", "/api/protected/roles/:id"},
		{fiber.MethodPut, "/api/protected/homepage-banners/bulk", "/api/protected/homepage-banners/:id"},
		{fiber.MethodPut, "/api/protected/member-divisions/bulk", "/api/protected/member-divisions/:id"},
	}

	for _, testCase := range cases {
		paths := registeredPaths(app, testCase.method)
		literalIndex := indexOf(paths, testCase.literal)
		paramIndex := indexOf(paths, testCase.param)

		if literalIndex == -1 {
			t.Errorf("%s %s belum terdaftar", testCase.method, testCase.literal)
			continue
		}
		if paramIndex == -1 {
			continue
		}
		if literalIndex > paramIndex {
			t.Errorf("%s %s terdaftar setelah %s, sehingga tidak akan pernah tercapai",
				testCase.method, testCase.literal, testCase.param)
		}
	}
}

func TestNewRoutesAreRegistered(t *testing.T) {
	app := setupTestApp(t)

	expected := map[string][]string{
		fiber.MethodGet: {
			"/api/galleries/categories",
			"/api/members/:id/divisions",
			"/api/divisions/:id/members",
			"/api/protected/import/:entity/template",
		},
		fiber.MethodPost: {
			"/api/forgot-password",
			"/api/reset-password",
			"/api/protected/auth/change-password",
			"/api/protected/member-divisions",
			"/api/protected/members/bulk",
			"/api/protected/divisions/bulk",
			"/api/protected/member-divisions/bulk",
			"/api/protected/events/bulk",
			"/api/protected/galleries/bulk",
			"/api/protected/gallery-items/bulk",
			"/api/protected/homepage-banners/bulk",
			"/api/protected/roles/bulk",
			"/api/protected/import/:entity",
			"/api/protected/import/:entity/preview",
		},
		fiber.MethodPut: {
			"/api/protected/member-divisions/:id",
			"/api/protected/gallery-items/bulk",
		},
		fiber.MethodDelete: {
			"/api/protected/member-divisions/:id",
		},
	}

	for method, paths := range expected {
		registered := registeredPaths(app, method)
		for _, path := range paths {
			if indexOf(registered, path) == -1 {
				t.Errorf("%s %s belum terdaftar", method, path)
			}
		}
	}
}
