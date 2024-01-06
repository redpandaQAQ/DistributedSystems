# Hand-in4
To run the program, you need to run the following commands in 3 different terminals:

- go run node/node.go -port 5455 -id 2 -nextport 5456 -nextid 3
- go run node/node.go -port 5456 -id 3 -nextport 5454 -nextid 1
- go run node/node.go -port 5454 -id 1 -nextport 5455 -nextid 2 (node with id = 1 will start the token ring)

after starting the nodes, you can type "x" in the terminal to request the access to the critical section.