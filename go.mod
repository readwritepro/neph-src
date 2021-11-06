module github.com/readwritepro/neph

go 1.16

replace github.com/readwritepro/error-handler => ../error-handler

replace github.com/readwritepro/figtree => ../figtree

replace github.com/readwritepro/compare-test-results => ../compare-test-results

require (
	github.com/pkg/errors v0.8.1 // indirect
	github.com/pkg/sftp v1.13.4 // indirect
	github.com/readwritepro/compare-test-results v0.0.0-00010101000000-000000000000
	github.com/readwritepro/error-handler v0.0.0-00010101000000-000000000000
	github.com/readwritepro/figtree v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
)
