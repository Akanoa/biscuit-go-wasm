use biscuit_auth::Algorithm;

mod keypair;
mod private_key;
mod public_key;

#[derive(Default)]
#[repr(C)]
#[allow(non_camel_case_types)]
pub enum SignatureAlgorithm {
    #[default]
    Ed25519 = 0,
    Secp256r1 = 1,
}

impl From<SignatureAlgorithm> for Algorithm {
    fn from(algorithm: SignatureAlgorithm) -> Self {
        match algorithm {
            SignatureAlgorithm::Ed25519 => Algorithm::Ed25519,
            SignatureAlgorithm::Secp256r1 => Algorithm::Secp256r1,
        }
    }
}