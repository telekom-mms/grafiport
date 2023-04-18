module grafana-exporter

go 1.19

require (
	github.com/gosimple/slug v1.1.1
	github.com/grafana/grafana-api-golang-client v0.18.4
)

require (
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/rainycape/unidecode v0.0.0-20150907023854-cb7f23ec59be // indirect
)

replace github.com/grafana/grafana-api-golang-client v0.18.4 => github.com/flkhndlr/grafana-api-golang-client v0.0.0-20230418085834-3a31f6f74a09
