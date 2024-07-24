#[macro_export]
macro_rules! tracer {
    () => {{
        tracer!("vald")
    }};

    ($name:expr) => {{
        global::tracer($name)
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
        global::meter($name)
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
