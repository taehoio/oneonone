package handler

import (
	"context"
	"database/sql"
	"fmt"

	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
)

type GetRandomQuestionByCategoryIdHandlerFunc func(ctx context.Context, req *oneononev1.GetRandomQuestionByCategoryIdRequest) (*oneononev1.GetRandomQuestionByCategoryIdResponse, error)

func GetRandomQuestionByCategoryId(
	db *sql.DB,
	cqr oneononeddlv1.CategoryQuestionRecorder,
	qr oneononeddlv1.QuestionRecorder,
) GetRandomQuestionByCategoryIdHandlerFunc {
	return func(ctx context.Context, req *oneononev1.GetRandomQuestionByCategoryIdRequest) (*oneononev1.GetRandomQuestionByCategoryIdResponse, error) {
		ddlCategoryQuestions, err := cqr.FindByCategoryId(db, req.CategoryId)
		if err != nil {
			return nil, err
		}

		questionIDs := listQuestionIDs(ddlCategoryQuestions)
		ddlQuestions, err := qr.FindByIDs(db, questionIDs)
		if err != nil {
			return nil, err
		}

		ddlQuestion, err := getRandomQuestion(ddlQuestions)
		if err != nil {
			return nil, err
		}

		return &oneononev1.GetRandomQuestionByCategoryIdResponse{
			Question: &oneononev1.Question{
				Id:       fmt.Sprintf("%d", ddlQuestion.Id),
				Question: ddlQuestion.Question,
			},
		}, nil
	}
}
