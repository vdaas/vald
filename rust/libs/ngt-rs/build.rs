fn main() -> miette::Result<()> {let current_dir = std::env::current_dir().unwrap();
    println!("cargo:rustc-link-search=native={}", current_dir.display());
    println!("cargo:rustc-link-search=native={}", current_dir.display());

    cxx_build::bridge("src/lib.rs")
        .file("src/input.cpp")
        .flag_if_supported("-std=c++20")
        .compile("ngt-rs");

    println!("cargo:rustc-link-lib=static=ngt");
    println!("cargo:rustc-link-lib=static=gomp");
    println!("cargo:rerun-if-changed=src/*");
    Ok(())
}