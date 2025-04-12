package main

import (
    "fmt"
    "internal/app/model/admin/dawonlod"
)

func main() {
    data := [][]string{
        {"Name", "Age", "City"},
        {"Alice", "30", "New York"},
        {"Bob", "25", "Los Angeles"},
    }

    err := dawonlod.WriteToCSV("output.csv", data)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("CSV file created successfully!")
    }
}










type Role_CDE struct {
	ID   int    `json:"id" db:"id"`
	Superviser_Name string `json:"name" db:"name"`
	Phone_number string `json:"phone_number" db:"phone_number"`
}

type Role_CATI struct {
	ID   int    `json:"id" db:"id"`
	Superviser_Name string `json:"name" db:"name"`
	Phone_number string `json:"phone_number" db:"phone_number"`
}

type Role_incotaion struct {
	ID   int    `json:"id" db:"id"`
	Superviser_Name string `json:"name" db:"name"`
	Phone_number string `json:"phone_number" db:"phone_number"`
}

type Program struct {
	ID   int    `json:"id" db:"id"`
	Status string `json:"status" db:"status"`
}