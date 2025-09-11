use core::mem;

mod authorizer_builder;
mod biscuit_builder;

/// Apply a closure to a builder on the heap.
/// If the closure succeeds, the builder is modified in place, slot was not reallocated, this
/// function behaves just like it mutably borrowed slot.
/// If the closure fails, the builder is dropped and the slot freed, just as if it took ownership
/// of slot.
///
/// This is not intended to be used from Rust, but called from extern program
fn apply_in_place<B, F, E>(mut slot: Box<B>, f: F) -> Result<(), E>
where
    F: FnOnce(B) -> Result<B, E>,
{
    #[allow(clippy::uninit_assumed_init)]
    let tmp = unsafe { mem::MaybeUninit::uninit().assume_init() };

    // Take the inner builder and replace it with temporary uninitialized data
    let builder = mem::replace(&mut *slot, tmp);

    // Execute the closure on the inner builder
    match f(builder) {
        Ok(result) => {
            // Put back the inner builder at its original place
            let _tmp = mem::replace(&mut *slot, result);
            mem::forget(slot);
            Ok(())
        }
        Err(error) => {
            // Get back the temporary value, do not run drop on its content since it is not
            // initialized, drop the Box
            let tmp: B = *slot;
            mem::forget(tmp);
            Err(error)
        }
    }
}
