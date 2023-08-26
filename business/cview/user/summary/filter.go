package usersummary

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/qcbit/service/business/sys/validate"
)

// QueryFilter holds the available files a query can be filtered by.
type QueryFilter struct {
	UserID   *uuid.UUID `validate:"omitempty,uuid4"`
	UserName *string    `validate:"omitempty,min=3"`
}

// Validate checks the data in the model.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("invalid query filter: %w", err)
	}
	return nil
}

// WithUserID adds a user id to the filter.
func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.UserID = &userID
}

// WithUserName adds a user name to the filter.
func (qf *QueryFilter) WithUserName(userName string) {
	qf.UserName = &userName
}
