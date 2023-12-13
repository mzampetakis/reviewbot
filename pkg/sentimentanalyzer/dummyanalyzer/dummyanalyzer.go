package dummyganalyzer

import (
	"context"
	"reviewbot/pkg/sentimentanalyzer"
	"reviewbot/pkg/sentimentanalyzer/sentimentanalysistypes"
)

type DummyAnalyzer struct{}

func NewDummyAnalyzer() sentimentanalyzer.SentimentAnalyze {
	return &DummyAnalyzer{}
}

func (da *DummyAnalyzer) Process(ctx context.Context, input string) (sentimentanalysistypes.SentimentAnalysisResult,
	error) {
	inputLen := len(input)
	scorePercentage := inputLen%3 - 1
	analysisResult := sentimentanalysistypes.SentimentAnalysisResult{
		SentimentScore: int64(scorePercentage),
	}
	return analysisResult, nil
}
