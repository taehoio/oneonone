package handler

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
)

func TestListCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cr := oneononeddlv1.NewMockCategoryRecorder(ctrl)
	cr.
		EXPECT().
		List(nil, nil, true, int64(100)).
		Return([]*oneononeddlv1.Category{
			{
				Id:        0,
				Name:      "About Manager",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			{
				Id:        1,
				Name:      "Career development",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
		}, nil).
		Times(1)

	ctx := context.Background()
	resp, err := ListCategories(nil, cr)(ctx, &oneononev1.ListCategoriesRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, &oneononev1.ListCategoriesResponse{
		Categories: []*oneononev1.Category{
			{
				Id:   "0",
				Name: "About Manager",
			},
			{
				Id:   "1",
				Name: "Career development",
			},
		},
	}, resp)
}
