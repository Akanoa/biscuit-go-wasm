package data

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type BiscuitError struct {
	FailedLogic *struct {
		Unauthorized *struct {
			Policy struct {
				Allow int `json:"Allow"`
			} `json:"policy"`
			Checks []struct {
				Block struct {
					BlockID int    `json:"block_id"`
					CheckID int    `json:"check_id"`
					Rule    string `json:"rule"`
				} `json:"Block"`
			} `json:"checks"`
		} `json:"Unauthorized"`
		InvalidBlockRule []any `json:"InvalidBlockRule"`
	} `json:"FailedLogic,omitempty"`
	Format *struct {
		Signature *struct {
			InvalidSignature string `json:"InvalidSignature"`
		} `json:"Signature"`
	} `json:"Format,omitempty"`
	Execution *string `json:"Execution,omitempty"`
	Raw       string  `json:"raw,omitempty"`
}

func (biscuitError BiscuitError) FromString(value string) (BiscuitError, error) {

	tmp := BiscuitError{}
	err := json.Unmarshal([]byte(value), &tmp)
	if err != nil {
		return BiscuitError{}, err
	}
	return tmp, nil
}

func (biscuitError BiscuitError) Equal(other BiscuitError) bool {
	// Compare FailedLogic
	if biscuitError.FailedLogic == nil && other.FailedLogic == nil {
		// Both are nil, continue to Format comparison
	} else if biscuitError.FailedLogic == nil || other.FailedLogic == nil {
		return false // One is nil, other is not
	} else {
		// Both are not nil, compare their contents

		// Compare Unauthorized
		if biscuitError.FailedLogic.Unauthorized == nil && other.FailedLogic.Unauthorized == nil {
			// Both are nil, continue
		} else if biscuitError.FailedLogic.Unauthorized == nil || other.FailedLogic.Unauthorized == nil {
			return false // One is nil, other is not
		} else {
			// Compare Policy
			if biscuitError.FailedLogic.Unauthorized.Policy.Allow != other.FailedLogic.Unauthorized.Policy.Allow {
				return false
			}

			// Compare Checks slices
			if len(biscuitError.FailedLogic.Unauthorized.Checks) != len(other.FailedLogic.Unauthorized.Checks) {
				return false
			}
			for i, check1 := range biscuitError.FailedLogic.Unauthorized.Checks {
				check2 := other.FailedLogic.Unauthorized.Checks[i]
				if check1.Block.BlockID != check2.Block.BlockID ||
					check1.Block.CheckID != check2.Block.CheckID ||
					check1.Block.Rule != check2.Block.Rule {
					return false
				}
			}
		}

		// Compare InvalidBlockRule slices
		if !reflect.DeepEqual(biscuitError.FailedLogic.InvalidBlockRule, other.FailedLogic.InvalidBlockRule) {
			return false
		}
	}

	// Compare Format
	if biscuitError.Format == nil && other.Format == nil {
		return true // Both are nil
	} else if biscuitError.Format == nil || other.Format == nil {
		return false // One is nil, other is not
	} else {
		// Both are not nil, compare their contents

		// Compare Signature
		if biscuitError.Format.Signature == nil && other.Format.Signature == nil {
			return true // Both are nil
		} else if biscuitError.Format.Signature == nil || other.Format.Signature == nil {
			return false // One is nil, other is not
		} else {
			// Compare InvalidSignature
			return biscuitError.Format.Signature.InvalidSignature == other.Format.Signature.InvalidSignature
		}
	}
}

func (biscuitError BiscuitError) Error() string {
	if biscuitError.Raw != "" {
		return biscuitError.Raw
	}
	if biscuitError.Format != nil {
		x, err := json.Marshal(biscuitError)
		if err != nil {
			return err.Error()
		}
		return string(x)
	}
	if biscuitError.FailedLogic != nil {
		x, err := json.Marshal(biscuitError)
		if err != nil {
			return err.Error()
		}
		return string(x)
	}

	if biscuitError.Execution != nil {
		x, err := json.Marshal(biscuitError)
		if err != nil {
			return err.Error()
		}
		return string(x)
	}

	fmt.Println("Not a biscuit biscuitError", biscuitError.Raw, biscuitError.Format, biscuitError.FailedLogic)

	return "Not a biscuit biscuitError"
}

func FromErrorAsBiscuitError(err error) BiscuitError {

	if err == nil {
		return BiscuitError{}
	}

	biscuitError, err := BiscuitError{}.FromString(err.Error())
	if err == nil {
		return biscuitError
	}

	panic("not a biscuit error")
}
