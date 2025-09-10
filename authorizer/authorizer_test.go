package authorizer_test

import (
	"biscuit-wasm-go/builder"
	error2 "biscuit-wasm-go/error"
	"biscuit-wasm-go/factory"
	"testing"
)

func TestAuthorizer_authorize(t *testing.T) {
	token, err := factory.MakeBiscuit(env, "user(2)")
	if err != nil {
		t.Error(err)
		return
	}

	authorizerBuilder, err := builder.AuthorizerBuilder{}.New(env)
	if err != nil {
		t.Error(err)
		return
	}

	authorizeCode := "deny if user(1); allow if true"

	err = authorizerBuilder.AddCode(authorizeCode)
	if err != nil {
		t.Error(err)
		return
	}

	authorizer, err := authorizerBuilder.Build(token.Token)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = authorizer.Authorize()
	// No error is expected
	if err != nil {
		t.Error(err)
		return
	}

}

// Test that the authorizer returns an error when the policy is not satisfied
// the authorize code is "deny if user(1); allow if true" and the world is "user(1)"
func TestAuthorizer_authorize_fail(t *testing.T) {
	token, err := factory.MakeBiscuit(env, "user(1)")
	if err != nil {
		t.Error(err)
		return
	}

	authorizerBuilder, err := builder.AuthorizerBuilder{}.New(env)
	if err != nil {
		t.Error(err)
		return
	}

	authorizeCode := "deny if user(1); allow if true"

	err = authorizerBuilder.AddCode(authorizeCode)
	if err != nil {
		t.Error(err)
		return
	}

	authorizer, err := authorizerBuilder.Build(token.Token)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = authorizer.Authorize()

	if err == nil {
		t.Errorf("The deny rule should have been triggered")
		return
	}

	errorFromAuthorizer := error2.FromErrorAsBiscuitError(err)

	expectedErrorString := `{"FailedLogic":{"Unauthorized":{"policy":{"Deny":0},"checks":[]}}}`
	expectedError, err := error2.BiscuitError{}.FromString(expectedErrorString)
	if err != nil {
		t.Error(err)
		return
	}

	if !expectedError.Equal(errorFromAuthorizer) {
		t.Errorf("The error returned by the authorizer is not the expected one")
		return
	}
}
