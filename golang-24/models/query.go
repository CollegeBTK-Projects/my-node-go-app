package models

import "errors"

type UserQuery struct {
	Page      *int   `form:"page"`
	Limit     *int   `form:"limit"`
	Sort      string `form:"sort"       binding:"omitempty,oneof=name -name email -email age -age created_at -created_at"`
	GroupName string `form:"group_name" binding:"omitempty,group"`
	Course    *int   `form:"course"`
	AgeMin    *int   `form:"age_min"`
	AgeMax    *int   `form:"age_max"`
	Search    string `form:"search"     binding:"omitempty,max=100"`
}

func (q UserQuery) Validate() error {
	if q.Page != nil && *q.Page < 1 {
		return errors.New("page должен быть целым числом >= 1")
	}
	if q.Limit != nil && (*q.Limit < 1 || *q.Limit > 100) {
		return errors.New("limit должен быть от 1 до 100")
	}
	if q.Course != nil && (*q.Course < 1 || *q.Course > 4) {
		return errors.New("course должен быть от 1 до 4")
	}
	if q.AgeMin != nil && *q.AgeMin < 0 {
		return errors.New("age_min не может быть отрицательным")
	}
	if q.AgeMin != nil && q.AgeMax != nil && *q.AgeMin > *q.AgeMax {
		return errors.New("age_min не может быть больше age_max")
	}
	return nil
}

func (q UserQuery) PageLimit() (page, limit int) {
	page, limit = 1, 10
	if q.Page != nil {
		page = *q.Page
	}
	if q.Limit != nil {
		limit = *q.Limit
	}
	return
}

type Meta struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Pages int64 `json:"pages"`
}

type UserList struct {
	Data []User `json:"data"`
	Meta Meta   `json:"meta"`
}

type Bucket struct {
	Key   any   `bson:"_id" json:"key"`
	Count int64 `bson:"count" json:"count"`
}

type Stats struct {
	Total    int64    `json:"total"`
	AvgAge   float64  `json:"avg_age"`
	ByGroup  []Bucket `json:"by_group"`
	ByCourse []Bucket `json:"by_course"`
}
