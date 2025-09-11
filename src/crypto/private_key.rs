use crate::wasm_export;
use crate::wasm_result::WasmResult;
use biscuit_auth::PrivateKey;
use core::str::FromStr;

// drop the private key and free the memory allocated for it
// Input:
// keypair is a pointer to the private key allocated in the wasm memory
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is 0
// data_len is 0
// is_ok is 1 because the function never fails
wasm_export!(
    fn private_key_drop(private_key: Box<PrivateKey>) {
        drop(private_key);
    }
);

// Format the private key as a hexadecimal string
// Input:
// private_key is a pointer to the private key allocated in the wasm memory
// Output:
// returnArea { data, data_len, is_ok = 1 }
// data is the private key in hex format
// data_len is the length of the hexadecimal representation of the private key in bytes
// is_ok is 1 because the function never fails
wasm_export!(
    fn private_key_to_hex(private_key: &PrivateKey) -> String {
        private_key.to_prefixed_string()
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
    fn private_key_from_hex(data: &str) -> Result<Box<PrivateKey>, biscuit_auth::error::Format> {
        Ok(Box::new(PrivateKey::from_str(data)?))
    }
);
