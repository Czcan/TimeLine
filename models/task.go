package models

import "gorm.io/gorm"

type Task struct {
	ID          int            `json:"id"`
	UserID      int            `json:"-"`
	FolderID    int            `json:"folder_id"`
	Content     string         `json:"content"`
	Description string         `json:"description"`
	Status      int            `json:"status"`
	StartAt     int            `json:"start_at"`
	UpdatedAt   int            `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func CreateTask(db *gorm.DB, userID int, folderID int, content string, description string, date int) ([]Task, error) {
	task := &Task{UserID: userID, FolderID: folderID, Content: content, Description: description, StartAt: date}
	if err := db.Save(&task).Error; err != nil {
		return nil, err
	}
	tasks := TaskList(db, folderID)
	return tasks, nil
}

func DeleteTask(db *gorm.DB, id int, folderID int) ([]Task, error) {
	if err := db.Where("id = ?", id).Delete(&Task{}).Error; err != nil {
		return nil, err
	}
	tasks := TaskList(db, folderID)
	return tasks, nil
}

func TaskList(db *gorm.DB, folderID int) []Task {
	tasks := []Task{}
	db.Where("folder_id = ?", folderID).Order("start_at desc").Find(&tasks)
	return tasks
}
