// This file contains the code to handle GET, POST, and PUT requests, and properly
// route the requests to the correct functions.

package main 

import (
	"github.com/gin-gonic/gin"
	"net/http"
)



// handle the routing of GET, POST, and PUT calls and listen on localhost port 8080
func main(){
	router := gin.Default()
	router.GET("/balance", getBalance)
	router.POST("/transaction", postTransaction)
	router.PUT("/spendPoints", spendPoints)

	router.Run("localhost:8080")
}

// A GET request for the users current balance
// RETURNS: the total number of points remaining for each Payer in the following format
// [ 
//    { "payer": "Payer1", "points": 1000}
// 	  { "payer": "Payer2", "points": 0}
//    { "payer": "Payer3", "points": 5300}
// ] 
func getBalance(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, balance)
}

// Post one transaction to the user's profile
// Expected transaction POST format
//  {"payer": "DANNON", "points": 1000, "timestamp": "2020-11-02T14:00:00Z"}
// RETURNS: A POST success message if properly formatted
//    If not properly formatted, return a failure message
func postTransaction(c *gin.Context) {
	var newTransaction transaction

	if err := c.BindJSON(&newTransaction); err !=nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "Incorrect transaction format."})
		return
	}
	errMessage :=isTransactionInvalid(newTransaction)
	if errMessage != "" {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": errMessage})

	} else if addTransaction(newTransaction){
		c.IndentedJSON(http.StatusCreated, newTransaction)
	} else {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "Transaction duplicates are not allowed."})
	}
}


// Try to spend a certain amount of points via a PUT request  
// Expected format for request
//  {"points":500}
// RETURNS: 
//   If the user tries to spend more points than they have or negative points, return failure message and spend no points
//   Return the following information for a valid spend PUT request
//   [
//		{ "payer": "DANNON", "points": -100 },
//		{ "payer": "UNILEVER", "points": -200 },
//		{ "payer": "MILLER COORS", "points": -4,700 }
//	 ]
func spendPoints(c *gin.Context) {
	var newSpendRequest spendRequest

	if err := c.BindJSON(&newSpendRequest); err !=nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "Incorrect spend request format."})
		return
	}
	pointsSpent :=trySpendPoints(newSpendRequest.Points)
	if pointsSpent != nil{
		c.IndentedJSON(http.StatusOK, pointsSpent)
	} else {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "You cannot spend more than you have or a negative amount of points."})
	}

}