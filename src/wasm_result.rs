use std::mem;

/// Struct to return bytes from a function
/// which returns a Result<Box<T>, E>
/// where E implements ToString
#[repr(C)]
pub struct WasmResult {
    /// Pointer to the bytes
    ptr: *const u8,
    /// Length of bytes
    ///
    /// - length is 0 when ptr is Box<T>
    /// - length != 0 when ptr is String or Vec<u8>
    len: usize,
    /// Whether the bytes are about data or an error
    /// true: Box<T>
    /// false: String
    is_ok: bool,
}

/// Trait to data into WasmResult
pub trait IntoWasmResult {
    fn into_wasm_result(self, ret: &mut WasmResult);
}
impl IntoWasmResult for String {
    fn into_wasm_result(self, ret: &mut WasmResult) {
        ret.ptr = self.as_ptr();
        ret.len = self.len();
        ret.is_ok = true;
        // leak the string dropped by the caller
        mem::forget(self);
    }
}

impl IntoWasmResult for Vec<u8> {
    fn into_wasm_result(self, ret: &mut WasmResult) {
        ret.ptr = self.as_ptr();
        ret.len = self.len();
        ret.is_ok = true;
        // leak the vec dropped by the caller
        mem::forget(self);
    }
}

impl<T> IntoWasmResult for Box<T> {
    fn into_wasm_result(self, ret: &mut WasmResult) {
        // Box<T> is a pointer to the data
        // T is consumed by the Box creation
        ret.ptr = Box::into_raw(self) as *const u8;
        ret.len = 0;
        ret.is_ok = true;
    }
}

impl IntoWasmResult for () {
    fn into_wasm_result(self, ret: &mut WasmResult) {
        ret.ptr = std::ptr::null();
        ret.len = 0;
        ret.is_ok = true;
    }
}

impl<T: IntoWasmResult, E: ToString> IntoWasmResult for Result<T, E> {
    fn into_wasm_result(self, ret: &mut WasmResult) {
        match self {
            // Return the data as bytes
            Ok(ok) => {
                ok.into_wasm_result(ret)
            }
            // Return the error as a string
            Err(err) => {
                let msg = err.to_string();
                ret.ptr = msg.as_ptr();
                ret.len = msg.len();
                ret.is_ok = false;
                mem::forget(msg);
            }
        }
    }
}

impl WasmResult {
    /// Captures a value and converts it into WebAssembly-compatible bytes.
    ///
    /// This function takes a value that implements the `IntoWasmResult` trait
    /// and converts it into WebAssembly-compatible bytes. The conversion is performed
    /// by calling the `into_wasm_bytes` method of the `IntoWasmResult` trait. The
    /// resulting bytes are then stored or used by the current context (`self`).
    ///
    /// # Type Parameters
    /// - `T`: A type that implements the `IntoWasmResult` trait, allowing it
    ///   to be converted into WebAssembly-compatible bytes.
    ///
    /// # Parameters
    /// - `value`: The value to be captured and converted into WebAssembly-compatible bytes.
    pub fn capture<T: IntoWasmResult>(&mut self, value: T) {
        value.into_wasm_result(self)
    }
}