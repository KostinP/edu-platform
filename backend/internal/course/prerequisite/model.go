package prerequisite

import "github.com/google/uuid"

type Prerequisite struct {
	CourseID       uuid.UUID `json:"course_id"`
	PrerequisiteID uuid.UUID `json:"prerequisite_id"`
}
