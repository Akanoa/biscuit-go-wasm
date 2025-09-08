use crate::print;
use crate::print_wasm;
use std::mem;
use std::mem::MaybeUninit;

mod biscuit_builder;
mod authorizer_builder;

pub struct Builder<T> (T);

impl<B> Builder<B> {

    fn chain<F, E>(&mut self, f: F) -> Result<(), E> where F: FnOnce(B) -> Result<B, E> {

        print_wasm!("Dark magic");
        let zero = unsafe {MaybeUninit::uninit().assume_init()};
        print_wasm!("Dark magic 1.5");
        let builder = mem::replace(&mut self.0, zero);

        print_wasm!("Dark magic 2");

        let result = f(builder)?;
        print_wasm!("Dark magic 3");
        let _zeroed = mem::replace(&mut self.0, result);
        print_wasm!("Dark magic 4");

        Ok(())
    }


    fn apply<F, T, E>(self, f: F) -> Result<T, E> where F: FnOnce(B) -> Result<T, E> {

        f(self.0)
    }
}

impl<B> From<B> for Builder<B> {
    fn from(value: B) -> Self {
        Builder(value)
    }
}