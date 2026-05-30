package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/truongle2004/campus_marketplace/internal/database"
	"github.com/truongle2004/campus_marketplace/internal/modules/campus"
	"github.com/truongle2004/campus_marketplace/internal/modules/category"
)

func main() {
	db := database.NewDatabase()
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	domain := "hcmut.edu.vn"
	campuses := []campus.Campus{
		{Name: "HCMUT", Slug: "hcmut", Domain: &domain, Country: "Vietnam", City: "Ho Chi Minh City", IsActive: true},
		{Name: "VNU-HCM", Slug: "vnu-hcm", Country: "Vietnam", City: "Ho Chi Minh City", IsActive: true},
	}

	for i := range campuses {
		var count int64
		db.Model(&campus.Campus{}).Where("slug = ?", campuses[i].Slug).Count(&count)
		if count == 0 {
			if err := db.Create(&campuses[i]).Error; err != nil {
				log.Fatalf("seed campus: %v", err)
			}
		}
	}

	electronicsID := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	booksID := uuid.MustParse("00000000-0000-4000-8000-000000000002")

	categories := []category.Category{
		{ID: electronicsID, Name: "Electronics", Slug: "electronics", SortOrder: 1, IsActive: true},
		{ID: booksID, Name: "Books", Slug: "books", SortOrder: 2, IsActive: true},
		{Name: "Phones", Slug: "phones", ParentID: &electronicsID, SortOrder: 1, IsActive: true},
		{Name: "Laptops", Slug: "laptops", ParentID: &electronicsID, SortOrder: 2, IsActive: true},
	}

	for i := range categories {
		var count int64
		db.Model(&category.Category{}).Where("slug = ?", categories[i].Slug).Count(&count)
		if count == 0 {
			if err := db.Create(&categories[i]).Error; err != nil {
				log.Fatalf("seed category: %v", err)
			}
		}
	}

	log.Println("seed completed")
}
