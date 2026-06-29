package pipeline

import (
	"encoding/json"
	"fmt"
)

func Decode(data []byte) (Finding, error) {
	var f Finding
	if err := json.Unmarshal(data, &f); err != nil {
		return Finding{}, fmt.Errorf("decode finding: %w", err)
	}
	return f, nil
}
