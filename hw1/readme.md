# Homework 1

**Student ID: 1004875**

**Name: Xiang Siqi**

## Table of Content

* Introduction
* Compilation and Execution
* Assignment 1
  * Client-server architecture
  * Lamport's logical clock
  * Vector clock
* Assignment 2

## Introduction

This homework contains two assignments, each focusing on specific concepts in distributed systems. This README will provide you with a clear understanding of this project implementation, how to compile and execute the assignments, the different types of parameters, and what each assignment entails. Each assignment's implementation will be explained in detail.

The structure of this homework is:

├─hw1
|  ├─Programming_HW1_2023.pdf
|  ├─readme.md
|  ├─assignment1
|  |      ├─main.go
|  |      ├─vectorclock
|  |      |      ├─client.go
|  |      |      ├─message.go
|  |      |      └server.go
|  |      ├─logger
|  |      |   └logger.go
|  |      ├─lamportclock
|  |      |      ├─client.go
|  |      |      ├─message.go
|  |      |      └server.go
|  ├─assignment2
|  |      ├─main.go
|  |      ├─bully
|  |      |   ├─message.go
|  |      |   ├─server.go
|  |      |   └util.go


## Compilation and Execution

To run each assignment, following these steps:

#### Assignment 1

1. Navigate to the `hw1/assignment1` directory.
2. Open a terminal in this directory.
3. Run the following command to start: ```go run main.go```.
4. Note that you can change the type of clock by setting `algorithm` to either `VECTOR_CLOCK` or `LAMPORT_CLOCK` in `main.go`.

#### Assignment 2

1. Navigate to the `hw1/assignment2` directory.
2. Open a terminal in this directory.
3. Run the following command to start: `go run main.go`.

**In this project, all logs generated during the execution will be written and saved in a file called `log.log` under `hw1` folder. There will be no log printed out in the command shell during execution.**

## Assignment 1

### Objects

1. Simulate the behavior of both the server and the registered clients via GO routines.
2. Use Lamport’s logical clock to determine a total order of all the messages received at all the registered clients. Subsequently, present (i.e., print) this order for all registered clients to know the order in which the messages should be read.
3. Use Vector clock to redo the assignment. Implement the detection of causality violation and print any such detected causality violation.

### Implementations

The configuration of this assignment can be customized through the following parameters defined in `main.go` (line 26 - line 28):

* `numOfClients`: Specifies the number of clients participating in this simulation. Set to 10 by default.
* `timeInterval`: Determines the time interval for a client to send a message to the server, measured in second. The default value is 1.
* `algorithm`: You can choose the clock algorithm to use, either `LAMPORT_CLOCK`, or `VECTOR_CLOCK`. The default setting is `VECTOR_CLOCK`.

You can experiment with different parameter values to observe the simulation's behavior.

#### Lamport Clock

The architecture of Lamport Clock simulation is illustrated in Figure 1. The Server and Clients are simulated via GO routines. Note that in this implementation, Each of the server and clients has only one channel.

![fig1_lamport_architecture](assignment1\doc\fig1_lamport_architecture.jpg)

<center style="font-size:14px; color:#C0C0C0">Figure 1: Lamport Clock Architecture</center>

When you run `go run main.go` with the `LAMPORT_CLOCK` algorithm, the program will generate log outputs in `log.log`. Each log entry includes information about the printer, its current clock value, and the operation being performed. Here is a sample log.

```
assignment 1: 2023/10/18 15:13:49 [ Server ] -- Clock 0 -- Server starts listening
assignment 1: 2023/10/18 15:13:49 [Client 0] -- Clock 0 -- client 0 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 1] -- Clock 0 -- client 1 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 2] -- Clock 0 -- client 2 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 3] -- Clock 0 -- client 3 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 4] -- Clock 0 -- client 4 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 5] -- Clock 0 -- client 5 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 6] -- Clock 0 -- client 6 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 7] -- Clock 0 -- client 7 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 8] -- Clock 0 -- client 8 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:49 [Client 9] -- Clock 0 -- client 9 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:13:50 [Client 9] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 3] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 2 -- receive message from client 9
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 3 -- discard message from client 9
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 4 -- receive message from client 3
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 5 -- broadcast message from client 3
assignment 1: 2023/10/18 15:13:50 [Client 5] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 8] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 7] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 6 -- receive message from client 5
assignment 1: 2023/10/18 15:13:50 [Client 6] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 2] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 1] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 0] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 4] -- Clock 1 -- send message to server
assignment 1: 2023/10/18 15:13:50 [Client 9] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 0] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 1] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 2] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 4] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 5] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 6] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 7] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [Client 8] -- Clock 6 -- receive server's broadcast message, originally from client 3
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 7 -- broadcast message from client 5
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 8 -- receive message from client 8
assignment 1: 2023/10/18 15:13:50 [ Server ] -- Clock 9 -- broadcast message from client 8
assignment 1: 2023/10/18 15:13:50 [Client 3] -- Clock 8 -- receive server's broadcast message, originally from client 5
```

#### Vector Clock

The architecture of Vector Clock simulation is illustrated in Figure 2. The Server and Clients are simulated via GO routines. Note that in this implementation, the server will create a channel for each client.

![fig2_vector_architecture](assignment1\doc\fig2_vector_architecture.jpg)

<center style="font-size:14px; color:#C0C0C0">Figure 2: Vector Clock Architecture</center>

