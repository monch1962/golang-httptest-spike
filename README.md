[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/monch1962/golang-httptest-spike)

[![Build Status](https://dev.azure.com/monch1962/monch1962/_apis/build/status/monch1962.golang-httptest-spike?branchName=master)](https://dev.azure.com/monch1962/monch1962/_build/latest?definitionId=10&branchName=master)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=monch1962_golang-httptest-spike&metric=alert_status)](https://sonarcloud.io/dashboard?id=monch1962_golang-httptest-spike)

# golang-httptest-spike

## To output JUnit format

`$ go get -u github.com/jstemmer/go-junit-report`

`$ BASE_URL=https://jsonplaceholder.typicode.com go test -v 2>&1 | go-junit-report`

## To compile tests into a standalone executable file

`$ go test -c -o main_test`

Then
`$ BASE_URL=https://jsonplaceholder.typicode.com HTTP_TIMEOUT=1000 ./main_test -test.v` will run the tests in verbose mode, and `$ BASE_URL=https://jsonplaceholder.typicode.com HTTP_TIMEOUT=1000 ./main_test` will run them without verbose mode.

Finally, to generate Junit reports from a compiled test file, 

`$ BASE_URL=https://jsonplaceholder.typicode.com HTTP_TIMEOUT=1000 ./main_test -test.v | go-junit-report`


`$ BASE_URL=https://jsonplaceholder.typicode.com go ./main_test | go-junit-report`