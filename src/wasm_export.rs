
#[macro_export]
macro_rules! wasm_export {
    (fn $f:ident($($arg_name:ident: $arg_type:ty),* $(,)?) $(-> $ret:ty)? $body:block) => {
        #[unsafe(no_mangle)]
        pub fn $f(__ret: &mut WasmResult, $($arg_name: $arg_type),*) {
            #[allow(unused_mut)]
            let mut ret = (|| $(-> $ret)? {
                $body
            });
            __ret.capture(ret())
        }
    };
}