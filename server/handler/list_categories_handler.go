package handler

import (
	"context"
	"database/sql"
	"fmt"

	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
)

type ListCategoriesHandlerFunc func(ctx context.Context, req *oneononev1.ListCategoriesRequest) (*oneononev1.ListCategoriesResponse, error)

func ListCategories(db *sql.DB, cr oneononeddlv1.CategoryRecorder) ListCategoriesHandlerFunc {
	return func(ctx context.Context, req *oneononev1.ListCategoriesRequest) (*oneononev1.ListCategoriesResponse, error) {
		ddlCategories, err := cr.List(db, nil, true, 100)
		if err != nil {
			return nil, err
		}

		var idlCategories []*oneononev1.Category
		for _, ddlCategory := range ddlCategories {
			idlCategory := &oneononev1.Category{
				Id:   fmt.Sprintf("%d", ddlCategory.GetId()),
				Name: ddlCategory.GetName(),
			}
			idlCategories = append(idlCategories, idlCategory)
		}

		return &oneononev1.ListCategoriesResponse{
			Categories: idlCategories,
		}, nil
	}
}
