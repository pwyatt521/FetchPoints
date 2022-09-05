package main 

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	//"fmt"
)

// each transaction has a payer, points, and timestamp associated with it
type transaction struct{
	Payer string `json:"payer"`
	Points int  `json:"points"`
	Timestamp string `json:"timestamp"`
}

// matches an amount of points to one payer
type payerPoints struct {
	Payer string `json:"payer"`
	Points int `json:"points"`
}

// the number of points to spend in a spend request
type spendRequest struct {
	Points int `json:"points"`
}

// stores the balance for each payer
var balance = map[string]payerPoints {}

// stores a list of all transactions sorted by timeline
var transactions = []transaction {}

// stores the total number of points the user currently has
var totalPoints = 0

// used to parse strings into a timestamp
const timeString = "2006-01-02T15:04:05Z07:00"

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

	if isTransactionInvalid(newTransaction){
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "Transactions cannot have a negative point value."})

	} else if addTransaction(newTransaction){
		c.IndentedJSON(http.StatusCreated, newTransaction)
	} else {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "Transaction duplicates not allowed."})
	}
}

// Check if the new transaction is valid
// The transaction cannot have a negative points value
// The transaction must have a correctly formatted timestamp
// RETURN: True if the transaction is invalid, and false if the transaction is valid
func isTransactionInvalid(newTransaction transaction) bool {
	if newTransaction.Points <= 0 {
		return true
	}
	_,error := time.Parse(timeString,newTransaction.Timestamp)

	if error != nil {
		return true
	}
	return false
}

// Given a transaction, check if it already exists in the transaction list, and if it does not, add it to the transaction list
// RETURNS: true if added, otherwise false
func addTransaction(newTransaction transaction) bool {
	if len(transactions) == 0 {
		transactions = append(transactions, newTransaction)
		updateTotals(newTransaction.Payer, newTransaction.Points)
		return true
	}
	parsedTimestamp,error := time.Parse(timeString, newTransaction.Timestamp)

	if error != nil {
		//fmt.Println(error)
		return false;
	}

	// binary search to see if duplicate, and if not duplicate insert into correct position
	low := 0
	high := len(transactions) -1
	mid := 0
	for low <= high {
		mid = (high+low)/2
		temp,_ := time.Parse(timeString, transactions[mid].Timestamp)
		if parsedTimestamp.After(temp){
			low = mid+1
		} else if parsedTimestamp.Before(temp) {
			high = mid-1
		} else {
			if transactions[mid].Payer == newTransaction.Payer {
				if transactions[mid].Points == newTransaction.Points {
					return false
				}
			}
			low=mid
			high=-1
		}
	}
	if low > len(transactions)-1 {
		transactions= append(transactions, newTransaction)
	} else {
		transactions = append(transactions[:low+1], transactions[mid:]...)
		transactions[mid] = newTransaction
	}
	
	updateTotals(newTransaction.Payer, newTransaction.Points)

	return true
}

// Given a payer string and a points integer, update the total points and the balance for the payer
func updateTotals(payer string, points int) {
	totalPoints += points
	if val, ok := balance[payer]; ok {
		balance[payer] = payerPoints{ payer, val.Points + points}
	} else {
		balance[payer] = payerPoints{payer,points}
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

// Try to spend points given an integer number of points to spend
// RETURNS: if points is greater than the total points available to spend or points is less than or equal to 0, return nil
//  otherwise, return a list of payerPoints that shows how the points will be spent
func trySpendPoints(points int) map[string]payerPoints {
	if (points > totalPoints)||(points <= 0){
		return nil
	}

	pointsSpent := map[string]payerPoints {}
	temp := 0
	tempPoints := points

	for i := 0; (i < len(transactions)) && (tempPoints > 0); i++ {
		temp = 0
		if _, ok := pointsSpent[transactions[i].Payer]; ok {
			temp = pointsSpent[transactions[i].Payer].Points
		}
		if tempPoints > transactions[i].Points {
			pointsSpent[transactions[i].Payer] = payerPoints { transactions[i].Payer, temp + transactions[i].Points }
			tempPoints -= transactions[i].Points
			updateTotals(transactions[i].Payer, 0-transactions[i].Points)
			transactions[i].Points = 0
		} else {
			pointsSpent[transactions[i].Payer] = payerPoints { transactions[i].Payer, temp + tempPoints }
			transactions[i].Points -= tempPoints
			updateTotals(transactions[i].Payer, 0-tempPoints)
			tempPoints = 0
		}
	}
	return pointsSpent
}