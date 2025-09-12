package data

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBiscuitError_Equal(t *testing.T) {
	// Test case 1: Both errors are empty (all nil fields)
	err1 := BiscuitError{}
	err2 := BiscuitError{}
	if !err1.Equal(err2) {
		t.Error("Empty BiscuitErrors should be equal")
	}

	// Test case 2: One has FailedLogic, other doesn't
	err3 := BiscuitError{
		FailedLogic: &struct {
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
		}{},
	}
	err4 := BiscuitError{}
	if err3.Equal(err4) {
		t.Error("BiscuitErrors with different FailedLogic should not be equal")
	}

	// Test case 3: Same FailedLogic with Unauthorized
	err5 := BiscuitError{
		FailedLogic: &struct {
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
		}{
			Unauthorized: &struct {
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
			}{
				Policy: struct {
					Allow int `json:"Allow"`
				}{Allow: 1},
				Checks: []struct {
					Block struct {
						BlockID int    `json:"block_id"`
						CheckID int    `json:"check_id"`
						Rule    string `json:"rule"`
					} `json:"Block"`
				}{
					{
						Block: struct {
							BlockID int    `json:"block_id"`
							CheckID int    `json:"check_id"`
							Rule    string `json:"rule"`
						}{
							BlockID: 1,
							CheckID: 2,
							Rule:    "test rule",
						},
					},
				},
			},
		},
	}

	err6 := BiscuitError{
		FailedLogic: &struct {
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
		}{
			Unauthorized: &struct {
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
			}{
				Policy: struct {
					Allow int `json:"Allow"`
				}{Allow: 1},
				Checks: []struct {
					Block struct {
						BlockID int    `json:"block_id"`
						CheckID int    `json:"check_id"`
						Rule    string `json:"rule"`
					} `json:"Block"`
				}{
					{
						Block: struct {
							BlockID int    `json:"block_id"`
							CheckID int    `json:"check_id"`
							Rule    string `json:"rule"`
						}{
							BlockID: 1,
							CheckID: 2,
							Rule:    "test rule",
						},
					},
				},
			},
		},
	}

	if !err5.Equal(err6) {
		t.Error("BiscuitErrors with same FailedLogic should be equal")
	}

	// Test case 4: Different Policy Allow values
	err7 := err5
	err8 := err6
	err8.FailedLogic.Unauthorized.Policy.Allow = 2
	if err7.Equal(err8) {
		t.Error("BiscuitErrors with different Policy Allow values should not be equal")
	}

	// Test case 5: Format comparison
	err9 := BiscuitError{
		Format: &struct {
			Signature *struct {
				InvalidSignature string `json:"InvalidSignature"`
			} `json:"Signature"`
		}{
			Signature: &struct {
				InvalidSignature string `json:"InvalidSignature"`
			}{
				InvalidSignature: "test signature",
			},
		},
	}

	err10 := BiscuitError{
		Format: &struct {
			Signature *struct {
				InvalidSignature string `json:"InvalidSignature"`
			} `json:"Signature"`
		}{
			Signature: &struct {
				InvalidSignature string `json:"InvalidSignature"`
			}{
				InvalidSignature: "test signature",
			},
		},
	}

	if !err9.Equal(err10) {
		t.Error("BiscuitErrors with same Format should be equal")
	}

	// Test case 6: Different InvalidSignature values
	err11 := err9
	err12 := err10
	err12.Format.Signature.InvalidSignature = "different signature"
	if err11.Equal(err12) {
		t.Error("BiscuitErrors with different InvalidSignature should not be equal")
	}

	// Test case 7: InvalidBlockRule comparison
	err13 := BiscuitError{
		FailedLogic: &struct {
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
		}{
			InvalidBlockRule: []any{"rule1", "rule2"},
		},
	}

	err14 := BiscuitError{
		FailedLogic: &struct {
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
		}{
			InvalidBlockRule: []any{"rule1", "rule2"},
		},
	}

	if !err13.Equal(err14) {
		t.Error("BiscuitErrors with same InvalidBlockRule should be equal")
	}

	// Test case 8: Different InvalidBlockRule
	err15 := err13
	err16 := err14
	err16.FailedLogic.InvalidBlockRule = []any{"rule1", "rule3"}
	if err15.Equal(err16) {
		t.Error("BiscuitErrors with different InvalidBlockRule should not be equal")
	}
}

type BiscuitError2 struct {
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
	Raw string `json:"raw,omitempty"`
}

func TestBiscuitError_FromString(t *testing.T) {
	data := `{"FailedLogic":{"Unauthorized":{"policy":{"Allow":0},"checks":[{"Block":{"block_id":1,"check_id":0,"rule":"check if resource($0), operation(\"read\"), right($0, \"read\")"}}]}}}`
	tmp := BiscuitError2{}
	err := json.Unmarshal([]byte(data), &tmp)
	fmt.Println(err)
	fmt.Printf("%+v\n", tmp)
}

func TestBiscuitError_FromString2(t *testing.T) {
	data := `{"FailedLogic":{"Unauthorized":{"policy":{"Allow":0},"checks":[{"Block":{"block_id":1,"check_id":0,"rule":"check if resource($0), operation(\"read\"), right($0, \"read\")"}}]}}}`
	tmp := BiscuitError{}
	tmp, err := tmp.FromString(data)
	fmt.Println(err)
	fmt.Printf("%+v\n", tmp)
}
