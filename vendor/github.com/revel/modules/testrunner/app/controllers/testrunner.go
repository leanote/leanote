// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"sort"
	"strings"

	"github.com/revel/revel"
	"github.com/revel/revel/testing"
)

// TestRunner is a controller which is used for running application tests in browser.
type TestRunner struct {
	*revel.Controller
}

// TestSuiteDesc is used for storing information about a single test suite.
// This structure is required by revel test cmd.
type TestSuiteDesc struct {
	Name  string
	Tests []TestDesc

	// Elem is reflect.Type which can be used for accessing methods
	// of the test suite.
	Elem reflect.Type
}

// TestDesc is used for describing a single test of some test suite.
// This structure is required by revel test cmd.
type TestDesc struct {
	Name string
}

// TestSuiteResult stores the results the whole test suite.
// This structure is required by revel test cmd.
type TestSuiteResult struct {
	Name    string
	Passed  bool
	Results []TestResult
}

// TestResult represents the results of running a single test of some test suite.
// This structure is required by revel test cmd.
type TestResult struct {
	Name         string
	Passed       bool
	ErrorHTML    template.HTML
	ErrorSummary string
}

var (
	testSuites []TestSuiteDesc // A list of all available tests.

	none = []reflect.Value{} // It is used as input for reflect call in a few places.

	// registeredTests simplifies the search of test suites by their name.
	// "TestSuite.TestName" is used as a key. Value represents index in testSuites.
	registeredTests map[string]int
)

/*
	Controller's action methods are below.
*/

// Index is an action which renders the full list of available test suites and their tests.
func (c TestRunner) Index() revel.Result {
	c.ViewArgs["suiteFound"] = len(testSuites) > 0
	return c.Render(testSuites)
}

// Suite method allows user to navigate to individual Test Suite and their tests
func (c TestRunner) Suite(suite string) revel.Result {
	var foundTestSuites []TestSuiteDesc
	for _, testSuite := range testSuites {
		if strings.EqualFold(testSuite.Name, suite) {
			foundTestSuites = append(foundTestSuites, testSuite)
		}
	}

	c.ViewArgs["testSuites"] = foundTestSuites
	c.ViewArgs["suiteFound"] = len(foundTestSuites) > 0
	c.ViewArgs["suiteName"] = suite

	return c.RenderTemplate("TestRunner/Index.html")
}

// Run runs a single test, given by the argument.
func (c TestRunner) Run(suite, test string) revel.Result {
	// Check whether requested test exists.
	suiteIndex, ok := registeredTests[suite+"."+test]
	if !ok {
		return c.NotFound("Test %s.%s does not exist", suite, test)
	}

	result := TestResult{Name: test}

	// Found the suite, create a new instance and run the named method.
	t := testSuites[suiteIndex].Elem
	v := reflect.New(t)
	func() {
		// When the function stops executing try to recover from panic.
		defer func() {
			if err := recover(); err != nil {
				// If panic error is empty, exit.
				panicErr := revel.NewErrorFromPanic(err)
				if panicErr == nil {
					return
				}

				// Otherwise, prepare and format the response of server if possible.
				testSuite := v.Elem().FieldByName("TestSuite").Interface().(testing.TestSuite)
				res := formatResponse(testSuite)

				// Render the error and save to the result structure.
				var buffer bytes.Buffer
				tmpl, _ := revel.MainTemplateLoader.TemplateLang("TestRunner/FailureDetail.html", "")
				_ = tmpl.Render(&buffer, map[string]interface{}{
					"error":    panicErr,
					"response": res,
					"postfix":  suite + "_" + test,
				})
				result.ErrorSummary = errorSummary(panicErr)
				result.ErrorHTML = template.HTML(buffer.String())
			}
		}()

		// Initialize the test suite with a NewTestSuite()
		testSuiteInstance := v.Elem().FieldByName("TestSuite")
		testSuiteInstance.Set(reflect.ValueOf(testing.NewTestSuite()))

		// Make sure After method will be executed at the end.
		if m := v.MethodByName("After"); m.IsValid() {
			defer m.Call(none)
		}

		// Start from running Before method of test suite if exists.
		if m := v.MethodByName("Before"); m.IsValid() {
			m.Call(none)
		}

		// Start the test method itself.
		v.MethodByName(test).Call(none)

		// No panic means success.
		result.Passed = true
	}()

	return c.RenderJSON(result)
}

// List returns a JSON list of test suites and tests.
// It is used by revel test command line tool.
func (c TestRunner) List() revel.Result {
	return c.RenderJSON(testSuites)
}

/*
	Below are helper functions.
*/

