## KZG Commitments using go-mcl

KZG commitments using go-mcl and go-kzg
- [go-mcl](https://github.com/alinush/go-mcl/)
- [go-kzg](https://github.com/protolambda/go-kzg)

```bash
time go test ./... -bench=. -run=^a -benchtime=100x -timeout 240m
time go test ./... -bench=BenchmarkPoly -run=^a -benchtime=100x  -timeout 240m
```
