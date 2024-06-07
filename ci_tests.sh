./ci_tests/gophermarttest-darwin-arm64 \
  -test.v -test.run=^TestGophermart$ \
  -gophermart-binary-path=cmd/gophermart/gophermart \
  -gophermart-host=localhost \
  -gophermart-port=8080 \
  -gophermart-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" \
  -accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
  -accrual-host=localhost \
  -accrual-port=$(./ci_tests/random-darwin-arm64 unused-port) \
  -accrual-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" &> ./ci_tests/test.log