// describeSuite expects testsuite interface as input parameter
// and returns its description in a form of TestSuiteDesc structure.
func describeSuite(testSuite interface{}) TestSuiteDesc {
	t := reflect.TypeOf(testSuite)

	// Get a list of methods of the embedded test type.
	// It will be used to make sure the same tests are not included in multiple test suites.
	super := t.Elem().Field(0).Type
	superMethods := map[string]bool{}
	for i := 0; i < super.NumMethod(); i++ {
		// Save the current method's name.
		superMethods[super.Method(i).Name] = true
	}

	// Get a list of methods on the test suite that take no parameters, return
	// no results, and were not part of the embedded type's method set.
	var tests []TestDesc
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		mt := m.Type

		// Make sure the test method meets the criterias:
		// - method of testSuite without input parameters;
		// - nothing is returned;
		// - has "Test" prefix;
		// - doesn't belong to the embedded structure.
		methodWithoutParams := (mt.NumIn() == 1 && mt.In(0) == t)
		nothingReturned := (mt.NumOut() == 0)
		hasTestPrefix := (strings.HasPrefix(m.Name, "Test"))
		if methodWithoutParams && nothingReturned && hasTestPrefix && !superMethods[m.Name] {
			// Register the test suite's index so we can quickly find it by test's name later.
			registeredTests[t.Elem().Name()+"."+m.Name] = len(testSuites)

			// Add test to the list of tests.
			tests = append(tests, TestDesc{m.Name})
		}
	}

	return TestSuiteDesc{
		Name:  t.Elem().Name(),
		Tests: tests,
		Elem:  t.Elem(),
	}
}

// errorSummary gets an error and returns its summary in human readable format.
func errorSummary(err *revel.Error) (message string) {
	expectedPrefix := "(expected)"
	actualPrefix := "(actual)"
	errDesc := err.Description
	//strip the actual/expected stuff to provide more condensed display.
	if strings.Index(errDesc, expectedPrefix) == 0 {
		errDesc = errDesc[len(expectedPrefix):]
	}
	if strings.LastIndex(errDesc, actualPrefix) > 0 {
		errDesc = errDesc[0 : len(errDesc)-len(actualPrefix)]
	}

	errFile := err.Path
	slashIdx := strings.LastIndex(errFile, "/")
	if slashIdx > 0 {
		errFile = errFile[slashIdx+1:]
	}

	message = fmt.Sprintf("%s %s#%d", errDesc, errFile, err.Line)

	/*
		// If line of error isn't known return the message as is.
		if err.Line == 0 {
			return
		}

		// Otherwise, include info about the line number and the relevant
		// source code lines.
		message += fmt.Sprintf(" (around line %d): ", err.Line)
		for _, line := range err.ContextSource() {
			if line.IsError {
				message += line.Source
			}
		}
	*/

	return
}

// formatResponse gets *revel.TestSuite as input parameter and
// transform response related info into a readable format.
func formatResponse(t testing.TestSuite) map[string]string {
	if t.Response == nil {
		return map[string]string{}
	}

	// Since Go 1.6 http.Request struct contains `Cancel <-chan struct{}` which
	// results in `json: unsupported type: <-chan struct {}`
	// So pull out required things for Request and Response
	req := map[string]interface{}{
		"Method":        t.Response.Request.Method,
		"URL":           t.Response.Request.URL,
		"Proto":         t.Response.Request.Proto,
		"ContentLength": t.Response.Request.ContentLength,
		"Header":        t.Response.Request.Header,
		"Form":          t.Response.Request.Form,
		"PostForm":      t.Response.Request.PostForm,
	}

	resp := map[string]interface{}{
		"Status":           t.Response.Status,
		"StatusCode":       t.Response.StatusCode,
		"Proto":            t.Response.Proto,
		"Header":           t.Response.Header,
		"ContentLength":    t.Response.ContentLength,
		"TransferEncoding": t.Response.TransferEncoding,
	}

	// Beautify the response JSON to make it human readable.
	respBytes, err := json.MarshalIndent(
		map[string]interface{}{
			"Response": resp,
			"Request":  req,
		},
		"",
		"   ")
	if err != nil {
		fmt.Println(err)
	}

	// Remove extra new line symbols so they do not take too much space on a result page.
	// Allow no more than 1 line break at a time.
	body := strings.Replace(string(t.ResponseBody), "\n\n", "\n", -1)
	body = strings.Replace(body, "\r\n\r\n", "\r\n", -1)

	return map[string]string{
		"Headers": string(respBytes),
		"Body":    strings.TrimSpace(body),
	}
}

//sortbySuiteName sorts the testsuites by name.
type sortBySuiteName []interface{}

func (a sortBySuiteName) Len() int      { return len(a) }
func (a sortBySuiteName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortBySuiteName) Less(i, j int) bool {
	return reflect.TypeOf(a[i]).Elem().Name() < reflect.TypeOf(a[j]).Elem().Name()
}

func init() {
	// Every time app is restarted convert the list of available test suites
	// provided by the revel testing package into a format which will be used by
	// the testrunner module and revel test cmd.
	revel.OnAppStart(func() {
		// Extracting info about available test suites from revel/testing package.
		registeredTests = map[string]int{}
		sort.Sort(sortBySuiteName(testing.TestSuites))
		for _, testSuite := range testing.TestSuites {
			testSuites = append(testSuites, describeSuite(testSuite))
		}
	})
}
