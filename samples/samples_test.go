// Copyright (c) 2019 Titanous, daeMOn63 and Contributors to the Eclipse Foundation.
// SPDX-License-Identifier: Apache-2.0

package biscuittest

import (
	"biscuit-wasm-go/builder"
	"biscuit-wasm-go/keypair"
	token2 "biscuit-wasm-go/token"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

type Samples struct {
	RootPrivateKey string     `json:"root_private_key"`
	RootPublicKey  string     `json:"root_public_key"`
	TestCases      []TestCase `json:"testcases"`
}

type TestCase struct {
	Title       string                `json:"title"`
	Filename    string                `json:"filename"`
	Token       []Block               `json:"token"`
	Validations map[string]Validation `json:"validations"`
}

type Block struct {
	Symbols     []string `json:"symbols"`
	PublicKeys  []any    `json:"public_keys"`
	ExternalKey any      `json:"external_key"`
	Code        string   `json:"code"`
}

type Result struct {
	Ok  *int          `json:"Ok"`
	Err *BiscuitError `json:"Err"`
}

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
}

type World struct {
	Facts    []ScopedFact `json:"facts"`
	Rules    []ScopedRule `json:"rules"`
	Checks   []string     `json:"checks"`
	Policies []string     `json:"policies"`
}

type ScopedFact struct {
	Fact  string
	Scope []*int32
}

func (sf *ScopedFact) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&sf.Fact, &sf.Scope}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in ScopedFact: %d != %d", g, e)
	}
	return nil
}

type ScopedRule struct {
	Rule  string
	Scope *int32
}

func (sr *ScopedRule) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&sr.Rule, &sr.Scope}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in ScopedRule: %d != %d", g, e)
	}
	return nil
}

func (w World) String() string {
	facts := []string{}
	for _, f := range w.Facts {
		visible := true
		for _, s := range f.Scope {

			if s != nil && *s != 0 {
				visible = false
				break
			}
		}

		if visible {
			facts = append(facts, f.Fact)
		}
	}
	sort.Strings(facts)
	rules := []string{}
	for _, r := range w.Rules {
		if r.Scope == nil || *r.Scope == 0 {
			rules = append(rules, r.Rule)
		}
	}
	sort.Strings(rules)

	return fmt.Sprintf("World {{\n\tfacts: %v\n\trules: %v\n}}", facts, rules)
}

type Validation struct {
	World          World    `json:"world"`
	Result         Result   `json:"result"`
	AuthorizerCode string   `json:"authorizer_code"`
	RevocationIds  []string `json:"revocation_ids"`
}

func CheckSample(root_key keypair.PublicKey, c TestCase, t *testing.T) {
	// all these contain v4 blocks, which are not supported yet
	if c.Filename == "test024_third_party.bc" ||
		c.Filename == "test025_check_all.bc" ||
		c.Filename == "test026_public_keys_interning.bc" ||
		c.Filename == "test027_integer_wraparound.bc" ||
		c.Filename == "test028_expressions_v4.bc" {
		t.SkipNow()
	}
	fmt.Printf("Checking sample %s\n", c.Filename)
	b, err := os.ReadFile("./data/current/" + c.Filename)
	require.NoError(t, err)
	token, err := token2.Biscuit{}.FromBytes(env, b, root_key)

	if err == nil {
		fmt.Printf("  Parsed file %s\n", c.Filename)
		// this sample uses a tampered biscuit file on purpose
		if c.Filename != "test006_reordered_blocks.bc" {
			//CompareBlocks(token, c.Token, t)
		}

		i := 0
		for _, v := range c.Validations {
			fmt.Printf("  Checking validation %d\n", i)
			CompareResult(c.Filename, token, v, t)
			i++
		}

	} else {
		fmt.Println(err)
		fmt.Println("  Parsing failed, all validations must be errors")
		for _, v := range c.Validations {
			require.Nil(t, v.Result.Ok)
		}
	}
}

//func CompareBlocks(token token2.Biscuit, blocks []Block, t *testing.T) {
//	sample := token.Code()
//	p := parser.New()
//
//	rng := rand.Reader
//	_, privateRoot, _ := ed25519.GenerateKey(rng)
//	authority, err := p.Block(blocks[0].Code, nil)
//	require.NoError(t, err)
//	builder := biscuit.NewBuilder(privateRoot)
//	builder.AddBlock(authority)
//	r, err := builder.Build()
//	require.NoError(t, err)
//	rebuilt := *r
//
//	for _, b := range blocks[1:] {
//		parsed, err := p.Block(b.Code, nil)
//		require.NoError(t, err)
//		builder := rebuilt.CreateBlock()
//		builder.AddBlock(parsed)
//		r, err := rebuilt.Append(rng, builder.Build())
//		require.NoError(t, err)
//		rebuilt = *r
//	}
//
//	require.Equal(t, sample, rebuilt.Code())
//}

func CompareResult(filename string, token token2.Biscuit, v Validation, t *testing.T) {

	authorizerBuilder, err := builder.AuthorizerBuilder{}.New(env)
	require.NoError(t, err)

	fmt.Println(v.AuthorizerCode)

	err = authorizerBuilder.AddCode(v.AuthorizerCode)
	require.NoError(t, err)

	authorizer, err := authorizerBuilder.Build(token)

	fmt.Println(authorizer.ToString())

	if err != nil {
		CompareError(err, v.Result.Err, t)
	} else {
		_, err = authorizer.Authorize()
		fmt.Println(err)
		if err != nil {
			CompareError(err, v.Result.Err, t)
		} else {
			require.NotNil(t, v.Result.Ok)
		}

		world, err := authorizer.ToString()
		require.NoError(t, err)

		require.Equal(t, v.World.String(), world,
			"World mismatch for sample %s", filename,
		)
	}
}

func CompareError(authorizationError error, sampleError *BiscuitError, t *testing.T) {
	errorString := authorizationError.Error()
	if sampleError.Format != nil {
		require.Equal(t, errorString, "biscuit: invalid signature")
	} else if sampleError.FailedLogic != nil {
		if sampleError.FailedLogic.Unauthorized != nil {
			// todo check the block and check ids (if there is a single failed check, because the lib only reports one)
			require.Regexp(t, "^biscuit: verification failed: failed to verify", errorString)
		} else if sampleError.FailedLogic.InvalidBlockRule != nil {
			// todo extract the block number
			require.Regexp(t, "^biscuit: verification failed: failed to verify", errorString)
		} else {
			require.Fail(t, errorString)
		}
	} else {
		fmt.Println(sampleError)
		require.Fail(t, errorString)
	}
}

func TestReadSamples(t *testing.T) {
	b, err := os.ReadFile("./data/current/samples.json")
	require.NoError(t, err)
	var samples Samples
	err = json.Unmarshal(b, &samples)
	require.NoError(t, err)

	rootKey, err := keypair.PublicKey{}.FromString(env, samples.RootPublicKey, keypair.Ed25519)
	require.NoError(t, err)
	fmt.Printf("Checking %d samples\n", len(samples.TestCases))
	for _, v := range samples.TestCases {
		t.Run(v.Filename, func(t *testing.T) { CheckSample(rootKey, v, t) })
	}

}
