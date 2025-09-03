package error

import "fmt"

type FailedLogicError struct {
	Data map[string]interface{}
}

func (e FailedLogicError) Error() string {
	var out string

	// Unauthorized
	for errorKey, errorValue := range e.Data {
		out += "Failed Logic\n"
		out += errorKey + ": \n"

		// expect a map[string]interface{} at this level; skip if not
		errorLevel1, ok := errorValue.(map[string]interface{})
		if !ok {
			continue
		}

		// policy, checks
		for errorItemKey, errorItemValue := range errorLevel1 {
			out += "\t" + errorItemKey + ": \n"

			// deepest level can be either map[string]uint64 or map[string]interface{}
			switch v := errorItemValue.(type) {
			case map[string]uint64:
				for errorDataKey, errorDataValue := range v {
					out += fmt.Sprintf("\t\t%s -> %d\n", errorDataKey, errorDataValue)
				}
			case map[string]int:
				for errorDataKey, errorDataValue := range v {
					out += fmt.Sprintf("\t\t%s -> %d\n", errorDataKey, errorDataValue)
				}
			case map[string]float64:
				for errorDataKey, errorDataValue := range v {
					out += fmt.Sprintf("\t\t%s -> %.0f\n", errorDataKey, errorDataValue)
				}
			case map[string]interface{}:
				for errorDataKey, errorDataValue := range v {
					// try to render numeric values as integers when possible
					switch n := errorDataValue.(type) {
					case int:
						out += fmt.Sprintf("\t\t%s -> %d\n", errorDataKey, n)
					case int64:
						out += fmt.Sprintf("\t\t%s -> %d\n", errorDataKey, n)
					case uint64:
						out += fmt.Sprintf("\t\t%s -> %d\n", errorDataKey, n)
					case float64:
						out += fmt.Sprintf("\t\t%s -> %.0f\n", errorDataKey, n)
					case string:
						out += fmt.Sprintf("\t\t%s -> %s\n", errorDataKey, n)
					default:
						out += fmt.Sprintf("\t\t%s -> %v\n", errorDataKey, n)
					}
				}
			default:
				// if it's a scalar value (not a map), print it directly under the item
				out += fmt.Sprintf("\t\t%v\n", v)
			}
		}
	}

	return out
}
