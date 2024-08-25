package leetcodeproblems

import (
	"testing"
)

// go test -v -run TestQuestionDatabase_Random
func TestQuestionDatabase_Random(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		database := NewQuestionDatabase()
		_ = database.Random()
	})
}

// go test -v -run TestQuestionDatabase_PrintAllQuestionsFinishSummary
func TestQuestionDatabase_PrintAllQuestionsFinishSummary(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		database := NewQuestionDatabase()
		_ = database.PrintAllQuestionsFinishSummary()
	})
}
