package main

import (
	"testing"
	"strconv"
	//"fmt"
)

// Test to make sure that invalid transactions are not allowed, and valid ones are allowed
func TestIsTransactionInvalid(t *testing.T) {
	testTransaction1 := transaction {"Payer1", 300,  "2020-11-02T14:00:00Z"}
 	testTransaction2 := transaction {"Payer2", -300,  "2020-11-02T14:00:00Z"}
 	testTransaction3 := transaction {"Payer2", 300,  "2020 Nov 02 14:00:00Z"}

	if isTransactionInvalid(testTransaction1) {
		t.Fatalf("testTransaction1 is valid, but isTransactionInvalid returns true.")
	}
	if !isTransactionInvalid(testTransaction2) {
		t.Fatalf("testTransaction2 is invalid, but isTransactionInvalid returns false.")
	}
	if !isTransactionInvalid(testTransaction3) {
		t.Fatalf("testTransaction3 is invalid, but isTransactionInvalid returns false.")
	}
}

// Test addTransaction and updateTotals to make sure everything is sorted correctly
// Then spend some points and make sure the correct transactions are subtracted from
func TestAddTransaction(t *testing.T) {
	testTransaction1 := transaction {"Payer1", 300,  "2020-11-02T14:00:00Z"}
	testTransaction3 := transaction {"Payer3", 200,  "2020-12-03T14:00:00Z"}
	testTransaction4 := transaction {"Payer4", 600,  "2021-11-02T14:00:00Z"}
	testTransaction5 := transaction {"Payer1", 100,  "2019-11-02T14:00:00Z"}
	if !addTransaction(testTransaction3) {
		t.Fatalf("testTransaction1 should have been added.")
	}
	if !addTransaction(testTransaction1) {
		t.Fatalf("testTransaction3 should have been added.")
	}
	if addTransaction(testTransaction1) {
		t.Fatalf("testTransaction1 should not have been added.")
	}
	if totalPoints != 500 {
		t.Fatalf("totalPoints should have the value of 500, but instead has a value of " + strconv.Itoa(totalPoints))
	}
	if len(transactions) != 2 {
		t.Fatalf("transactations should have a length of 2")
	}
	if transactions[0].Payer != "Payer1" {
		t.Fatalf("transactions should be sorted")
	}

	addTransaction(testTransaction4)
	addTransaction(testTransaction5)

	if balance["Payer1"].Points != 400 {
		t.Fatalf("Payer1 should have 400 points, but instead has a value of " + strconv.Itoa(balance["Payer1"].Points))
	}

	if totalPoints!= 1200 {
		t.Fatalf("totalPoints is expected to be 1200, but instead is "+ strconv.Itoa(totalPoints))
	}

	if trySpendPoints(1300) != nil {
		t.Fatalf("spending more points than the totalPoints should return nil")
	}

	if trySpendPoints(0) != nil {
		t.Fatalf("spending 0 or less points is not allowed, and should return nil")
	}

	spentPoints := trySpendPoints(500)

	if spentPoints["Payer3"].Points != 100{
		t.Fatalf("Payer3 points should have been 100, but was actually "+ strconv.Itoa(spentPoints["Payer3"].Points))
	}

	if totalPoints != 700 {
		t.Fatalf("totalPoints should have been 700 after spending 500, but was actually "+ strconv.Itoa(totalPoints))
	}

	if spentPoints["Payer4"].Points != 0 {
		t.Fatalf("Payer4 points should have been 0, but was actually "+ strconv.Itoa(spentPoints["Payer4"].Points))
	}

	if spentPoints["Payer1"].Points != 400 {
		t.Fatalf("Payer1 points should have been 400, but was actually "+ strconv.Itoa(spentPoints["Payer1"].Points))
	}
}