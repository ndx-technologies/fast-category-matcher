Fast Catagory Matching

This service provides fast category matching in taxonomy trees for natural language.
Inference is done `<5ms` which allows for real-time normalisation of approximate categories (e.g. LLM, OCR).
Zero memory allocations allows to run this in a tight loop without memory pressure.
This is accomplished by using state-of-the-art (`2025-08-14`) algorithms SIFT4[^1] and Snowball[^2].

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

[^1]: https://github.com/ndx-technologies/sift4
[^2]: https://github.com/kljensen/snowball
