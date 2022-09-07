// This file holds the unit tests for the fetchPoints file

package main

import (
	"testing"
	"strconv"
)

// Test to make sure that invalid transactions are not allowed, and valid ones are allowed
func TestIsTransactionInvalid(t *testing.T) {
	testTransaction1 := transaction {"Payer1", 300,  "2020-11-02T14:00:00Z"} // Valid transaction
 	testTransaction2 := transaction {"Payer2", -300,  "2020-11-02T14:00:00Z"} // Invalid number of points
 	testTransaction3 := transaction {"Payer2", 300,  "2020 Nov 02 14:00:00Z"}  // Invalid date format

	if isTransactionInvalid(testTransaction1) != "" {
		t.Fatalf("testTransaction1 is valid, but isTransactionInvalid returns an error message.")
	}
	if isTransactionInvalid(testTransaction2) == "" {
		t.Fatalf("testTransaction2 is invalid, but isTransactionInvalid does not return an error.")
	}
	if isTransactionInvalid(testTransaction3) == "" {
		t.Fatalf("testTransaction3 is invalid, but isTransactionInvalid does not return an error.")
	}
}

// Test addTransaction and updateTotals to make sure everything is sorted correctly
// Then spend some points and make sure the correct transactions are subtracted from
// The tag is large to easily test a growing transaction list and correctly interact with it multiple times
func TestAddTransaction(t *testing.T) {
	testTransaction1 := transaction {"Payer1", 300,  "2020-11-02T14:00:00Z"}
	testTransaction2 := transaction {"Payer2", 100,  "2019-11-02T14:00:00Z"}
	testTransaction3 := transaction {"Payer3", 200,  "2020-12-03T14:00:00Z"}
	testTransaction4 := transaction {"Payer4", 600,  "2021-11-02T14:00:00Z"}
	testTransaction5 := transaction {"Payer1", 100,  "2019-11-02T14:00:00Z"}
	testTransaction6 := transaction {"Payer2", 200,  "2018-11-02T14:00:00Z"}
	testTransaction7 := transaction {"Payer2", 300,  "2019-10-02T14:00:00Z"}
	testTransaction8 := transaction {"Payer4", 1000,  "2021-09-02T14:00:00Z"}

	// Test transactions are appropriately added
	if !addTransaction(testTransaction3) {
		t.Fatalf("testTransaction1 should have been added.")
	}
	if !addTransaction(testTransaction1) {
		t.Fatalf("testTransaction3 should have been added.")
	}
	if addTransaction(testTransaction1) {
		t.Fatalf("testTransaction1 should not have been added.")
	}

	// Test totals and transactions lists are as expected
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

	// Test balance is correct
	if balance["Payer1"].Points != 400 {
		t.Fatalf("Payer1 should have 400 points, but instead has a value of " + strconv.Itoa(balance["Payer1"].Points))
	}

	if totalPoints!= 1200 {
		t.Fatalf("totalPoints is expected to be 1200, but instead is "+ strconv.Itoa(totalPoints))
	}

	// Test spending points
	if trySpendPoints(1300) != nil {
		t.Fatalf("spending more points than the totalPoints should return nil")
	}

	if trySpendPoints(0) != nil {
		t.Fatalf("spending 0 or less points is not allowed, and should return nil")
	}

	spentPoints := trySpendPoints(500)

	if spentPoints["Payer3"].Points != -100{
		t.Fatalf("Payer3 points should have been -100, but was actually "+ strconv.Itoa(spentPoints["Payer3"].Points))
	}

	if totalPoints != 700 {
		t.Fatalf("totalPoints should have been 700 after spending 500, but was actually "+ strconv.Itoa(totalPoints))
	}

	if spentPoints["Payer4"].Points != 0 {
		t.Fatalf("Payer4 points should have been 0, but was actually "+ strconv.Itoa(spentPoints["Payer4"].Points))
	}

	if spentPoints["Payer1"].Points != -400 {
		t.Fatalf("Payer1 points should have been -400, but was actually "+ strconv.Itoa(spentPoints["Payer1"].Points))
	}

	// Test corner case of binary search where two transactions have the same timestamp and points, but different payers
	if !addTransaction(testTransaction2) {
		t.Fatalf("Payer2 should be added even though it has the same timestamp and points as Payer 3")
	}

	addTransaction(testTransaction6)
	addTransaction(testTransaction7)
	addTransaction(testTransaction8)

	// Test that the points are not spent twice, but that partially spent transactions are still correctly spent
	spentPoints= nil
	spentPoints = trySpendPoints(1000)
	if spentPoints["Payer2"].Points!=-600 {
		t.Fatalf("Payer2 points should have been -600, but was actually "+ strconv.Itoa(spentPoints["Payer2"].Points))
	}
	if spentPoints["Payer3"].Points!=-100 {
		t.Fatalf("Payer3 points should have been -100, but was actually "+ strconv.Itoa(spentPoints["Payer3"].Points))
	}
	if spentPoints["Payer4"].Points!=-300 {
		t.Fatalf("Payer4 points should have been -300, but was actually "+ strconv.Itoa(spentPoints["Payer4"].Points))
	}
	if spentPoints["Payer1"].Points!=0 {
		t.Fatalf("Payer1 points should have been 0, but was actually "+ strconv.Itoa(spentPoints["Payer1"].Points))
	}
}