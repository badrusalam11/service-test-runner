package db

import (
	"log"
)

// Project represents a row in the tbl_project table.
type TblProjects struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"unique;not null"`
	URL  string `gorm:"not null"`
}

// LoadProjects queries the tbl_project table and returns a mapping of project name to URL.
func LoadProjects() (map[string]string, error) {
	var projects []TblProjects
	result := DB.Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}

	projectMap := make(map[string]string)
	for _, p := range projects {
		projectMap[p.Name] = p.URL
	}
	log.Printf("Loaded projects from DB: %+v", projectMap)
	return projectMap, nil
}
