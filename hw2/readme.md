# Homework 2

**Student ID: 1004875**

**Name: Xiang Siqi**

## Table of Content

* Introduction
* Compilation and Execution
  * Configuration
  * Execution
  * Implementation

* Performance
* Appendix

## Introduction

This homework implements three types of distributed mutual exclusion protocols, and compares their performance in terms of time. These three types of protocols are:

1. Lamport's shared priority queue without Ricart and Agrawala's optimization.
2. Lamport's shared priority queue with Ricart and Agrawala's optimization.
3. Voting Protocol with deadlock avoidance.

The structure of this homework is:

├─hw2 
│  │  main.go 
│  │  Programming_HW2_2023.pdf 
│  │  readme.md 
│  │  
│  ├─doc 
│  ├─logger 
│  │      logger.go 
│  │      
│  ├─optimizedsharedpriorityqueue 
│  │      cluster.go 
│  │      server.go 
│  │      
│  ├─sharedpriorityqueue 
│  │      cluster.go 
│  │      server.go 
│  │      
│  ├─util 
│  │      message.go 
│  │      msg_priority_queue.go 
│  │      
│  └─voting 
│          cluster.go 
│          server.go 

## Compilation and Execution

### Configuration

This homework does not use any external package for GO. The custom logger used for this homework is implemented manually and resides in the `hw2/logger` directory.

The entry of the simulation is `main.go` under the `hw2` directory (the main directory). To simulate different scenarios, you can adjust the configuration variables defined in `main.go`.

* The configuration variables are found from line 14 to line 17 in `main.go`. The variables include:
  * **runningMode**: Specify the running mode of the simulation. It is of type `RunningMode` (iota)
    * ***PERFORMANCE_COMPARING_MODE***: Runs all three algorithms and output their performance to a file named `performance.log` under the `hw2` folder.
    * ***SINGLE_PERFORMANCE_MODE***: Runs a specific algorithm (specified by the `algorithm` variable) and outputs its performance to the `performance.log` under the `hw2` folder.
    * ***SINGLE_RUNNING_MODE***: Runs a single algorithm continuously. The requesters will periodically send critical section access requests.
  * **algorithm**: Specifies the algorithm to be run. It is of type `Algorithm` (iota)
    * ***SHARED_PRIORITY_QUEUE***: Lamport's shared priority queue without Ricart and Agrawala's optimization.
    * ***OPTIMIZED_SHARED_PRIORITY_QUEUE***: Lamport's shared priority queue with Ricart and Agrawala's optimization.
    * ***VOTING***: Voting Protocol with deadlock avoidance.
  * **numOfServers**: Specifies the number of servers in the simulation. It is of type `int`.
  * **numOfRequesters**: Specifies the number of servers that will request access to the critical section. This value must be less than or equal to `numOfServers`. It is of type `int`.

### Execution

To start the simulation, following these steps:

1. Navigate to `hw2` directory.

2. Open a terminal in this directory.

3. Run command `go run main.go` to start.

4. By default, the configuration variables are set to:

   ```go
   // line 14 to line 17 in main.go
   runningMode     = PERFORMANCE_COMPARING_MODE
   algorithm       = OPTIMIZED_SHARED_PRIORITY_QUEUE
   numOfServers    = 10
   numOfRequesters = 10
   ```

### Implementation

Some additional information that might be useful are given here:

1. **Log files**

   In this project, all logs generated during the execution will be written and saved in four (3 + 1) files:

   * The running log of SHARED_PRIORITY_QUEUE will be stored in `shared_priority_queue.log` under `hw2` folder.
   * The running log of OPTIMIZED_SHARED_PRIORITY_QUEUE will be stored in `optimizede_queue.log` under `hw2` folder.
   * The running log of VOTING will be stored in `voting_algorithm.log` under `hw2` folder.
   * If `runningMode` is set to `PERFORMANCE_COMPARING_MODE` or `SINGLE_PERFORMANCE_MODE`, an additional log called `performance.log` will be generated under `hw2` folder. It stores the running performance of this simulation.

   There will be no log printed out in the command shell during execution.

