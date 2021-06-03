package data

import (
	"fmt"

	"gorm.io/gorm"
)

type Course struct {
	Code           int32  `gorm:"code"`
	Title          string `gorm:"title"`
	DepartmentCode int32  `gorm:"department_code"`
	Description    string `gorm:"description"`
}

type CourseData struct {
	Db *gorm.DB
}

func NewCourseData(db *gorm.DB) *CourseData {
	return &CourseData{Db: db}
}

func (u CourseData) Add(course Course) (int32, error) {
	result := u.Db.Create(&course)
	if result.Error != nil {
		return -1, fmt.Errorf("can't create course, error: %w", result.Error)
	}
	return course.Code, nil
}

func (u CourseData) Read(code int32) (Course, error) {
	var course Course
	result := u.Db.Find(&course, code)
	if result.Error != nil {
		return course, fmt.Errorf("can't read course with given id, error: %w", result.Error)
	}
	return course, nil
}

func (u CourseData) ReadAll() ([]Course, error) {
	var courses []Course
	result := u.Db.Find(&courses)
	if result.Error != nil {
		return nil, fmt.Errorf("can't read courses from database, error: %w", result.Error)
	}
	return courses, nil
}

func (u CourseData) ChangeDescription(code int32, description string) (int32, error) {
	result := u.Db.Model(Course{}).Where("code = ?", code).Update("description", description)
	if result.Error != nil {
		return -1, fmt.Errorf("can't update course description, error: %w", result.Error)
	}
	return code, nil
}

func (u CourseData) Delete(code int32) error {
	result := u.Db.Delete(&Course{}, code)
	if result.Error != nil {
		return fmt.Errorf("can't delete course from database, error: %w", result.Error)
	}
	return nil
}

func (u CourseData) GetDepartmentName(code int32) (string, error) {
	var departmentName string
	result := u.Db.Model(&Course{}).
		Select("departments.name").
		Joins("join departments on department_code = departments.code").
		Where("courses.code = ?", code).
		Scan(&departmentName)
	if result.Error != nil {
		return "", fmt.Errorf("can't get department name from database, error: %w", result.Error)
	}
	return departmentName, nil
}
