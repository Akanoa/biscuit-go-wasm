package error

import (
	"testing"
)

func TestFailedLogicError_Error(t *testing.T) {

	data := map[string]interface{}{
		"Unauthorized": map[string]interface{}{
			"policy": map[string]interface{}{
				"deny": 1,
			},
		},
	}

	t.Log(FailedLogicError{data})
}
