use crate::wasm_result::WasmResult;
use std::mem;

#[unsafe(no_mangle)]
/// Allocates a new area of memory and returns a pointer to it.
/// *Parameters*
/// - `size`: allocated size in bytes
/// - `align`: alignment of the allocated area
pub fn malloc(size: usize, align: usize) -> *mut u8 {
    unsafe { alloc::alloc::alloc(alloc::alloc::Layout::from_size_align(size, align).unwrap()) }
}
#[unsafe(no_mangle)]
/// Frees a previously allocated area
/// *Parameters*
/// - `ptr`: allocated address
/// - `size`: allocated size
/// - `align`: alignment of the allocated area
pub fn free(ptr: *mut u8, size: usize, align: usize) {
    unsafe {
        alloc::alloc::dealloc(
            ptr,
            alloc::alloc::Layout::from_size_align(size, align).unwrap(),
        )
    }
}

#[unsafe(no_mangle)]
/// Allocates a WasmResult area and returns a pointer to it.
/// Rust doesn't handle the allocated region, user is responsible for
/// freeing after usage
pub fn get_return_area() -> Box<WasmResult> {
    Box::new(unsafe { mem::zeroed() })
}
