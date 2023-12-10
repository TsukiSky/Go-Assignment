# Homework 3

**Student ID: 1004875**

**Name: Xiang Siqi**

## Table of Content

* Introduction
* Implementation
* Compilation and Execution
  * Configuration
  * Execution
  * Additional Information
* Performance
* Appendix

## Introduction

This homework implements the original Ivy architecture discussed in the class, and further implements a fault tolerant version of Ivy.

The structure of this homework is:

├─hw3  
│  │  main.go  
│  │  readme.md  
│  │  
│  ├─doc  
│  │      Programming_HW3_2023.pdf  
│  │      
│  ├─logger  
│  │      logger.go  
│  │      
│  ├─ivy  
│  │      central_manager.go  
│  │      processor.go  
│  │      
│  ├─ivyfaulttolerant  
│  │      central_manager.go  
│  │      processor.go  
│  │      
│  ├─util  
│  │      heartbeat.go  
│  │      record.go  
│  │      message.go  

## Implementation

1. The Ivy architecture is implemented under `hw3/ivy`. There are two components: the central manager and the processor.

2. The fault-tolerant Ivy is implemented under `hw3/ivyfaulttolerant`. Similarly, there are two components: the manager and the processor. I added a heartbeat synchronization mechanism to guarantee the fault-tolerance feature of the system. Periodically, the primary manager will send its heartbeat along with `Pagetable` metadata to the backup manager. By this means, the backup manager is able to monitor the state of the primary server and synchronize its page table with the primary server.

3. My fault-tolerant version of Ivy still **preserves sequential consistency**. This is because:

   * All requests sent to the central manager (both the primary and the backup) follow the FIFO rule. This guarantees the sequential property of `<Write>` operations. It further indicates that all writing results will be observed by all processors according to some total order.
   * Any operations in one processor are still executed in order of timestamp.

   These two features allow the fault tolerant version of Ivy to still keep the sequential consistency property.

## Compilation and Execution

### Configuration

This homework does not use any external package for GO. The custom logger used for this homework is implemented manually and resides in the `hw3/logger` directory.

The entry of the simulation is `main.go` under the `hw3` directory (the main directory). To simulate different scenarios, you can adjust the configuration variables defined in `main.go`.

The configuration variables are found from line 13 to line 23 in `main.go`. These variables include:

`numOfProcessor`: The number of processors. By default, it is 10.

`numOfPage`: The number of pages. By defaults, it is set 50.

`readRequestInterval`: The intervals of  `<read, pageId>` requests in one processor, measured in second.

`writeRequestInterval`: The intervals of  `<write, pageId>` requests in one processor. It is worth noting that by adjust the ratio of `readRequestInterval` and `writeRequestInterval`, you can decide the  read/write ratio in the system, measured in second.

`isFaulty`: Set to true if you want to simulate the faulty scenarios, otherwise set to false.

`syncInterval`: The synchronization (heartbeat) interval between the primary central manager and the backup central manager.

`primaryDownTime`: The failing time of the primary central manager, starting from the beginning of the simulation, measured in second.

`primaryFailCount`: The down and restart times for the primary central manager.

`primaryRestartInterval`: The restarting interval of the primary central manager, starting from the shutdown of the primary central manager, measured in second.

`terminateReadNum`: Number of read requests sent by a processor before it stops.

`terminateWriteNum`: Number of write requests sent by a processor before it stops.

### Execution

To start the simulation, following these steps:

1. Navigate to `hw3` directory.

2. Open a terminal in this directory.

3. Run command `go run main.go` to start.

4. By default, the configuration variables are set to:

   ``````go
   // line 13 to line 23 in main.go
   	numOfProcessor         = 10
   	numOfPage              = 50
   	readRequestInterval    = 2 // seconds
   	writeRequestInterval   = 6 // seconds
   	isFaulty               = false
   	syncInterval           = 2 // seconds
   	primaryDownTime        = 5 // seconds
   	primaryDownCount       = 1
   	primaryRestartInterval = 8
   	terminateReadNum       = 15
   	terminateWriteNum      = 5
   ``````

### Additional Information

Some additional information that might be useful are given here:

**Log files**

In this project, all logs generated during the execution will be written and saved in four (3 + 1) files:

* The running log of Ivy will be stored in `assignment_1.log` under `hw3` folder.
* The running log of fault tolerant Ivy will be stored in `assignment_2.log` under `hw3` folder.
* The running performance will be store in `performance.log` under `hw3` folder.

There will be no log printed out in the command shell during execution.

# Performance

The performance evaluation is shown as follows. Note that among all these experiments, some settings hold consistently:

``````go
// consistent factors in experiments
numOfProcessor         = 10
numOfPage              = 50
readRequestInterval    = 1  // seconds
writeRequestInterval   = 3 // seconds
syncInterval           = 2  // seconds
primaryRestartInterval = 8 // seconds
``````

1. **No Fault Ivy Execution Times**

   | Total Number of Requests | Write/Read Ratio | Total Time Used |
   | ------------------------ | ---------------- | --------------- |
   | 200                      | 1/2              | 18.098 s        |
   | 200                      | 1/3              | 17.187 s        |
   | 200                      | 1/4              | 16.208 s        |

   It is apparently that as the number of write requests decreases, the overall execution time becomes shorter.

2. **Fault-tolerant Ivy Execution Times**

   | Total Number of Requests | Write/Read Ratio | Total Time Used |
   | ------------------------ | ---------------- | --------------- |
   | 200                      | 1/2              | 18.191 s        |
   | 200                      | 1/3              | 17.201 s        |
   | 200                      | 1/4              | 16.154 s        |

   Comparing with the results of **No Fault Ivy Execution Times**, the fault-tolerant Ivy has almost the same total execution time. This is because the adding of backup central manager does not affect the overall execution performance.

3. **Fault-tolerant Ivy & The PM Fails Once**

   | Total Number of Requests | Write/Read Ratio | Total Time Used |
   | ------------------------ | ---------------- | --------------- |
   | 200                      | 1/2              | 24.117 s        |
   | 200                      | 1/3              | 17.188 s        |
   | 200                      | 1/4              | 16.114 s        |

   In my experiment, the PM fails at 5 seconds. Please note that I did not implement stopping mechanism in Fault-tolerant cases, which means the program will not stop after reaching the number of requests limit. However, you can check the perfomance.log for the total time used.

4. **Fault-tolerant Ivy & The PM Fails Multiple Times**

   | Total Number of Requests | Write/Read Ratio | PM Fail Counts | Total Time Used |
   | ------------------------ | ---------------- | -------------- | --------------- |
   | 200                      | 1/4              | 1              | 16.114 s        |
   | 200                      | 1/4              | 2              | 17.164 s        |
   | 200                      | 1/4              | 3              | 17.166 s        |
   | 200                      | 1/4              | 4              | 17.117 s        |

   This experiment result indicates that the main factor of execution speed is the write/read ratio. The fail of PM only cost a little overhead in the execution.

## Appendix

Here shows an example of `performance.log`.

```
performance:2023/12/10 20:20:48 Processor 5 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 7 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 2 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 9 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 3 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 6 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 8 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 1 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 0 terminates at 17.1178149 seconds
performance:2023/12/10 20:20:48 Processor 4 terminates at 17.1178149 seconds
```

Here shows an example of `assignment_1.log`.

