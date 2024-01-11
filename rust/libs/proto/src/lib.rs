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

pub mod vald {
    pub mod v1 {
        include!("vald.v1.tonic.rs");
    }
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {}
}
