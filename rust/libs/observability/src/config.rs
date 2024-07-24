//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
use std::collections::HashMap;
use std::time::Duration;

#[derive(Clone, Debug)]
pub struct Config {
    enabled: bool,
    attributes: HashMap<String, String>,
    tracer: Tracer,
    meter: Meter,
}

#[derive(Clone, Debug)]
pub struct Tracer {
    enabled: bool,
    endpoint: String,
}

#[derive(Clone, Debug)]
pub struct Meter {
    enabled: bool,
    endpoint: String,
    export_duration: Duration,
}

impl Config {
    pub fn new() -> Self {
        Self {
            enabled: false,
            attributes: HashMap::new(),
            tracer: Tracer::default(),
            meter: Meter::default(),
        }
    }

    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }

    pub fn attributes(mut self, attrs: HashMap<String, String>) -> Self {
        self.attributes = attrs;
        self
    }

    pub fn attribute(mut self, key: String, value: String) -> Self {
        self.attributes.insert(key, value);
        self
    }

    pub fn tracer(mut self, cfg: Tracer) -> Self {
        self.tracer = cfg;
        self
    }

    pub fn meter(mut self, cfg: Meter) -> Self {
        self.meter = cfg;
        self
    }
}

impl Tracer {
    pub fn new() -> Self {
        Self {
            enabled: false,
            endpoint: "".to_string(),
        }
    }

    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }
}

impl Default for Tracer {
    fn default() -> Self {
        Self::new()
    }
}

impl Meter {
    pub fn new() -> Self {
        Self {
            enabled: false,
            endpoint: "".to_string(),
            export_duration: Duration::from_secs(1),
        }
    }

    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }

    pub fn export_duration(mut self, dur: Duration) -> Self {
        self.export_duration = dur;
        self
    }
}

impl Default for Meter {
    fn default() -> Self {
        Self::new()
    }
}
