package main

/*
Parties invloved

1. Insured
2. Insurer
3. Medical Provider/Service Provider
4. Employer

*/

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Claim is the data structure to hold claim records in the ledger
type Claim struct {
	claimId string `json:"claimId"`
	customerId string `json:"customerId"`
	firstName string `json:"firstName"`
	lastName string `json:"lastName"`
	address string `json:"address"`
	policyNumber string `json:"policyNumber"`
}	

/*

Claim implements chaincode interface which contains three
functions Chaincode.invoke, Chaincode.init and Chaincode.query
init will be called at the initialization of chain code
stub shim.ChaincodeStubInterface is used to put transactions into ledger 
args contains genesis record
error contains any errors generated
*/
func (t *Claim) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	fmt.Println("Init function start")
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Check if the table already exists, _, called blank identifier helps ignoring other returned values except for the one we are concenrned about
	_, err := stub.GetTable("Claim")
	
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}	
	
	// Create Claim table
	err = stub.CreateTable("Claim", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "claimId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "customerId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "firstName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "lastName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "policyNumber", Type: shim.ColumnDefinition_STRING, Key: false},
	})	
	
	if err != nil {
		return nil, errors.New("Failed creating ApplicationTable.")
	}

//	err := stub.PutState("user_type1_1", []byte("Insured"))
//	err := stub.PutState("user_type1_2", []byte("Insurer"))
//	err := stub.PutState("user_type1_3", []byte("Provider"))
//	err := stub.PutState("user_type1_4", []byte("Employer"))
	
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Init function end")
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *Claim) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *Claim) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {											//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error
	return nil, errors.New("Received unknown function query: " + function)
}

func (t *Claim) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	
	fmt.Println("runnin write()")
	key = args[0]
	value = args[1]

	err = stub.PutState(key, []byte(value))
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number od arguments.");
	}
	
	if err != nil {
        return nil, err
    }
    return nil, nil
}

func (t *Claim) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }
	valAsbytes, err := stub.GetState(key)
	
	if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil

}

func main() {
	err := shim.Start(new(Claim))
	if err != nil {
		fmt.Printf("Error starting Claim: %s", err)
	}
} 

