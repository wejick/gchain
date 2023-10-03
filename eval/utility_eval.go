package eval

import "encoding/json"

// ValidJson will check if input is valid json
type ValidJson struct {
}

// NewValidJson return valid json evaluator
func NewValidJson() *ValidJson {
	return &ValidJson{}
}

// Evaluate will return true if input is valid json
func (V *ValidJson) Evaluate(input string) (bool, error) {
	var v interface{}
	err := json.Unmarshal([]byte(input), v)
	if err != nil {
		return false, err
	}
	return true, nil
}
