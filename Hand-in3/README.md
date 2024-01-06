## Hand-in3
### Discuss, whether you are going to use server-side streaming, client-side streaming, or bidirectional streaming? 
- Use server-side streaming to receive messages from the server. because the client need to keep listening to the server, and the server need to broadcast all the the messages to all the clients.
- Use simple RPC to send messsge from client to server. So only when a client need to send a message, this function will be called. This will save some resources.
### Describe your system architecture - do you have a server-client architecture, peer-to-peer, or something else?
Use server-client architecture. The client send messages to the server, and the server broadcast the messages to all the clients. 
### Describe what RPC methods are implemented, of what type, and what messages types are used for communication
- message ChatMessage: The message type used for communication between server and client. It contains the following fields:
    - senderId: the id of the sender
    - timestamp: the timestamp of the message
    - messageText: the content of the message
- message ServerResponse: The message type used for giving response to the client. It contains the following fields:
    - timestamp: the timestamp of the message
    - result: true if the operation is successful, false otherwise
- message ClientRequest: The message type used for sending request(join or leave) to the server. It contains the following fields:
    - clientId: the id of the sender
    - timestamp: the timestamp of the message
- rpc ClientSend(ChatMessage) returns (ServerResponse): Simple RPC
    - ClientSend is used for publishing messages to the server. It returns a ServerResponse to indicate whether the operation is successful.
- rpc ClientReceive(ClientRequest) returns (stream ChatMessage): Server-side streaming RPC 
    - ClientReceive is used for receiving messages from the server. It returns a stream of ChatMessage. Server-side streaming RPC 
- rpc Join(ClientRequest) returns (ServerResponse): Simple RPC
    - Join is used for sending join request to the server. 
- rpc Leave(ClientRequest) returns (ServerResponse): Simple RPC
    - Leave is used for sending leave request to the server. 
### Describe how you have implemented the calculation of the Lamport timestamps
- At the beginning the clock of server and all the client are set to 0.
- When a client join the server/leave the server, send a message, it will increase the clock by 1. Then attach the current client clock to the TimeStamp.
- When a client receive a message, it will compare the clock of the message and its own clock, and set the clock to the max of the two clocks. Then increase the clock by 1.
- When a server broadcast a message to the clients, it will increase the clock by 1. Then attach the current server clock to the ChatMessage.TimeStamp.
- When a server receive a message, it will compare the clock of the message and its own clock, and set the clock to the max of the two clocks. Then increase the clock by 1.
### Provide a diagram, that traces a sequence of RPC calls together with the Lamport timestamps, that corresponds to a chosen sequence of interactions: Client X joins, Client X Publishes, ..., Client X leaves. Include documentation (system logs) in your appendix.
![alt text](https://github.com/redpandaQAQ/DistributedSystems/blob/main/Hand-in3/diagram.jpg)
- see diagram.log
### Provide a link to a Git repo with your source code in the report
https://github.itu.dk/HelloWorld/Hand-in3
### Include system logs, that document the requirements are met, in the appendix of your report
- see diagram.log
### Include a readme.md file that describes how to run your program. 
- Start server: go run server/server.go
- Start client x and join the chat: go run client/client.go -cid x (replace x with an integer)
- in the client's terminial,type any messsge to publish
- type "leave" to leave the chat
