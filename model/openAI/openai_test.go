package _openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOpenAIClientConfig(t *testing.T) {
	authToken := "my-auth-token"
	opts := OpenAIOption{
		BaseURL:    "https://my-base-url.com",
		OrgID:      "my-org-id",
		APIVersion: "v2",
	}

	t.Run("DefaultConfig", func(t *testing.T) {
		clientConfig := newOpenAIClientConfig(authToken, OpenAIOption{})
		assert.Equal(t, "https://api.openai.com/v1", clientConfig.BaseURL)
		assert.Equal(t, "", clientConfig.OrgID)
		assert.Equal(t, "", clientConfig.APIVersion)
	})

	t.Run("DefaultAzureConfig", func(t *testing.T) {
		clientConfig := newOpenAIClientConfig(authToken, opts)
		assert.Equal(t, opts.BaseURL, clientConfig.BaseURL)
		assert.Equal(t, opts.OrgID, clientConfig.OrgID)
		assert.Equal(t, "v2", clientConfig.APIVersion)
		assert.Equal(t, "AZURE", string(clientConfig.APIType))
	})
}
