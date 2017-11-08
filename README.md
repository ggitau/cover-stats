# cover-stats
utility that gets the average (total) coverage of golang package tests. You can integrate this tool with your
CI build steps if you require your golang tests to have a particular coverage threshold.
## Installation
`go get github.com/ggitau/cover-stats`
## Usage
* If you want to just get the average coverage
`cover-stats`
* If you want to test that your coverage meets a threshold run the command below. If the threshold is not met the program
will exit with a non zero status.
`cover-stats -threshold=50.5`
* The program uses the location /usr/local/ as GOROOT by default. You can override this by setting a value for the GOROOT
 environment variable before running the program. The program will then look for the go executable
 in the location $GOROOT/bin/go
