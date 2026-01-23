package database

import "category-management-api/models"

// In-memory storage
var Categories []models.Category
var NextID = 1

// Helper functions
func GetAllCategories() []models.Category {
	return Categories
}

func GetCategoryByID(id int) (*models.Category, bool) {
	for i := range Categories {
		if Categories[i].ID == id {
			return &Categories[i], true
		}
	}
	return nil, false
}

func AddCategory(cat models.Category) models.Category {
	cat.ID = NextID
	NextID++
	Categories = append(Categories, cat)
	return cat
}

func UpdateCategory(id int, cat models.Category) (*models.Category, bool) {
	for i := range Categories {
		if Categories[i].ID == id {
			Categories[i].Name = cat.Name
			Categories[i].Description = cat.Description
			return &Categories[i], true
		}
	}
	return nil, false
}

func DeleteCategory(id int) bool {
	for i := range Categories {
		if Categories[i].ID == id {
			// Hapus kategori dari slice
			Categories = append(Categories[:i], Categories[i+1:]...)
			return true
		}
	}
	return false
}
