// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/revel/cmd/harness"
	"github.com/revel/modules/testrunner/app/controllers"
	"github.com/revel/revel"
)

var cmdTest = &Command{
	UsageLine: "test [import path] [run mode] [suite.method]",
	Short:     "run all tests from the command-line",
	Long: `
Run all tests for the Revel app named by the given import path.

For example, to run the booking sample application's tests:

    revel test github.com/revel/examples/booking dev

The run mode is used to select which set of app.conf configuration should
apply and may be used to determine logic in the application itself.

Run mode defaults to "dev".

You can run a specific suite (and function) by specifying a third parameter.
For example, to run all of UserTest:

    revel test outspoken test UserTest

or one of UserTest's methods:

    revel test outspoken test UserTest.Test1
`,
}

func init() {
	cmdTest.Run = testApp
}

func testApp(args []string) {
	var err error
	if len(args) == 0 {
		errorf("No import path given.\nRun 'revel help test' for usage.\n")
	}

	mode := DefaultRunMode
	if len(args) >= 2 {
		mode = args[1]
	}

	// Find and parse app.conf
	revel.Init(mode, args[0], "")

	// Ensure that the testrunner is loaded in this mode.
	checkTestRunner()

	// Create a directory to hold the test result files.
	resultPath := filepath.Join(revel.BasePath, "test-results")
	if err = os.RemoveAll(resultPath); err != nil {
		errorf("Failed to remove test result directory %s: %s", resultPath, err)
	}
	if err = os.Mkdir(resultPath, 0777); err != nil {
		errorf("Failed to create test result directory %s: %s", resultPath, err)
	}

	// Direct all the output into a file in the test-results directory.
	file, err := os.OpenFile(filepath.Join(resultPath, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		errorf("Failed to create test result log file: %s", err)
	}

	app, reverr := harness.Build()
	if reverr != nil {
		errorf("Error building: %s", reverr)
	}
	cmd := app.Cmd()
	cmd.Stderr = io.MultiWriter(cmd.Stderr, file)
	cmd.Stdout = io.MultiWriter(cmd.Stderr, file)

	// Start the app...
	if err := cmd.Start(); err != nil {
		errorf("%s", err)
	}
	defer cmd.Kill()
	revel.INFO.Printf("Testing %s (%s) in %s mode\n", revel.AppName, revel.ImportPath, mode)

	var httpAddr = revel.HTTPAddr
	if httpAddr == "" {
		httpAddr = "127.0.0.1"
	}

	var httpProto = "http"
	if revel.HTTPSsl {
		httpProto = "https"
	}

	// Get a list of tests
	var baseURL = fmt.Sprintf("%s://%s:%d", httpProto, httpAddr, revel.HTTPPort)
	testSuites, _ := getTestsList(baseURL)

	// If a specific TestSuite[.Method] is specified, only run that suite/test
	if len(args) == 3 {
		testSuites = filterTestSuites(testSuites, args[2])
	}
	testSuiteCount := len(*testSuites)
	fmt.Printf("\n%d test suite%s to run.\n", testSuiteCount, pluralize(testSuiteCount, "", "s"))
	fmt.Println()

	// Run each suite.
	failedResults, overallSuccess := runTestSuites(baseURL, resultPath, testSuites)

	fmt.Println()
	if overallSuccess {
		writeResultFile(resultPath, "result.passed", "passed")
		fmt.Println("All Tests Passed.")
	} else {
		for _, failedResult := range *failedResults {
			fmt.Printf("Failures:\n")
			for _, result := range failedResult.Results {
				if !result.Passed {
					fmt.Printf("%s.%s\n", failedResult.Name, result.Name)
					fmt.Printf("%s\n\n", result.ErrorSummary)
				}
			}
		}
		writeResultFile(resultPath, "result.failed", "failed")
		errorf("Some tests failed.  See file://%s for results.", resultPath)
	}
}

func writeResultFile(resultPath, name, content string) {
	if err := ioutil.WriteFile(filepath.Join(resultPath, name), []byte(content), 0666); err != nil {
		errorf("Failed to write result file %s: %s", filepath.Join(resultPath, name), err)
	}
}

func pluralize(num int, singular, plural string) string {
	if num == 1 {
		return singular
	}
	return plural
}

// Filters test suites and individual tests to match
// the parsed command line parameter
func filterTestSuites(suites *[]controllers.TestSuiteDesc, suiteArgument string) *[]controllers.TestSuiteDesc {
	var suiteName, testName string
	argArray := strings.Split(suiteArgument, ".")
	suiteName = argArray[0]
	if suiteName == "" {
		return suites
	}
	if len(argArray) == 2 {
		testName = argArray[1]
	}
	for _, suite := range *suites {
		if suite.Name != suiteName {
			continue
		}
		if testName == "" {
			return &[]controllers.TestSuiteDesc{suite}
		}
		// Only run a particular test in a suite
		for _, test := range suite.Tests {
			if test.Name != testName {
				continue
			}
			return &[]controllers.TestSuiteDesc{
				{
					Name:  suite.Name,
					Tests: []controllers.TestDesc{test},
				},
			}
		}
		errorf("Couldn't find test %s in suite %s", testName, suiteName)
	}
	errorf("Couldn't find test suite %s", suiteName)
	return nil
}

func checkTestRunner() {
	testRunnerFound := false
	for _, module := range revel.Modules {
		if module.ImportPath == revel.Config.StringDefault("module.testrunner", "github.com/revel/modules/testrunner") {
			testRunnerFound = true
			break
		}
	}

	if !testRunnerFound {
		errorf(`Error: The testrunner module is not running.

You can add it to a run mode configuration with the following line:

	module.testrunner = github.com/revel/modules/testrunner

`)
	}
}

// Get a list of tests from server.
// Since this is the first request to the server, retry/sleep a couple times
// in case it hasn't finished starting up yet.
func getTestsList(baseURL string) (*[]controllers.TestSuiteDesc, error) {
	var (
		err        error
		resp       *http.Response
		testSuites []controllers.TestSuiteDesc
	)
	for i := 0; ; i++ {
		if resp, err = http.Get(baseURL + "/@tests.list"); err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		if i < 3 {
			time.Sleep(3 * time.Second)
			continue
		}
		if err != nil {
			errorf("Failed to request test list: %s", err)
		} else {
			errorf("Failed to request test list: non-200 response")
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&testSuites)

	return &testSuites, err
}

func runTestSuites(baseURL, resultPath string, testSuites *[]controllers.TestSuiteDesc) (*[]controllers.TestSuiteResult, bool) {
	// Load the result template, which we execute for each suite.
	module, _ := revel.ModuleByName("testrunner")
	TemplateLoader := revel.NewTemplateLoader([]string{filepath.Join(module.Path, "app", "views")})
	if err := TemplateLoader.Refresh(); err != nil {
		errorf("Failed to compile templates: %s", err)
	}
	resultTemplate, err := TemplateLoader.Template("TestRunner/SuiteResult.html")
	if err != nil {
		errorf("Failed to load suite result template: %s", err)
	}

	var (
		overallSuccess = true
		failedResults  []controllers.TestSuiteResult
	)
	for _, suite := range *testSuites {
		// Print the name of the suite we're running.
		name := suite.Name
		if len(name) > 22 {
			name = name[:19] + "..."
		}
		fmt.Printf("%-22s", name)

		// Run every test.
		startTime := time.Now()
		suiteResult := controllers.TestSuiteResult{Name: suite.Name, Passed: true}
		for _, test := range suite.Tests {
			testURL := baseURL + "/@tests/" + suite.Name + "/" + test.Name
			resp, err := http.Get(testURL)
			if err != nil {
				errorf("Failed to fetch test result at url %s: %s", testURL, err)
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			var testResult controllers.TestResult
			err = json.NewDecoder(resp.Body).Decode(&testResult)
			if err == nil && !testResult.Passed {
				suiteResult.Passed = false
			}
			suiteResult.Results = append(suiteResult.Results, testResult)
		}
		overallSuccess = overallSuccess && suiteResult.Passed

		// Print result.  (Just PASSED or FAILED, and the time taken)
		suiteResultStr, suiteAlert := "PASSED", ""
		if !suiteResult.Passed {
			suiteResultStr, suiteAlert = "FAILED", "!"
			failedResults = append(failedResults, suiteResult)
		}
		fmt.Printf("%8s%3s%6ds\n", suiteResultStr, suiteAlert, int(time.Since(startTime).Seconds()))
		// Create the result HTML file.
		suiteResultFilename := filepath.Join(resultPath,
			fmt.Sprintf("%s.%s.html", suite.Name, strings.ToLower(suiteResultStr)))
		suiteResultFile, err := os.Create(suiteResultFilename)
		if err != nil {
			errorf("Failed to create result file %s: %s", suiteResultFilename, err)
		}
		if err = resultTemplate.Render(suiteResultFile, suiteResult); err != nil {
			errorf("Failed to render result template: %s", err)
		}
	}

	return &failedResults, overallSuccess
}
