package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type TestData struct {
	Description     string                 `json:"description"`
	GraphQLEndpoint string                 `json:"graphql_endpoint"`
	PlaygroundURL   string                 `json:"playground_url"`
	Mutations       map[string]TestCase    `json:"mutations"`
	Queries         map[string]TestCase    `json:"queries"`
	TestFlow        TestFlow               `json:"test_flow"`
}

type TestCase struct {
	Description      string      `json:"description"`
	Query            string      `json:"query"`
	ExpectedResponse interface{} `json:"expected_response"`
}

type TestFlow struct {
	Description string     `json:"description"`
	Steps       []TestStep `json:"steps"`
}

type TestStep struct {
	Step     int    `json:"step"`
	Action   string `json:"action"`
	Mutation string `json:"mutation,omitempty"`
	Query    string `json:"query,omitempty"`
	Note     string `json:"note"`
}

type GraphQLRequest struct {
	Query string `json:"query"`
}

type TestResponse struct {
	TestName    string      `json:"test_name"`
	Description string      `json:"description"`
	Query       string      `json:"query"`
	Response    interface{} `json:"response"`
	Status      string      `json:"status"`
	Error       string      `json:"error,omitempty"`
	Timestamp   string      `json:"timestamp"`
}

type TestResults struct {
	Summary struct {
		TotalTests   int    `json:"total_tests"`
		PassedTests  int    `json:"passed_tests"`
		FailedTests  int    `json:"failed_tests"`
		TestRunTime  string `json:"test_run_time"`
	} `json:"summary"`
	Results []TestResponse `json:"results"`
}

func main() {
	fmt.Println("🧪 Starting GraphQL API Test Runner...")
	
	// Load test data
	testData, err := loadTestData("test_data.json")
	if err != nil {
		fmt.Printf("❌ Error loading test data: %v\n", err)
		return
	}

	fmt.Printf("📊 Loaded test suite: %s\n", testData.Description)
	fmt.Printf("🎯 Target endpoint: %s\n", testData.GraphQLEndpoint)

	var results TestResults
	results.Results = []TestResponse{}
	startTime := time.Now()

	// Test mutations first (they create data)
	fmt.Println("\n🔧 Testing Mutations...")
	for name, testCase := range testData.Mutations {
		result := runTest(name, testCase, testData.GraphQLEndpoint)
		results.Results = append(results.Results, result)
		if result.Status == "PASSED" {
			results.Summary.PassedTests++
		} else {
			results.Summary.FailedTests++
		}
		results.Summary.TotalTests++
		
		// Small delay between tests
		time.Sleep(500 * time.Millisecond)
	}

	// Test queries
	fmt.Println("\n🔍 Testing Queries...")
	for name, testCase := range testData.Queries {
		result := runTest(name, testCase, testData.GraphQLEndpoint)
		results.Results = append(results.Results, result)
		if result.Status == "PASSED" {
			results.Summary.PassedTests++
		} else {
			results.Summary.FailedTests++
		}
		results.Summary.TotalTests++
		
		// Small delay between tests
		time.Sleep(500 * time.Millisecond)
	}

	results.Summary.TestRunTime = time.Since(startTime).String()

	// Save results
	err = saveTestResults(results, "test_response.json")
	if err != nil {
		fmt.Printf("❌ Error saving test results: %v\n", err)
		return
	}

	// Print summary
	fmt.Println("\n📋 Test Summary:")
	fmt.Printf("   Total Tests: %d\n", results.Summary.TotalTests)
	fmt.Printf("   ✅ Passed: %d\n", results.Summary.PassedTests)
	fmt.Printf("   ❌ Failed: %d\n", results.Summary.FailedTests)
	fmt.Printf("   ⏱️  Runtime: %s\n", results.Summary.TestRunTime)
	fmt.Printf("   📄 Results saved to: test_response.json\n")

	if results.Summary.FailedTests > 0 {
		fmt.Println("\n❌ Some tests failed. Check test_response.json for details.")
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All tests passed!")
	}
}

func loadTestData(filename string) (*TestData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var testData TestData
	err = json.Unmarshal(data, &testData)
	if err != nil {
		return nil, err
	}

	return &testData, nil
}

func runTest(testName string, testCase TestCase, endpoint string) TestResponse {
	fmt.Printf("  🧪 Running: %s\n", testName)
	
	result := TestResponse{
		TestName:    testName,
		Description: testCase.Description,
		Query:       testCase.Query,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	// Prepare GraphQL request
	reqBody := GraphQLRequest{
		Query: testCase.Query,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("Failed to marshal request: %v", err)
		return result
	}

	// Make HTTP request
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("HTTP request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("Failed to read response: %v", err)
		return result
	}

	// Parse JSON response
	var jsonResponse interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("Failed to parse JSON response: %v", err)
		result.Response = string(body)
		return result
	}

	result.Response = jsonResponse

	// Check if response contains errors
	if respMap, ok := jsonResponse.(map[string]interface{}); ok {
		if errors, hasErrors := respMap["errors"]; hasErrors {
			result.Status = "FAILED"
			result.Error = fmt.Sprintf("GraphQL errors: %v", errors)
			return result
		}
		
		// Check if response has data
		if _, hasData := respMap["data"]; hasData {
			result.Status = "PASSED"
		} else {
			result.Status = "FAILED"
			result.Error = "Response missing data field"
		}
	} else {
		result.Status = "FAILED"
		result.Error = "Invalid response format"
	}

	return result
}

func saveTestResults(results TestResults, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}