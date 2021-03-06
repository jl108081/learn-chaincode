/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 4 {
        return nil, errors.New("Incorrect number of arguments. Expecting 4")
    }

	err := stub.PutState("name", []byte(args[0]))
    err1 := stub.PutState("rewards", []byte(args[1]))
	err2 := stub.PutState("duration", []byte(args[2]))
	err3 := stub.PutState("target", []byte(args[3]))
    if err != nil {
        return nil, err
    }
	if err1 != nil {
        return nil, err1
    }
	if err2 != nil {
        return nil, err2
    }
	if err3 != nil {
        return nil, err3
    }

    return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "write" {
        return t.write(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    
    var err error
    fmt.Println("running write()")

    if len(args) != 4 {
        return nil, errors.New("Incorrect number of arguments. Expecting 4")
    }
                
    var projectName = args[0]
	err = stub.PutState("name", []byte(projectName))  //write the variable into the chaincode state
	                
    var projectRewards = args[1]
	err = stub.PutState("rewards", []byte(projectRewards))  //write the variable into the chaincode state
	                
    var projectDuration = args[2]
	err = stub.PutState("duration", []byte(projectDuration))  //write the variable into the chaincode state
	                 
    var projectTarget = args[3]
	err = stub.PutState("target", []byte(projectTarget))  //write the variable into the chaincode state
	
	
    if err != nil {
        return nil, err
    }
    return nil, nil
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