2. **Simulation Workload**

   In the simulation implementation, once a server successfully acquires the critical section, it undergoes a simulated workload by sleeping for 1 second before releasing the critical section.

3. **Clock**

   Following the modified instruction, all protocol implementations adopt scalar clocks instead of vector clocks.

## Performance

To obtain the performance for the three implementations, set the `runningMode` parameter to `PERFORMANCE_COMPARING_MODE` in the `main.go` file. Additionally, ensure that you properly set the values for `numOfServers (int)` and `numOfRequesters (int)`. Note that to simulate the workload, each server undergoes a 1 second sleep before releasing the critical section.

In the first experiment, the value of `numOfServers` is set to a constant value of 10. The value of `numOfRequesters` is ranging from 1 to 10. It means that a number of `numOfRequesters` servers will request access to the critical section once during the simulation.

Below are the performance results for the first experiment:

| numOfServers | numOfRequesters | Shared Priority Queue (s) | Optimized Shared Priority Queue (s) | Voting Protocol (s) |
| ------------ | --------------- | ------------------------- | ----------------------------------- | ------------------- |
| 10           | 1               | 1.0053259                 | 1.0042429                           | 1.0055147           |
| 10           | 2               | 2.0149021                 | 2.0103170                           | 2.0303858           |
| 10           | 3               | 3.0266128                 | 3.0306643                           | 3.0268125           |
| 10           | 4               | 4.0369802                 | 4.0279125                           | 4.0401182           |
| 10           | 5               | 5.0559047                 | 5.0455057                           | 5.0620421           |
| 10           | 6               | 6.0569362                 | 6.0563891                           | 6.0641657           |
| 10           | 7               | 7.0718667                 | 7.0518215                           | 7.0653808           |
| 10           | 8               | 8.0637985                 | 8.0538078                           | 8.0805836           |
| 10           | 9               | 9.0621716                 | 9.0596152                           | 9.0875623           |
| 10           | 10              | 10.0958771                | 10.0669812                          | 10.0557619          |

In the second experiment, the value of `numOfServers` is intentionally set to be the same as `numOfRequesters`, ranging from 2 to 11. The 1 server case is intentionally ignored, as it is impractical for a single server to request access. This setup ensures that all servers will request access to the critical section precisely once during the simulation.

Below are the performance result for the second experiment:

| numOfServers | numOfRequesters | Shared Priority Queue (s) | Optimized Shared Priority Queue (s) | Voting Protocol (s) |
| ------------ | --------------- | ------------------------- | ----------------------------------- | ------------------- |
| 2            | 2               | 2.0183889                 | 2.0049165                           | 2.0080983           |
| 3            | 3               | 3.0352952                 | 3.0331427                           | 3.0353020           |
| 4            | 4               | 4.0347881                 | 4.0186394                           | 4.0370095           |
| 5            | 5               | 5.0312824                 | 5.0275441                           | 5.0369574           |
| 6            | 6               | 6.0608842                 | 6.0582980                           | 6.0225319           |
| 7            | 7               | 7.0567044                 | 7.0549108                           | 7.0754125           |
| 8            | 8               | 8.0412546                 | 8.0582339                           | 8.0537618           |
| 9            | 9               | 9.0537163                 | 9.0520255                           | 9.0652227           |
| 10           | 10              | 10.1284541                | 10.1165193                          | 10.1071885          |
| 11           | 11              | 11.0831693                | 11.0816846                          | 11.0940093          |

## Appendix

Here shows an example of `performance.log`.

