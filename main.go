package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type flags struct {
	url    string
	file   string
	output string
	method string
}

func isFlaggedPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func paramFileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func CheckHTTP(url string) bool {
	r := regexp.MustCompile(`^(?i)https?:\/\/`)
	if r.MatchString(url) {
		return true
	} else {
		return false
	}
}

func readParamFile(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("[+] Unable to to open %v: %s", f, err)
		os.Exit(0)

	}
	defer f.Close()

	paramsOut := []string{}
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Text()
		paramsOut = append(paramsOut, line)
	}
	return paramsOut
}

func requestParam(method, url, param string) {
	preForwardSlash := regexp.MustCompile(`^\/`)
	if preForwardSlash.MatchString(param) {
		//fmt.Printf("[+] Removing slash from parameter: %q\n", param)
		param = preForwardSlash.ReplaceAllString(param, "")
	}
	fullUrl := url + "/" + param
	r, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		fmt.Print(err)
	}

	req, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("[+] %v : %v : %v\n", req.Request.Method, fullUrl, req.StatusCode)
}

func main() {

	m := new(flags)
	flag.StringVar(&m.url, "u", "", "Url to target")
	flag.StringVar(&m.output, "o", "/tmp/hoopla", "Output file to save params")
	flag.StringVar(&m.file, "f", "", "File to read for fuzzing params")
	flag.StringVar(&m.method, "x", "GET", "HTTP Method to apply")
	flag.Parse()

	if !isFlaggedPassed("u") {
		fmt.Print("[+] No url passed. Check -h for help. ***\n")
		os.Exit(0)
	}

	if !isFlaggedPassed("f") {
		fmt.Print("[+] No param test file. Check -h for help. ***\n")
		os.Exit(0)
	}

	if !CheckHTTP(m.url) {
		fmt.Print("[+] Make sure http/s is added. Exiting ***\n")
		os.Exit(0)
	}

	if !paramFileExists(m.file) {
		fmt.Printf("[+] File %q not found.\n", m.file)
		os.Exit(0)
	}

	params := readParamFile(m.file)
	for _, i := range params {
		requestParam(m.method, m.url, i)
	}

}
