package model

import (
	"github.com/satori/uuid"
)

func NextID() string {
	id := uuid.NewV4()
	return id.String()
}
