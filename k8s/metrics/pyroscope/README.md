# Pyroscope

This is the mafests to deploy pyroscope server and pyroscope agent.

The pyroscope server scrapes pprof data.

Which pod data to scrape is determined by `podAnnotations`.

The following is an example of scraping data of `vald-agent-ngt`.

```
# valdrelease.yaml
agent:
  podAnnotations:
    pyroscope.io/scrape: "true"
    pyroscope.io/application-name: "vald-agent-ngt"
    pyroscope.io/profile-cpu-enabled: "true"
    pyroscope.io/profile-mem-enabled: "true"
    pyroscope.io/port: "6060" # pprof server port
```
