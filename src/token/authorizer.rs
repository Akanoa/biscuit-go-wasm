use crate::wasm_result::WasmResult;
use crate::wasm_export;
use biscuit_auth::datalog::RunLimits;
use biscuit_auth::Authorizer;
use std::time::Duration;

// create a new authorizer builder
// Input:
// max_time: the maximum time in milliseconds for the authorizer to run
// Output:
// returnArea { data, data_len=0, is_ok = 1 }
// data is a pointer to the run limits allocated in the wasm memory
// data_len is 0 because of the opaque type
// is_ok is 1 because the function never fails
wasm_export!(
    fn authorizer_run_limits(max_time: u64) -> Box<RunLimits> {
        Box::new(RunLimits {
            max_time: Duration::from_millis(max_time),
            ..Default::default()
        })
    }
);

// create a new authorizer builder
// Input:
// authorizer_builder: a pointer to the authorizer builder allocated in the wasm memory
// run_limits: a pointer to the run limits allocated in the wasm memory
// Output:
// returnArea { data, data_len, is_ok }
//
// if is_ok = 1
// data is the index of the first matching policy
// data_len is 0 because of a number type
//
// if is_ok = 0
// data is the pointer to the error message
// data_len is the length of the error message
wasm_export!(
    fn authorizer_authorize(authorizer: &mut Authorizer,run_limits: Box<RunLimits>) -> Result<u32, biscuit_auth::error::Token> {
        Ok(authorizer.authorize_with_limits(*run_limits)? as u32)
    }
);

// print the world
// Input:
// authorizer: a pointer to the authorizer allocated in the wasm memory
// Output:
// returnArea { data, data_len, is_ok=1 }
//
// data is the world in string format
// data_len is the length of the world
// is_ok is 1 because the function never fails
wasm_export!(
    fn authorizer_print_world(authorizer: &Authorizer) -> String {
        authorizer.print_world()
    }
);

// drop a authorizer
// Input:
// authorizer: a pointer to an authorizer allocated in the wasm memory
// Output:
// returnArea { data=0, data_len=0, is_ok=1 }
//
// data is 0 because no data is returned
// data_len is 0 because no data is returned
// is_ok is 1 because the function never fails
wasm_export!(
    fn authorizer_drop(authorizer: Box<Authorizer>) {
        drop(authorizer);
    }
);