use crate::builder::apply_in_place;
use crate::wasm_export;
use crate::wasm_result::WasmResult;
use biscuit_auth::{Authorizer, AuthorizerBuilder, Biscuit};

// create a new authorizer builder
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the authorizer builder allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn authorizer_builder_new() -> Box<AuthorizerBuilder> {
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
    fn authorizer_builder_drop(builder: Box<AuthorizerBuilder>) {
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
    fn authorizer_builder_build(
        builder: Box<AuthorizerBuilder>,
        token: &Biscuit,
    ) -> Result<Box<Authorizer>, biscuit_auth::error::Token> {
        let authorizer = builder.build(token)?;
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
    fn authorizer_builder_add_code(
        builder: Box<AuthorizerBuilder>,
        code: &str,
    ) -> Result<(), biscuit_auth::error::Token> {
        apply_in_place(builder, |builder| builder.code(code))
    }
);
