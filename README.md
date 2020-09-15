[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/monch1962/golang-httptest-spike)

[![Build Status](https://dev.azure.com/monch1962/monch1962/_apis/build/status/monch1962.golang-httptest-spike?branchName=master)](https://dev.azure.com/monch1962/monch1962/_build/latest?definitionId=10&branchName=master)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=monch1962_golang-httptest-spike&metric=alert_status)](https://sonarcloud.io/dashboard?id=monch1962_golang-httptest-spike)

# golang-httptest-spike

Intention is to produce a Go/httptest _template_ that can be used by Pactical (https://github.com/monch1962/pactical) to generate Go test cases from Pacts. My thinking is:
- Go tests execute particularly quickly, and consume very few hardware resources
- Go code compiles into very small EXEs with no external dependencies
- Go code can be cross-compiled for just about any platform
- it's fairly easy to turn Go code into serverless code that will execute on AWS, Azure, GCP and KNative platforms

Combine all these features, and you get a particularly flexible API testing framework that's suitable for many different test platforms. It'd be very nice to have a single framework that can execute with zero dependencies on just about any infrastructure.


## Simplest case

To run without compiling and get test results in Go test format:

`$ BASE_URL=https://jsonplaceholder.typicode.com go test`

## To output JUnit format

`$ go get -u github.com/jstemmer/go-junit-report`

`$ BASE_URL=https://jsonplaceholder.typicode.com go test -v 2>&1 | go-junit-report`

## To compile tests into a standalone executable file

`$ go test -c -o main_test`

will compile all the tests into a standalone EXE called `main_test`, which can be moved and executed on other hardware

`$ BASE_URL=https://jsonplaceholder.typicode.com HTTP_TIMEOUT=1000 ./main_test -test.v` will run the tests in verbose mode, and `$ BASE_URL=https://jsonplaceholder.typicode.com HTTP_TIMEOUT=1000 ./main_test` will run them without verbose mode.

Finally, to generate Junit reports from a compiled test file, 

`$ BASE_URL=https://jsonplaceholder.typicode.com HTTP_TIMEOUT=1000 ./main_test -test.v | go-junit-report`
