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

func TestListQuestionsByCategoryId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cqr := oneononeddlv1.NewMockCategoryQuestionRecorder(ctrl)
	cqr.
		EXPECT().
		FindByCategoryId(nil, "0").
		Return([]*oneononeddlv1.CategoryQuestion{
			{
				Id:         0,
				CategoryId: 0,
				QuestionId: 1,
				CreatedAt:  timestamppb.Now(),
				UpdatedAt:  timestamppb.Now(),
			},
			{
				Id:         1,
				CategoryId: 0,
				QuestionId: 2,
				CreatedAt:  timestamppb.Now(),
				UpdatedAt:  timestamppb.Now(),
			},
		}, nil).
		Times(1)

	qr := oneononeddlv1.NewMockQuestionRecorder(ctrl)
	qr.
		EXPECT().
		FindByIDs(nil, []uint64{1, 2}).
		Return([]*oneononeddlv1.Question{
			{
				Id:        1,
				Question:  "What is your favorite color?",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			{
				Id:        2,
				Question:  "What is your favorite food?",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
		}, nil).
		Times(1)

	ctx := context.Background()
	resp, err := ListQuestionsByCategoryId(nil, cqr, qr)(ctx, &oneononev1.ListQuestionsByCategoryIdRequest{
		CategoryId: "0",
	})
	assert.NoError(t, err)
	assert.Equal(t, &oneononev1.ListQuestionsByCategoryIdResponse{
		Questions: []*oneononev1.Question{
			{
				Id:       "1",
				Question: "What is your favorite color?",
			},
			{
				Id:       "2",
				Question: "What is your favorite food?",
			},
		},
	}, resp)
}
