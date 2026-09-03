package requests

import "github.com/google/uuid"

// DivisionProgram adalah satu program kerja divisi yang ditampilkan di halaman
// detail divisi pada situs publik.
type DivisionProgram struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

type CreateDivisionRequest struct {
	Name             string            `json:"name" validate:"required"`
	Slug             string            `json:"slug" validate:"required"`
	Subtitle         string            `json:"subtitle" validate:"omitempty"`
	Description      string            `json:"description" validate:"omitempty"`
	// Tipe divisi: internal (di bawah Wakil Ketua Internal) atau external
	// (di bawah Wakil Ketua External). Kosong dianggap internal.
	DivisionType     string            `json:"division_type" validate:"omitempty,oneof=internal external"`
	Responsibilities []string          `json:"responsibilities" validate:"omitempty,dive,required"`
	Programs         []DivisionProgram `json:"programs" validate:"omitempty,dive"`
	CoordinatorID    uuid.UUID         `json:"coordinator_id" validate:"omitempty,uuid4"`
}

type UpdateDivisionRequest struct {
	Name             string            `json:"name" validate:"omitempty"`
	Slug             string            `json:"slug" validate:"omitempty"`
	Subtitle         string            `json:"subtitle" validate:"omitempty"`
	Description      string            `json:"description" validate:"omitempty"`
	DivisionType     string            `json:"division_type" validate:"omitempty,oneof=internal external"`
	Responsibilities []string          `json:"responsibilities" validate:"omitempty,dive,required"`
	Programs         []DivisionProgram `json:"programs" validate:"omitempty,dive"`
	CoordinatorID    uuid.UUID         `json:"coordinator_id" validate:"omitempty,uuid4"`
	// ClearCoordinator melepas koordinator (coordinator_id di-set NULL).
	// Diperlukan karena coordinator_id kosong berarti "tidak diubah".
	ClearCoordinator bool  `json:"clear_coordinator" validate:"omitempty"`
	IsActive         *bool `json:"is_active" validate:"omitempty"`
}
