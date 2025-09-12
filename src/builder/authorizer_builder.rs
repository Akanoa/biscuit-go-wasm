use crate::builder::Builder;
use crate::wasm_export;
use crate::wasm_result::WasmResult;
use biscuit_auth::{Authorizer, AuthorizerBuilder, Biscuit};

// wrapper around the biscuit builder to allow using the builder pattern without reallocating the builder every time
// refcell is used to perform interior mutability (safe because the builder is not exposed to the user)
pub type AuthorizerBuilderWrapper = Builder<AuthorizerBuilder>;


// create a new authorizer builder
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the authorizer builder allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn authorizer_builder_new() -> Box<AuthorizerBuilderWrapper> {
        Box::new(AuthorizerBuilder::new().into())
    }
);

// drop the authorizer builder and free the memory allocated for it
// Input:
// builder is a pointer to the authorizer builder allocated in the wasm memory
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is 0
// data_len is 0
// is_ok is 1 because the function never fails
wasm_export!(
    fn authorizer_builder_drop(builder: Box<AuthorizerBuilderWrapper>) {
        drop(builder);
    }
);

// build the authorizer
// Input:
// builder is a pointer to the authorizer builder allocated in the wasm memory
// token is a pointer to the biscuit token
// Output:
// returnArea { data, data_len, is_ok}
//
// if is_ok is 0
// data is a pointer to the authorizer allocated in the wasm memory
// data_len is 0 because of the opaque type
//
// if is_ok is 1
// data is the pointer to the authorizer error message allocated in the wasm memory
// data_len is the length of the error message
wasm_export!(
    fn authorizer_builder_build(builder: Box<AuthorizerBuilderWrapper>, token: &Biscuit) -> Result<Box<Authorizer>, biscuit_auth::error::Token> {
        let authorizer = builder.0.build(token)?;
        Ok(Box::new(authorizer))
    }
);

// add a code to the authorizer builder
// Input:
// builder is a pointer to the authorizer builder allocated in the wasm memory
// code {ptr, len} is a string containing the code to add
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok is 0
// data is 0 because no data is returned
// data_len is 0 because no data is returned
//
// if is_ok is 1
// data is the pointer to the authorizer error message allocated in the wasm memory
// data_len is the length of the error message
wasm_export!(
    fn authorizer_builder_add_code(builder: &mut AuthorizerBuilderWrapper, code: &str) -> Result<(), biscuit_auth::error::Token> {
        builder.apply(|builder|builder.code(code))
    }
);