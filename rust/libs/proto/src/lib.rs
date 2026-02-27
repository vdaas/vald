//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
#[allow(clippy::all)]
pub mod core {
    pub mod v1 {
        include!("core/v1/core.v1.tonic.rs");
    }
}
#[allow(clippy::all)]
pub mod discoverer {
    pub mod v1 {
        include!("discoverer/v1/discoverer.v1.tonic.rs");
    }
}
#[allow(clippy::all)]
pub mod filter {
    pub mod egress {
        pub mod v1 {
            include!("filter/egress/v1/filter.egress.v1.tonic.rs");
        }
    }
    pub mod ingress {
        pub mod v1 {
            include!("filter/ingress/v1/filter.ingress.v1.tonic.rs");
        }
    }
}
#[allow(clippy::all)]
pub mod google {
    pub mod protobuf {
        include!(concat!(env!("OUT_DIR"), "/google.protobuf.rs"));
    }
    pub mod rpc {
        include!("google/rpc/status.rs");
    }
}
#[allow(clippy::all)]
pub mod meta {
    pub mod v1 {
        include!("meta/v1/meta.v1.tonic.rs");
    }
}
#[allow(clippy::all)]
pub mod mirror {
    pub mod v1 {
        include!("mirror/v1/mirror.v1.tonic.rs");
    }
}
#[allow(clippy::all)]
pub mod payload {
    pub mod v1 {
        include!("payload/v1/payload.v1.rs");
    }
}
#[allow(clippy::all)]
pub mod rpc {
    pub mod v1 {
        include!("rpc/v1/rpc.v1.rs");
        include!("rpc/v1/rpc.v1.tonic.rs");
    }
}
#[allow(clippy::all)]
pub mod sidecar {
    pub mod v1 {
        include!("sidecar/v1/sidecar.v1.tonic.rs");
    }
}
#[allow(clippy::all)]
pub mod vald {
    pub mod v1 {
        include!("vald/v1/vald.v1.tonic.rs");
    }
}
