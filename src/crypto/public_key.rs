use crate::crypto::SignatureAlgorithm;
use crate::wasm_export;
use crate::wasm_result::WasmResult;
use biscuit_auth::PublicKey;

// Format the public key as a hexadecimal string
// Input:
// private_key is a pointer to the public key allocated in the wasm memory
// Output:
// returnArea { data, data_len, is_ok = 1 }
// data is the private key in hex format
// data_len is the length of the hexadecimal representation of the private key in bytes
// is_ok is 1 because the function never fails
wasm_export!(
    fn public_key_to_hex(public_key: &PublicKey) -> String {
        public_key.to_string()
    }
);

// Create a private key from a hexadecimal string
// Input:
// data is the private key in hex format
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok = 1
// data is a pointer to the private key allocated in the wasm memory
// data_len is 0 because of the opaque type
//
// if is_ok = 0
// data is the pointer to the error message allocated in the wasm memory
// data_len is the length of the error message in bytes
wasm_export!(
    fn public_key_from_hex(data: &str, algorithm: SignatureAlgorithm) -> Result<Box<PublicKey>, biscuit_auth::error::Format> {
        Ok(Box::new(PublicKey::from_bytes_hex(data, algorithm.into())?))
    }
);