```
performance:2023/11/19 14:34:09 ########################################################################################
performance:2023/11/19 14:34:09 [Algorithm]: Shared Priority Queue
performance:2023/11/19 14:34:09 [Number of Servers]: 11
performance:2023/11/19 14:34:09 [Number of Requesters]: 11
performance:2023/11/19 14:34:21 [Time (s)]: 11.0831693s
performance:2023/11/19 14:34:21 ########################################################################################
performance:2023/11/19 14:34:21 [Algorithm]: Optimized Shared Priority Queue (Ricart and Agrawala’s Optimization)
performance:2023/11/19 14:34:21 [Number of Servers]: 11
performance:2023/11/19 14:34:21 [Number of Requesters]: 11
performance:2023/11/19 14:34:32 [Time (s)]: 11.0816846s
performance:2023/11/19 14:34:32 ########################################################################################
performance:2023/11/19 14:34:32 [Algorithm]: Voting Protocol
performance:2023/11/19 14:34:32 [Number of Servers]: 11
performance:2023/11/19 14:34:32 [Number of Requesters]: 11
performance:2023/11/19 14:34:43 [Time (s)]: 11.0940093s
performance:2023/11/19 14:34:43 ########################################################################################
```



Here shows an example of `shared_priority_queue.log`.

```
shared priority queue:2023/11/19 14:39:27 [Cluster ] Server 0 added to the cluster
shared priority queue:2023/11/19 14:39:27 [Cluster ] Server 1 added to the cluster
shared priority queue:2023/11/19 14:39:27 [Cluster ] Server 2 added to the cluster
shared priority queue:2023/11/19 14:39:27 [Cluster ] Server 3 added to the cluster
shared priority queue:2023/11/19 14:39:27 [Cluster ] Server 4 added to the cluster
shared priority queue:2023/11/19 14:39:27 [Server 0] Activated as One-time Requester
shared priority queue:2023/11/19 14:39:27 [Server 1] Activated as One-time Requester
shared priority queue:2023/11/19 14:39:27 [Server 2] Activated as One-time Requester
shared priority queue:2023/11/19 14:39:27 [Server 3] Activated as Listener
shared priority queue:2023/11/19 14:39:27 [Server 4] Activated as Listener
shared priority queue:2023/11/19 14:39:27 [Server 2] Sent a request to access the critical section
shared priority queue:2023/11/19 14:39:27 [Server 1] Sent a request to access the critical section
shared priority queue:2023/11/19 14:39:27 [Server 0] Received a request from server 2
shared priority queue:2023/11/19 14:39:27 [Server 2] Received a request from server 1
shared priority queue:2023/11/19 14:39:27 [Server 2] Replied to server 1
shared priority queue:2023/11/19 14:39:27 [Server 0] Sent a request to access the critical section
shared priority queue:2023/11/19 14:39:27 [Server 1] Received a request from server 2
shared priority queue:2023/11/19 14:39:27 [Server 1] Received reply from 2
shared priority queue:2023/11/19 14:39:27 [Server 1] Replied to server 2
shared priority queue:2023/11/19 14:39:27 [Server 1] Received a request from server 0
shared priority queue:2023/11/19 14:39:27 [Server 1] Replied to server 0
shared priority queue:2023/11/19 14:39:27 [Server 3] Received a request from server 2
shared priority queue:2023/11/19 14:39:27 [Server 3] Replied to server 2
shared priority queue:2023/11/19 14:39:27 [Server 3] Received a request from server 1
shared priority queue:2023/11/19 14:39:27 [Server 3] Replied to server 1
shared priority queue:2023/11/19 14:39:27 [Server 2] Received a request from server 0
shared priority queue:2023/11/19 14:39:27 [Server 2] Replied to server 0
shared priority queue:2023/11/19 14:39:27 [Server 2] Received reply from 1
shared priority queue:2023/11/19 14:39:27 [Server 0] Replied to server 2
shared priority queue:2023/11/19 14:39:27 [Server 4] Received a request from server 2
shared priority queue:2023/11/19 14:39:27 [Server 4] Replied to server 2
shared priority queue:2023/11/19 14:39:27 [Server 3] Received a request from server 0
shared priority queue:2023/11/19 14:39:27 [Server 3] Replied to server 0
shared priority queue:2023/11/19 14:39:27 [Server 1] Received reply from 3
shared priority queue:2023/11/19 14:39:27 [Server 2] Received reply from 3
shared priority queue:2023/11/19 14:39:27 [Server 2] Received reply from 0
shared priority queue:2023/11/19 14:39:27 [Server 0] Received reply from 1
shared priority queue:2023/11/19 14:39:27 [Server 0] Received reply from 2
shared priority queue:2023/11/19 14:39:27 [Server 0] Received a request from server 1
shared priority queue:2023/11/19 14:39:27 [Server 0] Replied to server 1
shared priority queue:2023/11/19 14:39:27 [Server 0] Received reply from 3
shared priority queue:2023/11/19 14:39:27 [Server 4] Received a request from server 1
shared priority queue:2023/11/19 14:39:27 [Server 4] Replied to server 1
shared priority queue:2023/11/19 14:39:27 [Server 2] Received reply from 4
shared priority queue:2023/11/19 14:39:27 [Server 1] Received reply from 0
shared priority queue:2023/11/19 14:39:27 [Server 1] Received reply from 4
shared priority queue:2023/11/19 14:39:27 [Server 4] Received a request from server 0
shared priority queue:2023/11/19 14:39:27 [Server 4] Replied to server 0
shared priority queue:2023/11/19 14:39:27 [Server 0] Received reply from 4
shared priority queue:2023/11/19 14:39:27 [Server 0] Executing the critical section
shared priority queue:2023/11/19 14:39:28 [Server 0] Released the critical section
shared priority queue:2023/11/19 14:39:28 [Server 4] Received release from 0
shared priority queue:2023/11/19 14:39:28 [Server 2] Received release from 0
shared priority queue:2023/11/19 14:39:28 [Server 1] Received release from 0
shared priority queue:2023/11/19 14:39:28 [Server 1] Executing the critical section
shared priority queue:2023/11/19 14:39:28 [Server 3] Received release from 0
shared priority queue:2023/11/19 14:39:29 [Server 1] Released the critical section
shared priority queue:2023/11/19 14:39:29 [Server 4] Received release from 1
shared priority queue:2023/11/19 14:39:29 [Server 0] Received release from 1
shared priority queue:2023/11/19 14:39:29 [Server 2] Received release from 1
shared priority queue:2023/11/19 14:39:29 [Server 2] Executing the critical section
shared priority queue:2023/11/19 14:39:29 [Server 3] Received release from 1
shared priority queue:2023/11/19 14:39:30 [Server 2] Released the critical section
shared priority queue:2023/11/19 14:39:30 [Server 1] Received release from 2
shared priority queue:2023/11/19 14:39:30 [Server 0] Received release from 2
shared priority queue:2023/11/19 14:39:30 [Server 3] Received release from 2
shared priority queue:2023/11/19 14:39:30 [Server 4] Received release from 2
```



