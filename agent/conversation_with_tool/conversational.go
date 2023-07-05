package conversation_with_tool

import (
	"context"

	"github.com/wejick/gchain/agent"
	"github.com/wejick/gchain/memory"
	"github.com/wejick/gchain/model"
)

type ConversationalAgent struct {
	model  *model.ChatModel
	memory memory.Memory
}

func (ChatModel *ConversationalAgent) Plan(ctx context.Context, userPrompt string, actionTaken []agent.Action) (plan agent.Action, err error) {

	return
}
