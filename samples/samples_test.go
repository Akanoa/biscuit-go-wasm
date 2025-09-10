// Copyright (c) 2019 Titanous, daeMOn63 and Contributors to the Eclipse Foundation.
// SPDX-License-Identifier: Apache-2.0

package biscuittest

import (
	"biscuit-wasm-go/builder"
	biscuitData "biscuit-wasm-go/data"
	"biscuit-wasm-go/keypair"
	biscuitModule "biscuit-wasm-go/token"
	"encoding/json"
	"fmt"
	"os"
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
	Ok  *int                      `json:"Ok"`
	Err *biscuitData.BiscuitError `json:"Err"`
}

type Validation struct {
	World          biscuitData.World `json:"world"`
	Result         Result            `json:"result"`
	AuthorizerCode string            `json:"authorizer_code"`
	RevocationIds  []string          `json:"revocation_ids"`
}

func CheckSample(rootKey keypair.PublicKey, testCase TestCase, t *testing.T) {
	// FFI is not supported yet by golang
	if testCase.Filename == "test035_ffi.bc" {
		t.SkipNow()
	}
	t.Logf("Checking sample %s\n", testCase.Filename)

	// Reading biscuit file
	biscuitBytes, err := os.ReadFile("./data/current/" + testCase.Filename)
	require.NoError(t, err)

	// Decoded biscuit file
	token, err := biscuitModule.Biscuit{}.FromBytes(env, biscuitBytes, rootKey)

	// Decoding as failed
	if err != nil {
		for _, validation := range testCase.Validations {
			// searching for error in validation if one of the validations is successful, then
			// there is no decoding error expected
			require.Nil(t, validation.Result.Ok)
		}
		return
	}

	// Decoding as successful, then checking validations
	for validationName, validation := range testCase.Validations {
		t.Logf("Checking validation %s\n", validationName)
		CompareResult(testCase.Filename, token, validation, t)
	}
}

//func CompareBlocks(token biscuitModule.Biscuit, blocks []Block, t *testing.T) {
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

func CompareResult(filename string, token biscuitModule.Biscuit, v Validation, t *testing.T) {

	// create a new authorizer
	authorizerBuilder, err := builder.AuthorizerBuilder{}.New(env)
	require.NoError(t, err)

	// add the code-authorizing
	err = authorizerBuilder.AddCode(v.AuthorizerCode)
	require.NoError(t, err)

	// build the authorizer with the biscuit token
	authorizer, err := authorizerBuilder.Build(token)

	// The authorizer may fail to build, in which case the test fails
	if err != nil {
		// If the authorizer failed to build, verify that the error is as expected
		require.Nil(t, v.Result.Ok)
		// Compare the returned error with the expected one
		CompareError(err, v.Result.Err, t)
		return
	}

	// Authorize the token
	_, err = authorizer.Authorize()

	// If the authorization failed, compare the returned error with the expected one
	if err != nil {
		CompareError(err, v.Result.Err, t)
		return
	}

	// If the authorization is successful, verify that's the expected result
	require.NotNil(t, v.Result.Ok)

	// If the authorization is successful, verify that the world is as expected
	world, err := authorizer.ToString()
	require.NoError(t, err)
	// Compare the world with the expected one
	require.Equal(t, world, v.World.String(),
		"World mismatch for sample %s", filename,
	)
}

func CompareError(authorizationError error, sampleError *biscuitData.BiscuitError, t *testing.T) {
	biscuitAuthorizationError := biscuitData.FromErrorAsBiscuitError(authorizationError)

	if !biscuitAuthorizationError.Equal(*sampleError) {
		require.Fail(t, "BiscuitErrors are not equal", sampleError.Error(), biscuitAuthorizationError.Error())
	}
}

func TestReadSamples(t *testing.T) {
	sampleFileBytes, err := os.ReadFile("./data/current/samples.json")
	require.NoError(t, err)
	var samples Samples
	err = json.Unmarshal(sampleFileBytes, &samples)
	require.NoError(t, err)

	rootKey, err := keypair.PublicKey{}.FromString(env, samples.RootPublicKey, keypair.Ed25519)
	require.NoError(t, err)
	fmt.Printf("Checking %d samples\n", len(samples.TestCases))
	for _, v := range samples.TestCases {
		t.Run(v.Filename, func(t *testing.T) { CheckSample(rootKey, v, t) })
	}

}
