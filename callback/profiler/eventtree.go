package profiler

import (
	"github.com/wejick/gchain/callback"
)

type EventNode struct {
	Parent   *EventNode
	Data     *callback.CallbackData
	Childern []*EventNode
}

type EventTree struct {
	Root *EventNode
}

func buildTree(events []*callback.CallbackData) (tree EventTree) {
	mapIDNode := make(map[string]*EventNode)

	for _, event := range events {
		if event.ParentID == "" {
			tree.Root = &EventNode{
				Data:     event,
				Childern: make([]*EventNode, 0),
			}
			continue
		}

		node, exist := mapIDNode[event.ID]
		if !exist {
			node = &EventNode{
				Data: event,
			}
			mapIDNode[event.ID] = node
		} else {
			node.Data = event
		}

		parentNode, parentExist := mapIDNode[event.ParentID]
		if !parentExist {
			parentNode = &EventNode{
				Childern: make([]*EventNode, 0),
			}
			mapIDNode[event.ParentID] = parentNode
		}
	}

	return
}
