# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	example.com/nightguide/cmd/nightguide	[no test files]
ok  	example.com/nightguide/internal/app	0.052s
ok  	example.com/nightguide/internal/booking	0.086s
ok  	example.com/nightguide/internal/catalog	0.001s
ok  	example.com/nightguide/internal/domain	0.001s
?   	example.com/nightguide/internal/guide	[no test files]
ok  	example.com/nightguide/internal/httpapi	0.071s
?   	example.com/nightguide/internal/itinerary	[no test files]
?   	example.com/nightguide/internal/ledger	[no test files]
?   	example.com/nightguide/internal/operations	[no test files]
?   	example.com/nightguide/internal/policy	[no test files]
?   	example.com/nightguide/internal/quality	[no test files]
?   	example.com/nightguide/internal/report	[no test files]
ok  	example.com/nightguide/internal/store	0.063s
--- FAIL: TestNightTourConfirmationUsesOwnMeetingPoint (0.06s)
    night_tour_test.go:23: meeting=Old Street Gate
FAIL
FAIL	example.com/nightguide/internal/tours	0.059s
ok  	example.com/nightguide/internal/workflow	0.100s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/nightguide): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/nightguide): exit `0`
- Frontend build (web): exit `0`
