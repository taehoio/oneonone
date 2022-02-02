package handler

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetRandomQuestionByCategoryId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cqr := oneononeddlv1.NewMockCategoryQuestionRecorder(ctrl)
	cqr.
		EXPECT().
		FindByCategoryId(nil, "1").
		Return([]*oneononeddlv1.CategoryQuestion{
			{
				Id:         5,
				CategoryId: 1,
				QuestionId: 2,
				CreatedAt:  timestamppb.Now(),
				UpdatedAt:  timestamppb.Now(),
			},
		}, nil).
		Times(1)

	qr := oneononeddlv1.NewMockQuestionRecorder(ctrl)
	qr.
		EXPECT().
		FindByIDs(nil, []uint64{2}).
		Return([]*oneononeddlv1.Question{
			{
				Id:        2,
				Question:  "What is your favorite food?",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
		}, nil).
		Times(1)

	ctx := context.Background()
	resp, err := GetRandomQuestionByCategoryId(nil, cqr, qr)(ctx, &oneononev1.GetRandomQuestionByCategoryIdRequest{
		CategoryId: "1",
	})
	assert.NoError(t, err)
	assert.Equal(t, &oneononev1.GetRandomQuestionByCategoryIdResponse{
		Question: &oneononev1.Question{
			Id:       "2",
			Question: "What is your favorite food?",
		},
	}, resp)
}
