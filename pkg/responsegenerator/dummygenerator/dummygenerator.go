package dummygenerator

import (
	"context"
	"reviewbot/pkg/responsegenerator"
	"reviewbot/pkg/responsegenerator/responsegeneratortypes"
	"reviewbot/pkg/sentimentanalyzer/sentimentanalysistypes"
)

type DummyGenerator struct{}

func NewDummyGenerator() responsegenerator.ResponseGenerator {
	return &DummyGenerator{}
}

func (dg *DummyGenerator) Generate(ctx context.Context, sentimentAnalysisScore sentimentanalysistypes.
	SentimentAnalysisResult, context string) (responsegeneratortypes.ResponseGeneratorResult, error) {
	responseResult := responsegeneratortypes.ResponseGeneratorResult{
		Response: "",
	}
	if sentimentAnalysisScore.SentimentScore > 0 {
		responseResult.Response = "Happy to hear that!"
	} else {
		responseResult.Response = "Sorry to hear that."
	}
	return responseResult, nil
}