Here shows an example of `optimized_shared_priority_queue.log`.

```
optimized shared priority queue:2023/11/19 14:39:30 [Cluster ] Server 0 added to the cluster
optimized shared priority queue:2023/11/19 14:39:30 [Cluster ] Server 1 added to the cluster
optimized shared priority queue:2023/11/19 14:39:30 [Cluster ] Server 2 added to the cluster
optimized shared priority queue:2023/11/19 14:39:30 [Cluster ] Server 3 added to the cluster
optimized shared priority queue:2023/11/19 14:39:30 [Cluster ] Server 4 added to the cluster
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Activated as One-time Requester
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Activated as One-time Requester
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Activated as One-time Requester
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Sent a request to access the critical section
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Received a request from server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Replied to server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Received a request from server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Replied to server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 3] Activated as Listener
optimized shared priority queue:2023/11/19 14:39:30 [Server 4] Activated as Listener
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Sent a request to access the critical section
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Sent a request to access the critical section
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Received reply from 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Received reply from 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Received a request from server 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 3] Received a request from server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 4] Received a request from server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Received a request from server 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 3] Replied to server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 4] Replied to server 0
optimized shared priority queue:2023/11/19 14:39:30 [Server 4] Received a request from server 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 4] Replied to server 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 4] Received a request from server 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 3] Received a request from server 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 3] Replied to server 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Received reply from 3
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Received a request from server 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 4] Replied to server 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Received reply from 4
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Received reply from 3
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Received a request from server 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Replied to server 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 2] Received reply from 4
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Received reply from 2
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Received reply from 4
optimized shared priority queue:2023/11/19 14:39:30 [Server 0] Executing the critical section
optimized shared priority queue:2023/11/19 14:39:30 [Server 3] Received a request from server 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 3] Replied to server 1
optimized shared priority queue:2023/11/19 14:39:30 [Server 1] Received reply from 3
optimized shared priority queue:2023/11/19 14:39:31 [Server 0] Finished executing the critical section
optimized shared priority queue:2023/11/19 14:39:31 [Server 0] Replied to server 2
optimized shared priority queue:2023/11/19 14:39:31 [Server 0] Replied to server 1
optimized shared priority queue:2023/11/19 14:39:31 [Server 1] Received reply from 0
optimized shared priority queue:2023/11/19 14:39:31 [Server 1] Executing the critical section
optimized shared priority queue:2023/11/19 14:39:31 [Server 2] Received reply from 0
optimized shared priority queue:2023/11/19 14:39:32 [Server 1] Finished executing the critical section
optimized shared priority queue:2023/11/19 14:39:32 [Server 1] Replied to server 2
optimized shared priority queue:2023/11/19 14:39:32 [Server 2] Received reply from 1
optimized shared priority queue:2023/11/19 14:39:32 [Server 2] Executing the critical section
optimized shared priority queue:2023/11/19 14:39:33 [Server 2] Finished executing the critical section
```



