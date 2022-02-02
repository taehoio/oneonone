package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
)

var (
	ErrInvalidQuestionID = status.Error(codes.InvalidArgument, "invalid question id")
	ErrNotFound          = status.Error(codes.NotFound, "not found")
)

type GetQuestionHandlerFunc func(ctx context.Context, req *oneononev1.GetQuestionRequest) (*oneononev1.GetQuestionResponse, error)

func GetQuestion(
	db *sql.DB,
	qr oneononeddlv1.QuestionRecorder,
) GetQuestionHandlerFunc {
	return func(ctx context.Context, req *oneononev1.GetQuestionRequest) (*oneononev1.GetQuestionResponse, error) {
		questionID, err := strconv.ParseUint(req.GetId(), 10, 64)
		if err != nil {
			return nil, ErrInvalidQuestionID
		}

		ddlQuestion, err := qr.Get(db, questionID)
		if err != nil {
			return nil, err
		}

		if ddlQuestion == nil {
			return nil, ErrNotFound
		}

		return &oneononev1.GetQuestionResponse{
			Question: &oneononev1.Question{
				Id:       fmt.Sprintf("%d", ddlQuestion.GetId()),
				Question: ddlQuestion.GetQuestion(),
			},
		}, nil
	}
}
