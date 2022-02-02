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

func TestGetQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qr := oneononeddlv1.NewMockQuestionRecorder(ctrl)
	qr.
		EXPECT().
		Get(nil, uint64(2)).
		Return(&oneononeddlv1.Question{
			Id:        2,
			Question:  "What is your favorite food?",
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		}, nil).
		Times(1)

	ctx := context.Background()
	resp, err := GetQuestion(nil, qr)(ctx, &oneononev1.GetQuestionRequest{
		Id: "2",
	})
	assert.NoError(t, err)
	assert.Equal(t, &oneononev1.GetQuestionResponse{
		Question: &oneononev1.Question{
			Id:       "2",
			Question: "What is your favorite food?",
		},
	}, resp)
}
