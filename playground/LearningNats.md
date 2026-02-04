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