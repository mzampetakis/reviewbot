package responsegenerator

import (
	"context"
	"reviewbot/pkg/responsegenerator/responsegeneratortypes"
	"reviewbot/pkg/sentimentanalyzer/sentimentanalysistypes"
)

// ResponseGenerator interface for generating responses based on sentiment analysis
type ResponseGenerator interface {
	// Generate generates a response based on the given sentiment analysis score
	// sentimentAnalysisScore is the score of the sentiment analysis
	// context provides more context about the given prompt
	Generate(ctx context.Context, sentimentAnalysisScore sentimentanalysistypes.SentimentAnalysisResult,
		context string) (responsegeneratortypes.ResponseGeneratorResult, error)
}
