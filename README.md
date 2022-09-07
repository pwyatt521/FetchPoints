# FetchPoints
Fetch Rewards Coding Exercise - Backend Software Engineering

I wrote my project in go.  This is my first time writing in go, so I had to set it up on my compueter.
To download and setup go, use https://go.dev/doc/install

To run the project, in the project folder in a terminal window, run the following:
go run .

To run the unit tests for the project, in the project folder in a terminal window, run the following:
go test

In a separate terminal, you can send requests to localhost port 8080. See below for more information on the available requests.

To get the point balance for each payer, you can use the following GET request:
curl http://localhost:8080/balance

To post a new transaction, you can use the following POST request format:
curl http://localhost:8080/transaction --inclue --header "Content-Type: application/json" --request "POST" --data '{"payer":"Payer3", "points":100,  "timestamp":"2019-01-02T14:00:00Z"}'

NOTE: I added some input validation for adding transactions.
- The points cannot be 0 or negative when adding a transaction. I made this decision for the following reasons:
	* While using fetch, I have never been punished points for adding a transaction and I could not think of a reason to do so, as it would punish users for uploading a transaction.  This would disincentivizing the user from using the app.
	* Allowing negative values in specific transactions but not the payer balances would require looping over all transactions each time a transaction is added for a payer with a 0 balance to correct for negative values.  
- The timestamp must follow the expected format of "2006-01-02T15:04:05Z07:00".  This is the format that best matches the examples given in the instructions, and allows for copy and pasting the json objects from the instructions to help create transactions. 
- You cannot enter the same timestamp twice.  At least one of the properties must be different.  If the payer, points, and timestamp are all the exact same, an error message will be returned to the user.

To spend points, you can use the following PUT request format:
curl http://localhost:8080/spendPoints --include --header "Content-Type: application/json" --request "PUT" --data '{"points":500}' e --header "Content-Type: application
HTTP/1.1 200 OKt "PUT" --data '{"points":500}'

- You cannot spend more points then you have.  If you have 500 total points, and attempt to spend 1000, an error will be given and no change will be made to the current point balance.
- You cannot spend 0 or a negative amount of points.

Future Enhancements:
- Improve security by not trusing all proxies (https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies)
- Add a user property, so different users can make transactions 
- Add session tokens
- Separate out unit tests so each tag only tests one case.  It is a little annoying to fix the only test that failed case just to have the next one fail.  I kept them together for simplicity interacting with the local memory storage and not having to add extra set up tags.
