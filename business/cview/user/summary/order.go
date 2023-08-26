package usersummary

import "github.com/qcbit/service/business/data/order"

// DefaultOrderBy represents the default order for results.
var DefaultOrderBy = order.NewBy(OrderByUserID, order.ASC)

// Set of fields that the results can be ordered by. These are the names
// that should be used by the application layer.
const (
	OrderByUserID   = "userid"
	OrderByUserName = "username"
)
