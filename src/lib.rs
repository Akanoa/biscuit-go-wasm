extern crate alloc;
extern crate core;

mod allocation;
mod builder;
mod crypto;
mod token;
mod wasm_export;
mod wasm_result;

#[allow(unused)]
pub(crate) fn make_rng() -> rand::rngs::StdRng {
    let mut data = [0u8; 8];
    getrandom::getrandom(&mut data[..]).unwrap();
    rand::SeedableRng::seed_from_u64(u64::from_le_bytes(data))
}

#[cfg(feature = "print_wasm")]
unsafe extern "C" {
    pub fn print(ptr: *const u8, len: usize);
}

#[macro_export]
macro_rules! print_wasm {
    ($($args:tt)*) => {
        #[cfg(feature = "print_wasm")]
        {
            let msg = format!($($args)*);
            unsafe { $crate::print(msg.as_ptr(), msg.len()) };
        }
    };
}
