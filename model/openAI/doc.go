/*
To use ChatStreaming, call it as go routine and then listen to the channel until it's closed

Example of ChatStreaming usage :

	Chat := NewOpenAIChatModel(token,"",modelname)
	streamingChannel := make(chan model.ChatMessage,100) // make it buffered
	go Chat.ChatStreaming(ctx,messages, model.WithIsStreaming(true), model.WithStreamingChannel(streamingChannel))
	for {
		value, ok := <-streamingChannel
		if ok && !model.IsStreamFinished(value) {
			fmt.Print(value.Content)
		} else {
			fmt.Println("")
			break
		}
	}
*/
package _openai
