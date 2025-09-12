extern crate alloc;
extern crate core;

mod wasm_result;
mod allocation;
mod wasm_export;
mod crypto;
mod builder;
mod token;

#[allow(unused)]
pub(crate) fn make_rng() -> rand::rngs::StdRng {
    let mut data = [0u8; 8];
    getrandom::getrandom(&mut data[..]).unwrap();
    rand::SeedableRng::seed_from_u64(u64::from_le_bytes(data))
}

unsafe extern "C" {
    pub fn print(ptr: *const u8, len: usize);
}

#[macro_export]
macro_rules! print_wasm {
    ($($args:tt)*) => {
        let msg = format!($($args)*);
        unsafe { print(msg.as_ptr(), msg.len()) };
    };
}