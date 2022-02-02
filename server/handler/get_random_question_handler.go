package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"

	oneononeddlv1 "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1"
	oneononev1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/oneonone/v1"
)

type GetRandomQuestionHandlerFunc func(ctx context.Context, req *oneononev1.GetRandomQuestionRequest) (*oneononev1.GetRandomQuestionResponse, error)

func GetRandomQuestion(
	db *sql.DB,
	cr oneononeddlv1.CategoryRecorder,
	cqr oneononeddlv1.CategoryQuestionRecorder,
	qr oneononeddlv1.QuestionRecorder,
) GetRandomQuestionHandlerFunc {
	return func(ctx context.Context, req *oneononev1.GetRandomQuestionRequest) (*oneononev1.GetRandomQuestionResponse, error) {
		ddlCategories, err := cr.List(db, nil, true, 1000)
		if err != nil {
			return nil, err
		}

		ddlCategory, err := getRandomCategory(ddlCategories)
		if err != nil {
			return nil, err
		}

		ddlCategoryQuestions, err := cqr.FindByCategoryId(db, ddlCategory.Id)
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

		return &oneononev1.GetRandomQuestionResponse{
			Question: &oneononev1.Question{
				Id:       fmt.Sprintf("%d", ddlQuestion.Id),
				Question: ddlQuestion.Question,
			},
		}, nil
	}
}

func getRandomCategory(categories []*oneononeddlv1.Category) (*oneononeddlv1.Category, error) {
	if len(categories) == 0 {
		return nil, errors.New("no categories")
	}

	return categories[rand.Intn(len(categories))], nil
}

func listQuestionIDs(categoryQuestions []*oneononeddlv1.CategoryQuestion) []uint64 {
	var questionIDs []uint64
	for _, categoryQuestion := range categoryQuestions {
		questionIDs = append(questionIDs, categoryQuestion.GetQuestionId())
	}

	return questionIDs
}

func getRandomQuestion(questions []*oneononeddlv1.Question) (*oneononeddlv1.Question, error) {
	if len(questions) == 0 {
		return nil, errors.New("no questions")
	}

	return questions[rand.Intn(len(questions))], nil
}