```
assignment 1: 2023/12/10 19:34:51 [Central Manager] Central Manager activated
assignment 1: 2023/12/10 19:34:52 [Processor 6] -- Send <<<Read Request>>> for Page 49
assignment 1: 2023/12/10 19:34:52 [Processor 5] -- Send <<<Read Request>>> for Page 10
assignment 1: 2023/12/10 19:34:52 [Processor 4] -- Send <<<Read Request>>> for Page 2
assignment 1: 2023/12/10 19:34:52 [Processor 1] -- Send <<<Read Request>>> for Page 41
assignment 1: 2023/12/10 19:34:52 [Processor 9] -- Send <<<Read Request>>> for Page 21
assignment 1: 2023/12/10 19:34:52 [Processor 7] -- Send <<<Read Request>>> for Page 7
assignment 1: 2023/12/10 19:34:52 [Processor 8] -- Send <<<Read Request>>> for Page 13
assignment 1: 2023/12/10 19:34:52 [Processor 3] -- Send <<<Read Request>>> for Page 2
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 49 from Processor 6
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 10 from Processor 5
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 2 from Processor 4
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 41 from Processor 1
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 21 from Processor 9
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 7 from Processor 7
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 13 from Processor 8
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 2 from Processor 3
assignment 1: 2023/12/10 19:34:52 [Processor 9] -- Receive <<<Page>>> Page 21
assignment 1: 2023/12/10 19:34:52 [Processor 7] -- Receive <<<Page>>> Page 7
assignment 1: 2023/12/10 19:34:52 [Processor 5] -- Receive <<<Page>>> Page 10
assignment 1: 2023/12/10 19:34:52 [Processor 2] -- Send <<<Read Request>>> for Page 14
assignment 1: 2023/12/10 19:34:52 [Processor 0] -- Send <<<Read Request>>> for Page 7
assignment 1: 2023/12/10 19:34:52 [Processor 6] -- Receive <<<Page>>> Page 49
assignment 1: 2023/12/10 19:34:52 [Processor 4] -- Receive <<<Page>>> Page 2
assignment 1: 2023/12/10 19:34:52 [Processor 4] -- Receive <<<Read Forward>>> for Page 2 to Processor 3
assignment 1: 2023/12/10 19:34:52 [Processor 1] -- Receive <<<Page>>> Page 41
assignment 1: 2023/12/10 19:34:52 [Processor 8] -- Receive <<<Page>>> Page 13
assignment 1: 2023/12/10 19:34:52 [Processor 3] -- Receive <<<Page>>> Page 2
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 14 from Processor 2
assignment 1: 2023/12/10 19:34:52 [Central Manager] Receive <<<Read Request>>> for Page 7 from Processor 0
assignment 1: 2023/12/10 19:34:52 [Processor 7] -- Receive <<<Read Forward>>> for Page 7 to Processor 0
assignment 1: 2023/12/10 19:34:52 [Processor 2] -- Receive <<<Page>>> Page 14
assignment 1: 2023/12/10 19:34:52 [Processor 0] -- Receive <<<Page>>> Page 7
assignment 1: 2023/12/10 19:34:53 [Processor 6] -- Send <<<Read Request>>> for Page 12
assignment 1: 2023/12/10 19:34:53 [Processor 9] -- Send <<<Read Request>>> for Page 45
assignment 1: 2023/12/10 19:34:53 [Processor 1] -- Send <<<Read Request>>> for Page 48
assignment 1: 2023/12/10 19:34:53 [Processor 0] -- Send <<<Read Request>>> for Page 5
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 12 from Processor 6
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 45 from Processor 9
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 48 from Processor 1
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 5 from Processor 0
assignment 1: 2023/12/10 19:34:53 [Processor 0] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:34:53 [Processor 1] -- Receive <<<Page>>> Page 48
assignment 1: 2023/12/10 19:34:53 [Processor 9] -- Receive <<<Page>>> Page 45
assignment 1: 2023/12/10 19:34:53 [Processor 4] -- Send <<<Read Request>>> for Page 15
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 15 from Processor 4
assignment 1: 2023/12/10 19:34:53 [Processor 7] -- Send <<<Read Request>>> for Page 34
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 7
assignment 1: 2023/12/10 19:34:53 [Processor 8] -- Send <<<Read Request>>> for Page 17
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 17 from Processor 8
assignment 1: 2023/12/10 19:34:53 [Processor 2] -- Send <<<Read Request>>> for Page 27
assignment 1: 2023/12/10 19:34:53 [Processor 8] -- Receive <<<Page>>> Page 17
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 27 from Processor 2
assignment 1: 2023/12/10 19:34:53 [Processor 2] -- Receive <<<Page>>> Page 27
assignment 1: 2023/12/10 19:34:53 [Processor 3] -- Send <<<Read Request>>> for Page 30
assignment 1: 2023/12/10 19:34:53 [Processor 5] -- Send <<<Read Request>>> for Page 7
assignment 1: 2023/12/10 19:34:53 [Processor 6] -- Receive <<<Page>>> Page 12
assignment 1: 2023/12/10 19:34:53 [Processor 4] -- Receive <<<Page>>> Page 15
assignment 1: 2023/12/10 19:34:53 [Processor 7] -- Receive <<<Page>>> Page 34
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 30 from Processor 3
assignment 1: 2023/12/10 19:34:53 [Central Manager] Receive <<<Read Request>>> for Page 7 from Processor 5
assignment 1: 2023/12/10 19:34:53 [Processor 3] -- Receive <<<Page>>> Page 30
assignment 1: 2023/12/10 19:34:53 [Processor 7] -- Receive <<<Read Forward>>> for Page 7 to Processor 5
assignment 1: 2023/12/10 19:34:53 [Processor 5] -- Receive <<<Page>>> Page 7
assignment 1: 2023/12/10 19:34:54 [Processor 1] -- Send <<<Write Request>>> for Page 48
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 48 from Processor 1
assignment 1: 2023/12/10 19:34:54 [Processor 4] -- Send <<<Write Request>>> for Page 2
assignment 1: 2023/12/10 19:34:54 [Processor 9] -- Send <<<Write Request>>> for Page 12
assignment 1: 2023/12/10 19:34:54 [Processor 5] -- Send <<<Write Request>>> for Page 19
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 2 from Processor 4
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 12 from Processor 9
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 19 from Processor 5
assignment 1: 2023/12/10 19:34:54 [Processor 5] -- Receive <<<Page>>> Page 19
assignment 1: 2023/12/10 19:34:54 [Processor 7] -- Send <<<Write Request>>> for Page 32
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 32 from Processor 7
assignment 1: 2023/12/10 19:34:54 [Processor 7] -- Receive <<<Page>>> Page 32
assignment 1: 2023/12/10 19:34:54 [Processor 3] -- Receive <<<Invalidate>>> Page 2
assignment 1: 2023/12/10 19:34:54 [Processor 6] -- Receive <<<Write Forward>>> for Page 12 to Processor 9
assignment 1: 2023/12/10 19:34:54 [Processor 9] -- Receive <<<Page>>> Page 12
assignment 1: 2023/12/10 19:34:54 [Processor 6] -- Receive <<<Invalidate>>> Page 12
assignment 1: 2023/12/10 19:34:54 [Processor 3] -- Send <<<Write Request>>> for Page 13
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 13 from Processor 3
assignment 1: 2023/12/10 19:34:54 [Processor 6] -- Send <<<Write Request>>> for Page 3
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 3 from Processor 6
assignment 1: 2023/12/10 19:34:54 [Processor 6] -- Receive <<<Page>>> Page 3
assignment 1: 2023/12/10 19:34:54 [Processor 8] -- Receive <<<Write Forward>>> for Page 13 to Processor 3
assignment 1: 2023/12/10 19:34:54 [Processor 8] -- Receive <<<Invalidate>>> Page 13
assignment 1: 2023/12/10 19:34:54 [Processor 8] -- Send <<<Write Request>>> for Page 23
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 23 from Processor 8
assignment 1: 2023/12/10 19:34:54 [Processor 8] -- Receive <<<Page>>> Page 23
assignment 1: 2023/12/10 19:34:54 [Processor 3] -- Receive <<<Page>>> Page 13
assignment 1: 2023/12/10 19:34:54 [Processor 2] -- Send <<<Write Request>>> for Page 9
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 9 from Processor 2
assignment 1: 2023/12/10 19:34:54 [Processor 2] -- Receive <<<Page>>> Page 9
assignment 1: 2023/12/10 19:34:54 [Processor 0] -- Send <<<Write Request>>> for Page 5
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Write Request>>> for Page 5 from Processor 0
assignment 1: 2023/12/10 19:34:54 [Processor 5] -- Send <<<Read Request>>> for Page 5
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 5 from Processor 5
assignment 1: 2023/12/10 19:34:54 [Processor 0] -- Receive <<<Read Forward>>> for Page 5 to Processor 5
assignment 1: 2023/12/10 19:34:54 [Processor 5] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:34:54 [Processor 9] -- Send <<<Read Request>>> for Page 16
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 16 from Processor 9
assignment 1: 2023/12/10 19:34:54 [Processor 9] -- Receive <<<Page>>> Page 16
assignment 1: 2023/12/10 19:34:54 [Processor 0] -- Send <<<Read Request>>> for Page 11
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 11 from Processor 0
assignment 1: 2023/12/10 19:34:54 [Processor 0] -- Receive <<<Page>>> Page 11
assignment 1: 2023/12/10 19:34:54 [Processor 3] -- Send <<<Read Request>>> for Page 17
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 17 from Processor 3
assignment 1: 2023/12/10 19:34:54 [Processor 6] -- Send <<<Read Request>>> for Page 8
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 8 from Processor 6
assignment 1: 2023/12/10 19:34:54 [Processor 6] -- Receive <<<Page>>> Page 8
assignment 1: 2023/12/10 19:34:54 [Processor 1] -- Send <<<Read Request>>> for Page 2
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 2 from Processor 1
assignment 1: 2023/12/10 19:34:54 [Processor 4] -- Receive <<<Read Forward>>> for Page 2 to Processor 1
assignment 1: 2023/12/10 19:34:54 [Processor 1] -- Receive <<<Page>>> Page 2
assignment 1: 2023/12/10 19:34:54 [Processor 8] -- Receive <<<Read Forward>>> for Page 17 to Processor 3
assignment 1: 2023/12/10 19:34:54 [Processor 3] -- Receive <<<Page>>> Page 17
assignment 1: 2023/12/10 19:34:54 [Processor 4] -- Send <<<Read Request>>> for Page 43
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 43 from Processor 4
assignment 1: 2023/12/10 19:34:54 [Processor 4] -- Receive <<<Page>>> Page 43
assignment 1: 2023/12/10 19:34:54 [Processor 8] -- Send <<<Read Request>>> for Page 7
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 7 from Processor 8
assignment 1: 2023/12/10 19:34:54 [Processor 7] -- Receive <<<Read Forward>>> for Page 7 to Processor 8
assignment 1: 2023/12/10 19:34:54 [Processor 8] -- Receive <<<Page>>> Page 7
assignment 1: 2023/12/10 19:34:54 [Processor 2] -- Send <<<Read Request>>> for Page 39
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 2
assignment 1: 2023/12/10 19:34:54 [Processor 2] -- Receive <<<Page>>> Page 39
assignment 1: 2023/12/10 19:34:54 [Processor 7] -- Send <<<Read Request>>> for Page 35
assignment 1: 2023/12/10 19:34:54 [Central Manager] Receive <<<Read Request>>> for Page 35 from Processor 7
assignment 1: 2023/12/10 19:34:54 [Processor 7] -- Receive <<<Page>>> Page 35
assignment 1: 2023/12/10 19:34:55 [Processor 3] -- Send <<<Read Request>>> for Page 41
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 41 from Processor 3
assignment 1: 2023/12/10 19:34:55 [Processor 2] -- Send <<<Read Request>>> for Page 20
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 20 from Processor 2
assignment 1: 2023/12/10 19:34:55 [Processor 2] -- Receive <<<Page>>> Page 20
assignment 1: 2023/12/10 19:34:55 [Processor 8] -- Send <<<Read Request>>> for Page 45
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 45 from Processor 8
assignment 1: 2023/12/10 19:34:55 [Processor 5] -- Send <<<Read Request>>> for Page 25
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 25 from Processor 5
assignment 1: 2023/12/10 19:34:55 [Processor 5] -- Receive <<<Page>>> Page 25
assignment 1: 2023/12/10 19:34:55 [Processor 9] -- Send <<<Read Request>>> for Page 30
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 30 from Processor 9
assignment 1: 2023/12/10 19:34:55 [Processor 3] -- Receive <<<Read Forward>>> for Page 30 to Processor 9
assignment 1: 2023/12/10 19:34:55 [Processor 0] -- Send <<<Read Request>>> for Page 40
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 40 from Processor 0
assignment 1: 2023/12/10 19:34:55 [Processor 4] -- Send <<<Read Request>>> for Page 36
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 36 from Processor 4
assignment 1: 2023/12/10 19:34:55 [Processor 6] -- Send <<<Read Request>>> for Page 9
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 9 from Processor 6
assignment 1: 2023/12/10 19:34:55 [Processor 2] -- Receive <<<Read Forward>>> for Page 9 to Processor 6
assignment 1: 2023/12/10 19:34:55 [Processor 1] -- Send <<<Read Request>>> for Page 49
assignment 1: 2023/12/10 19:34:55 [Processor 7] -- Send <<<Read Request>>> for Page 4
assignment 1: 2023/12/10 19:34:55 [Processor 1] -- Receive <<<Read Forward>>> for Page 41 to Processor 3
assignment 1: 2023/12/10 19:34:55 [Processor 9] -- Receive <<<Read Forward>>> for Page 45 to Processor 8
assignment 1: 2023/12/10 19:34:55 [Processor 9] -- Receive <<<Page>>> Page 30
assignment 1: 2023/12/10 19:34:55 [Processor 0] -- Receive <<<Page>>> Page 40
assignment 1: 2023/12/10 19:34:55 [Processor 4] -- Receive <<<Page>>> Page 36
assignment 1: 2023/12/10 19:34:55 [Processor 6] -- Receive <<<Page>>> Page 9
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 49 from Processor 1
assignment 1: 2023/12/10 19:34:55 [Central Manager] Receive <<<Read Request>>> for Page 4 from Processor 7
assignment 1: 2023/12/10 19:34:55 [Processor 3] -- Receive <<<Page>>> Page 41
assignment 1: 2023/12/10 19:34:55 [Processor 8] -- Receive <<<Page>>> Page 45
assignment 1: 2023/12/10 19:34:55 [Processor 7] -- Receive <<<Page>>> Page 4
assignment 1: 2023/12/10 19:34:55 [Processor 6] -- Receive <<<Read Forward>>> for Page 49 to Processor 1
assignment 1: 2023/12/10 19:34:55 [Processor 1] -- Receive <<<Page>>> Page 49
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 3 from Processor 6
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 19 from Processor 5
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 48 from Processor 1
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 13 from Processor 3
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 9 from Processor 2
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 5 from Processor 0
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 32 from Processor 7
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 23 from Processor 8
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 2 from Processor 4
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Write Ack>>> for Page 12 from Processor 9
assignment 1: 2023/12/10 19:34:56 [Processor 4] -- Send <<<Read Request>>> for Page 24
assignment 1: 2023/12/10 19:34:56 [Processor 2] -- Send <<<Read Request>>> for Page 35
assignment 1: 2023/12/10 19:34:56 [Processor 1] -- Send <<<Read Request>>> for Page 11
assignment 1: 2023/12/10 19:34:56 [Processor 3] -- Send <<<Read Request>>> for Page 22
assignment 1: 2023/12/10 19:34:56 [Processor 9] -- Send <<<Read Request>>> for Page 48
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 4
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 35 from Processor 2
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 11 from Processor 1
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 22 from Processor 3
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 48 from Processor 9
assignment 1: 2023/12/10 19:34:56 [Processor 7] -- Receive <<<Read Forward>>> for Page 35 to Processor 2
assignment 1: 2023/12/10 19:34:56 [Processor 4] -- Receive <<<Page>>> Page 24
assignment 1: 2023/12/10 19:34:56 [Processor 1] -- Receive <<<Read Forward>>> for Page 48 to Processor 9
assignment 1: 2023/12/10 19:34:56 [Processor 9] -- Receive <<<Page>>> Page 48
assignment 1: 2023/12/10 19:34:56 [Processor 7] -- Send <<<Read Request>>> for Page 11
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 11 from Processor 7
assignment 1: 2023/12/10 19:34:56 [Processor 5] -- Send <<<Read Request>>> for Page 12
assignment 1: 2023/12/10 19:34:56 [Processor 0] -- Send <<<Read Request>>> for Page 46
assignment 1: 2023/12/10 19:34:56 [Processor 6] -- Read Page 8 from local page table
assignment 1: 2023/12/10 19:34:56 [Processor 8] -- Send <<<Read Request>>> for Page 43
assignment 1: 2023/12/10 19:34:56 [Processor 0] -- Receive <<<Read Forward>>> for Page 11 to Processor 1
assignment 1: 2023/12/10 19:34:56 [Processor 0] -- Receive <<<Read Forward>>> for Page 11 to Processor 7
assignment 1: 2023/12/10 19:34:56 [Processor 3] -- Receive <<<Page>>> Page 22
assignment 1: 2023/12/10 19:34:56 [Processor 2] -- Receive <<<Page>>> Page 35
assignment 1: 2023/12/10 19:34:56 [Processor 7] -- Receive <<<Page>>> Page 11
assignment 1: 2023/12/10 19:34:56 [Processor 1] -- Receive <<<Page>>> Page 11
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 12 from Processor 5
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 46 from Processor 0
assignment 1: 2023/12/10 19:34:56 [Central Manager] Receive <<<Read Request>>> for Page 43 from Processor 8
assignment 1: 2023/12/10 19:34:56 [Processor 4] -- Receive <<<Read Forward>>> for Page 43 to Processor 8
assignment 1: 2023/12/10 19:34:56 [Processor 8] -- Receive <<<Page>>> Page 43
assignment 1: 2023/12/10 19:34:56 [Processor 9] -- Receive <<<Read Forward>>> for Page 12 to Processor 5
assignment 1: 2023/12/10 19:34:56 [Processor 5] -- Receive <<<Page>>> Page 12
assignment 1: 2023/12/10 19:34:56 [Processor 0] -- Receive <<<Page>>> Page 46
assignment 1: 2023/12/10 19:34:57 [Processor 0] -- Send <<<Write Request>>> for Page 7
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 7 from Processor 0
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Send <<<Write Request>>> for Page 41
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 41 from Processor 5
assignment 1: 2023/12/10 19:34:57 [Processor 8] -- Send <<<Write Request>>> for Page 33
assignment 1: 2023/12/10 19:34:57 [Processor 1] -- Receive <<<Write Forward>>> for Page 41 to Processor 5
assignment 1: 2023/12/10 19:34:57 [Processor 1] -- Receive <<<Invalidate>>> Page 41
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 33 from Processor 8
assignment 1: 2023/12/10 19:34:57 [Processor 8] -- Receive <<<Page>>> Page 33
assignment 1: 2023/12/10 19:34:57 [Processor 1] -- Send <<<Write Request>>> for Page 25
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 25 from Processor 1
assignment 1: 2023/12/10 19:34:57 [Processor 9] -- Send <<<Write Request>>> for Page 1
assignment 1: 2023/12/10 19:34:57 [Processor 6] -- Send <<<Write Request>>> for Page 2
assignment 1: 2023/12/10 19:34:57 [Processor 7] -- Send <<<Write Request>>> for Page 28
assignment 1: 2023/12/10 19:34:57 [Processor 3] -- Send <<<Write Request>>> for Page 47
assignment 1: 2023/12/10 19:34:57 [Processor 7] -- Receive <<<Write Forward>>> for Page 7 to Processor 0
assignment 1: 2023/12/10 19:34:57 [Processor 7] -- Receive <<<Invalidate>>> Page 7
assignment 1: 2023/12/10 19:34:57 [Processor 2] -- Send <<<Write Request>>> for Page 9
assignment 1: 2023/12/10 19:34:57 [Processor 0] -- Receive <<<Page>>> Page 7
assignment 1: 2023/12/10 19:34:57 [Processor 4] -- Send <<<Write Request>>> for Page 0
assignment 1: 2023/12/10 19:34:57 [Processor 3] -- Receive <<<Invalidate>>> Page 41
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Receive <<<Page>>> Page 41
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Receive <<<Write Forward>>> for Page 25 to Processor 1
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Receive <<<Invalidate>>> Page 7
assignment 1: 2023/12/10 19:34:57 [Processor 1] -- Receive <<<Page>>> Page 25
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 1 from Processor 9
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 2 from Processor 6
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Receive <<<Invalidate>>> Page 25
assignment 1: 2023/12/10 19:34:57 [Processor 8] -- Receive <<<Invalidate>>> Page 7
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 28 from Processor 7
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 47 from Processor 3
assignment 1: 2023/12/10 19:34:57 [Processor 9] -- Receive <<<Page>>> Page 1
assignment 1: 2023/12/10 19:34:57 [Processor 4] -- Receive <<<Write Forward>>> for Page 2 to Processor 6
assignment 1: 2023/12/10 19:34:57 [Processor 6] -- Receive <<<Page>>> Page 2
assignment 1: 2023/12/10 19:34:57 [Processor 1] -- Receive <<<Invalidate>>> Page 2
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 9 from Processor 2
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Write Request>>> for Page 0 from Processor 4
assignment 1: 2023/12/10 19:34:57 [Processor 7] -- Receive <<<Page>>> Page 28
assignment 1: 2023/12/10 19:34:57 [Processor 3] -- Receive <<<Page>>> Page 47
assignment 1: 2023/12/10 19:34:57 [Processor 4] -- Receive <<<Page>>> Page 0
assignment 1: 2023/12/10 19:34:57 [Processor 6] -- Receive <<<Invalidate>>> Page 9
assignment 1: 2023/12/10 19:34:57 [Processor 7] -- Send <<<Read Request>>> for Page 46
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 46 from Processor 7
assignment 1: 2023/12/10 19:34:57 [Processor 8] -- Send <<<Read Request>>> for Page 5
assignment 1: 2023/12/10 19:34:57 [Processor 6] -- Send <<<Read Request>>> for Page 32
assignment 1: 2023/12/10 19:34:57 [Processor 0] -- Receive <<<Read Forward>>> for Page 46 to Processor 7
assignment 1: 2023/12/10 19:34:57 [Processor 7] -- Receive <<<Page>>> Page 46
assignment 1: 2023/12/10 19:34:57 [Processor 4] -- Send <<<Read Request>>> for Page 8
assignment 1: 2023/12/10 19:34:57 [Processor 9] -- Send <<<Read Request>>> for Page 19
assignment 1: 2023/12/10 19:34:57 [Processor 3] -- Send <<<Read Request>>> for Page 19
assignment 1: 2023/12/10 19:34:57 [Processor 0] -- Send <<<Read Request>>> for Page 33
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Send <<<Read Request>>> for Page 2
assignment 1: 2023/12/10 19:34:57 [Processor 1] -- Send <<<Read Request>>> for Page 39
assignment 1: 2023/12/10 19:34:57 [Processor 2] -- Send <<<Read Request>>> for Page 37
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 5 from Processor 8
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 32 from Processor 6
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 8 from Processor 4
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 19 from Processor 9
assignment 1: 2023/12/10 19:34:57 [Processor 0] -- Receive <<<Read Forward>>> for Page 5 to Processor 8
assignment 1: 2023/12/10 19:34:57 [Processor 8] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 19 from Processor 3
assignment 1: 2023/12/10 19:34:57 [Processor 7] -- Receive <<<Read Forward>>> for Page 32 to Processor 6
assignment 1: 2023/12/10 19:34:57 [Processor 6] -- Receive <<<Read Forward>>> for Page 8 to Processor 4
assignment 1: 2023/12/10 19:34:57 [Processor 6] -- Receive <<<Page>>> Page 32
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Receive <<<Read Forward>>> for Page 19 to Processor 9
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Receive <<<Read Forward>>> for Page 19 to Processor 3
assignment 1: 2023/12/10 19:34:57 [Processor 4] -- Receive <<<Page>>> Page 8
assignment 1: 2023/12/10 19:34:57 [Processor 3] -- Receive <<<Page>>> Page 19
assignment 1: 2023/12/10 19:34:57 [Processor 9] -- Receive <<<Page>>> Page 19
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 33 from Processor 0
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 2 from Processor 5
assignment 1: 2023/12/10 19:34:57 [Processor 8] -- Receive <<<Read Forward>>> for Page 33 to Processor 0
assignment 1: 2023/12/10 19:34:57 [Processor 6] -- Receive <<<Read Forward>>> for Page 2 to Processor 5
assignment 1: 2023/12/10 19:34:57 [Processor 5] -- Receive <<<Page>>> Page 2
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 1
assignment 1: 2023/12/10 19:34:57 [Central Manager] Receive <<<Read Request>>> for Page 37 from Processor 2
assignment 1: 2023/12/10 19:34:57 [Processor 0] -- Receive <<<Page>>> Page 33
assignment 1: 2023/12/10 19:34:57 [Processor 2] -- Receive <<<Read Forward>>> for Page 39 to Processor 1
assignment 1: 2023/12/10 19:34:57 [Processor 2] -- Receive <<<Page>>> Page 37
assignment 1: 2023/12/10 19:34:57 [Processor 1] -- Receive <<<Page>>> Page 39
assignment 1: 2023/12/10 19:34:58 [Processor 0] -- Send <<<Read Request>>> for Page 34
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 0
assignment 1: 2023/12/10 19:34:58 [Processor 7] -- Receive <<<Read Forward>>> for Page 34 to Processor 0
assignment 1: 2023/12/10 19:34:58 [Processor 0] -- Receive <<<Page>>> Page 34
assignment 1: 2023/12/10 19:34:58 [Processor 7] -- Send <<<Read Request>>> for Page 47
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 47 from Processor 7
assignment 1: 2023/12/10 19:34:58 [Processor 3] -- Receive <<<Read Forward>>> for Page 47 to Processor 7
assignment 1: 2023/12/10 19:34:58 [Processor 7] -- Receive <<<Page>>> Page 47
assignment 1: 2023/12/10 19:34:58 [Processor 3] -- Send <<<Read Request>>> for Page 24
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 3
assignment 1: 2023/12/10 19:34:58 [Processor 4] -- Receive <<<Read Forward>>> for Page 24 to Processor 3
assignment 1: 2023/12/10 19:34:58 [Processor 3] -- Receive <<<Page>>> Page 24
assignment 1: 2023/12/10 19:34:58 [Processor 9] -- Send <<<Read Request>>> for Page 8
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 8 from Processor 9
assignment 1: 2023/12/10 19:34:58 [Processor 4] -- Send <<<Read Request>>> for Page 28
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 28 from Processor 4
assignment 1: 2023/12/10 19:34:58 [Processor 7] -- Receive <<<Read Forward>>> for Page 28 to Processor 4
assignment 1: 2023/12/10 19:34:58 [Processor 4] -- Receive <<<Page>>> Page 28
assignment 1: 2023/12/10 19:34:58 [Processor 6] -- Receive <<<Read Forward>>> for Page 8 to Processor 9
assignment 1: 2023/12/10 19:34:58 [Processor 9] -- Receive <<<Page>>> Page 8
assignment 1: 2023/12/10 19:34:58 [Processor 2] -- Read Page 9 from local page table
assignment 1: 2023/12/10 19:34:58 [Processor 5] -- Send <<<Read Request>>> for Page 1
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 1 from Processor 5
assignment 1: 2023/12/10 19:34:58 [Processor 1] -- Read Page 11 from local page table
assignment 1: 2023/12/10 19:34:58 [Processor 6] -- Send <<<Read Request>>> for Page 13
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 13 from Processor 6
assignment 1: 2023/12/10 19:34:58 [Processor 3] -- Receive <<<Read Forward>>> for Page 13 to Processor 6
assignment 1: 2023/12/10 19:34:58 [Processor 8] -- Send <<<Read Request>>> for Page 8
assignment 1: 2023/12/10 19:34:58 [Central Manager] Receive <<<Read Request>>> for Page 8 from Processor 8
assignment 1: 2023/12/10 19:34:58 [Processor 9] -- Receive <<<Read Forward>>> for Page 1 to Processor 5
assignment 1: 2023/12/10 19:34:58 [Processor 5] -- Receive <<<Page>>> Page 1
assignment 1: 2023/12/10 19:34:58 [Processor 6] -- Receive <<<Page>>> Page 13
assignment 1: 2023/12/10 19:34:58 [Processor 6] -- Receive <<<Read Forward>>> for Page 8 to Processor 8
assignment 1: 2023/12/10 19:34:58 [Processor 8] -- Receive <<<Page>>> Page 8
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 25 from Processor 1
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 41 from Processor 5
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 2 from Processor 6
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 9 from Processor 2
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 33 from Processor 8
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 7 from Processor 0
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 1 from Processor 9
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 47 from Processor 3
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 28 from Processor 7
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Write Ack>>> for Page 0 from Processor 4
assignment 1: 2023/12/10 19:34:59 [Processor 6] -- Send <<<Read Request>>> for Page 21
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 21 from Processor 6
assignment 1: 2023/12/10 19:34:59 [Processor 8] -- Send <<<Read Request>>> for Page 44
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 44 from Processor 8
assignment 1: 2023/12/10 19:34:59 [Processor 9] -- Receive <<<Read Forward>>> for Page 21 to Processor 6
assignment 1: 2023/12/10 19:34:59 [Processor 6] -- Receive <<<Page>>> Page 21
assignment 1: 2023/12/10 19:34:59 [Processor 1] -- Send <<<Read Request>>> for Page 6
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 6 from Processor 1
assignment 1: 2023/12/10 19:34:59 [Processor 1] -- Receive <<<Page>>> Page 6
assignment 1: 2023/12/10 19:34:59 [Processor 4] -- Send <<<Read Request>>> for Page 21
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 21 from Processor 4
assignment 1: 2023/12/10 19:34:59 [Processor 9] -- Receive <<<Read Forward>>> for Page 21 to Processor 4
assignment 1: 2023/12/10 19:34:59 [Processor 4] -- Receive <<<Page>>> Page 21
assignment 1: 2023/12/10 19:34:59 [Processor 5] -- Read Page 41 from local page table
assignment 1: 2023/12/10 19:34:59 [Processor 0] -- Send <<<Read Request>>> for Page 48
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 48 from Processor 0
assignment 1: 2023/12/10 19:34:59 [Processor 1] -- Receive <<<Read Forward>>> for Page 48 to Processor 0
assignment 1: 2023/12/10 19:34:59 [Processor 0] -- Receive <<<Page>>> Page 48
assignment 1: 2023/12/10 19:34:59 [Processor 7] -- Send <<<Read Request>>> for Page 20
assignment 1: 2023/12/10 19:34:59 [Processor 3] -- Send <<<Read Request>>> for Page 20
assignment 1: 2023/12/10 19:34:59 [Processor 9] -- Send <<<Read Request>>> for Page 5
assignment 1: 2023/12/10 19:34:59 [Processor 2] -- Send <<<Read Request>>> for Page 0
assignment 1: 2023/12/10 19:34:59 [Processor 8] -- Receive <<<Page>>> Page 44
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 20 from Processor 7
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 20 from Processor 3
assignment 1: 2023/12/10 19:34:59 [Processor 2] -- Receive <<<Read Forward>>> for Page 20 to Processor 7
assignment 1: 2023/12/10 19:34:59 [Processor 2] -- Receive <<<Read Forward>>> for Page 20 to Processor 3
assignment 1: 2023/12/10 19:34:59 [Processor 7] -- Receive <<<Page>>> Page 20
assignment 1: 2023/12/10 19:34:59 [Processor 3] -- Receive <<<Page>>> Page 20
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 5 from Processor 9
assignment 1: 2023/12/10 19:34:59 [Central Manager] Receive <<<Read Request>>> for Page 0 from Processor 2
assignment 1: 2023/12/10 19:34:59 [Processor 4] -- Receive <<<Read Forward>>> for Page 0 to Processor 2
assignment 1: 2023/12/10 19:34:59 [Processor 2] -- Receive <<<Page>>> Page 0
assignment 1: 2023/12/10 19:34:59 [Processor 0] -- Receive <<<Read Forward>>> for Page 5 to Processor 9
assignment 1: 2023/12/10 19:34:59 [Processor 9] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:35:00 [Processor 4] -- Send <<<Write Request>>> for Page 19
assignment 1: 2023/12/10 19:35:00 [Processor 9] -- Send <<<Write Request>>> for Page 12
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 19 from Processor 4
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 12 from Processor 9
assignment 1: 2023/12/10 19:35:00 [Processor 0] -- Send <<<Write Request>>> for Page 9
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 9 from Processor 0
assignment 1: 2023/12/10 19:35:00 [Processor 3] -- Send <<<Write Request>>> for Page 27
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 27 from Processor 3
assignment 1: 2023/12/10 19:35:00 [Processor 9] -- Receive <<<Invalidate>>> Page 19
assignment 1: 2023/12/10 19:35:00 [Processor 3] -- Receive <<<Invalidate>>> Page 19
assignment 1: 2023/12/10 19:35:00 [Processor 6] -- Send <<<Write Request>>> for Page 39
assignment 1: 2023/12/10 19:35:00 [Processor 7] -- Send <<<Write Request>>> for Page 12
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Send <<<Write Request>>> for Page 5
assignment 1: 2023/12/10 19:35:00 [Processor 5] -- Send <<<Write Request>>> for Page 5
assignment 1: 2023/12/10 19:35:00 [Processor 1] -- Send <<<Write Request>>> for Page 16
assignment 1: 2023/12/10 19:35:00 [Processor 5] -- Receive <<<Write Forward>>> for Page 19 to Processor 4
assignment 1: 2023/12/10 19:35:00 [Processor 5] -- Receive <<<Invalidate>>> Page 12
assignment 1: 2023/12/10 19:35:00 [Processor 4] -- Receive <<<Page>>> Page 19
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Write Forward>>> for Page 9 to Processor 0
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Write Forward>>> for Page 27 to Processor 3
assignment 1: 2023/12/10 19:35:00 [Processor 8] -- Send <<<Write Request>>> for Page 48
assignment 1: 2023/12/10 19:35:00 [Processor 3] -- Receive <<<Page>>> Page 27
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 39 from Processor 6
assignment 1: 2023/12/10 19:35:00 [Processor 0] -- Receive <<<Page>>> Page 9
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Invalidate>>> Page 27
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Write Forward>>> for Page 39 to Processor 6
assignment 1: 2023/12/10 19:35:00 [Processor 6] -- Receive <<<Page>>> Page 39
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 12 from Processor 7
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 5 from Processor 2
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 5 from Processor 5
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 16 from Processor 1
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Write Request>>> for Page 48 from Processor 8
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Invalidate>>> Page 39
assignment 1: 2023/12/10 19:35:00 [Processor 0] -- Receive <<<Write Forward>>> for Page 5 to Processor 2
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:35:00 [Processor 8] -- Receive <<<Invalidate>>> Page 5
assignment 1: 2023/12/10 19:35:00 [Processor 9] -- Receive <<<Write Forward>>> for Page 16 to Processor 1
assignment 1: 2023/12/10 19:35:00 [Processor 9] -- Receive <<<Invalidate>>> Page 48
assignment 1: 2023/12/10 19:35:00 [Processor 9] -- Receive <<<Invalidate>>> Page 5
assignment 1: 2023/12/10 19:35:00 [Processor 9] -- Receive <<<Invalidate>>> Page 16
assignment 1: 2023/12/10 19:35:00 [Processor 1] -- Receive <<<Write Forward>>> for Page 48 to Processor 8
assignment 1: 2023/12/10 19:35:00 [Processor 1] -- Receive <<<Invalidate>>> Page 39
assignment 1: 2023/12/10 19:35:00 [Processor 1] -- Receive <<<Page>>> Page 16
assignment 1: 2023/12/10 19:35:00 [Processor 8] -- Receive <<<Page>>> Page 48
assignment 1: 2023/12/10 19:35:00 [Processor 5] -- Receive <<<Invalidate>>> Page 5
assignment 1: 2023/12/10 19:35:00 [Processor 0] -- Receive <<<Invalidate>>> Page 48
assignment 1: 2023/12/10 19:35:00 [Processor 9] -- Read Page 1 from local page table
assignment 1: 2023/12/10 19:35:00 [Processor 1] -- Send <<<Read Request>>> for Page 14
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 14 from Processor 1
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Send <<<Read Request>>> for Page 3
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 3 from Processor 2
assignment 1: 2023/12/10 19:35:00 [Processor 4] -- Send <<<Read Request>>> for Page 33
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 33 from Processor 4
assignment 1: 2023/12/10 19:35:00 [Processor 6] -- Send <<<Read Request>>> for Page 22
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 22 from Processor 6
assignment 1: 2023/12/10 19:35:00 [Processor 3] -- Receive <<<Read Forward>>> for Page 22 to Processor 6
assignment 1: 2023/12/10 19:35:00 [Processor 0] -- Send <<<Read Request>>> for Page 19
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 19 from Processor 0
assignment 1: 2023/12/10 19:35:00 [Processor 4] -- Receive <<<Read Forward>>> for Page 19 to Processor 0
assignment 1: 2023/12/10 19:35:00 [Processor 0] -- Receive <<<Page>>> Page 19
assignment 1: 2023/12/10 19:35:00 [Processor 5] -- Send <<<Read Request>>> for Page 31
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 31 from Processor 5
assignment 1: 2023/12/10 19:35:00 [Processor 7] -- Read Page 32 from local page table
assignment 1: 2023/12/10 19:35:00 [Processor 3] -- Send <<<Read Request>>> for Page 25
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 25 from Processor 3
assignment 1: 2023/12/10 19:35:00 [Processor 1] -- Receive <<<Read Forward>>> for Page 25 to Processor 3
assignment 1: 2023/12/10 19:35:00 [Processor 8] -- Send <<<Read Request>>> for Page 42
assignment 1: 2023/12/10 19:35:00 [Central Manager] Receive <<<Read Request>>> for Page 42 from Processor 8
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Read Forward>>> for Page 14 to Processor 1
assignment 1: 2023/12/10 19:35:00 [Processor 6] -- Receive <<<Read Forward>>> for Page 3 to Processor 2
assignment 1: 2023/12/10 19:35:00 [Processor 6] -- Receive <<<Page>>> Page 22
assignment 1: 2023/12/10 19:35:00 [Processor 2] -- Receive <<<Page>>> Page 3
assignment 1: 2023/12/10 19:35:00 [Processor 8] -- Receive <<<Read Forward>>> for Page 33 to Processor 4
assignment 1: 2023/12/10 19:35:00 [Processor 8] -- Receive <<<Page>>> Page 42
assignment 1: 2023/12/10 19:35:00 [Processor 5] -- Receive <<<Page>>> Page 31
assignment 1: 2023/12/10 19:35:00 [Processor 3] -- Receive <<<Page>>> Page 25
assignment 1: 2023/12/10 19:35:00 [Processor 1] -- Receive <<<Page>>> Page 14
assignment 1: 2023/12/10 19:35:00 [Processor 4] -- Receive <<<Page>>> Page 33
assignment 1: 2023/12/10 19:35:01 [Processor 0] -- Read Page 7 from local page table
assignment 1: 2023/12/10 19:35:01 [Processor 2] -- Read Page 5 from local page table
assignment 1: 2023/12/10 19:35:01 [Processor 5] -- Send <<<Read Request>>> for Page 25
assignment 1: 2023/12/10 19:35:01 [Central Manager] Receive <<<Read Request>>> for Page 25 from Processor 5
assignment 1: 2023/12/10 19:35:01 [Processor 1] -- Receive <<<Read Forward>>> for Page 25 to Processor 5
assignment 1: 2023/12/10 19:35:01 [Processor 5] -- Receive <<<Page>>> Page 25
assignment 1: 2023/12/10 19:35:01 [Processor 9] -- Send <<<Read Request>>> for Page 36
assignment 1: 2023/12/10 19:35:01 [Central Manager] Receive <<<Read Request>>> for Page 36 from Processor 9
assignment 1: 2023/12/10 19:35:01 [Processor 4] -- Receive <<<Read Forward>>> for Page 36 to Processor 9
assignment 1: 2023/12/10 19:35:01 [Processor 9] -- Receive <<<Page>>> Page 36
assignment 1: 2023/12/10 19:35:01 [Processor 8] -- Read Page 43 from local page table
assignment 1: 2023/12/10 19:35:01 [Processor 1] -- Read Page 16 from local page table
assignment 1: 2023/12/10 19:35:01 [Processor 4] -- Read Page 36 from local page table
assignment 1: 2023/12/10 19:35:01 [Processor 7] -- Read Page 35 from local page table
assignment 1: 2023/12/10 19:35:01 [Processor 6] -- Send <<<Read Request>>> for Page 33
assignment 1: 2023/12/10 19:35:01 [Central Manager] Receive <<<Read Request>>> for Page 33 from Processor 6
assignment 1: 2023/12/10 19:35:01 [Processor 8] -- Receive <<<Read Forward>>> for Page 33 to Processor 6
assignment 1: 2023/12/10 19:35:01 [Processor 6] -- Receive <<<Page>>> Page 33
assignment 1: 2023/12/10 19:35:01 [Processor 3] -- Send <<<Read Request>>> for Page 26
assignment 1: 2023/12/10 19:35:01 [Central Manager] Receive <<<Read Request>>> for Page 26 from Processor 3
assignment 1: 2023/12/10 19:35:01 [Processor 3] -- Receive <<<Page>>> Page 26
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 12 from Processor 9
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 27 from Processor 3
assignment 1: 2023/12/10 19:35:02 [Processor 9] -- Receive <<<Write Forward>>> for Page 12 to Processor 7
assignment 1: 2023/12/10 19:35:02 [Processor 7] -- Receive <<<Page>>> Page 12
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 16 from Processor 1
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 5 from Processor 2
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 9 from Processor 0
assignment 1: 2023/12/10 19:35:02 [Processor 2] -- Receive <<<Write Forward>>> for Page 5 to Processor 5
assignment 1: 2023/12/10 19:35:02 [Processor 5] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 39 from Processor 6
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 48 from Processor 8
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Write Ack>>> for Page 19 from Processor 4
assignment 1: 2023/12/10 19:35:02 [Processor 1] -- Send <<<Read Request>>> for Page 48
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 48 from Processor 1
assignment 1: 2023/12/10 19:35:02 [Processor 8] -- Receive <<<Read Forward>>> for Page 48 to Processor 1
assignment 1: 2023/12/10 19:35:02 [Processor 1] -- Receive <<<Page>>> Page 48
assignment 1: 2023/12/10 19:35:02 [Processor 5] -- Send <<<Read Request>>> for Page 15
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 15 from Processor 5
assignment 1: 2023/12/10 19:35:02 [Processor 4] -- Receive <<<Read Forward>>> for Page 15 to Processor 5
assignment 1: 2023/12/10 19:35:02 [Processor 2] -- Send <<<Read Request>>> for Page 7
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 7 from Processor 2
assignment 1: 2023/12/10 19:35:02 [Processor 0] -- Send <<<Read Request>>> for Page 42
assignment 1: 2023/12/10 19:35:02 [Processor 0] -- Receive <<<Read Forward>>> for Page 7 to Processor 2
assignment 1: 2023/12/10 19:35:02 [Processor 2] -- Receive <<<Page>>> Page 7
assignment 1: 2023/12/10 19:35:02 [Processor 5] -- Receive <<<Page>>> Page 15
assignment 1: 2023/12/10 19:35:02 [Processor 4] -- Send <<<Read Request>>> for Page 37
assignment 1: 2023/12/10 19:35:02 [Processor 8] -- Send <<<Read Request>>> for Page 4
assignment 1: 2023/12/10 19:35:02 [Processor 6] -- Send <<<Read Request>>> for Page 18
assignment 1: 2023/12/10 19:35:02 [Processor 9] -- Read Page 36 from local page table
assignment 1: 2023/12/10 19:35:02 [Processor 7] -- Send <<<Read Request>>> for Page 39
assignment 1: 2023/12/10 19:35:02 [Processor 3] -- Send <<<Read Request>>> for Page 38
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 42 from Processor 0
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 37 from Processor 4
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 4 from Processor 8
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 18 from Processor 6
assignment 1: 2023/12/10 19:35:02 [Processor 8] -- Receive <<<Read Forward>>> for Page 42 to Processor 0
assignment 1: 2023/12/10 19:35:02 [Processor 0] -- Receive <<<Page>>> Page 42
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 7
assignment 1: 2023/12/10 19:35:02 [Processor 2] -- Receive <<<Read Forward>>> for Page 37 to Processor 4
assignment 1: 2023/12/10 19:35:02 [Processor 4] -- Receive <<<Page>>> Page 37
assignment 1: 2023/12/10 19:35:02 [Processor 7] -- Receive <<<Read Forward>>> for Page 4 to Processor 8
assignment 1: 2023/12/10 19:35:02 [Processor 8] -- Receive <<<Page>>> Page 4
assignment 1: 2023/12/10 19:35:02 [Processor 6] -- Receive <<<Page>>> Page 18
assignment 1: 2023/12/10 19:35:02 [Processor 6] -- Receive <<<Read Forward>>> for Page 39 to Processor 7
assignment 1: 2023/12/10 19:35:02 [Processor 7] -- Receive <<<Page>>> Page 39
assignment 1: 2023/12/10 19:35:02 [Central Manager] Receive <<<Read Request>>> for Page 38 from Processor 3
assignment 1: 2023/12/10 19:35:02 [Processor 3] -- Receive <<<Page>>> Page 38
assignment 1: 2023/12/10 19:35:03 [Processor 6] -- Send <<<Read Request>>> for Page 31
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 31 from Processor 6
assignment 1: 2023/12/10 19:35:03 [Processor 0] -- Send <<<Read Request>>> for Page 24
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 0
assignment 1: 2023/12/10 19:35:03 [Processor 5] -- Receive <<<Read Forward>>> for Page 31 to Processor 6
assignment 1: 2023/12/10 19:35:03 [Processor 1] -- Send <<<Read Request>>> for Page 32
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 32 from Processor 1
assignment 1: 2023/12/10 19:35:03 [Processor 5] -- Send <<<Read Request>>> for Page 48
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 48 from Processor 5
assignment 1: 2023/12/10 19:35:03 [Processor 8] -- Receive <<<Read Forward>>> for Page 48 to Processor 5
assignment 1: 2023/12/10 19:35:03 [Processor 9] -- Read Page 30 from local page table
assignment 1: 2023/12/10 19:35:03 [Processor 4] -- Send <<<Read Request>>> for Page 46
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 46 from Processor 4
assignment 1: 2023/12/10 19:35:03 [Processor 0] -- Receive <<<Read Forward>>> for Page 46 to Processor 4
assignment 1: 2023/12/10 19:35:03 [Processor 7] -- Send <<<Read Request>>> for Page 3
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 3 from Processor 7
assignment 1: 2023/12/10 19:35:03 [Processor 8] -- Send <<<Read Request>>> for Page 6
assignment 1: 2023/12/10 19:35:03 [Processor 3] -- Read Page 13 from local page table
assignment 1: 2023/12/10 19:35:03 [Processor 2] -- Send <<<Read Request>>> for Page 41
assignment 1: 2023/12/10 19:35:03 [Processor 4] -- Receive <<<Read Forward>>> for Page 24 to Processor 0
assignment 1: 2023/12/10 19:35:03 [Processor 4] -- Receive <<<Page>>> Page 46
assignment 1: 2023/12/10 19:35:03 [Processor 0] -- Receive <<<Page>>> Page 24
assignment 1: 2023/12/10 19:35:03 [Processor 6] -- Receive <<<Page>>> Page 31
assignment 1: 2023/12/10 19:35:03 [Processor 7] -- Receive <<<Read Forward>>> for Page 32 to Processor 1
assignment 1: 2023/12/10 19:35:03 [Processor 5] -- Receive <<<Page>>> Page 48
assignment 1: 2023/12/10 19:35:03 [Processor 6] -- Receive <<<Read Forward>>> for Page 3 to Processor 7
assignment 1: 2023/12/10 19:35:03 [Processor 7] -- Receive <<<Page>>> Page 3
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 6 from Processor 8
assignment 1: 2023/12/10 19:35:03 [Processor 1] -- Receive <<<Page>>> Page 32
assignment 1: 2023/12/10 19:35:03 [Processor 1] -- Receive <<<Read Forward>>> for Page 6 to Processor 8
assignment 1: 2023/12/10 19:35:03 [Processor 8] -- Receive <<<Page>>> Page 6
assignment 1: 2023/12/10 19:35:03 [Central Manager] Receive <<<Read Request>>> for Page 41 from Processor 2
assignment 1: 2023/12/10 19:35:03 [Processor 5] -- Receive <<<Read Forward>>> for Page 41 to Processor 2
assignment 1: 2023/12/10 19:35:03 [Processor 2] -- Receive <<<Page>>> Page 41
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Write Ack>>> for Page 5 from Processor 5
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Write Ack>>> for Page 12 from Processor 7
assignment 1: 2023/12/10 19:35:04 [Processor 3] -- Send <<<Read Request>>> for Page 10
assignment 1: 2023/12/10 19:35:04 [Processor 2] -- Send <<<Read Request>>> for Page 28
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 10 from Processor 3
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 28 from Processor 2
assignment 1: 2023/12/10 19:35:04 [Processor 7] -- Receive <<<Read Forward>>> for Page 28 to Processor 2
assignment 1: 2023/12/10 19:35:04 [Processor 2] -- Receive <<<Page>>> Page 28
assignment 1: 2023/12/10 19:35:04 [Processor 4] -- Send <<<Read Request>>> for Page 42
assignment 1: 2023/12/10 19:35:04 [Processor 8] -- Read Page 17 from local page table
assignment 1: 2023/12/10 19:35:04 [Processor 9] -- Send <<<Read Request>>> for Page 20
assignment 1: 2023/12/10 19:35:04 [Processor 0] -- Read Page 42 from local page table
assignment 1: 2023/12/10 19:35:04 [Processor 5] -- Send <<<Read Request>>> for Page 17
assignment 1: 2023/12/10 19:35:04 [Processor 5] -- Receive <<<Read Forward>>> for Page 10 to Processor 3
assignment 1: 2023/12/10 19:35:04 [Processor 3] -- Receive <<<Page>>> Page 10
assignment 1: 2023/12/10 19:35:04 [Processor 1] -- Send <<<Read Request>>> for Page 27
assignment 1: 2023/12/10 19:35:04 [Processor 6] -- Send <<<Read Request>>> for Page 35
assignment 1: 2023/12/10 19:35:04 [Processor 7] -- Send <<<Read Request>>> for Page 31
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 42 from Processor 4
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 20 from Processor 9
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 17 from Processor 5
assignment 1: 2023/12/10 19:35:04 [Processor 8] -- Receive <<<Read Forward>>> for Page 42 to Processor 4
assignment 1: 2023/12/10 19:35:04 [Processor 8] -- Receive <<<Read Forward>>> for Page 17 to Processor 5
assignment 1: 2023/12/10 19:35:04 [Processor 4] -- Receive <<<Page>>> Page 42
assignment 1: 2023/12/10 19:35:04 [Processor 2] -- Receive <<<Read Forward>>> for Page 20 to Processor 9
assignment 1: 2023/12/10 19:35:04 [Processor 9] -- Receive <<<Page>>> Page 20
assignment 1: 2023/12/10 19:35:04 [Processor 5] -- Receive <<<Page>>> Page 17
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 27 from Processor 1
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 35 from Processor 6
assignment 1: 2023/12/10 19:35:04 [Processor 3] -- Receive <<<Read Forward>>> for Page 27 to Processor 1
assignment 1: 2023/12/10 19:35:04 [Central Manager] Receive <<<Read Request>>> for Page 31 from Processor 7
assignment 1: 2023/12/10 19:35:04 [Processor 5] -- Receive <<<Read Forward>>> for Page 31 to Processor 7
assignment 1: 2023/12/10 19:35:04 [Processor 7] -- Receive <<<Read Forward>>> for Page 35 to Processor 6
assignment 1: 2023/12/10 19:35:04 [Processor 7] -- Receive <<<Page>>> Page 31
assignment 1: 2023/12/10 19:35:04 [Processor 1] -- Receive <<<Page>>> Page 27
assignment 1: 2023/12/10 19:35:04 [Processor 6] -- Receive <<<Page>>> Page 35
assignment 1: 2023/12/10 19:35:05 [Processor 1] -- Read Page 27 from local page table
assignment 1: 2023/12/10 19:35:05 [Processor 9] -- Send <<<Read Request>>> for Page 11
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 11 from Processor 9
assignment 1: 2023/12/10 19:35:05 [Processor 0] -- Receive <<<Read Forward>>> for Page 11 to Processor 9
assignment 1: 2023/12/10 19:35:05 [Processor 9] -- Receive <<<Page>>> Page 11
assignment 1: 2023/12/10 19:35:05 [Processor 6] -- Send <<<Read Request>>> for Page 11
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 11 from Processor 6
assignment 1: 2023/12/10 19:35:05 [Processor 0] -- Receive <<<Read Forward>>> for Page 11 to Processor 6
assignment 1: 2023/12/10 19:35:05 [Processor 6] -- Receive <<<Page>>> Page 11
assignment 1: 2023/12/10 19:35:05 [Processor 3] -- Send <<<Read Request>>> for Page 31
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 31 from Processor 3
assignment 1: 2023/12/10 19:35:05 [Processor 5] -- Receive <<<Read Forward>>> for Page 31 to Processor 3
assignment 1: 2023/12/10 19:35:05 [Processor 3] -- Receive <<<Page>>> Page 31
assignment 1: 2023/12/10 19:35:05 [Processor 2] -- Send <<<Read Request>>> for Page 33
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 33 from Processor 2
assignment 1: 2023/12/10 19:35:05 [Processor 8] -- Receive <<<Read Forward>>> for Page 33 to Processor 2
assignment 1: 2023/12/10 19:35:05 [Processor 2] -- Receive <<<Page>>> Page 33
assignment 1: 2023/12/10 19:35:05 [Processor 8] -- Send <<<Read Request>>> for Page 19
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 19 from Processor 8
assignment 1: 2023/12/10 19:35:05 [Processor 0] -- Send <<<Read Request>>> for Page 14
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 14 from Processor 0
assignment 1: 2023/12/10 19:35:05 [Processor 2] -- Receive <<<Read Forward>>> for Page 14 to Processor 0
assignment 1: 2023/12/10 19:35:05 [Processor 5] -- Send <<<Read Request>>> for Page 34
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 5
assignment 1: 2023/12/10 19:35:05 [Processor 7] -- Receive <<<Read Forward>>> for Page 34 to Processor 5
assignment 1: 2023/12/10 19:35:05 [Processor 7] -- Read Page 46 from local page table
assignment 1: 2023/12/10 19:35:05 [Processor 4] -- Send <<<Read Request>>> for Page 5
assignment 1: 2023/12/10 19:35:05 [Central Manager] Receive <<<Read Request>>> for Page 5 from Processor 4
assignment 1: 2023/12/10 19:35:05 [Processor 4] -- Receive <<<Read Forward>>> for Page 19 to Processor 8
assignment 1: 2023/12/10 19:35:05 [Processor 8] -- Receive <<<Page>>> Page 19
assignment 1: 2023/12/10 19:35:05 [Processor 0] -- Receive <<<Page>>> Page 14
assignment 1: 2023/12/10 19:35:05 [Processor 5] -- Receive <<<Page>>> Page 34
assignment 1: 2023/12/10 19:35:05 [Processor 5] -- Receive <<<Read Forward>>> for Page 5 to Processor 4
assignment 1: 2023/12/10 19:35:05 [Processor 4] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:35:06 [Processor 4] -- Read Page 19 from local page table
assignment 1: 2023/12/10 19:35:06 [Processor 1] -- Send <<<Read Request>>> for Page 42
assignment 1: 2023/12/10 19:35:06 [Processor 6] -- Send <<<Read Request>>> for Page 5
assignment 1: 2023/12/10 19:35:06 [Processor 3] -- Send <<<Read Request>>> for Page 45
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 42 from Processor 1
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 5 from Processor 6
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 45 from Processor 3
assignment 1: 2023/12/10 19:35:06 [Processor 9] -- Receive <<<Read Forward>>> for Page 45 to Processor 3
assignment 1: 2023/12/10 19:35:06 [Processor 3] -- Receive <<<Page>>> Page 45
assignment 1: 2023/12/10 19:35:06 [Processor 5] -- Receive <<<Read Forward>>> for Page 5 to Processor 6
assignment 1: 2023/12/10 19:35:06 [Processor 8] -- Receive <<<Read Forward>>> for Page 42 to Processor 1
assignment 1: 2023/12/10 19:35:06 [Processor 1] -- Receive <<<Page>>> Page 42
assignment 1: 2023/12/10 19:35:06 [Processor 0] -- Send <<<Read Request>>> for Page 4
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 4 from Processor 0
assignment 1: 2023/12/10 19:35:06 [Processor 7] -- Receive <<<Read Forward>>> for Page 4 to Processor 0
assignment 1: 2023/12/10 19:35:06 [Processor 0] -- Receive <<<Page>>> Page 4
assignment 1: 2023/12/10 19:35:06 [Processor 5] -- Send <<<Read Request>>> for Page 42
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 42 from Processor 5
assignment 1: 2023/12/10 19:35:06 [Processor 7] -- Send <<<Read Request>>> for Page 38
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 38 from Processor 7
assignment 1: 2023/12/10 19:35:06 [Processor 3] -- Receive <<<Read Forward>>> for Page 38 to Processor 7
assignment 1: 2023/12/10 19:35:06 [Processor 2] -- Send <<<Read Request>>> for Page 25
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 25 from Processor 2
assignment 1: 2023/12/10 19:35:06 [Processor 1] -- Receive <<<Read Forward>>> for Page 25 to Processor 2
assignment 1: 2023/12/10 19:35:06 [Processor 2] -- Receive <<<Page>>> Page 25
assignment 1: 2023/12/10 19:35:06 [Processor 8] -- Send <<<Read Request>>> for Page 26
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 26 from Processor 8
assignment 1: 2023/12/10 19:35:06 [Processor 3] -- Receive <<<Read Forward>>> for Page 26 to Processor 8
assignment 1: 2023/12/10 19:35:06 [Processor 9] -- Send <<<Read Request>>> for Page 29
assignment 1: 2023/12/10 19:35:06 [Central Manager] Receive <<<Read Request>>> for Page 29 from Processor 9
assignment 1: 2023/12/10 19:35:06 [Processor 9] -- Receive <<<Page>>> Page 29
assignment 1: 2023/12/10 19:35:06 [Processor 6] -- Receive <<<Page>>> Page 5
assignment 1: 2023/12/10 19:35:06 [Processor 8] -- Receive <<<Read Forward>>> for Page 42 to Processor 5
assignment 1: 2023/12/10 19:35:06 [Processor 8] -- Receive <<<Page>>> Page 26
assignment 1: 2023/12/10 19:35:06 [Processor 7] -- Receive <<<Page>>> Page 38
assignment 1: 2023/12/10 19:35:06 [Processor 5] -- Receive <<<Page>>> Page 42
assignment 1: 2023/12/10 19:35:07 [Processor 7] -- Send <<<Read Request>>> for Page 49
assignment 1: 2023/12/10 19:35:07 [Central Manager] Receive <<<Read Request>>> for Page 49 from Processor 7
assignment 1: 2023/12/10 19:35:07 [Processor 6] -- Receive <<<Read Forward>>> for Page 49 to Processor 7
assignment 1: 2023/12/10 19:35:07 [Processor 7] -- Receive <<<Page>>> Page 49
assignment 1: 2023/12/10 19:35:07 [Processor 3] -- Read Page 45 from local page table
assignment 1: 2023/12/10 19:35:07 [Processor 2] -- Send <<<Read Request>>> for Page 12
assignment 1: 2023/12/10 19:35:07 [Central Manager] Receive <<<Read Request>>> for Page 12 from Processor 2
assignment 1: 2023/12/10 19:35:07 [Processor 7] -- Receive <<<Read Forward>>> for Page 12 to Processor 2
assignment 1: 2023/12/10 19:35:07 [Processor 6] -- Send <<<Read Request>>> for Page 15
assignment 1: 2023/12/10 19:35:07 [Processor 8] -- Send <<<Read Request>>> for Page 36
assignment 1: 2023/12/10 19:35:07 [Processor 2] -- Receive <<<Page>>> Page 12
assignment 1: 2023/12/10 19:35:07 [Processor 1] -- Send <<<Read Request>>> for Page 44
assignment 1: 2023/12/10 19:35:07 [Processor 4] -- Read Page 43 from local page table
assignment 1: 2023/12/10 19:35:07 [Processor 9] -- Send <<<Read Request>>> for Page 23
assignment 1: 2023/12/10 19:35:07 [Processor 0] -- Read Page 40 from local page table
assignment 1: 2023/12/10 19:35:07 [Processor 5] -- Send <<<Read Request>>> for Page 40
assignment 1: 2023/12/10 19:35:07 [Central Manager] Receive <<<Read Request>>> for Page 15 from Processor 6
assignment 1: 2023/12/10 19:35:07 [Central Manager] Receive <<<Read Request>>> for Page 36 from Processor 8
assignment 1: 2023/12/10 19:35:07 [Processor 4] -- Receive <<<Read Forward>>> for Page 15 to Processor 6
assignment 1: 2023/12/10 19:35:07 [Processor 4] -- Receive <<<Read Forward>>> for Page 36 to Processor 8
assignment 1: 2023/12/10 19:35:07 [Processor 8] -- Receive <<<Page>>> Page 36
assignment 1: 2023/12/10 19:35:07 [Processor 6] -- Receive <<<Page>>> Page 15
assignment 1: 2023/12/10 19:35:07 [Central Manager] Receive <<<Read Request>>> for Page 44 from Processor 1
assignment 1: 2023/12/10 19:35:07 [Central Manager] Receive <<<Read Request>>> for Page 23 from Processor 9
assignment 1: 2023/12/10 19:35:07 [Processor 8] -- Receive <<<Read Forward>>> for Page 44 to Processor 1
assignment 1: 2023/12/10 19:35:07 [Processor 8] -- Receive <<<Read Forward>>> for Page 23 to Processor 9
assignment 1: 2023/12/10 19:35:07 [Processor 9] -- Receive <<<Page>>> Page 23
assignment 1: 2023/12/10 19:35:07 [Central Manager] Receive <<<Read Request>>> for Page 40 from Processor 5
assignment 1: 2023/12/10 19:35:07 [Processor 1] -- Receive <<<Page>>> Page 44
assignment 1: 2023/12/10 19:35:07 [Processor 0] -- Receive <<<Read Forward>>> for Page 40 to Processor 5
assignment 1: 2023/12/10 19:35:07 [Processor 5] -- Receive <<<Page>>> Page 40
assignment 1: 2023/12/10 19:35:08 [Processor 4] -- Read Page 19 from local page table
assignment 1: 2023/12/10 19:35:08 [Processor 2] -- Send <<<Read Request>>> for Page 19
assignment 1: 2023/12/10 19:35:08 [Central Manager] Receive <<<Read Request>>> for Page 19 from Processor 2
assignment 1: 2023/12/10 19:35:08 [Processor 4] -- Receive <<<Read Forward>>> for Page 19 to Processor 2
assignment 1: 2023/12/10 19:35:08 [Processor 2] -- Receive <<<Page>>> Page 19
assignment 1: 2023/12/10 19:35:08 [Processor 8] -- Send <<<Read Request>>> for Page 7
assignment 1: 2023/12/10 19:35:08 [Central Manager] Receive <<<Read Request>>> for Page 7 from Processor 8
assignment 1: 2023/12/10 19:35:08 [Processor 0] -- Receive <<<Read Forward>>> for Page 7 to Processor 8
assignment 1: 2023/12/10 19:35:08 [Processor 8] -- Receive <<<Page>>> Page 7
assignment 1: 2023/12/10 19:35:08 [Processor 0] -- Send <<<Read Request>>> for Page 39
assignment 1: 2023/12/10 19:35:08 [Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 0
assignment 1: 2023/12/10 19:35:08 [Processor 6] -- Receive <<<Read Forward>>> for Page 39 to Processor 0
assignment 1: 2023/12/10 19:35:08 [Processor 1] -- Send <<<Read Request>>> for Page 38
assignment 1: 2023/12/10 19:35:08 [Central Manager] Receive <<<Read Request>>> for Page 38 from Processor 1
assignment 1: 2023/12/10 19:35:08 [Processor 9] -- Send <<<Read Request>>> for Page 34
assignment 1: 2023/12/10 19:35:08 [Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 9
assignment 1: 2023/12/10 19:35:08 [Processor 7] -- Receive <<<Read Forward>>> for Page 34 to Processor 9
assignment 1: 2023/12/10 19:35:08 [Processor 5] -- Send <<<Read Request>>> for Page 36
assignment 1: 2023/12/10 19:35:08 [Central Manager] Receive <<<Read Request>>> for Page 36 from Processor 5
assignment 1: 2023/12/10 19:35:08 [Processor 7] -- Send <<<Read Request>>> for Page 43
assignment 1: 2023/12/10 19:35:09 [Central Manager] Receive <<<Read Request>>> for Page 43 from Processor 7
assignment 1: 2023/12/10 19:35:08 [Processor 3] -- Send <<<Read Request>>> for Page 48
assignment 1: 2023/12/10 19:35:08 [Processor 6] -- Send <<<Read Request>>> for Page 45
assignment 1: 2023/12/10 19:35:08 [Processor 3] -- Receive <<<Read Forward>>> for Page 38 to Processor 1
assignment 1: 2023/12/10 19:35:09 [Processor 1] -- Receive <<<Page>>> Page 38
assignment 1: 2023/12/10 19:35:08 [Processor 0] -- Receive <<<Page>>> Page 39
assignment 1: 2023/12/10 19:35:08 [Processor 9] -- Receive <<<Page>>> Page 34
assignment 1: 2023/12/10 19:35:09 [Processor 4] -- Receive <<<Read Forward>>> for Page 36 to Processor 5
assignment 1: 2023/12/10 19:35:09 [Processor 4] -- Receive <<<Read Forward>>> for Page 43 to Processor 7
assignment 1: 2023/12/10 19:35:09 [Processor 7] -- Receive <<<Page>>> Page 43
assignment 1: 2023/12/10 19:35:09 [Central Manager] Receive <<<Read Request>>> for Page 48 from Processor 3
assignment 1: 2023/12/10 19:35:09 [Processor 5] -- Receive <<<Page>>> Page 36
assignment 1: 2023/12/10 19:35:09 [Central Manager] Receive <<<Read Request>>> for Page 45 from Processor 6
assignment 1: 2023/12/10 19:35:09 [Processor 9] -- Receive <<<Read Forward>>> for Page 45 to Processor 6
assignment 1: 2023/12/10 19:35:09 [Processor 8] -- Receive <<<Read Forward>>> for Page 48 to Processor 3
assignment 1: 2023/12/10 19:35:09 [Processor 3] -- Receive <<<Page>>> Page 48
assignment 1: 2023/12/10 19:35:09 [Processor 6] -- Receive <<<Page>>> Page 45
```

