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
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
// Struct for storing name / password and balance as value
type User struct {
	Name			string `json:"name"`
	Password	string `json:"password"`
	Balance		int 	 `json:"balance"`

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
		return nil, errors.New("Incorrect number of arguments. Expecting 10")
	}
	// first 4 parameters are for the project
	err := stub.PutState("name", []byte(args[0]))
	err1 := stub.PutState("reward", []byte(args[1]))
	err2 := stub.PutState("funds", []byte(args[2]))
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
	// next 6 arguments are split for user 1 (name/passwd/balance) and user 2 *(name/passwd/balance)

	var usersArray []string   // creating a array of users
	// first user
	var userone User
	userone.Name = args[4]
	userone.Password = args [5]
	balance, err := strconv.Atoi(args[6])           // string convert from integer
	if err != nil{
		return nil, errors.New("Expecting integer value for your balance at the 7th place")
		}

		userone.Balance = balance

		b, err := json.Marshal(userone)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for userone")
		}

		err = stub.PutState(args[4], b)
		if err != nil {
			return nil, err
		}
		// second user
		userone.Name = args[7]
		userone.Password =args[8]
		balance, err = strconv.Atoi(args[9])  // string convert from integer
		if err != nil {
			return nil, errors.New("Expecting integer value for your balance at the 10th place")
		}

		userone.Balance = balance

		b, err = json.Marshal(userone)        // convert object to one suitable for storage
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for userone")
		}

		err = stub.PutState(args[7], b)
		if err != nil {
			return nil, err
		}

		usersArray = append(usersArray, args[4]) // put the first user at the end of the usersArray
		usersArray = append(usersArray, args[7]) // put the second user at the end of the usersArray

		b, err = json.Marshal(usersArray)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Errors while creating json string for usertwo")
		}

		err = stub.PutState("users", b)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Invoke isur entry point to invoke a chaincode function
	func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
		fmt.Println("invoke is running " + function)

		// Handle different functions
		if function == "init" {
			return t.Init(stub, "init", args)
			} else if function == "write" {
				return t.write(stub, args)
				} else if function == "transaction"{
					return t.Transaction(stub, args)
					}else if function == "createUser"{
						return t.CreateUser(stub, args)
					}
				fmt.Println("invoke did not find func: " + function)
				return nil, errors.New("Received unknown function invocation: " + function)
			}

			// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
				fmt.Println("query is running " + function)

				// Handle different functions
				if function == "read" { //read a variable
					return t.read(stub, args)
					} else if function == "listUsers"{
						return t.listUsers(stub, args)
					}
					fmt.Println("query did not find func: " + function)

					return nil, errors.New("Received unknown function query: " + function)
				}

				// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
					var projectName, projectReward, projectFunds, projectTarget string
					var err error
					fmt.Println("running write()")

					if len(args) != 4 {
						return nil, errors.New("Incorrect number of arguments. Expecting 4 name of the key and value to set")
					}

					projectName = args[0] //rename for funsies
					projectReward = args[1]
					projectFunds = args[2]
					projectTarget = args[3]
					err = stub.PutState("name", []byte(projectName)) //write the variable into the chaincode state
					err = stub.PutState("rewards", []byte(projectReward)) //write the variable into the chaincode state
					err = stub.PutState("funds", []byte(projectFunds))
					err = stub.PutState("target", []byte(projectTarget))
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
func (t *SimpleChaincode)CreateUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
					if len(args) != 3 {
						return nil, errors.New("Incorrect number of arguments. Expecting name, password, balance to create user")
					}
					usersArray, err := stub.GetState("users")    // retrieve the list of users from the ledger
					if err != nil {
						return nil, err
					}

					var users [] string

					err = json.Unmarshal(usersArray,&users)     // revert the data back to a object

					if err != nil {
						return nil, err
					}

					users = append(users, args[0])

					b, err :=json.Marshal(users)
					if err != nil {
						fmt.Println(err)
						return nil, errors.New("Errors while creating json string for usertwo ")
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
								return nil, errors.New("Expecting integer value for asset holding at place 3")
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
func (t *SimpleChaincode) Transaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

							var tx int // transaction value
							var err error

							if len(args) != 3 {
								return nil, errors.New("incorrect number of arguments. Expecting 3 NameOfSender/ NameOfReceiver / transactionAmount")
							}
							// retrieve sender information from the ledger
							Sendervalbytes, err := stub.GetState(args[0])
							if err != nil{
								return nil, errors.New("Failed to get senderState")
							}
							var Sender User
							err = json.Unmarshal(Sendervalbytes, &Sender)
							if err != nil {
								return nil, errors.New("Failed to marshal string to struct of Sender")
							}
							// retrieve reciever information from the ledger
							Receivervalbytes, err := stub.GetState(args[1])
							if err != nil{
								return nil, errors.New("Failed to get receiverState")
							}
							var Receiver User
							err = json.Unmarshal(Receivervalbytes, &Receiver)
							if err != nil {
								return nil, errors.New("Failed to marshal string to struct of userB")
							}
							// perform the transaction

							tx, err = strconv.Atoi(args[2])
								if err != nil {
									return nil, errors.New("Third argument needs to be a integer")
								}

							Sender.Balance = Sender.Balance - tx
							Receiver.Balance = Receiver.Balance + tx
							fmt.Printf("SenderValue = %d, ReceiverValue = %d\n, Sender.Balance, Receiver.Balance")

							b, err := json.Marshal(Sender)
							if err != nil{
								fmt.Println(err)
								return nil, errors.New("Errors while creating json string for Sender")
							}
							// Write the state back to the ledger
							err = stub.PutState(Sender.Name, b)
							if err != nil {
								return nil, err
							}

							b, err = json.Marshal(Receiver)
							if err != nil {
								fmt.Println(err)
								return nil, errors.New("Errors while creating json string for Receiver")
							}

							err = stub.PutState(Receiver.Name, b)
							if err != nil {
									return nil, err
							}
							return nil, nil
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