Here shows an example of `voting_algorithm.log`.

```
voting algorithm:2023/11/19 14:39:33 [Cluster ] Server 0 added to the cluster
voting algorithm:2023/11/19 14:39:33 [Cluster ] Server 1 added to the cluster
voting algorithm:2023/11/19 14:39:33 [Cluster ] Server 2 added to the cluster
voting algorithm:2023/11/19 14:39:33 [Cluster ] Server 3 added to the cluster
voting algorithm:2023/11/19 14:39:33 [Cluster ] Server 4 added to the cluster
voting algorithm:2023/11/19 14:39:33 [Server 0] Activated as One-time Requester
voting algorithm:2023/11/19 14:39:33 [Server 1] Activated as One-time Requester
voting algorithm:2023/11/19 14:39:33 [Server 2] Activated as One-time Requester
voting algorithm:2023/11/19 14:39:33 [Server 3] Activated as Listener
voting algorithm:2023/11/19 14:39:33 [Server 4] Activated as Listener
voting algorithm:2023/11/19 14:39:33 [Server 0] Sent a request to access the critical section
voting algorithm:2023/11/19 14:39:33 [Server 2] Received a vote request from server 0
voting algorithm:2023/11/19 14:39:33 [Server 2] Vote for server 0
voting algorithm:2023/11/19 14:39:33 [Server 0] Received a vote request from server 0
voting algorithm:2023/11/19 14:39:33 [Server 0] Vote for server 0
voting algorithm:2023/11/19 14:39:33 [Server 0] Received a vote from server 2
voting algorithm:2023/11/19 14:39:33 [Server 3] Received a vote request from server 0
voting algorithm:2023/11/19 14:39:33 [Server 3] Vote for server 0
voting algorithm:2023/11/19 14:39:33 [Server 0] Received a vote from server 3
voting algorithm:2023/11/19 14:39:33 [Server 4] Received a vote request from server 0
voting algorithm:2023/11/19 14:39:33 [Server 4] Vote for server 0
voting algorithm:2023/11/19 14:39:33 [Server 1] Received a vote request from server 0
voting algorithm:2023/11/19 14:39:33 [Server 1] Vote for server 0
voting algorithm:2023/11/19 14:39:33 [Server 1] Sent a request to access the critical section
voting algorithm:2023/11/19 14:39:33 [Server 2] Sent a request to access the critical section
voting algorithm:2023/11/19 14:39:33 [Server 2] Received a vote request from server 2
voting algorithm:2023/11/19 14:39:33 [Server 0] Received a vote from server 0
voting algorithm:2023/11/19 14:39:33 [Server 0] Executing the critical section
voting algorithm:2023/11/19 14:39:34 [Server 0] Release vote to server 2
voting algorithm:2023/11/19 14:39:34 [Server 0] Release vote to server 3
voting algorithm:2023/11/19 14:39:34 [Server 0] Release vote to server 0
voting algorithm:2023/11/19 14:39:34 [Server 0] Received a vote from server 4
voting algorithm:2023/11/19 14:39:34 [Server 0] Release vote to server 4
voting algorithm:2023/11/19 14:39:34 [Server 0] Received a vote from server 1
voting algorithm:2023/11/19 14:39:34 [Server 0] Release vote to server 1
voting algorithm:2023/11/19 14:39:34 [Server 0] Received a vote request from server 1
voting algorithm:2023/11/19 14:39:34 [Server 0] Received a vote request from server 2
voting algorithm:2023/11/19 14:39:34 [Server 4] Received a vote request from server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a release from server 0
voting algorithm:2023/11/19 14:39:34 [Server 2] Vote for server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a vote from server 2
voting algorithm:2023/11/19 14:39:34 [Server 0] Received a release from server 0
voting algorithm:2023/11/19 14:39:34 [Server 0] Vote for server 1
voting algorithm:2023/11/19 14:39:34 [Server 4] Received a release from server 0
voting algorithm:2023/11/19 14:39:34 [Server 4] Vote for server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a vote from server 4
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a vote request from server 1
voting algorithm:2023/11/19 14:39:34 [Server 2] Rescind vote to server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a rescind request from server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Release rescind vote to server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a release from server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Vote for server 1
voting algorithm:2023/11/19 14:39:34 [Server 1] Received a vote request from server 2
voting algorithm:2023/11/19 14:39:34 [Server 1] Received a vote from server 0
voting algorithm:2023/11/19 14:39:34 [Server 1] Received a release from server 0
voting algorithm:2023/11/19 14:39:34 [Server 1] Vote for server 2
voting algorithm:2023/11/19 14:39:34 [Server 1] Received a vote from server 2
voting algorithm:2023/11/19 14:39:34 [Server 3] Received a vote request from server 2
voting algorithm:2023/11/19 14:39:34 [Server 3] Received a release from server 0
voting algorithm:2023/11/19 14:39:34 [Server 3] Vote for server 2
voting algorithm:2023/11/19 14:39:34 [Server 3] Received a vote request from server 1
voting algorithm:2023/11/19 14:39:34 [Server 3] Rescind vote to server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a rescind request from server 3
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a vote from server 1
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a vote from server 3
voting algorithm:2023/11/19 14:39:34 [Server 2] Release rescind vote to server 3
voting algorithm:2023/11/19 14:39:34 [Server 3] Received a release from server 2
voting algorithm:2023/11/19 14:39:34 [Server 3] Vote for server 1
voting algorithm:2023/11/19 14:39:34 [Server 1] Received a vote from server 3
voting algorithm:2023/11/19 14:39:34 [Server 1] Executing the critical section
voting algorithm:2023/11/19 14:39:34 [Server 4] Received a vote request from server 1
voting algorithm:2023/11/19 14:39:34 [Server 4] Rescind vote to server 2
voting algorithm:2023/11/19 14:39:34 [Server 2] Received a rescind request from server 4
voting algorithm:2023/11/19 14:39:34 [Server 2] Release rescind vote to server 4
voting algorithm:2023/11/19 14:39:34 [Server 4] Received a release from server 2
voting algorithm:2023/11/19 14:39:34 [Server 4] Vote for server 1
voting algorithm:2023/11/19 14:39:35 [Server 1] Release vote to server 0
voting algorithm:2023/11/19 14:39:35 [Server 1] Release vote to server 2
voting algorithm:2023/11/19 14:39:35 [Server 1] Release vote to server 3
voting algorithm:2023/11/19 14:39:35 [Server 1] Received a vote request from server 1
voting algorithm:2023/11/19 14:39:35 [Server 1] Rescind vote to server 2
voting algorithm:2023/11/19 14:39:35 [Server 1] Received a vote from server 4
voting algorithm:2023/11/19 14:39:35 [Server 1] Release vote to server 4
voting algorithm:2023/11/19 14:39:35 [Server 4] Received a release from server 1
voting algorithm:2023/11/19 14:39:35 [Server 2] Received a release from server 1
voting algorithm:2023/11/19 14:39:35 [Server 2] Vote for server 2
voting algorithm:2023/11/19 14:39:35 [Server 0] Received a release from server 1
voting algorithm:2023/11/19 14:39:35 [Server 0] Vote for server 2
voting algorithm:2023/11/19 14:39:35 [Server 3] Received a release from server 1
voting algorithm:2023/11/19 14:39:35 [Server 3] Vote for server 2
voting algorithm:2023/11/19 14:39:35 [Server 4] Vote for server 2
voting algorithm:2023/11/19 14:39:35 [Server 2] Received a rescind request from server 1
voting algorithm:2023/11/19 14:39:35 [Server 2] Release rescind vote to server 1
voting algorithm:2023/11/19 14:39:35 [Server 2] Received a vote from server 2
voting algorithm:2023/11/19 14:39:35 [Server 2] Received a vote from server 0
voting algorithm:2023/11/19 14:39:35 [Server 2] Received a vote from server 3
voting algorithm:2023/11/19 14:39:35 [Server 2] Executing the critical section
voting algorithm:2023/11/19 14:39:35 [Server 1] Received a release from server 2
voting algorithm:2023/11/19 14:39:35 [Server 1] Vote for server 1
voting algorithm:2023/11/19 14:39:35 [Server 1] Received a vote from server 1
voting algorithm:2023/11/19 14:39:35 [Server 1] Release vote to server 1
voting algorithm:2023/11/19 14:39:35 [Server 1] Received a release from server 1
voting algorithm:2023/11/19 14:39:35 [Server 1] Vote for server 2
voting algorithm:2023/11/19 14:39:36 [Server 2] Release vote to server 2
voting algorithm:2023/11/19 14:39:36 [Server 2] Release vote to server 0
voting algorithm:2023/11/19 14:39:36 [Server 2] Release vote to server 3
voting algorithm:2023/11/19 14:39:36 [Server 2] Received a vote from server 4
voting algorithm:2023/11/19 14:39:36 [Server 2] Release vote to server 4
voting algorithm:2023/11/19 14:39:36 [Server 2] Received a vote from server 1
voting algorithm:2023/11/19 14:39:36 [Server 2] Release vote to server 1
voting algorithm:2023/11/19 14:39:36 [Server 1] Received a release from server 2
voting algorithm:2023/11/19 14:39:36 [Server 2] Received a release from server 2
voting algorithm:2023/11/19 14:39:36 [Server 4] Received a release from server 2
voting algorithm:2023/11/19 14:39:36 [Server 0] Received a release from server 2
voting algorithm:2023/11/19 14:39:36 [Server 3] Received a release from server 2
```