pub use biscuit_wasm::*;
use serde::{Deserialize, Serialize};
use wasm_bindgen::prelude::wasm_bindgen;
use wasm_bindgen::JsValue;


#[derive(Deserialize, Serialize, Debug, Default)]
pub struct RunLimits {
    pub max_facts: Option<u64>,
    pub max_iterations: Option<u64>,
    pub max_time_micro: Option<u64>,
}

#[wasm_bindgen(js_name = publickey_ToString)]
pub fn public_key_to_string(public_key: &PublicKey) -> Result<String, JsValue> {
    Ok(public_key.to_string())
}

#[wasm_bindgen(js_name = authorizer_getRunLimits)]
pub fn run_limits(max_execution_time: u64) -> JsValue {
    let run_limits = RunLimits {
        max_time_micro: Some(max_execution_time),
        ..RunLimits::default()
    };
    serde_wasm_bindgen::to_value(&run_limits).unwrap()
}