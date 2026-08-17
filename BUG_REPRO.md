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
?   	weddingguide/cmd/weddingguide	[no test files]
?   	weddingguide/internal/clock	[no test files]
?   	weddingguide/internal/config	[no test files]
?   	weddingguide/internal/fixture	[no test files]
?   	weddingguide/internal/ids	[no test files]
ok  	weddingguide/internal/audit	0.006s
ok  	weddingguide/internal/domain	0.001s
--- FAIL: TestBatchImportProcessesAllRecords (0.00s)
    importer_test.go:39: batch import returned error: visitor import stopped
FAIL
FAIL	weddingguide/internal/importer	0.008s
ok  	weddingguide/internal/service	0.010s
ok  	weddingguide/internal/store	0.008s
ok  	weddingguide/internal/web	0.012s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/weddingguide): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/weddingguide): exit `0`
