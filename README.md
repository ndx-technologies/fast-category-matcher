Fast Catagory Matching

This service provides fast category matching in taxonomy trees for natural language.
Inference is done `<5ms` which allows for real-time normalisation of approximate categories (e.g. LLM, OCR).
Zero memory alloc allows this inference to be done in tight loop without memory pressure.

```bash
$ go test -bench=. -benchmem .
goos: darwin
goarch: arm64
pkg: github.com/ndx-technologies/fast-category-matcher
cpu: Apple M3 Max
BenchmarkApproxStrGoogleTaxonomyMatcher-16           314           3823583 ns/op             922 B/op         29 allocs/op
PASS
ok      github.com/ndx-technologies/fast-category-matcher       2.036s
```
