//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
#[macro_export]
macro_rules! tracer {
    () => {{
        tracer!("vald")
    }};

    ($name:expr) => {{
        opentelemetry::global::tracer($name)
    }};
}

#[macro_export]
macro_rules! ctx_span {
    ($ctx:expr, $name:expr) => {{
        ctx_span!($ctx, $name, opentelemetry::trace::SpanKind::Internal)
    }};

    ($ctx:expr, $name:expr, $kind:expr) => {{
        let tracer = tracer!();
        let parent_ctx: &opentelemetry::Context = $ctx;
        let span = tracer
            .span_builder($name)
            .with_kind($kind)
            .start_with_context(&tracer, parent_ctx);
        opentelemetry::Context::current_with_span(span)
    }};
}

#[macro_export]
macro_rules! meter {
    () => {{
        meter!("vald")
    }};

    ($name:expr) => {{
        opentelemetry::global::meter($name)
    }};
}

#[derive(Debug, PartialEq)]
pub enum InstrumentKind {
    UpdownCounter,
    Counter,
    Histogram,
    Gauge,
}

#[macro_export]
macro_rules! instrument {
    (InstrumentKind::Counter, $typ:ty, $name:expr, $disc:expr, $unit:expr) => {{
        let meter = meter!();
        paste::paste! {
          meter
            .[<$typ _counter>]($name) // typ = f64 or u64
            .with_description($disc)
            .with_unit($unit)
            .init()
        }
    }};

    (InstrumentKind::Counter, $typ:ty, $name:expr, $disc:expr, $unit:expr, $measurement:expr, $key_value:expr) => {{
        let meter = meter!();
        paste::paste! {
          meter
            .[<$typ _observable_counter>]($name) // typ = f64 or u64
            .with_description($disc)
            .with_unit($unit)
            .with_callback(|observe| {
                observe.observe($measurement, $key_value);
            })
            .init()
        }
    }};

    (InstrumentKind::UpdownCounter, $typ:ty, $name:expr, $disc:expr, $unit:expr) => {{
        let meter = meter!();
        paste::paste! {
          meter
            .[<$typ _up_down_counter>]($name) // typ = f64 or i64
            .with_description($disc)
            .with_unit($unit)
            .init()
        }
    }};

    (InstrumentKind::UpdownCounter, $typ:ty, $name:expr, $disc:expr, $unit:expr, $measurement:expr, $key_value:expr) => {{
        let meter = meter!();
        paste::paste! {
          meter
            .[<$typ _observable_up_down_counter>]($name) // typ = f64 or i64
            .with_description($disc)
            .with_unit($unit)
            .with_callback(|observe| {
                observe.observe($measurement, $key_value);
            })
            .init()
        }
    }};

    (InstrumentKind::Histogram, $typ:ty, $name:expr, $disc:expr, $unit:expr) => {{
        let meter = meter!();
        paste::paste! {
          meter
            .[<$typ _histogram>]($name) // typ = f64 or i64
            .with_description($disc)
            .with_unit($unit)
            .init()
        }
    }};

    (InstrumentKind::Gauge, $typ:ty, $name:expr, $disc:expr, $unit:expr) => {{
        let meter = meter!();
        paste::paste! {
          meter
            .[<$typ _gauge>]($name) // typ = f64 or i64 or u64
            .with_description($disc)
            .with_unit($unit)
            .init()
        }
    }};

    (InstrumentKind::Gauge, $typ:ty, $name:expr, $disc:expr, $unit:expr, $measurement:expr, $key_value:expr) => {{
        let meter = meter!();
        paste::paste! {
          meter
            .[<$typ _observable_gauge>]($name) // typ = f64 or u64 or u64
            .with_description($disc)
            .with_unit($unit)
            .with_callback(|observe| {
                observe.observe($measurement, $key_value);
            })
            .init()
        }
    }};
}
