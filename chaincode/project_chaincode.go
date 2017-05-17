/*
Chaincode created for Oracle hackathon

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var projectsIndexStr = "GE:::ABCConsulting"
var proj1Rates = "proj1:::rates"
var proj2Rates = "proj2:::rates"

//Data elements

// user entered time entry

type TimeEntry struct {
	ProjectName     string `json:"projectname"`
	TaskName        string `json:"taskname"`
	PersonName      string `json:"personname"`
	QuantityInHours string `json:"quantityinhours"`
	ExpenseType		string `json:"expensetype"`
	DerivedAmount   string `json:"derivedamount"`
}

type ProjectTimeEntry struct {
	ProjectName     string 		`json:"projectname"`
	TotalAmount		string		`json:"totalamount"`
	TimeEntries 	[]TimeEntry	`json:"timeentries"`
}

type ProjectMilestone struct {
	ProjectName   string `json:"projectname"`
	MilestoneName string `json:"milestonename"`
	User          string `json:"user"`
	Amount        string `json:"amount"`
}

// list of project milestones , as  example P1 --> M1 1000
type AllProjectMilestones struct {
	ProjectMileStones []ProjectMilestone `json:"project_milestones"`
}

type UserRate struct {
	User string `json:"user"`
	Rate string `json:"rate"`
}

// ============================================================================================================================
//  Main - main - Starts up the chaincode
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("test", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	_, err = t.initializeData(stub, args)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_time_entry" { //initialize the chaincode state, used as reset
		return t.create_time_entry(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" { //read a variable
		fmt.Println("hi there " + function) //error
		return nil, nil
	} else if function == "read" {
		return t.read(stub, args)
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = strings.ToLower(args[0])
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//Initilizing project Data

func (t *SimpleChaincode) initializeData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//Initilizing the sample projects
	consultingProjects := []string{"proj1", "proj2"}

	jsonAsBytes, _ := json.Marshal(consultingProjects)
	err := stub.PutState(projectsIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	//Rates for proj1
	str := `{"chandra" : "90", "shiva" : "80", "sudheer" : "70"}`;
	err = stub.PutState(proj1Rates, []byte(str))
	if err != nil {
		return nil, err
	}

	//Rates for proj2
	str = `{"chandra" : "120", "shiva" : "110", "sudheer" : "100"}`;
	err = stub.PutState(proj2Rates, []byte(str))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Create Time Entry - create a new time entry, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) create_time_entry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//   0          1         2        3    4
	// "Proj1", "Task11", "chandra", "40", "ST"
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	//input sanitation
	fmt.Println("- start create time entry")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	projectname := strings.ToLower(args[0])
	taskname := strings.ToLower(args[1])
	personname := args[2]
	quantity, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	expensetype := strings.ToLower(args[4])
	
	fmt.Println("Time Entry : ", projectname, taskname, personname, quantity, expensetype)
	var projRates string
	if projectname == "proj1" {
		projRates = proj1Rates
	} else if projectname == "proj2" {
		projRates = proj2Rates
	}
	fmt.Println("projRates : "+projRates);

	//get the project time entries
	projectRatesAsBytes, err := stub.GetState(projRates)
	if err != nil {
		return nil, errors.New("Failed to get rates for "+ projRates)
	}
	if projectRatesAsBytes == nil {
		return nil, errors.New("Rates are not defined for "+ projectname)
	}
	
	var projRatesObj map[string]interface{}
	json.Unmarshal(projectRatesAsBytes, &projRatesObj)
	fmt.Println("projRatesObj : ", projRatesObj);
	
	userRate := projRatesObj[personname].(string)
	fmt.Println("userRate : "+userRate);

	userRateInt, err := strconv.Atoi(userRate)
	if err != nil {
		return nil, errors.New("Rate is not an integer")
	}
	fmt.Println("userRateInt : ", userRateInt);
	
	derivedAmount := userRateInt * quantity
	fmt.Println("derivedAmount : ", derivedAmount);

	//build the time entry json string manually
	timeEntryAsBytes := []byte(`{"projectname": "` + projectname + `", "taskname": "` + taskname + `", "personname": "` + personname + `", "quantityinhours": "` + strconv.Itoa(quantity) + `", "expensetype": "` + expensetype + `", "derivedamount": "` + strconv.Itoa(derivedAmount) + `"}`)

	fmt.Println("timeEntryAsBytes : ", timeEntryAsBytes);
	fmt.Println("timeEntryAsBytes as a string : ", string(timeEntryAsBytes));
	var timeEntry TimeEntry
	json.Unmarshal(timeEntryAsBytes, &timeEntry);
	fmt.Println("timeEntry : ", timeEntry);

	//get the project time entries
	projectTimeEntryAsBytes, err := stub.GetState(projectname)
	if err != nil {
		return nil, errors.New("Failed to get time entries for "+ projectname)
	}
	var projectTimeEntry ProjectTimeEntry
	json.Unmarshal(projectTimeEntryAsBytes, &projectTimeEntry)								//un stringify it aka JSON.parse()
	fmt.Println("Project Time Entry (Before append) : ", projectTimeEntry);
	
	var totalAmountInt int
	if len(projectTimeEntry.TimeEntries) > 0 {
		fmt.Println("No. of time entries : ", len(projectTimeEntry.TimeEntries));
		totalAmountInt, err = strconv.Atoi(projectTimeEntry.TotalAmount)
		if err != nil {
			return nil, errors.New("Total Amount in Project Time Entry is not an integer")
		}
		fmt.Println("totalAmountInt (Before): ", totalAmountInt);
		totalAmountInt = totalAmountInt + derivedAmount
		fmt.Println("totalAmountInt (After): ", totalAmountInt);
	} else {
		projectTimeEntry.ProjectName = projectname
		totalAmountInt = derivedAmount
	}

	fmt.Println("totalAmountInt : ", totalAmountInt);
	fmt.Println("strconv.Itoa(totalAmountInt) : ", strconv.Itoa(totalAmountInt));
	projectTimeEntry.TotalAmount = strconv.Itoa(totalAmountInt)
	fmt.Println("projectTimeEntry.TotalAmount : ", projectTimeEntry.TotalAmount);
	timeEntries := projectTimeEntry.TimeEntries
	timeEntries = append(timeEntries, timeEntry)
	projectTimeEntry.TimeEntries = timeEntries
	fmt.Println("! Project Time Entry (After append) : ", projectTimeEntry)
	jsonAsBytes, _ := json.Marshal(projectTimeEntry)
	err = stub.PutState(projectname, jsonAsBytes)

	fmt.Println("- end create time entry")
	return nil, nil
}