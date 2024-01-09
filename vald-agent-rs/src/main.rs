pub mod google {
    pub mod rpc {
        pub type Status = tonic_types::Status;
    }
}

pub mod payload {
    pub mod v1 {
        include!("payload.v1.rs");
    }
}

pub mod core {
    pub mod v1 {
        include!("core.v1.tonic.rs");
    }
}

pub mod vald {
    pub mod v1 {
        include!("vald.v1.tonic.rs");
    }
}

fn main() {
    println!("Hello, world!");
}
