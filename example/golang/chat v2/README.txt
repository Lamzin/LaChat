http://raycompstuff.blogspot.com/2009/12/simpler-chat-server-and-client-in.html


SERVER:

Server has three part which are running as goroutines and communicate via channels.
1) handlingINOUT() simple wait for input of clientreceiver() and send to all clientsender() which are in the list.
2) clientreceiver() wait for his data from client via networkconnection and send it to a inputchannel to handlingINOUT
3) clientsender() wait for data from channel and send it to client

every client connection get a his own clientreceiver/sender and a list entry. on disconnection the list entry will be deleted.



CLIENT:

client start two goroutines: for getting input from stdin and send it to network. and get from network and print it out.
