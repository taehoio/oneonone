package handler

import (
	"context"
	"database/sql"
	"fmt"

	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
)

type ListQuestionsByCategoryIdHandlerFunc func(ctx context.Context, req *oneononev1.ListQuestionsByCategoryIdRequest) (*oneononev1.ListQuestionsByCategoryIdResponse, error)

func ListQuestionsByCategoryId(
	db *sql.DB,
	cqr oneononeddlv1.CategoryQuestionRecorder,
	qr oneononeddlv1.QuestionRecorder,
) ListQuestionsByCategoryIdHandlerFunc {
	return func(ctx context.Context, req *oneononev1.ListQuestionsByCategoryIdRequest) (*oneononev1.ListQuestionsByCategoryIdResponse, error) {
		ddlCategoryQuestions, err := cqr.FindByCategoryId(db, req.GetCategoryId())
		if err != nil {
			return nil, err
		}

		var questionIDs []uint64
		for _, cq := range ddlCategoryQuestions {
			questionIDs = append(questionIDs, cq.QuestionId)
		}

		ddlQustions, err := qr.FindByIDs(db, questionIDs)
		if err != nil {
			return nil, err
		}

		var idlQuestions []*oneononev1.Question
		for _, ddlQuestion := range ddlQustions {
			idlQuestion := &oneononev1.Question{
				Id:       fmt.Sprintf("%d", ddlQuestion.GetId()),
				Question: ddlQuestion.GetQuestion(),
			}
			idlQuestions = append(idlQuestions, idlQuestion)
		}

		return &oneononev1.ListQuestionsByCategoryIdResponse{
			Questions: idlQuestions,
		}, nil
	}
}
