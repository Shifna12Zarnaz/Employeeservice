// Package controller contains ...
package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	//	"github.com/SigNoz/sample-golang-app/controllers"
	"github.com/connect2naga/logger/logging"
	"github.com/gorilla/mux"
)

/*
Author : Nagarjuna S
Date : 30-04-2022 18:18
Project : sample-http-service
File : endpoint-controller.go
*/

type EmployeeDetails struct {
	Id        string
	Name      string
	Locations string
}

type EndpointHandler struct {
	logger          logging.Logger
	EmployeeDetails map[string]EmployeeDetails
	// EmployeeDetails map[EmployeeDetails]string
}

func NewEndpointHandler(logger logging.Logger) *EndpointHandler {
	return &EndpointHandler{logger: logger, EmployeeDetails: make(map[string]EmployeeDetails)}
	// return &EndpointHandler{logger: logger, EmployeeDetails: make(map[EmployeeDetails]string)}
}
func (e *EndpointHandler) Status(w http.ResponseWriter, r *http.Request) {
	e.logger.Infof(context.Background(), "endpoint hit......")
	w.WriteHeader(http.StatusOK)
}

func (e *EndpointHandler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	e.logger.Infof(context.Background(), "GetAllEmployees hit......")
	data, err := json.Marshal(e.EmployeeDetails)
	if err != nil {
		fmt.Printf("failed to marshl...")
		w.Write([]byte("error while fetching data"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
func (e *EndpointHandler) PutEmployees(w http.ResponseWriter, r *http.Request) {
	var emp = []EmployeeDetails{}
	p := EmployeeDetails{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decodng response ", http.StatusBadRequest)
		return
	}
	emp = append(emp, p)
	fmt.Println(p)
	fmt.Println((emp))
	//	sample := e.EmployeeDetails
	//fmt.Println(sample)
	//  sample=map["1"]p

	e.EmployeeDetails = map[string]EmployeeDetails{"1": p}
	//	sample := map[EmployeeDetails]int{p: 1}
	fmt.Println("Map is ", e.EmployeeDetails)
	response, err := json.Marshal(&p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response ", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	//fmt.Fprintf(w, "Success")

}
func (e *EndpointHandler) PutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range e.EmployeeDetails {
		//	fmt.Println(index, item)
		//	fmt.Println(item.Id)
		fmt.Println(index)

		if item.Id == params["id"] {
			delete(e.EmployeeDetails, index)
			var emp EmployeeDetails
			_ = json.NewDecoder(r.Body).Decode(&emp)

			//rand.Seed(time.Now().UnixNano())
			//key := strconv.Itoa(rand.Intn(1000000000))
			key := index

			e.EmployeeDetails[key] = emp
			w.Write([]byte("data updated successfully"))
		}

	}

}

/*func (e *EndpointHandler) PutEmployees(w http.ResponseWriter, r *http.Request) {
	e.logger.Infof(context.Background(), "PushAllEmployees hit......")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	var emp EmployeeDetails
	er := json.Unmarshal(bodyBytes, &emp)
	if er != nil {
		//  if error is not nil
		//  print error
		fmt.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&emp)
	fmt.Println("EmployeeDetails is:", emp)
	json.NewEncoder(w).Encode(emp)
	w.WriteHeader(http.StatusOK)
}*/

func (e *EndpointHandler) GetAllEmployeeById(w http.ResponseWriter, r *http.Request) {
	e.logger.Infof(context.Background(), "GetAllEmployeeById hit......")

	vars := mux.Vars(r)
	empId := vars["id"]

	empDetails, ok := e.EmployeeDetails[empId]
	if !ok {
		fmt.Printf("no data availale...")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("given EmpID %s not found", empId)))
		return
	}

	data, err := json.Marshal(empDetails)
	if err != nil {
		fmt.Printf("failed to marshl...")
		w.Write([]byte("error while marshling data"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
func (e *EndpointHandler) DeleteEmp(w http.ResponseWriter, r *http.Request) {
	for index, item := range e.EmployeeDetails {
		if item.Id != "" {
			delete(e.EmployeeDetails, index)
		}
	}
}
