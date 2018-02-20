# Partial-Sum
 A GoLang multithreaded applications to compute the sum of integers stored in a file. 
 
## Problem Statement:
Your program will take two command-line input parameters: M and fname. Here M is an integer and fname is the pathname (relative or absolute) to the input data file.
1. The format of the input data file is: a sequence of integers separated by white space, written in ASCII decimal text notation; you may assume that each integer is the ASCII representation of a 64-bit signed integer (and has ~20 digits).1. 
2. Your program should spawn M workers threads and one coordinator thread. The main() of your program simply spawns the coordinator.
3. The workers and the coordinator are to be implemented as goroutines in the GoLang. Workers can communicate with the coordinator but do not communicate among themselves. Use GoLang channels for communication between the workers and the coordinators.
4. The coordinator partitions the data file in M equal-size contiguous fragments; each kth fragment will be given to the kth worker via a JSON message of the form that includes the datafile's filename and the start and end byte position of the kth fragment, eg "{datafile: fname, start: pos1 , end: pos2}"
5. Each worker upon receiving its assignment via a JSON message it computes a partial sum which is the sum of all the integers which are fully contained within its assigned fragment. It also identifies a prefix or suffix of its fragment that could be part of integers split between adjacent (neighboring) fragments. Upon completion, the worker communicates its response to the coordinator via a JSON message with the partial sum, suffix, and prefix of its fragment, eg: "{value: ddd, prefix: '1224 ', suffix: ' 678'}".
6. The coordinator, upon receving a response from each worker, accumulates all the workers's partial sums, as well as the sums of the couple of integers in the concatenation of kth suffix and (k+1)th prefix (received by the the workers assigned the kth and (k+1)th fragments respectively.

Your implementation should be correct and efficient by ensuring that
- The coordinator, upon completion, computes the correct sum of the integers in the datafile.
there are no race conditions.
- Has least amount of information (#bytes) communicated between the coordinator and the workers.
least synchronization delay (elapsed time of syncrhonization between workers and coordinator).
- Has least latency time for the coordinator (time difference between the last response received and the coordinator finishing).
- Has least response time for the coordinator (time difference between the coordinator starting and the receiving its 1st response from a worker).
- Has least elapsed time for the "slowest" worker (the time from when a worker starts and when it finishes).
- Has least elapsed time for the coordinator (the time from when the coordinator starts and finishes).

## Implementation:
For above problem, I created a method and a function. There are certain edge cases for computing the sum from partial sum where we need to deal with suffix, preffix and error( basically a chunk of integer which doesn't have neither suffix nor prefix). All of the above are handled using a continue_ flag which helps to understand if we need to return the chunk or concatenate them.

Method: Worker; which gets details about the chunk it needs to compute while taking care of prefix, suffix and exception cases.
- input: waitGroup, channel of jobs and results
- receives JSON via jobs channel
- sends JSON via results channel

Function: Co-ordinator; Its responsibility is to generate chunks for worker, assign it to specific worker via channels, receive the partial sum from the workers and do the sum of all partial sums
- inputs: Filename and Number of Chunks
- returns total sum

## References :
1. https://www.csee.umbc.edu/~kalpakis/Courses/621-sp18/homeworks/hw1c.php

