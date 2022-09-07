# FetchPoints
Fetch Rewards Coding Exercise - Backend Software Engineering

I wrote my project in go.  This is my first time writing in go, so I had to set it up on my compueter.
To download and setup go, use https://go.dev/doc/install

To run the project, in the project folder in a terminal window, run the following:
go run .

In a separate terminal, you can send requests to localhost port 8080.

To get the point balance for each payer, you can use the following GET request:
curl http://localhost:8080/balance

To post a new transaction, you can use the following POST request format:
curl http://localhost:8080/transaction --inclue --header "Content-Type: application/json" --request "POST" --data '{"payer":"Payer3", "points":100,  "timestamp":"2019-01-02T14:00:00Z"}'

To spend points, you can use the following PUT request format:
curl http://localhost:8080/spendPoints --include --header "Content-Type: application/json" --request "PUT" --data '{"points":500}' e --header "Content-Type: application
HTTP/1.1 200 OKt "PUT" --data '{"points":500}'

Future Enhancements:
- Improve security by not trusing all proxies (https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies)
- Add a user property, so different users can make transactions 
- Add session tokens