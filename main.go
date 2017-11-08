package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func exit(err error) {
	fmt.Printf("error:%v\n", err)
	os.Exit(1)
}

func Average(goTestOutput []string) (float64, error) {
	covRe, err := regexp.Compile("coverage:\\s\\d+\\.\\d+")
	if err != nil {
		return 0, fmt.Errorf("compile coverage regex:%v", err)
	}

	failRe, err := regexp.Compile("^FAIL")
	if err != nil {
		return 0, fmt.Errorf("compile fail regex:%v", err)
	}

	var totalCov float64
	var pkgCount float64
	for _, output := range goTestOutput {
		failMatch := failRe.FindString(output)
		if failMatch != "" {
			return 0, fmt.Errorf("tests failed")
		}

		match := covRe.FindString(output)
		if match == "" {
			continue
		}

		cov, err := strconv.ParseFloat(strings.Split(match, " ")[1], 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse coverage:%v", err)
		}

		totalCov += cov
		pkgCount = pkgCount + 1
	}

	if pkgCount == 0 {
		return 0, nil
	}
	return totalCov / pkgCount, nil
}
func goTest(verbose bool) ([]string, error) {
	var goRoot = "/usr/local"
	if root := os.Getenv("GOROOT"); root != "" {
		goRoot = os.Getenv("GOROOT")
	}
	cmd := exec.Command(filepath.Join(goRoot, "bin", "go"), "test", "-cover", "./...", ".")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not get stdin:%v", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start:%v", err)
	}

	line := make(chan string)
	var output []string
	b := bufio.NewScanner(stdout)
	go func() {
		for b.Scan() {
			if verbose {
				fmt.Println(b.Text())
			}
			line <- b.Text()
		}
		close(line)
	}()
	for l := range line {
		output = append(output, l)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("wait:%v", err)
	}
	return output, nil
}

func main() {
	var threshold float64
	var verbose bool
	flag.Float64Var(&threshold, "threshold", 0.0, "the minimum coverage")
	flag.BoolVar(&verbose, "verbose", false, "whether or not to print test output.")
	flag.Parse()

	lines, err := goTest(verbose)
	if err != nil {
		exit(fmt.Errorf("tests failed:%v", err))
	}
	avg, err := Average(lines)

	fmt.Printf("average coverage:%f\n", avg)

	if threshold == 0 {
		os.Exit(0)
	}

	if avg < threshold {
		fmt.Println("threshold not met")
		os.Exit(1)
	}

	fmt.Println("threshold met!")

}
