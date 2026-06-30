# pipeline-testing

A mock vulnerability-findings stream pipeline in Go. Built test-first, no real broker.

Decode -> validate -> route -> dedup -> produce, with bad input dead-lettered.

## Layout

```
findings-pipeline/
├── go.mod
├── finding.go        # Finding, Severity, Topic
├── route.go          # Route: severity -> topic
├── validate.go       # Finding.Validate + sentinel errors
├── decode.go         # Decode: json bytes -> Finding
├── producer.go       # Message + the narrow Producer interface
├── pipeline.go       # Pipeline.Handle: decode -> validate -> dedup -> route -> produce
├── dedup.go          # concurrency-safe Deduper
└── *_test.go         # tests
```

## Running

```bash
go test        # all unit tests

```