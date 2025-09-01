pub use biscuit_wasm::*;
use wasm_bindgen::prelude::wasm_bindgen;
use wasm_bindgen::JsValue;

#[wasm_bindgen(js_name = public_key_ToString)]
pub fn public_key_to_string(public_key: &PublicKey) -> Result<String, JsValue> {
    Ok(public_key.to_string())
}