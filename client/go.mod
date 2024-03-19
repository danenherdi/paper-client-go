module github.com/Paper-Book/paper-client-go/client

go 1.18

require internal/sheet_writer v1.0.0
replace internal/sheet_writer => ./internal/sheet_writer

require internal/sheet_reader v1.0.0
replace internal/sheet_reader => ./internal/sheet_reader

require internal/tcp_client v1.0.0
replace internal/tcp_client => ./internal/tcp_client

require internal/response v1.0.0
replace internal/response => ./internal/response
