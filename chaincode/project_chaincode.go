/*
Chaincode created for Oracle hackathon

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var projectsIndexStr = "GE::ABCConsulting"

//Data elements

// user entered time entry

type TimeEntry struct {
	ProjectName     string `json:"projectname"`
	TaskName        string `json:"taskname"`
	User            string `json:"user"`
	QuantityInHours string `json:"quantityhours"`
	TotalAmount     string `json:"totalamount"`
}

// stored as P1::User1
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
	consultingProjects := []string{"Proj1", "Proj2", "Proj3"}

	jsonAsBytes, _ := json.Marshal(consultingProjects)
	err := stub.PutState(projectsIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	//Initilizing the Project user rates p1 -->[ {u1 100}, {u2 200 }]
	var projectUserRates []UserRate
	userrate := UserRate{}
	userrate.User = "Chandra"
	userrate.Rate = "90"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Sudheer"
	userrate.Rate = "100"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Sanjay"
	userrate.Rate = "80"
	projectUserRates = append(projectUserRates, userrate)

	jsonAsBytes, _ = json.Marshal(projectUserRates)
	err = stub.PutState("Proj1", jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