Here shows an example of `assignment_2.log`.

```
assignment 2: 2023/12/10 20:20:31 [Primary Central Manager] Central Manager activated
assignment 2: 2023/12/10 20:20:31 [Backup Central Manager] Backup Central Manager activated
assignment 2: 2023/12/10 20:20:32 [Processor 6] -- Send <<<Read Request>>> for Page 4
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 4 from Processor 6
assignment 2: 2023/12/10 20:20:32 [Processor 2] -- Send <<<Read Request>>> for Page 10
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 10 from Processor 2
assignment 2: 2023/12/10 20:20:32 [Processor 2] -- Receive <<<Page>>> Page 10
assignment 2: 2023/12/10 20:20:32 [Processor 5] -- Send <<<Read Request>>> for Page 13
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 13 from Processor 5
assignment 2: 2023/12/10 20:20:32 [Processor 3] -- Send <<<Read Request>>> for Page 0
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 0 from Processor 3
assignment 2: 2023/12/10 20:20:32 [Processor 3] -- Receive <<<Page>>> Page 0
assignment 2: 2023/12/10 20:20:32 [Processor 9] -- Send <<<Read Request>>> for Page 22
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 22 from Processor 9
assignment 2: 2023/12/10 20:20:32 [Processor 0] -- Send <<<Read Request>>> for Page 30
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 30 from Processor 0
assignment 2: 2023/12/10 20:20:32 [Processor 0] -- Receive <<<Page>>> Page 30
assignment 2: 2023/12/10 20:20:32 [Processor 1] -- Send <<<Read Request>>> for Page 42
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 42 from Processor 1
assignment 2: 2023/12/10 20:20:32 [Processor 1] -- Receive <<<Page>>> Page 42
assignment 2: 2023/12/10 20:20:32 [Processor 4] -- Send <<<Read Request>>> for Page 20
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 20 from Processor 4
assignment 2: 2023/12/10 20:20:32 [Processor 7] -- Send <<<Read Request>>> for Page 35
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 35 from Processor 7
assignment 2: 2023/12/10 20:20:32 [Processor 7] -- Receive <<<Page>>> Page 35
assignment 2: 2023/12/10 20:20:32 [Processor 6] -- Receive <<<Page>>> Page 4
assignment 2: 2023/12/10 20:20:32 [Processor 8] -- Send <<<Read Request>>> for Page 3
assignment 2: 2023/12/10 20:20:32 [Primary Central Manager] Receive <<<Read Request>>> for Page 3 from Processor 8
assignment 2: 2023/12/10 20:20:32 [Processor 5] -- Receive <<<Page>>> Page 13
assignment 2: 2023/12/10 20:20:32 [Processor 9] -- Receive <<<Page>>> Page 22
assignment 2: 2023/12/10 20:20:32 [Processor 4] -- Receive <<<Page>>> Page 20
assignment 2: 2023/12/10 20:20:32 [Processor 8] -- Receive <<<Page>>> Page 3
assignment 2: 2023/12/10 20:20:33 [Processor 1] -- Send <<<Read Request>>> for Page 28
assignment 2: 2023/12/10 20:20:33 [Processor 4] -- Send <<<Read Request>>> for Page 10
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Send <<<Heartbeat>>> to Backup Central Manager
assignment 2: 2023/12/10 20:20:33 [Processor 6] -- Send <<<Read Request>>> for Page 34
assignment 2: 2023/12/10 20:20:33 [Processor 9] -- Send <<<Read Request>>> for Page 27
assignment 2: 2023/12/10 20:20:33 [Backup Central Manager] Receive <<<Heartbeat>>> from Primary Central Manager
assignment 2: 2023/12/10 20:20:33 [Processor 3] -- Send <<<Read Request>>> for Page 24
assignment 2: 2023/12/10 20:20:33 [Processor 8] -- Send <<<Read Request>>> for Page 31
assignment 2: 2023/12/10 20:20:33 [Processor 2] -- Send <<<Read Request>>> for Page 37
assignment 2: 2023/12/10 20:20:33 [Processor 5] -- Send <<<Read Request>>> for Page 30
assignment 2: 2023/12/10 20:20:33 [Processor 0] -- Send <<<Read Request>>> for Page 36
assignment 2: 2023/12/10 20:20:33 [Processor 7] -- Send <<<Read Request>>> for Page 39
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 28 from Processor 1
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 10 from Processor 4
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 6
assignment 2: 2023/12/10 20:20:33 [Processor 1] -- Receive <<<Page>>> Page 28
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 27 from Processor 9
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 3
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 31 from Processor 8
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 37 from Processor 2
assignment 2: 2023/12/10 20:20:33 [Processor 2] -- Receive <<<Read Forward>>> for Page 10 to Processor 4
assignment 2: 2023/12/10 20:20:33 [Processor 2] -- Receive <<<Page>>> Page 37
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 30 from Processor 5
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 36 from Processor 0
assignment 2: 2023/12/10 20:20:33 [Primary Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 7
assignment 2: 2023/12/10 20:20:33 [Processor 7] -- Receive <<<Page>>> Page 39
assignment 2: 2023/12/10 20:20:33 [Processor 4] -- Receive <<<Page>>> Page 10
assignment 2: 2023/12/10 20:20:33 [Processor 0] -- Receive <<<Page>>> Page 36
assignment 2: 2023/12/10 20:20:33 [Processor 9] -- Receive <<<Page>>> Page 27
assignment 2: 2023/12/10 20:20:33 [Processor 8] -- Receive <<<Page>>> Page 31
assignment 2: 2023/12/10 20:20:33 [Processor 0] -- Receive <<<Read Forward>>> for Page 30 to Processor 5
assignment 2: 2023/12/10 20:20:33 [Processor 6] -- Receive <<<Page>>> Page 34
assignment 2: 2023/12/10 20:20:33 [Processor 3] -- Receive <<<Page>>> Page 24
assignment 2: 2023/12/10 20:20:33 [Processor 5] -- Receive <<<Page>>> Page 30
assignment 2: 2023/12/10 20:20:34 [Processor 1] -- Send <<<Write Request>>> for Page 49
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 49 from Processor 1
assignment 2: 2023/12/10 20:20:34 [Processor 1] -- Receive <<<Page>>> Page 49
assignment 2: 2023/12/10 20:20:34 [Processor 1] -- Start <<<Writing>>> Page 49
assignment 2: 2023/12/10 20:20:34 [Processor 6] -- Send <<<Write Request>>> for Page 36
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 36 from Processor 6
assignment 2: 2023/12/10 20:20:34 [Processor 7] -- Send <<<Write Request>>> for Page 15
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 15 from Processor 7
assignment 2: 2023/12/10 20:20:34 [Processor 7] -- Receive <<<Page>>> Page 15
assignment 2: 2023/12/10 20:20:34 [Processor 7] -- Start <<<Writing>>> Page 15
assignment 2: 2023/12/10 20:20:34 [Processor 9] -- Send <<<Write Request>>> for Page 14
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 14 from Processor 9
assignment 2: 2023/12/10 20:20:34 [Processor 9] -- Receive <<<Page>>> Page 14
assignment 2: 2023/12/10 20:20:34 [Processor 9] -- Start <<<Writing>>> Page 14
assignment 2: 2023/12/10 20:20:34 [Processor 8] -- Send <<<Write Request>>> for Page 43
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 43 from Processor 8
assignment 2: 2023/12/10 20:20:34 [Processor 8] -- Receive <<<Page>>> Page 43
assignment 2: 2023/12/10 20:20:34 [Processor 8] -- Start <<<Writing>>> Page 43
assignment 2: 2023/12/10 20:20:34 [Processor 3] -- Send <<<Write Request>>> for Page 29
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 29 from Processor 3
assignment 2: 2023/12/10 20:20:34 [Processor 3] -- Receive <<<Page>>> Page 29
assignment 2: 2023/12/10 20:20:34 [Processor 3] -- Start <<<Writing>>> Page 29
assignment 2: 2023/12/10 20:20:34 [Processor 0] -- Receive <<<Write Forward>>> for Page 36 to Processor 6
assignment 2: 2023/12/10 20:20:34 [Processor 6] -- Receive <<<Page>>> Page 36
assignment 2: 2023/12/10 20:20:34 [Processor 6] -- Start <<<Writing>>> Page 36
assignment 2: 2023/12/10 20:20:34 [Processor 5] -- Send <<<Write Request>>> for Page 10
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 10 from Processor 5
assignment 2: 2023/12/10 20:20:34 [Processor 4] -- Receive <<<Invalidate>>> Page 10
assignment 2: 2023/12/10 20:20:34 [Processor 2] -- Send <<<Write Request>>> for Page 0
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 0 from Processor 2
assignment 2: 2023/12/10 20:20:34 [Processor 2] -- Receive <<<Write Forward>>> for Page 10 to Processor 5
assignment 2: 2023/12/10 20:20:34 [Processor 5] -- Receive <<<Page>>> Page 10
assignment 2: 2023/12/10 20:20:34 [Processor 5] -- Start <<<Writing>>> Page 10
assignment 2: 2023/12/10 20:20:34 [Processor 3] -- Receive <<<Write Forward>>> for Page 0 to Processor 2
assignment 2: 2023/12/10 20:20:34 [Processor 2] -- Receive <<<Page>>> Page 0
assignment 2: 2023/12/10 20:20:34 [Processor 2] -- Start <<<Writing>>> Page 0
assignment 2: 2023/12/10 20:20:34 [Processor 0] -- Send <<<Write Request>>> for Page 45
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 45 from Processor 0
assignment 2: 2023/12/10 20:20:34 [Processor 0] -- Receive <<<Page>>> Page 45
assignment 2: 2023/12/10 20:20:34 [Processor 0] -- Start <<<Writing>>> Page 45
assignment 2: 2023/12/10 20:20:34 [Processor 4] -- Send <<<Write Request>>> for Page 26
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Write Request>>> for Page 26 from Processor 4
assignment 2: 2023/12/10 20:20:34 [Processor 4] -- Receive <<<Page>>> Page 26
assignment 2: 2023/12/10 20:20:34 [Processor 4] -- Start <<<Writing>>> Page 26
assignment 2: 2023/12/10 20:20:34 [Processor 9] -- Read Page 14 from local page table
assignment 2: 2023/12/10 20:20:34 [Processor 1] -- Send <<<Read Request>>> for Page 10
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 10 from Processor 1
assignment 2: 2023/12/10 20:20:34 [Processor 6] -- Send <<<Read Request>>> for Page 24
assignment 2: 2023/12/10 20:20:34 [Processor 7] -- Send <<<Read Request>>> for Page 2
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 6
assignment 2: 2023/12/10 20:20:34 [Processor 4] -- Read Page 20 from local page table
assignment 2: 2023/12/10 20:20:34 [Processor 3] -- Receive <<<Read Forward>>> for Page 24 to Processor 6
assignment 2: 2023/12/10 20:20:34 [Processor 0] -- Send <<<Read Request>>> for Page 33
assignment 2: 2023/12/10 20:20:34 [Processor 8] -- Send <<<Read Request>>> for Page 0
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 2 from Processor 7
assignment 2: 2023/12/10 20:20:34 [Processor 3] -- Send <<<Read Request>>> for Page 7
assignment 2: 2023/12/10 20:20:34 [Processor 7] -- Receive <<<Page>>> Page 2
assignment 2: 2023/12/10 20:20:34 [Processor 5] -- Send <<<Read Request>>> for Page 11
assignment 2: 2023/12/10 20:20:34 [Processor 2] -- Send <<<Read Request>>> for Page 34
assignment 2: 2023/12/10 20:20:34 [Processor 6] -- Receive <<<Page>>> Page 24
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 33 from Processor 0
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 0 from Processor 8
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 7 from Processor 3
assignment 2: 2023/12/10 20:20:34 [Processor 0] -- Receive <<<Page>>> Page 33
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 11 from Processor 5
assignment 2: 2023/12/10 20:20:34 [Processor 3] -- Receive <<<Page>>> Page 7
assignment 2: 2023/12/10 20:20:34 [Primary Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 2
assignment 2: 2023/12/10 20:20:34 [Processor 5] -- Receive <<<Page>>> Page 11
assignment 2: 2023/12/10 20:20:34 [Processor 6] -- Receive <<<Read Forward>>> for Page 34 to Processor 2
assignment 2: 2023/12/10 20:20:34 [Processor 2] -- Receive <<<Page>>> Page 34
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 26 from Processor 4
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 29 from Processor 3
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 36 from Processor 6
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 43 from Processor 8
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 49 from Processor 1
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 15 from Processor 7
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 14 from Processor 9
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 45 from Processor 0
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 10 from Processor 5
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Write Ack>>> for Page 0 from Processor 2
assignment 2: 2023/12/10 20:20:35 [Processor 8] -- Send <<<Read Request>>> for Page 9
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 9 from Processor 8
assignment 2: 2023/12/10 20:20:35 [Processor 8] -- Receive <<<Page>>> Page 9
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Send <<<Heartbeat>>> to Backup Central Manager
assignment 2: 2023/12/10 20:20:35 [Backup Central Manager] Receive <<<Heartbeat>>> from Primary Central Manager
assignment 2: 2023/12/10 20:20:35 [Processor 5] -- Send <<<Read Request>>> for Page 44
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 44 from Processor 5
assignment 2: 2023/12/10 20:20:35 [Processor 5] -- Receive <<<Page>>> Page 44
assignment 2: 2023/12/10 20:20:35 [Processor 0] -- Send <<<Read Request>>> for Page 21
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 21 from Processor 0
assignment 2: 2023/12/10 20:20:35 [Processor 0] -- Receive <<<Page>>> Page 21
assignment 2: 2023/12/10 20:20:35 [Processor 4] -- Send <<<Read Request>>> for Page 37
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 37 from Processor 4
assignment 2: 2023/12/10 20:20:35 [Processor 2] -- Receive <<<Read Forward>>> for Page 37 to Processor 4
assignment 2: 2023/12/10 20:20:35 [Processor 4] -- Receive <<<Page>>> Page 37
assignment 2: 2023/12/10 20:20:35 [Processor 6] -- Send <<<Read Request>>> for Page 45
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 45 from Processor 6
assignment 2: 2023/12/10 20:20:35 [Processor 0] -- Receive <<<Read Forward>>> for Page 45 to Processor 6
assignment 2: 2023/12/10 20:20:35 [Processor 6] -- Receive <<<Page>>> Page 45
assignment 2: 2023/12/10 20:20:35 [Processor 1] -- Send <<<Read Request>>> for Page 38
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 38 from Processor 1
assignment 2: 2023/12/10 20:20:35 [Processor 1] -- Receive <<<Page>>> Page 38
assignment 2: 2023/12/10 20:20:35 [Processor 9] -- Send <<<Read Request>>> for Page 12
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 12 from Processor 9
assignment 2: 2023/12/10 20:20:35 [Processor 9] -- Receive <<<Page>>> Page 12
assignment 2: 2023/12/10 20:20:35 [Processor 3] -- Send <<<Read Request>>> for Page 31
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 31 from Processor 3
assignment 2: 2023/12/10 20:20:35 [Processor 8] -- Receive <<<Read Forward>>> for Page 31 to Processor 3
assignment 2: 2023/12/10 20:20:35 [Processor 3] -- Receive <<<Page>>> Page 31
assignment 2: 2023/12/10 20:20:35 [Processor 2] -- Send <<<Read Request>>> for Page 39
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 2
assignment 2: 2023/12/10 20:20:35 [Processor 7] -- Receive <<<Read Forward>>> for Page 39 to Processor 2
assignment 2: 2023/12/10 20:20:35 [Processor 2] -- Receive <<<Page>>> Page 39
assignment 2: 2023/12/10 20:20:35 [Processor 7] -- Send <<<Read Request>>> for Page 12
assignment 2: 2023/12/10 20:20:35 [Primary Central Manager] Receive <<<Read Request>>> for Page 12 from Processor 7
assignment 2: 2023/12/10 20:20:35 [Processor 9] -- Receive <<<Read Forward>>> for Page 12 to Processor 7
assignment 2: 2023/12/10 20:20:35 [Processor 7] -- Receive <<<Page>>> Page 12
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Primary is DOWN!
assignment 2: 2023/12/10 20:20:36 [Processor 7] -- Send <<<Read Request>>> for Page 13
assignment 2: 2023/12/10 20:20:36 [Processor 8] -- Send <<<Read Request>>> for Page 18
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 13 from Processor 7
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 18 from Processor 8
assignment 2: 2023/12/10 20:20:36 [Processor 8] -- Receive <<<Page>>> Page 18
assignment 2: 2023/12/10 20:20:36 [Processor 5] -- Send <<<Read Request>>> for Page 20
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 20 from Processor 5
assignment 2: 2023/12/10 20:20:36 [Processor 4] -- Receive <<<Read Forward>>> for Page 20 to Processor 5
assignment 2: 2023/12/10 20:20:36 [Processor 5] -- Receive <<<Page>>> Page 20
assignment 2: 2023/12/10 20:20:36 [Processor 5] -- Receive <<<Read Forward>>> for Page 13 to Processor 7
assignment 2: 2023/12/10 20:20:36 [Processor 2] -- Send <<<Read Request>>> for Page 4
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 4 from Processor 2
assignment 2: 2023/12/10 20:20:36 [Processor 1] -- Send <<<Read Request>>> for Page 39
assignment 2: 2023/12/10 20:20:36 [Processor 6] -- Receive <<<Read Forward>>> for Page 4 to Processor 2
assignment 2: 2023/12/10 20:20:36 [Processor 2] -- Receive <<<Page>>> Page 4
assignment 2: 2023/12/10 20:20:36 [Processor 3] -- Send <<<Read Request>>> for Page 44
assignment 2: 2023/12/10 20:20:36 [Processor 0] -- Send <<<Read Request>>> for Page 23
assignment 2: 2023/12/10 20:20:36 [Processor 4] -- Read Page 20 from local page table
assignment 2: 2023/12/10 20:20:36 [Processor 9] -- Send <<<Read Request>>> for Page 29
assignment 2: 2023/12/10 20:20:36 [Processor 6] -- Send <<<Read Request>>> for Page 47
assignment 2: 2023/12/10 20:20:36 [Processor 7] -- Receive <<<Page>>> Page 13
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 1
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 44 from Processor 3
assignment 2: 2023/12/10 20:20:36 [Processor 7] -- Receive <<<Read Forward>>> for Page 39 to Processor 1
assignment 2: 2023/12/10 20:20:36 [Processor 1] -- Receive <<<Page>>> Page 39
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 23 from Processor 0
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 29 from Processor 9
assignment 2: 2023/12/10 20:20:36 [Primary Central Manager] Receive <<<Read Request>>> for Page 47 from Processor 6
assignment 2: 2023/12/10 20:20:36 [Processor 0] -- Receive <<<Page>>> Page 23
assignment 2: 2023/12/10 20:20:36 [Processor 5] -- Receive <<<Read Forward>>> for Page 44 to Processor 3
assignment 2: 2023/12/10 20:20:36 [Processor 6] -- Receive <<<Page>>> Page 47
assignment 2: 2023/12/10 20:20:36 [Processor 3] -- Receive <<<Read Forward>>> for Page 29 to Processor 9
assignment 2: 2023/12/10 20:20:36 [Processor 3] -- Receive <<<Page>>> Page 44
assignment 2: 2023/12/10 20:20:36 [Processor 9] -- Receive <<<Page>>> Page 29
assignment 2: 2023/12/10 20:20:37 [Processor 2] -- Send <<<Write Request>>> for Page 44
assignment 2: 2023/12/10 20:20:37 [Processor 4] -- Send <<<Write Request>>> for Page 36
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Send <<<Write Request>>> for Page 7
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 44 from Processor 2
assignment 2: 2023/12/10 20:20:37 [Processor 1] -- Send <<<Write Request>>> for Page 43
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Receive <<<Invalidate>>> Page 44
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Send <<<Write Request>>> for Page 30
assignment 2: 2023/12/10 20:20:37 [Processor 7] -- Send <<<Write Request>>> for Page 13
assignment 2: 2023/12/10 20:20:37 [Processor 9] -- Send <<<Write Request>>> for Page 7
assignment 2: 2023/12/10 20:20:37 [Processor 6] -- Send <<<Write Request>>> for Page 36
assignment 2: 2023/12/10 20:20:37 [Processor 8] -- Send <<<Write Request>>> for Page 0
assignment 2: 2023/12/10 20:20:37 [Processor 0] -- Send <<<Write Request>>> for Page 35
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 36 from Processor 4
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 7 from Processor 5
assignment 2: 2023/12/10 20:20:37 [Processor 6] -- Receive <<<Write Forward>>> for Page 36 to Processor 4
assignment 2: 2023/12/10 20:20:37 [Processor 4] -- Receive <<<Page>>> Page 36
assignment 2: 2023/12/10 20:20:37 [Processor 4] -- Start <<<Writing>>> Page 36
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Receive <<<Write Forward>>> for Page 44 to Processor 2
assignment 2: 2023/12/10 20:20:37 [Processor 2] -- Receive <<<Page>>> Page 44
assignment 2: 2023/12/10 20:20:37 [Processor 2] -- Start <<<Writing>>> Page 44
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Receive <<<Write Forward>>> for Page 7 to Processor 5
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 43 from Processor 1
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Receive <<<Page>>> Page 7
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Start <<<Writing>>> Page 7
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 30 from Processor 3
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 13 from Processor 7
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 7 from Processor 9
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Receive <<<Write Forward>>> for Page 13 to Processor 7
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Receive <<<Invalidate>>> Page 30
assignment 2: 2023/12/10 20:20:37 [Processor 0] -- Receive <<<Write Forward>>> for Page 30 to Processor 3
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Receive <<<Page>>> Page 30
assignment 2: 2023/12/10 20:20:37 [Processor 8] -- Receive <<<Write Forward>>> for Page 43 to Processor 1
assignment 2: 2023/12/10 20:20:37 [Processor 1] -- Receive <<<Page>>> Page 43
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 36 from Processor 6
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 0 from Processor 8
assignment 2: 2023/12/10 20:20:37 [Processor 7] -- Receive <<<Invalidate>>> Page 13
assignment 2: 2023/12/10 20:20:37 [Processor 7] -- Receive <<<Page>>> Page 13
assignment 2: 2023/12/10 20:20:37 [Processor 7] -- Start <<<Writing>>> Page 13
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Start <<<Writing>>> Page 30
assignment 2: 2023/12/10 20:20:37 [Processor 1] -- Start <<<Writing>>> Page 43
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Write Request>>> for Page 35 from Processor 0
assignment 2: 2023/12/10 20:20:37 [Processor 2] -- Receive <<<Write Forward>>> for Page 0 to Processor 8
assignment 2: 2023/12/10 20:20:37 [Processor 8] -- Receive <<<Page>>> Page 0
assignment 2: 2023/12/10 20:20:37 [Processor 8] -- Start <<<Writing>>> Page 0
assignment 2: 2023/12/10 20:20:37 [Processor 7] -- Receive <<<Write Forward>>> for Page 35 to Processor 0
assignment 2: 2023/12/10 20:20:37 [Processor 0] -- Receive <<<Page>>> Page 35
assignment 2: 2023/12/10 20:20:37 [Processor 0] -- Start <<<Writing>>> Page 35
assignment 2: 2023/12/10 20:20:37 [Processor 7] -- Send <<<Read Request>>> for Page 28
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 28 from Processor 7
assignment 2: 2023/12/10 20:20:37 [Processor 1] -- Receive <<<Read Forward>>> for Page 28 to Processor 7
assignment 2: 2023/12/10 20:20:37 [Processor 7] -- Receive <<<Page>>> Page 28
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Send <<<Read Request>>> for Page 3
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 3 from Processor 3
assignment 2: 2023/12/10 20:20:37 [Processor 8] -- Receive <<<Read Forward>>> for Page 3 to Processor 3
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Receive <<<Page>>> Page 3
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Send <<<Read Request>>> for Page 27
assignment 2: 2023/12/10 20:20:37 [Processor 0] -- Send <<<Read Request>>> for Page 8
assignment 2: 2023/12/10 20:20:37 [Processor 1] -- Send <<<Read Request>>> for Page 24
assignment 2: 2023/12/10 20:20:37 [Processor 8] -- Send <<<Read Request>>> for Page 6
assignment 2: 2023/12/10 20:20:37 [Processor 4] -- Send <<<Read Request>>> for Page 10
assignment 2: 2023/12/10 20:20:37 [Processor 9] -- Read Page 27 from local page table
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 27 from Processor 5
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 8 from Processor 0
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 1
assignment 2: 2023/12/10 20:20:37 [Processor 0] -- Receive <<<Page>>> Page 8
assignment 2: 2023/12/10 20:20:37 [Processor 9] -- Receive <<<Read Forward>>> for Page 27 to Processor 5
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Receive <<<Page>>> Page 27
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 6 from Processor 8
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 10 from Processor 4
assignment 2: 2023/12/10 20:20:37 [Processor 5] -- Receive <<<Read Forward>>> for Page 10 to Processor 4
assignment 2: 2023/12/10 20:20:37 [Processor 4] -- Receive <<<Page>>> Page 10
assignment 2: 2023/12/10 20:20:37 [Processor 2] -- Send <<<Read Request>>> for Page 24
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 2
assignment 2: 2023/12/10 20:20:37 [Processor 8] -- Receive <<<Page>>> Page 6
assignment 2: 2023/12/10 20:20:37 [Processor 6] -- Send <<<Read Request>>> for Page 14
assignment 2: 2023/12/10 20:20:37 [Primary Central Manager] Receive <<<Read Request>>> for Page 14 from Processor 6
assignment 2: 2023/12/10 20:20:37 [Processor 9] -- Receive <<<Read Forward>>> for Page 14 to Processor 6
assignment 2: 2023/12/10 20:20:37 [Processor 6] -- Receive <<<Page>>> Page 14
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Receive <<<Read Forward>>> for Page 24 to Processor 1
assignment 2: 2023/12/10 20:20:37 [Processor 3] -- Receive <<<Read Forward>>> for Page 24 to Processor 2
assignment 2: 2023/12/10 20:20:37 [Processor 2] -- Receive <<<Page>>> Page 24
assignment 2: 2023/12/10 20:20:37 [Processor 1] -- Receive <<<Page>>> Page 24
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 44 from Processor 2
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 36 from Processor 4
assignment 2: 2023/12/10 20:20:38 [Processor 4] -- Receive <<<Write Forward>>> for Page 36 to Processor 6
assignment 2: 2023/12/10 20:20:38 [Processor 6] -- Receive <<<Page>>> Page 36
assignment 2: 2023/12/10 20:20:38 [Processor 6] -- Start <<<Writing>>> Page 36
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 13 from Processor 7
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 30 from Processor 3
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 0 from Processor 8
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 35 from Processor 0
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 7 from Processor 5
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Write Ack>>> for Page 43 from Processor 1
assignment 2: 2023/12/10 20:20:38 [Processor 5] -- Receive <<<Write Forward>>> for Page 7 to Processor 9
assignment 2: 2023/12/10 20:20:38 [Processor 9] -- Receive <<<Page>>> Page 7
assignment 2: 2023/12/10 20:20:38 [Processor 9] -- Start <<<Writing>>> Page 7
assignment 2: 2023/12/10 20:20:38 [Processor 9] -- Send <<<Read Request>>> for Page 11
assignment 2: 2023/12/10 20:20:38 [Processor 6] -- Send <<<Read Request>>> for Page 8
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 11 from Processor 9
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 8 from Processor 6
assignment 2: 2023/12/10 20:20:38 [Processor 3] -- Send <<<Read Request>>> for Page 49
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 49 from Processor 3
assignment 2: 2023/12/10 20:20:38 [Processor 1] -- Receive <<<Read Forward>>> for Page 49 to Processor 3
assignment 2: 2023/12/10 20:20:38 [Processor 1] -- Send <<<Read Request>>> for Page 1
assignment 2: 2023/12/10 20:20:38 [Processor 8] -- Send <<<Read Request>>> for Page 25
assignment 2: 2023/12/10 20:20:38 [Processor 2] -- Send <<<Read Request>>> for Page 8
assignment 2: 2023/12/10 20:20:38 [Processor 0] -- Send <<<Read Request>>> for Page 43
assignment 2: 2023/12/10 20:20:38 [Processor 4] -- Send <<<Read Request>>> for Page 36
assignment 2: 2023/12/10 20:20:38 [Processor 5] -- Send <<<Read Request>>> for Page 29
assignment 2: 2023/12/10 20:20:38 [Processor 7] -- Send <<<Read Request>>> for Page 38
assignment 2: 2023/12/10 20:20:38 [Processor 0] -- Receive <<<Read Forward>>> for Page 8 to Processor 6
assignment 2: 2023/12/10 20:20:38 [Processor 6] -- Receive <<<Page>>> Page 8
assignment 2: 2023/12/10 20:20:38 [Processor 5] -- Receive <<<Read Forward>>> for Page 11 to Processor 9
assignment 2: 2023/12/10 20:20:38 [Processor 9] -- Receive <<<Page>>> Page 11
assignment 2: 2023/12/10 20:20:38 [Processor 3] -- Receive <<<Page>>> Page 49
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 1 from Processor 1
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 25 from Processor 8
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 8 from Processor 2
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 43 from Processor 0
assignment 2: 2023/12/10 20:20:38 [Processor 0] -- Receive <<<Read Forward>>> for Page 8 to Processor 2
assignment 2: 2023/12/10 20:20:38 [Processor 2] -- Receive <<<Page>>> Page 8
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 36 from Processor 4
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 29 from Processor 5
assignment 2: 2023/12/10 20:20:38 [Primary Central Manager] Receive <<<Read Request>>> for Page 38 from Processor 7
assignment 2: 2023/12/10 20:20:38 [Processor 1] -- Receive <<<Page>>> Page 1
assignment 2: 2023/12/10 20:20:38 [Processor 1] -- Receive <<<Read Forward>>> for Page 43 to Processor 0
assignment 2: 2023/12/10 20:20:38 [Processor 1] -- Receive <<<Read Forward>>> for Page 38 to Processor 7
assignment 2: 2023/12/10 20:20:38 [Processor 7] -- Receive <<<Page>>> Page 38
assignment 2: 2023/12/10 20:20:38 [Processor 3] -- Receive <<<Read Forward>>> for Page 29 to Processor 5
assignment 2: 2023/12/10 20:20:38 [Processor 5] -- Receive <<<Page>>> Page 29
assignment 2: 2023/12/10 20:20:38 [Processor 8] -- Receive <<<Page>>> Page 25
assignment 2: 2023/12/10 20:20:38 [Processor 0] -- Receive <<<Page>>> Page 43
assignment 2: 2023/12/10 20:20:39 [Primary Central Manager] Receive <<<Write Ack>>> for Page 36 from Processor 6
assignment 2: 2023/12/10 20:20:39 [Primary Central Manager] Receive <<<Write Ack>>> for Page 7 from Processor 9
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Primary Central Manager is down, promote Backup Central Manager to Primary Central Manager
assignment 2: 2023/12/10 20:20:39 [Processor 9] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 4] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 0] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 1] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 2] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 3] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 6] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 5] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 7] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 8] -- Receive <<<Primary Down>>>
assignment 2: 2023/12/10 20:20:39 [Processor 5] -- Read Page 29 from local page table
assignment 2: 2023/12/10 20:20:39 [Processor 0] -- Send <<<Read Request>>> for Page 34
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 0
assignment 2: 2023/12/10 20:20:39 [Processor 6] -- Send <<<Read Request>>> for Page 39
assignment 2: 2023/12/10 20:20:39 [Processor 7] -- Send <<<Read Request>>> for Page 24
assignment 2: 2023/12/10 20:20:39 [Processor 2] -- Send <<<Read Request>>> for Page 32
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 39 from Processor 6
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 24 from Processor 7
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 32 from Processor 2
assignment 2: 2023/12/10 20:20:39 [Processor 7] -- Receive <<<Read Forward>>> for Page 39 to Processor 6
assignment 2: 2023/12/10 20:20:39 [Processor 8] -- Send <<<Read Request>>> for Page 44
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 44 from Processor 8
assignment 2: 2023/12/10 20:20:39 [Processor 8] -- Receive <<<Page>>> Page 44
assignment 2: 2023/12/10 20:20:39 [Processor 4] -- Send <<<Read Request>>> for Page 28
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 28 from Processor 4
assignment 2: 2023/12/10 20:20:39 [Processor 3] -- Send <<<Read Request>>> for Page 40
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 40 from Processor 3
assignment 2: 2023/12/10 20:20:39 [Processor 6] -- Receive <<<Read Forward>>> for Page 34 to Processor 0
assignment 2: 2023/12/10 20:20:39 [Processor 6] -- Receive <<<Page>>> Page 39
assignment 2: 2023/12/10 20:20:39 [Processor 9] -- Send <<<Read Request>>> for Page 17
assignment 2: 2023/12/10 20:20:39 [Processor 1] -- Send <<<Read Request>>> for Page 45
assignment 2: 2023/12/10 20:20:39 [Processor 2] -- Receive <<<Page>>> Page 32
assignment 2: 2023/12/10 20:20:39 [Processor 3] -- Receive <<<Read Forward>>> for Page 24 to Processor 7
assignment 2: 2023/12/10 20:20:39 [Processor 3] -- Receive <<<Page>>> Page 40
assignment 2: 2023/12/10 20:20:39 [Processor 1] -- Receive <<<Read Forward>>> for Page 28 to Processor 4
assignment 2: 2023/12/10 20:20:39 [Processor 4] -- Receive <<<Page>>> Page 28
assignment 2: 2023/12/10 20:20:39 [Processor 0] -- Receive <<<Page>>> Page 34
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 17 from Processor 9
assignment 2: 2023/12/10 20:20:39 [Backup Central Manager] Receive <<<Read Request>>> for Page 45 from Processor 1
assignment 2: 2023/12/10 20:20:39 [Processor 0] -- Receive <<<Read Forward>>> for Page 45 to Processor 1
assignment 2: 2023/12/10 20:20:39 [Processor 7] -- Receive <<<Page>>> Page 24
assignment 2: 2023/12/10 20:20:39 [Processor 1] -- Receive <<<Page>>> Page 45
assignment 2: 2023/12/10 20:20:39 [Processor 9] -- Receive <<<Page>>> Page 17
assignment 2: 2023/12/10 20:20:40 [Processor 5] -- Send <<<Write Request>>> for Page 43
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 43 from Processor 5
assignment 2: 2023/12/10 20:20:40 [Processor 3] -- Send <<<Write Request>>> for Page 5
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 5 from Processor 3
assignment 2: 2023/12/10 20:20:40 [Processor 3] -- Receive <<<Page>>> Page 5
assignment 2: 2023/12/10 20:20:40 [Processor 3] -- Start <<<Writing>>> Page 5
assignment 2: 2023/12/10 20:20:40 [Processor 0] -- Send <<<Write Request>>> for Page 7
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 7 from Processor 0
assignment 2: 2023/12/10 20:20:40 [Processor 7] -- Send <<<Write Request>>> for Page 15
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 15 from Processor 7
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Receive <<<Write Forward>>> for Page 43 to Processor 5
assignment 2: 2023/12/10 20:20:40 [Processor 3] -- Receive <<<Write Forward>>> for Page 7 to Processor 0
assignment 2: 2023/12/10 20:20:40 [Processor 7] -- Start <<<Writing>>> Page 15
assignment 2: 2023/12/10 20:20:40 [Processor 1] -- Send <<<Write Request>>> for Page 28
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 28 from Processor 1
assignment 2: 2023/12/10 20:20:40 [Processor 4] -- Receive <<<Invalidate>>> Page 28
assignment 2: 2023/12/10 20:20:40 [Processor 2] -- Send <<<Write Request>>> for Page 3
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 3 from Processor 2
assignment 2: 2023/12/10 20:20:40 [Processor 4] -- Send <<<Write Request>>> for Page 18
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 18 from Processor 4
assignment 2: 2023/12/10 20:20:40 [Processor 4] -- Receive <<<Page>>> Page 18
assignment 2: 2023/12/10 20:20:40 [Processor 4] -- Start <<<Writing>>> Page 18
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Send <<<Write Request>>> for Page 14
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 14 from Processor 8
assignment 2: 2023/12/10 20:20:40 [Processor 1] -- Start <<<Writing>>> Page 28
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Receive <<<Write Forward>>> for Page 3 to Processor 2
assignment 2: 2023/12/10 20:20:40 [Processor 2] -- Receive <<<Page>>> Page 3
assignment 2: 2023/12/10 20:20:40 [Processor 2] -- Start <<<Writing>>> Page 3
assignment 2: 2023/12/10 20:20:40 [Processor 9] -- Receive <<<Write Forward>>> for Page 14 to Processor 8
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Receive <<<Page>>> Page 14
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Start <<<Writing>>> Page 14
assignment 2: 2023/12/10 20:20:40 [Processor 6] -- Send <<<Write Request>>> for Page 31
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 31 from Processor 6
assignment 2: 2023/12/10 20:20:40 [Processor 9] -- Send <<<Write Request>>> for Page 46
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Write Request>>> for Page 46 from Processor 9
assignment 2: 2023/12/10 20:20:40 [Processor 9] -- Receive <<<Page>>> Page 46
assignment 2: 2023/12/10 20:20:40 [Processor 9] -- Start <<<Writing>>> Page 46
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Receive <<<Write Forward>>> for Page 31 to Processor 6
assignment 2: 2023/12/10 20:20:40 [Processor 6] -- Receive <<<Page>>> Page 31
assignment 2: 2023/12/10 20:20:40 [Processor 6] -- Start <<<Writing>>> Page 31
assignment 2: 2023/12/10 20:20:40 [Processor 1] -- Send <<<Read Request>>> for Page 32
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 32 from Processor 1
assignment 2: 2023/12/10 20:20:40 [Processor 2] -- Receive <<<Read Forward>>> for Page 32 to Processor 1
assignment 2: 2023/12/10 20:20:40 [Processor 1] -- Receive <<<Page>>> Page 32
assignment 2: 2023/12/10 20:20:40 [Processor 4] -- Send <<<Read Request>>> for Page 38
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 38 from Processor 4
assignment 2: 2023/12/10 20:20:40 [Processor 4] -- Receive <<<Page>>> Page 38
assignment 2: 2023/12/10 20:20:40 [Processor 5] -- Send <<<Read Request>>> for Page 34
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 5
assignment 2: 2023/12/10 20:20:40 [Processor 6] -- Receive <<<Read Forward>>> for Page 34 to Processor 5
assignment 2: 2023/12/10 20:20:40 [Processor 5] -- Receive <<<Page>>> Page 34
assignment 2: 2023/12/10 20:20:40 [Processor 6] -- Read Page 39 from local page table
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Send <<<Read Request>>> for Page 34
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 34 from Processor 8
assignment 2: 2023/12/10 20:20:40 [Processor 6] -- Receive <<<Read Forward>>> for Page 34 to Processor 8
assignment 2: 2023/12/10 20:20:40 [Processor 8] -- Receive <<<Page>>> Page 34
assignment 2: 2023/12/10 20:20:40 [Processor 7] -- Send <<<Read Request>>> for Page 33
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 33 from Processor 7
assignment 2: 2023/12/10 20:20:40 [Processor 0] -- Receive <<<Read Forward>>> for Page 33 to Processor 7
assignment 2: 2023/12/10 20:20:40 [Processor 7] -- Receive <<<Page>>> Page 33
assignment 2: 2023/12/10 20:20:40 [Processor 9] -- Send <<<Read Request>>> for Page 40
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 40 from Processor 9
assignment 2: 2023/12/10 20:20:40 [Processor 3] -- Receive <<<Read Forward>>> for Page 40 to Processor 9
assignment 2: 2023/12/10 20:20:40 [Processor 9] -- Receive <<<Page>>> Page 40
assignment 2: 2023/12/10 20:20:40 [Processor 0] -- Send <<<Read Request>>> for Page 27
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 27 from Processor 0
assignment 2: 2023/12/10 20:20:40 [Processor 9] -- Receive <<<Read Forward>>> for Page 27 to Processor 0
assignment 2: 2023/12/10 20:20:40 [Processor 0] -- Receive <<<Page>>> Page 27
assignment 2: 2023/12/10 20:20:40 [Processor 2] -- Send <<<Read Request>>> for Page 29
assignment 2: 2023/12/10 20:20:40 [Backup Central Manager] Receive <<<Read Request>>> for Page 29 from Processor 2
assignment 2: 2023/12/10 20:20:40 [Processor 3] -- Receive <<<Read Forward>>> for Page 29 to Processor 2
```