# cover-stats
utility that gets the average (total) coverage of golang package tests. You can integrate this tool with your
CI build steps if you require your golang tests to have a particular coverage threshold.
## Installation
`go get github.com/ggitau/cover-stats`
## Usage
* If you want to just get the average coverage
`go test -cover ./... . | cover-stats`
* If you want to test that your coverage meets a threshold run the command below. If the threshold is not met the program
will exit with a non zero status.
`go test -cover ./... . | cover-stats -threshold=50.5`
* You can also supply a file that contains the output of `go test cover ./... .` by running the `cover-stats` like this
`cover-stats -file=path-to-file`