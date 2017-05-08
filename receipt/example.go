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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

type Project struct {
	Name   	string  `json:"name"`
	Reward 	string  `json:"reward"`
	Funds  	int     `json:"funds"`
	Target 	int     `json:"target"`
	Stat	bool	`json:"stat"`
	Creator string 	`json:"creator"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 10 {
		return nil, errors.New("Incorrect number of arguments. Execting 10")
	}

	var usersArray []string
	var projectsArray []string

	var userone User
	userone.Name = args[0]
	userone.Password = args[1]
	balance, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	}
	
	userone.Balance = balance

	b, err := json.Marshal(userone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userone")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	userone.Name = args[3]
	userone.Password = args[4]
	balance, err = strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	}

	userone.Balance = balance

	b, err = json.Marshal(userone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userone")
	}

	err = stub.PutState(args[3], b)
	if err != nil {
		return nil, err
	}

	usersArray = append(usersArray, args[0])
	usersArray = append(usersArray, args[3])

	b, err = json.Marshal(usersArray)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for usertwo")
	}

	err = stub.PutState("users", b)
	if err != nil {
		return nil, err
	}
	
	var projectone Project
	
	projectone.Name = args[6]
	projectone.Reward = args[7]
	funds, err := strconv.Atoi(args[8])
	if err != nil {
		return nil, errors.New("Expecting integer value for the projectFunds at place 9")
	}
	target, err := strconv.Atoi(args[9])
	if err != nil {
		return nil, errors.New("Expecting integer value for the projectTarget at place 10")
	}
	projectone.Stat = false
	projectone.Creator = (args[0])
	
	projectone.Funds = funds
	projectone.Target = target
	
	b, err = json.Marshal(projectone)
	if err != nil{
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for projectone")
	}
	
	err = stub.PutState(args[6], b)
	if err != nil {
		return nil, err
	}
	
	projectsArray = append(projectsArray, args[6])
	
	b, err = json.Marshal(projectsArray)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for projectsarray")
	}
	
	err = stub.PutState("projects", b)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) Transaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var X int // Transaction value
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	var userA User
	err = json.Unmarshal(Avalbytes, &userA)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of userA")
	}

	Bvalbytes, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	var userB User
	err = json.Unmarshal(Bvalbytes, &userB)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of userB")
	}

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Third argument must be integer")
	}

	userA.Balance = userA.Balance - X
	userB.Balance = userB.Balance + X
	fmt.Printf("Aval = %d, Bval = %d\n", userA.Balance, userB.Balance)
	// valdidation
	if userA.Balance < 0 {
		userA.Balance = userA.Balance + X
		userB.Balance = userB.Balance - X
		return nil, errors.New("unsufficient balance please fund your account")
	}
	
	b, err := json.Marshal(userA)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for usera")
	}

	// Write the state back to the ledger
	err = stub.PutState(userA.Name, b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(userB)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userb")
	}

	err = stub.PutState(userB.Name, b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (t *SimpleChaincode) InvestProject(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. Name of the project, name of the investor and the amount")
	}
	
	var X int // investment value
	var err error
	
	// get the state from the ledger
	
	projectState, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	var projectX Project
	
	err = json.Unmarshal(projectState, &projectX)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of projectX")
	}
	
	userState, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	var userX User
	
	err = json.Unmarshal(userState, &userX)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of userX")
	}
	
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Third argument must be a integer")
	}
	
	userX.Balance = userX.Balance - X
	projectX.Funds = projectX.Funds + X
	fmt.Printf("Funds for project %d is %d and the %d balance is %d", projectX.Name, projectX.Funds, userX.Name, userX.Balance)
	
	if userX.Balance < 0 {
		userX.Balance = userX.Balance + X
		projectX.Funds = projectX.Funds - X
		return nil, errors.New("unsufficient balance please fund your account")
	}
	
	// self execution
	if projectX.Funds >= projectX.Target {
		creatorState, err := stub.GetState(projectX.Creator)
		if err != nil {
			return nil, errors.New("Failed to get creatorstate")
		}
		var creatorX User

		err = json.Unmarshal(creatorState, &creatorX)
		if err != nil {
			return nil, errors.New("Failed to marshal string to struct of creator")
		}
		// transfer all the funds to the creator from the project
		X = projectX.Funds
		projectX.Funds = projectX.Funds - X
		creatorX.Balance = creatorX.Balance + X
		projectX.Stat = true
		fmt.Println("The project has been succesfully funded the funds have been transferred to the Creator of the project")
		// write everything back to the ledger
		b, err := json.Marshal(creatorX)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for creatorX")
		}

		err = stub.PutState(creatorX.Name, b)
		if err != nil {
			return nil, err
		}

		b, err = json.Marshal(projectX)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for projectX")
		}

		err = stub.PutState(projectX.Name, b)
		if err != nil {
			return nil, err
		}
	}
	b, err := json.Marshal(projectX)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for projectX")
	}
	
	err = stub.PutState(projectX.Name, b)
	if err != nil {
		return nil, err
	}
	
	b, err = json.Marshal(userX)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userX")
	}
	
	err = stub.PutState(userX.Name, b)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

func (t *SimpleChaincode) RechargeBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. Name of the investor, the amount and the password")
	} // password is the sha256 hash of mendix
	if args[2] == "1274d60ff458da72bf3e58107cc2ebcf1f542b587b94c358eb65265f85c72cf5"{
		var X int // charge amount
		var err error
		// get the state from the user from the ledger
		userState, err := stub.GetState(args[0])
		if err != nil {
			return nil, errors.New("Failed to get user state")
		}

		var userX User
		err = json.Unmarshal(userState, &userX)
		if err != nil {
			return nil, errors.New("Failed to marshal string from user struct")
		}
		// perform the execution
		X, err = strconv.Atoi(args[1])
		if err !=  nil {
			return nil, errors.New("second argument must be a integer")
		}

		userX.Balance = userX.Balance + X
		fmt.Printf("new user balance is %d", userX.Balance)

		b, err := json.Marshal(userX)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for userX")
		}
		// write back to the ledger
		err = stub.PutState(userX.Name, b)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("the password is incorrect: tips: 'mendix' 'bitcoin secure'")
	}
	return nil, nil
}

func (t *SimpleChaincode) CreateUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name,password,balance to create user")
	}

	usersArray, err := stub.GetState("users")
	if err != nil {
		return nil, err
	}

	var users []string

	err = json.Unmarshal(usersArray, &users)

	if err != nil {
		return nil, err
	}

	users = append(users, args[0])

	b, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for usertwo")
	}

	err = stub.PutState("users", b)
	if err != nil {
		return nil, err
	}

	var userone User
	userone.Name = args[0]
	userone.Password = args[1]
	balance, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	}

	userone.Balance = balance

	b, err = json.Marshal(userone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userone")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (t *SimpleChaincode) CreateProject(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5. name, reward, funds, target and creator to create project")
	}
	
	projectsArray, err := stub.GetState("projects")
	if err != nil {
		return nil, err
	}
	
	var projects []string
	
	err = json.Unmarshal (projectsArray, &projects)
	
	if err != nil {
		return nil, err
	}
	
	projects = append(projects, args[0])
	
	b, err := json.Marshal(projects)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for projects")
	}
	
	err = stub.PutState("projects", b)
	if err != nil {
		return nil, err
	}
	
	var projectone Project
	projectone.Name = args[0]
	projectone.Reward = args[1]
	funds, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for the projectFunds at place 3")
	}
	target, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for the projectTarget at place 4")
	}
	projectone.Stat = false
	projectone.Creator = (args[4])
	
	projectone.Funds = funds
	projectone.Target = target
	
	b, err = json.Marshal(projectone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for projectone")
	}
	
	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// Invoke is your entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "transaction" {
		return t.Transaction(stub, args)
	} else if function == "create_user" {
		return t.CreateUser(stub, args)
	} else if function == "create_project" {
		return t.CreateProject(stub, args)
	} else if function == "investment" {
		return t.InvestProject(stub, args)
	} else if function == "recharge" {
		return t.RechargeBalance(stub, args)
	}
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "list_users" {
		return t.listUsers(stub, args)
	} else if function == "list_projects" {
		return t.listProjects(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
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

func (t *SimpleChaincode) listUsers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("users")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for users}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
func (t *SimpleChaincode) listProjects(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
						
	valAsbytes, err := stub.GetState("projects")
	if err != nil {
	jsonResp = "{\"Error\":\"Failed to get state for projects}"
	return nil, errors.New(jsonResp)
	}
		
	return valAsbytes, nil
	}
