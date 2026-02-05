This is just my own notes for how I am working through learning it, and debugging issues When
I do go wrong.

Before I begin, I needed to install NATS server and the NATS CLI on Docker Network.

I am currently having issues with accessing nats within that docker network from my host matchine, so I will be running the commands from within a container on that network. 

The docker network does not seem to be the issue. 

I had used the wrong command to run the NATS server.
'''bash
"-js -store_dir /data"
'''
This seems to be the correct command to run NATS server, and a data directory.

After the Container is running, I can exec into it, and try out the Various patterns.
'''bash
docker-compose exec toolbox sh
'''
It gives me access to a shell inside the containe, to access the nats CLI which is able to access the ports.

##Pattern 1: Publish and Subscribe
It is a decoupling pattern where it separates the message sender from the receiver, currently in the Project, the API is being called and returned towards the same port and are working togethe in sync. This pattern allows for asynchronous communication in the sense, that the publishers(source of messages, or API calls) are able to send messages without waiting for immediate responses, which allows for more flexibility and scalability in the system. The subscribers(services that process the messages) receive the messages, and manage them independently, allowing for better resource utilization and responsiveness.

Use Cases:
- Event-Driven Architectures: In systems where events trigger actions, such as user registrations or order placements, publishers can send event notifications to subscribers that handle these events asynchronously.
- Distributed Logging: Applications send log messages to a central service which collects and distributes them as required to systems. 


Now, How can we do this in the CLI?(I will get to the code implementation later)
First, I need to start a subscriber that listens to a specific subject(topic). In NATS, subjects are used to categorize messages. Let's say we want to listen to the subject "updates".
'''bash
nats sub updates
'''
This command starts a subscriber that listens for messages on the "updates" subject.
Next, in another terminal (or another session within the same container), I can publish a message to the "updates" subject.\
'''bash
nats pub updates "Hello, NATS!"
'''
When I run this command, the message "Hello, NATS!" is sent to the "updates" subject, and the subscriber that I started earlier should receive and display this message.


How I can use this in the current project:
Currently, This project is receving and listening on the same port synchronously, Which makes it tight for coupling, and not something that is good for a Routing sysem which can have even 10,000 drivers asking for a route. Having more ports is also not a good solution, as it can be hard to manage as well as expensive. We can change the project pattern to instead have a central NATS Server which can listen for the requests from the API, and then which can be utilized by the current worker pools to process the requests, which then get sent to a temporary inbox, which can then be given back to the Client through websockets or any other method(I dont know how to do this yet, but I will learn it later).


##Pattern 2: Queue Groups
Standard PubSub pattern generally leads to multiple workers receving the same message, which is not ideal for a worker pool. In JetStream, there is an inbuilt tool which is Queue Groups. It functions as a load balancer for messages, Messages are distributed among its members, ensuring that each message is processed by only one member of the group. 

Overall Picture: We can utilize this for the project in the sense that. NATS automatically handles the distribution of messages among the workers in the pool, ensuring that each request is processed by only one worker, which can help to improve the efficiency and scalability of the system. If we have a large number of requests, we can just add more workers to the queue group.

'''bash
nats sub updates --queue group1
'''

'''bash
nats sub updates --queue group2
'''

'''bash
nats pub orders "Job 1"

nats pub orders "Job 2"

nats pub orders "Job 3"

nats pub orders "Job 4"
'''
When I run the publish command, the message "Hello, NATS!" will be received by only one of the subscribers in either group1 or group2, demonstrating the load balancing effect of queue groups.

(All this will be executed in the container, code patterns will be shown and added later)