use std::mem;
use std::mem::MaybeUninit;
use std::ops::{Deref, DerefMut};

mod biscuit_builder;
mod authorizer_builder;

pub struct Builder<T> (T);

impl<B> Builder<B> {

    /// Get the inner builder from the heap and replace it by "zeroed" area
    /// Execute the closure on the inner builder
    /// Put back the inner builder at its original place
    fn apply<F, E>(&mut self, f: F) -> Result<(), E> where F: FnOnce(B) -> Result<B, E> {

        #[allow(clippy::uninit_assumed_init)]
        let zero = unsafe {MaybeUninit::uninit().assume_init()};

        // Take the inner builder and replace it with a zeroed area
        let builder = mem::replace(&mut self.0, zero);

        // Execute the closure on the inner builder
        let result = f(builder)?;

        // Put back the inner builder at its original place
        let _zeroed = mem::replace(&mut self.0, result);

        Ok(())
    }
    
    fn apply_no_return<F>(&mut self, f: F) where F: FnOnce(B) -> B {
        #[allow(clippy::uninit_assumed_init)]
        let zero = unsafe {MaybeUninit::uninit().assume_init()};

        // Take the inner builder and replace it with a zeroed area
        let builder = mem::replace(&mut self.0, zero);

        // Execute the closure on the inner builder
        let builder = f(builder);

        // Put back the inner builder at its original place
        let _zeroed = mem::replace(&mut self.0, builder);
        
    }
}

impl<B> From<B> for Builder<B> {
    fn from(value: B) -> Self {
        Builder(value)
    }
}

impl<B> Deref for Builder<B> {
    type Target = B;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

impl<B> DerefMut for Builder<B> {
    fn deref_mut(&mut self) -> &mut Self::Target {
        &mut self.0
    }
}