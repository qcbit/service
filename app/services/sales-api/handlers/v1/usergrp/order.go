package usergrp

import (
	"errors"
	"net/http"

	"github.com/qcbit/service/business/core/user"
	usersummary "github.com/qcbit/service/business/cview/user/summary"
	"github.com/qcbit/service/business/data/order"
	"github.com/qcbit/service/business/sys/validate"
)

var orderByFields = map[string]struct{}{
	user.OrderByID:      {},
	user.OrderByName:    {},
	user.OrderByEmail:   {},
	user.OrderByRoles:   {},
	user.OrderByEnabled: {},
}

func parseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, user.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError("orderBy", errors.New("invalid order by field"))
	}

	return orderBy, nil
}

// -----------------------------------------------------------------------------

var orderBySummaryFields = map[string]struct{}{
	usersummary.OrderByUserID:   {},
	usersummary.OrderByUserName: {},
}

func parseSummaryOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, usersummary.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderBySummaryFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}
