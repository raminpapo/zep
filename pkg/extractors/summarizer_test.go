package extractors

import (
	"testing"

	"github.com/getzep/zep/pkg/models"
	"github.com/getzep/zep/pkg/testutils"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func TestSummarize(t *testing.T) {
	windowSize := 10
	newMessageCountAfterSummary := windowSize / 2

	messages := make([]models.Message, len(testutils.TestMessages))
	err := copier.Copy(&messages, &testutils.TestMessages)
	assert.NoError(t, err)

	messages = messages[:windowSize+2]
	for i := range messages {
		messages[i].UUID = uuid.New()
	}

	newestMessageToSummarizeIndex := len(
		messages,
	) - newMessageCountAfterSummary - 1 // the seventh-oldest message, leaving 5 messages after it
	newSummaryPointUUID := messages[newestMessageToSummarizeIndex].UUID

	tests := []struct {
		name     string
		messages []models.Message
		summary  *models.Summary
	}{
		{
			name:     "With an existing summary",
			messages: messages,
			summary: &models.Summary{
				Content:    "Existing summary content",
				TokenCount: 10,
			},
		},
		{
			name:     "With a nil-value passed as the summary argument",
			messages: messages,
			summary:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newSummary, err := summarize(testCtx, appState, windowSize, tt.messages, tt.summary, 0)
			assert.NoError(t, err)

			assert.Equal(t, newSummaryPointUUID, newSummary.SummaryPointUUID)
			assert.NotEmpty(t, newSummary.Content)
			assert.True(t, newSummary.TokenCount > 0)
		})
	}
}
