package error

import (
	"encoding/json"
	"errors"
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
	} `json:"FailedLogic"`
	Format *struct {
		Signature *struct {
			InvalidSignature string `json:"InvalidSignature"`
		} `json:"Signature"`
	} `json:"Format"`
	Raw string `json:"-"`
}

func (error BiscuitError) FromString(value string) (BiscuitError, error) {
	data := BiscuitError{}
	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (error BiscuitError) Equal(error2 BiscuitError) bool {
	// Compare FailedLogic
	if error.FailedLogic == nil && error2.FailedLogic == nil {
		// Both are nil, continue to Format comparison
	} else if error.FailedLogic == nil || error2.FailedLogic == nil {
		return false // One is nil, other is not
	} else {
		// Both are not nil, compare their contents

		// Compare Unauthorized
		if error.FailedLogic.Unauthorized == nil && error2.FailedLogic.Unauthorized == nil {
			// Both are nil, continue
		} else if error.FailedLogic.Unauthorized == nil || error2.FailedLogic.Unauthorized == nil {
			return false // One is nil, other is not
		} else {
			// Compare Policy
			if error.FailedLogic.Unauthorized.Policy.Allow != error2.FailedLogic.Unauthorized.Policy.Allow {
				return false
			}

			// Compare Checks slices
			if len(error.FailedLogic.Unauthorized.Checks) != len(error2.FailedLogic.Unauthorized.Checks) {
				return false
			}
			for i, check1 := range error.FailedLogic.Unauthorized.Checks {
				check2 := error2.FailedLogic.Unauthorized.Checks[i]
				if check1.Block.BlockID != check2.Block.BlockID ||
					check1.Block.CheckID != check2.Block.CheckID ||
					check1.Block.Rule != check2.Block.Rule {
					return false
				}
			}
		}

		// Compare InvalidBlockRule slices
		if !reflect.DeepEqual(error.FailedLogic.InvalidBlockRule, error2.FailedLogic.InvalidBlockRule) {
			return false
		}
	}

	// Compare Format
	if error.Format == nil && error2.Format == nil {
		return true // Both are nil
	} else if error.Format == nil || error2.Format == nil {
		return false // One is nil, other is not
	} else {
		// Both are not nil, compare their contents

		// Compare Signature
		if error.Format.Signature == nil && error2.Format.Signature == nil {
			return true // Both are nil
		} else if error.Format.Signature == nil || error2.Format.Signature == nil {
			return false // One is nil, other is not
		} else {
			// Compare InvalidSignature
			return error.Format.Signature.InvalidSignature == error2.Format.Signature.InvalidSignature
		}
	}
}

func (error BiscuitError) Error() string {
	return error.Raw
}

func FromErrorAsBiscuitError(err error) BiscuitError {
	var biscuitError BiscuitError
	if errors.As(err, &biscuitError) {
		return biscuitError
	}
	panic("not a biscuit error")
}
