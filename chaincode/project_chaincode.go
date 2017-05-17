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
	TotalAmount     string `json:"totalamount"`
}

type AllProjectTimeEntry struct {
	ProjectTimeEntry []TimeEntry `json:"project_timeentry"`
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

	key = args[0]
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

	str := `{"chandra" : "90", "shiva" : "80", "sudheer" : "70"}`;
	
	//Initilizing the Project user rates proj1 -->[ {u1 100}, {u2 200 }]
	/*var projectUserRates []UserRate
	userrate := UserRate{}
	userrate.User = "chandra"
	userrate.Rate = "90"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "shiva"
	userrate.Rate = "80"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "sudheer"
	userrate.Rate = "70"
	projectUserRates = append(projectUserRates, userrate)

	jsonAsBytes, _ = json.Marshal(projectUserRates)*/
	err = stub.PutState(proj1Rates, []byte(str))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) create_time_entry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0          1          2         3     4
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
	
	var projRates string
	if projectname == "proj1" {
		projRates = proj1Rates
	} else if projectname == "proj2" {
		projRates = proj2Rates
	}

	//get the project time entries
	projectRatesAsBytes, err := stub.GetState(projRates)
	if err != nil {
		return nil, errors.New("Failed to get rates for "+ projRates)
	}
	
	var projRatesObj map[string]interface{}
	json.Unmarshal(projectRatesAsBytes, &projRatesObj)
	
	userRate := projRatesObj[personname].(string)
	userRateInt, err := strconv.Atoi(userRate)
	if err != nil {
		return nil, errors.New("Rate is not an integer")
	}
	
	totalAmount := userRateInt * quantity

	//check if marble already exists
	/*marbleAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get marble name")
	}
	res := TimeEntry{}
	json.Unmarshal(timeEntryAsBytes, &res)
	if res.Name == name{
		fmt.Println("This marble arleady exists: " + name)
		fmt.Println(res);
		return nil, errors.New("This time entry already exists")				//all stop a time entry by this name exists
	}*/
	
	//build the marble json string manually
	timeEntryAsBytes := []byte(`{"projectname": "` + projectname + `", "taskname": "` + taskname + `", "personname": "` + personname + `", "quantityinhours": ` + strconv.Itoa(quantity) + `", "expensetype": "` + expensetype + `", "totalamount": ` + strconv.Itoa(totalAmount) + `"}`)

	var timeEntry TimeEntry
	json.Unmarshal(timeEntryAsBytes, &timeEntry);

	//get the project time entries
	projectTimeEntriesAsBytes, err := stub.GetState(projectname)
	if err != nil {
		return nil, errors.New("Failed to get time entries for "+ projectname)
	}
	var projectTimeEntriesArray []TimeEntry
	json.Unmarshal(projectTimeEntriesAsBytes, &projectTimeEntriesArray)								//un stringify it aka JSON.parse()

	//append
	projectTimeEntriesArray = append(projectTimeEntriesArray, timeEntry)									//add marble name to index list
	fmt.Println("! Project Time Entries: ", projectTimeEntriesArray)
	jsonAsBytes, _ := json.Marshal(projectTimeEntriesArray)
	err = stub.PutState(projectname, jsonAsBytes)						//store name of marble

	fmt.Println("- end create time entry")
	return nil, nil
}