package repositories

import (
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// ── Generations ──

func (r *Repository) ListHofGenerations(c fiber.Ctx) ([]generated.HofGeneration, error) {
	return r.Queries.SelectHofGenerations(c)
}

func (r *Repository) CreateHofGeneration(params generated.InsertHofGenerationParams, c fiber.Ctx) (generated.HofGeneration, error) {
	return r.Queries.InsertHofGeneration(c, params)
}

func (r *Repository) UpdateHofGeneration(params generated.UpdateHofGenerationParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateHofGeneration(c, params)
	return err
}

func (r *Repository) DeleteHofGeneration(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteHofGeneration(c, id)
}

// ── People ──

func (r *Repository) ListHofPeople(c fiber.Ctx) ([]generated.SelectHofPeopleRow, error) {
	return r.Queries.SelectHofPeople(c)
}

func (r *Repository) CreateHofPerson(params generated.InsertHofPersonParams, c fiber.Ctx) (generated.HofPerson, error) {
	return r.Queries.InsertHofPerson(c, params)
}

func (r *Repository) UpdateHofPerson(params generated.UpdateHofPersonParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateHofPerson(c, params)
	return err
}

func (r *Repository) DeleteHofPerson(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteHofPerson(c, id)
}

// ── Achievements ──

func (r *Repository) ListHofAchievements(c fiber.Ctx) ([]generated.HofAchievement, error) {
	return r.Queries.SelectHofAchievements(c)
}

func (r *Repository) CreateHofAchievement(params generated.InsertHofAchievementParams, c fiber.Ctx) (generated.HofAchievement, error) {
	return r.Queries.InsertHofAchievement(c, params)
}

func (r *Repository) UpdateHofAchievement(params generated.UpdateHofAchievementParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateHofAchievement(c, params)
	return err
}

func (r *Repository) DeleteHofAchievement(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteHofAchievement(c, id)
}

// ── Timeline ──

func (r *Repository) ListHofTimelineEvents(c fiber.Ctx) ([]generated.HofTimelineEvent, error) {
	return r.Queries.SelectHofTimelineEvents(c)
}

func (r *Repository) CreateHofTimelineEvent(params generated.InsertHofTimelineEventParams, c fiber.Ctx) (generated.HofTimelineEvent, error) {
	return r.Queries.InsertHofTimelineEvent(c, params)
}

func (r *Repository) UpdateHofTimelineEvent(params generated.UpdateHofTimelineEventParams, c fiber.Ctx) error {
	_, err := r.Queries.UpdateHofTimelineEvent(c, params)
	return err
}

func (r *Repository) DeleteHofTimelineEvent(id uuid.UUID, c fiber.Ctx) error {
	return r.Queries.DeleteHofTimelineEvent(c, id)
}
