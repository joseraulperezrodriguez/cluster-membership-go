# Membership protocol and services

This project provides an implementation for offering cluster membership services with the aim of being useful and extensible in any case that need answer to the question: *who are the peers in the cluster*?
</br></br> 
Two main components are part of this system:
- [A gossip like protocol](#a-gossip-like-protocol)
- [Membership services](#membership-services)

## A gossip like protocol

Used for message propagation to the cluster and membership management. The protocol works as follows.

Every 3 seconds by default (configurable), each node picks one node at random to exchange messages with it, (uses Netty library for this communications between nodes). There are six types of messages, and each one belong to one of three categories:  

- ADD-TO-CLUSTER (RUMOR). A message containing a node to add to the cluster.
<br/><br/> 
- SUSPECT-DEAD (RUMOR). A message indicating to node "T" that node "D" can't be reached. If after some configurable period of time the node "T" don't receive a KEEP-ALIVE message (regarding "D"), the node "D" is removed from list, and is out of the cluster. This task is performed in every node and at the same time, because of time synchronization, and message priority, details later.
<br/><br/>
- KEEP-ALIVE (RUMOR). A message containing a node "D" and meaning it is alive, this avoid node "D" is removed from cluster.
<br/><br/>
- REMOVE-FROM-CLUSTER (RUMOR). A message to remove node "D" from cluster, the received will remove it from its node list.
<br/><br/>
- PROBE (PROBE). A message sent to node "T" with the only purpose of test if node "D" can be reached by the sender node "S". When the node "D" is the same as "T" the probe is called 'direct', in other case is an 'indirect' probe and has a different treatment in the protocol.
<br/><br/>
- UNSUBSCRIPTION (CLUSTER). A message send to node "T" meaning that node "D" wants to leave the cluster, after node "T" receives the message it creates a "REMOVE-FROM-CLUSTER" message and propagate it to other nodes in the cluster.
  
The message categories describes better each message type:

- PROBE, a message with the only purpose of test the other node connectivity.
<br/><br/>
- CLUSTER, is a message at cluster level, that means the sender just wants "some" node of the cluster receive the message, then it is processed by the received node taking specific actions. In the case of UNSUBSCRIPTION, the action implemented by the receiver is create another message type (REMOVE-FROM-CLUSTER) and spread it across the cluster.
<br/><br/>
- RUMOR, is a message that needs to be received to every node in the cluster, we achieve this using the gossip behaviour and continuous synchronization feature, explained later.

The cluster state are conformed by some data structures, described next.

- **nodes**, an ordered list (by node id) containing the nodes in the cluster.
<br/><br/>
- **suspectingNodesTimeout**, a treeset like data structure holding <Node, Time> tuples to store the nodes in suspecting state, and the time each nodes will expire if not KEEP-ALIVE message is received, sorted by Time ascending.
<br/><br/>
- **failed**, a treeset like data structure holding <Node, Time> tuples to store the failing nodes, and the fail time for each node. When the node selected for message exchange on each iteration can't be reached by the sender, is inserted in the "failed" list of the sender node. Then, on next iterations a node is polled from this set and send an 'indirect' probe to that node, if that probe also fails, the node is inserted in "suspectingNodesTimeout" and a SUSPECT-DEAD rumor is spread to the cluster so to every node is aware of such failure.
<br/><br/>
- **rumorsToSend**, a treeset like data structure with the message of type RUMOR pending to send to other nodes, every iteration send the messages stored in this set
<br/><br/>
- **receivedRumors**, a tree set like data structure to store a configurable number of RUMOR messages received (by default 10000), this is to make possible synchronization between nodes based on historic messages, the message are ordered by the generated time, the time the message was created by the creator node.

#### How to ensure a RUMOR message reach all the nodes in the cluster?

Basically in two ways

First, we have the concept of iterations, each message object, has an iteration field that stores how many times that message needs to propagate to other nodes until reach all nodes. The process is as follows, when a node creates a new message, lets say ADD-TO-CLUSTER it calculates the number of iterations "I" for it, using the formula Log2(N*2) + 1 where N is the total nodes in the cluster, then when the message is received by a node, it decrease "I" by one, because the fact that a node receives a message means that an iteration have been completed. To better illustrate how it works, image a cluster with 7 nodes.

1- The node "A" creates a RUMOR message "M" with "I"=4, when node "B" receive the message, 2 from 7 nodes knows the message, and remains 3 iterations to complete the propagation. "B" insert the message into 'rumorsToSend', because has 3 iterations remaining, and also in 'receivedRumors' because is the first time it receive that message and need to store it for synchronization purposed mentioned before.

2- Now, nodes "A" and "B" knows the message, in the next iterations both "A" and "B" picks randomly a node, lets say "A" picks "D" and "B" picks "C". A this time both nodes has the message M in the 'rumorsToSend' list, they read it and propagate M to "C" and "D", now "A","B","C" and "D" knows the message and the "I" field of message M is 3 in all nodes.

3- In this way, after "I" iterations  the message reach every node with high probability, depending of the random function implementation. Also is possible a node "E" receive the same RUMOR from different nodes, but that don't change the state of the node, because the RUMOR handling is idempotent.


The other way the system uses to RUMOR exchange is using continuous synchronization, this intend to work as a fall back of the previous process, and used time frames data to make queries to 'receivedRumors' set on a remote node to check both are synchronized. On every iteration, the sender node "A" sends to "D" a simple metadata object based on his 'receivedRumors' data. The metadata consist of a time frame; 'start', 'end' in milliseconds and a 'counter' value, meaning that in the time frame [start to end] it has received 'counter' rumors. The receiver node "D" makes a query to his own 'receivedRumors' set, in the same range and if the number is higher than the "A" one, it sends back to "A" as a response all the RUMORs received in that period, then is the job of node "A" to found the missing RUMOR and process it.

The time frame for this synchronization is carefully selected on every iteration, the main restriction here is that the frame must be a period that don't affects any current message based on propagation formula, to avoid synchronization attempts be too early. For example, for a 7 nodes cluster and an iteration interval of 3 seconds, all messages that started propagation more than 12 seconds ago (because number of iterations = 4, from previous case) is supposed have been reached all nodes, so a good time frame to synchronize could be:

start = (current time - 16 second)
end   = start + 4 seconds = (current time - 12 seconds)

because the message started in this frame has stopped already the propagation and not change is supposed to happen. Besides that, the frame is about 4 seconds, that is a period in which not a high number of RUMOR messages comes to the system, and is less resource consuming to query and process a short period than a larger one.

Implementation details in the source code documentation.

## Membership services

A REST API for consumer to get info about the cluster

It is implemented using Spring Boot REST API libraries. It has the following end points:

- **/membership/size**, simply get the size of the cluster
</br></br>
- **/membership/unsubscribe**, start an unsubscribe process in the node receiving the message, that starts generating a UNSUBSCRIBE message, and as is a CLUSTER level message type category, it needs to confirm that "some" one in the cluster received the message, after that, the node starts the shutdown.
</br></br>
- **/membership/nodes**, retrieve the list of nodes object that are in the cluster
</br></br>
- **/membership/subscribe**, it is the subscribe end point, when a node comes to live, the first call is to this end point, to get cluster state copy 


Debugging end points

- **/membership/debug/state-info**, is for describing the current state, useful for debugging mode
</br></br>
- **/membership/debug/shutdown**, is for failure simulation in debugging mode
</br></br>
- **/membership/debug/pause**, is for delaying simulation in debugging mode