When you run `go run main.go` with the `VECTOR_CLOCK` algorithm, the program will generate log outputs in `log.log`. Each log entry includes information about the printer, its current vector clock, and the operation being performed. Here is a sample log.

```
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  1 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Server Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- Server starts listening
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  2 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  3 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  4 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  5 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  6 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  7 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  8 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client  9 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:50 [Client Activate] -- Clock [0 0 0 0 0 0 0 0 0 0 0] -- client 10 starts listening and sends periodical messages
assignment 1: 2023/10/18 15:47:51 [Client  4] -- Clock [0 0 0 0 1 0 0 0 0 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client  5] -- Clock [0 0 0 0 0 1 0 0 0 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client  1] -- Clock [0 1 0 0 0 0 0 0 0 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client  2] -- Clock [0 0 1 0 0 0 0 0 0 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client  8] -- Clock [0 0 0 0 0 0 0 0 1 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client 10] -- Clock [0 0 0 0 0 0 0 0 0 0 1] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client  9] -- Clock [0 0 0 0 0 0 0 0 0 1 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client  6] -- Clock [0 0 0 0 0 0 1 0 0 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [Client  7] -- Clock [0 0 0 0 0 0 0 1 0 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [1 0 0 0 1 0 0 0 0 0 0] -- receive message from client  4
assignment 1: 2023/10/18 15:47:51 [Client  3] -- Clock [0 0 0 1 0 0 0 0 0 0 0] -- send message to server
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [2 0 0 0 1 0 0 0 0 0 0] -- discard message from client 4
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [3 0 0 0 1 1 0 0 0 0 0] -- receive message from client  5
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [4 0 0 0 1 1 0 0 0 0 0] -- broadcast message from client  5
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [5 1 0 0 1 1 0 0 0 0 0] -- receive message from client  1
assignment 1: 2023/10/18 15:47:51 [Client  7] -- Clock [4 0 0 0 1 1 0 2 0 0 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client  2] -- Clock [4 0 2 0 1 1 0 0 0 0 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client  6] -- Clock [4 0 0 0 1 1 2 0 0 0 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client  3] -- Clock [4 0 0 2 1 1 0 0 0 0 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client  4] -- Clock [4 0 0 0 2 1 0 0 0 0 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client  1] -- Clock [4 2 0 0 1 1 0 0 0 0 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client  8] -- Clock [4 0 0 0 1 1 0 0 2 0 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [6 1 0 0 1 1 0 0 0 0 0] -- discard message from client 1
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [7 1 1 0 1 1 0 0 0 0 0] -- receive message from client  2
assignment 1: 2023/10/18 15:47:51 [ Server  ] -- Clock [8 1 1 0 1 1 0 0 0 0 0] -- broadcast message from client  2
assignment 1: 2023/10/18 15:47:51 [Client  8] -- Clock [8 1 1 0 1 1 0 0 3 0 0] -- receive  2's message
assignment 1: 2023/10/18 15:47:51 [Client  9] -- Clock [4 0 0 0 1 1 0 0 0 2 0] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client  9] -- Clock [8 1 1 0 1 1 0 0 0 3 0] -- receive  2's message
assignment 1: 2023/10/18 15:47:51 [Client  1] -- Clock [8 3 1 0 1 1 0 0 0 0 0] -- receive  2's message
assignment 1: 2023/10/18 15:47:51 [Client  3] -- Clock [8 1 1 3 1 1 0 0 0 0 0] -- receive  2's message
assignment 1: 2023/10/18 15:47:51 [Client  4] -- Clock [8 1 1 0 3 1 0 0 0 0 0] -- receive  2's message
assignment 1: 2023/10/18 15:47:51 [Client  5] -- Clock [8 1 1 0 1 2 0 0 0 0 0] -- receive  2's message
assignment 1: 2023/10/18 15:47:51 [Client 10] -- Clock [4 0 0 0 1 1 0 0 0 0 2] -- receive  5's message
assignment 1: 2023/10/18 15:47:51 [Client 10] -- Clock [8 1 1 0 1 1 0 0 0 0 3] -- receive  2's message
```

This implementation does not contain any causality violations by default. This is because The FIFO (First-In-First-Out)  nature of GO channels. It ensures that messages are received in the order they are sent, thereby preventing any causality violations.

However, to introduce and demonstrate the detection of causality violations, I have implemented a `MadlyActive` method. This method instructs the client to intentionally send two messages with causality violations to the server. As a result, the server detects and handles these violations. Below, you can find the relevant code for this feature.

```go
c.incrementClock()
clockSmall := make([]int, len(c.vectorClock))
copy(clockSmall, c.vectorClock)
msgSmall := Message{senderId: c.Id, vectorClock: clockSmall} // construct a message with a smaller vector clock

c.incrementClock()
clockLarge := make([]int, len(c.vectorClock))
copy(clockLarge, c.vectorClock)
msgLarge := Message{senderId: c.Id, vectorClock: clockLarge} // construct a message with a larger vector clock
c.mu.Unlock()

// send the messages in a wrong order
c.sendMsg(msgLarge)
c.sendMsg(msgSmall)
```

There is an implement of "madly launching" in the `main.go`. You have to comment the normal launching code, and uncomment the madly launching code to see the simulation.

You will see something like:

```
assignment 1: 2023/10/18 16:04:25 [Potential Causality Violation Detected on Server when receiving  2's message]
-- Vector Clock on Server -- [14 3 4 3 0 1 0 0 0 1 2]
-- Vector Clock from client  2-- [4 0 3 0 0 0 0 0 0 0 2]
```

