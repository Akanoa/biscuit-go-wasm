use crate::crypto::SignatureAlgorithm;
use crate::wasm_result::WasmResult;
use crate::{make_rng, wasm_export};
use biscuit_auth::{KeyPair, PrivateKey, PublicKey};

// create a new keypair with the given signature algorithm
// Input:
// signature_algorithm is an u32 value
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the keypair allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn keypair_new(signature_algorithm: SignatureAlgorithm) -> Box<KeyPair> {
        Box::new(KeyPair::new_with_rng(signature_algorithm.into(), &mut make_rng()))
    }
);

// drop the keypair and free the memory allocated for it
// Input:
// keypair is a pointer to the keypair allocated in the wasm memory
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is 0
// data_len is 0
// is_ok is 1 because the function never fails
wasm_export!(
    fn keypair_drop(keypair: Box<KeyPair>) {
        drop(keypair);
    }
);

// get the public key from the keypair
// Input:
// keypair is a pointer to the keypair allocated in the wasm memory
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the public key allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn keypair_public_key(keypair: &KeyPair) -> Box<PublicKey> {
        Box::new(keypair.public())
    }
);

// get the private key from the keypair
// Input:
// keypair is a pointer to the keypair allocated in the wasm memory
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the private key allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn keypair_private_key(keypair: &KeyPair) -> Box<PrivateKey> {
        Box::new(keypair.private())
    }
);

// create a new keypair from a private key
// Input:
// private_key is a pointer to the private key allocated in the wasm memory
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the keypair allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn keypair_from_private_key(private_key: &PrivateKey) -> Box<KeyPair> {
        Box::new(KeyPair::from(private_key))
    }
);