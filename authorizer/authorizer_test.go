package authorizer_test

import (
	"biscuit-wasm-go/builder"
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

	authorizeCode := "allow if user(1); deny if user(2);"

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

	matchedPolicy, err := authorizer.Authorize()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(matchedPolicy)

}
