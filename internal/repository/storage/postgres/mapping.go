package postgres

import "github.com/evgeniy-dammer/marketplace-api/pkg/columncode"

var mappingSortUser = map[columncode.ColumnCode]string{
	"id":         "id",
	"phone":      "phone",
	"first_name": "first_name",
	"last_name":  "last_name",
	"status_id":  "status_id",
	"created_at": "created_at",
}

var mappingSortRole = map[columncode.ColumnCode]string{
	"id":   "id",
	"name": "name",
}

var mappingSortOrganization = map[columncode.ColumnCode]string{
	"id":         "id",
	"name":       "name",
	"user_id":    "address",
	"address":    "address",
	"phone":      "phone",
	"created_at": "created_at",
}

var mappingSortItem = map[columncode.ColumnCode]string{
	"id":              "id",
	"name_tm":         "name_tm",
	"name_ru":         "name_ru",
	"name_tr":         "name_tr",
	"name_en":         "name_en",
	"description_tm":  "description_tm",
	"description_ru":  "description_ru",
	"description_tr":  "description_tr",
	"description_en":  "description_en",
	"internal_id":     "internal_id",
	"price":           "price",
	"rating":          "rating",
	"comments_qty":    "comments_qty",
	"category_id":     "category_id",
	"organization_id": "organization_id",
	"brand_id":        "brand_id",
	"created_at":      "created_at",
}

var mappingSortTable = map[columncode.ColumnCode]string{
	"id":         "id",
	"name":       "name",
	"user_id":    "address",
	"address":    "address",
	"phone":      "phone",
	"created_at": "created_at",
}

var mappingSortImage = map[columncode.ColumnCode]string{
	"id":              "id",
	"object_id":       "object_id",
	"type":            "type",
	"origin":          "origin",
	"middle":          "middle",
	"small":           "small",
	"organization_id": "organization_id",
	"is_main":         "is_main",
}

var mappingSortCategory = map[columncode.ColumnCode]string{
	"id":              "id",
	"name_tm":         "name_tm",
	"name_ru":         "name_ru",
	"name_tr":         "name_tr",
	"name_en":         "name_en",
	"parent_id":       "parent_id",
	"level":           "level",
	"organization_id": "organization_id",
}

var mappingSortComment = map[columncode.ColumnCode]string{
	"id":              "id",
	"item_id":         "item_id",
	"organization_id": "organization_id",
	"content":         "content",
	"status_id":       "status_id",
	"rating":          "rating",
	"user_created":    "user_created",
	"created_at":      "created_at",
}

var mappingSortSpecification = map[columncode.ColumnCode]string{
	"id":              "id",
	"item_id":         "item_id",
	"organization_id": "organization_id",
	"name_tm":         "name_tm",
	"name_ru":         "name_ru",
	"name_tr":         "name_tr",
	"name_en":         "name_en",
	"description_tm":  "description_tm",
	"description_ru":  "description_ru",
	"description_tr":  "description_tr",
	"description_en":  "description_en",
	"value":           "value",
}

var mappingSortRule = map[columncode.ColumnCode]string{
	"id":    "id",
	"ptype": "ptype",
	"v0":    "v0",
	"v1":    "v1",
	"v2":    "v2",
	"v3":    "v3",
	"v4":    "v4",
	"v5":    "v5",
}

var mappingSortOrder = map[columncode.ColumnCode]string{
	"id":              "id",
	"user_id":         "user_id",
	"organization_id": "organization_id",
	"table_id":        "table_id",
	"status_id":       "status_id",
	"totalsum":        "totalsum",
	"created_at":      "created_at",
}
