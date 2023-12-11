package sentimentanalyzer

import (
	"context"
	"reviewbot/pkg/sentimentanalyzer/sentimentanalysistypes"
)

// SentimentAnalyze interface for sentiment analysis of a sentence
type SentimentAnalyze interface {
	// Process processes the given input string for sentiment analysis resulting with a SentimentAnalysisResult
	Process(context.Context, string) (sentimentanalysistypes.SentimentAnalysisResult, error)
}
