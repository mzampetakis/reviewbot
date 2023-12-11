package sentimentanalysistypes

// SentimentAnalysisResult is the given result of the sentiment analysis of a sentence
type SentimentAnalysisResult struct {
	// SentimentScore provides the sentiment score of the given sentence
	// positive values mean positive sentiment
	// negative values mean negative sentiment
	SentimentScore int64
}
