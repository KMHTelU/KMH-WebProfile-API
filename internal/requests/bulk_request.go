package requests

import "github.com/google/uuid"

// Setiap permintaan bulk memakai bentuk yang sama: sebuah objek dengan field
// "items". Batas 500 item dijaga di tingkat validasi supaya permintaan yang
// terlalu besar ditolak sebelum menyentuh basis data.
//
// Varian update membungkus request tunggal yang sudah ada dan menambahkan "id",
// sehingga aturan validasi field tetap satu sumber.

type BulkCreateMembersRequest struct {
	Items []CreateMemberRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateMemberItem struct {
	ID uuid.UUID `json:"id" validate:"required"`
	UpdateMemberRequest
}

type BulkUpdateMembersRequest struct {
	Items []BulkUpdateMemberItem `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkCreateDivisionsRequest struct {
	Items []CreateDivisionRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateDivisionItem struct {
	ID uuid.UUID `json:"id" validate:"required"`
	UpdateDivisionRequest
}

type BulkUpdateDivisionsRequest struct {
	Items []BulkUpdateDivisionItem `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkCreateMemberDivisionsRequest struct {
	Items []CreateMemberDivisionRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateMemberDivisionItem struct {
	ID uuid.UUID `json:"id" validate:"required"`
	UpdateMemberDivisionRequest
}

type BulkUpdateMemberDivisionsRequest struct {
	Items []BulkUpdateMemberDivisionItem `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkCreateEventsRequest struct {
	Items []CreateEventRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateEventItem struct {
	ID uuid.UUID `json:"id" validate:"required"`
	UpdateEventRequest
}

type BulkUpdateEventsRequest struct {
	Items []BulkUpdateEventItem `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkCreateGalleriesRequest struct {
	Items []CreateGalleryRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateGalleryItem struct {
	ID uuid.UUID `json:"id" validate:"required"`
	UpdateGalleryRequest
}

type BulkUpdateGalleriesRequest struct {
	Items []BulkUpdateGalleryItem `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkCreateGalleryItemsRequest struct {
	Items []CreateGalleryItemRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateGalleryItemEntry struct {
	ID uuid.UUID `json:"id" validate:"required"`
	UpdateGalleryItemRequest
}

type BulkUpdateGalleryItemsRequest struct {
	Items []BulkUpdateGalleryItemEntry `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkCreateHomepageBannersRequest struct {
	Items []HomepageBannerJSONRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateHomepageBannerItem struct {
	ID uuid.UUID `json:"id" validate:"required"`
	HomepageBannerJSONRequest
}

type BulkUpdateHomepageBannersRequest struct {
	Items []BulkUpdateHomepageBannerItem `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkCreateRolesRequest struct {
	Items []CreateRoleRequest `json:"items" validate:"required,min=1,max=500,dive"`
}

type BulkUpdateRoleItem struct {
	ID uuid.UUID `json:"id" validate:"required"`
	UpdateRoleRequest
}

type BulkUpdateRolesRequest struct {
	Items []BulkUpdateRoleItem `json:"items" validate:"required,min=1,max=500,dive"`
}
