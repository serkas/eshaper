# Event Shaper

Shaper can enforce limit on maximum number of code executions per instance of time.
It is safe to adjust the rate limit in the runtime.

It can be usable in cases when:
* the source of events produces bursts of load followed by periods of silence
* we need to flatten the load before next stages of processing
* we cannot drop the events. *If dropping events is ok, maybe check [cloudfoundry/go-diodes](https://github.com/cloudfoundry/go-diodes)*
* there is no back-pressure mechanism from the downstream stages


![Shaping of bursts](./chart.png)
