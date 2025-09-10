use crate::wasm_result::WasmResult;
use crate::wasm_export;
use biscuit_auth::{Authorizer, Biscuit, BiscuitBuilder, PublicKey};


// create a new biscuit builder
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the biscuit builder allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn biscuit_builder() -> Box<BiscuitBuilder> {
        Box::new(BiscuitBuilder::new())
    }
);

// create a new biscuit authorizer used for querying
// Input:
// biscuit: a pointer to a biscuit allocated in the wasm memory
//
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok = 0
// data is a pointer to the biscuit authorizer allocated in the wasm memory
// data_len is 0 because of the opaque type
//
// if is_ok = 1
// data is the pointer to the error message
// data_len is the length of the error message
wasm_export!(
    fn biscuit_authorizer(biscuit: &Biscuit) -> Result<Box<Authorizer>,biscuit_auth::error::Token> {
        Ok(Box::new(biscuit.authorizer()?))
    }
);

// create a new biscuit from bytes
// Input:
// data {ptr, len}: a pointer to the base64 encoded bytes allocated in the wasm memory
// root_public_key: a pointer to the root public key allocated in the wasm memory
//
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok = 0
// data is a pointer to the biscuit allocated in the wasm memory
// data_len is 0 because of the opaque type
//
// if is_ok = 1
// data is the pointer to the error message
// data_len is the length of the error message
wasm_export!(
    fn biscuit_from_bytes(data: &[u8], root_public_key: &PublicKey) -> Result<Box<Biscuit>,biscuit_auth::error::Token> {
        Ok(Box::new(Biscuit::from(data, root_public_key)?))
    }
);

// create a new biscuit from base64
// Input:
// data {ptr, len}: a pointer to the base64 encoded bytes allocated in the wasm memory
// root_public_key: a pointer to the root public key allocated in the wasm memory
//
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok = 0
// data is a pointer to the biscuit allocated in the wasm memory
// data_len is 0 because of the opaque type
//
// if is_ok = 1
// data is the pointer to the error message
// data_len is the length of the error message
wasm_export!(
    fn biscuit_from_base64(data: &str, root_public_key: &PublicKey) -> Result<Box<Biscuit>,biscuit_auth::error::Token> {
        let biscuit = Biscuit::from_base64(data, root_public_key);
        let biscuit = biscuit?;
        Ok(Box::new(biscuit))
    }
);

// drop a biscuit
// Input:
// biscuit: a pointer to a biscuit allocated in the wasm memory
// Output:
// returnArea { data=0, data_len=0, is_ok=1 }
//
// data is 0 because no data is returned
// data_len is 0 because no data is returned
// is_ok is 1 because the function never fails
wasm_export!(
    fn biscuit_drop(biscuit: Box<Biscuit>) {
        drop(biscuit);
    }
);

// convert a biscuit to base64
// Input:
// biscuit: a pointer to a biscuit allocated in the wasm memory
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok = 0
// data is the pointer to the base64 encoded bytes allocated in the wasm memory
// data_len is the length of the base64 encoded bytes
//
// if is_ok = 1
// data is the pointer to the error message
// data_len is the length of the error message
wasm_export!(
    fn biscuit_to_base64(biscuit: &Biscuit) -> Result<String,biscuit_auth::error::Token> {
        biscuit.to_base64()
    }
);

// convert a biscuit to bytes
// Input:
// biscuit: a pointer to a biscuit allocated in the wasm memory
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok = 0
// data is the pointer to the bytes allocated in the wasm memory
// data_len is the length of the bytes
//
// if is_ok = 1
// data is the pointer to the error message
// data_len is the length of the error message
wasm_export!(
    fn biscuit_to_bytes(biscuit: &Biscuit) -> Result<Vec<u8>,biscuit_auth::error::Token> {
        biscuit.to_vec()
    }
);