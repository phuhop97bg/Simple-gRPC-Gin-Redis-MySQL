# Simple gRPC project
**My gRPC project**<br />

**Service A:** 

-Receiving RESTful API request from Client<br />
-Handling request and call RPC function to Service B<br />
-Receiving RPC response from Service B and handling<br />
-Return Response to Client<br />


**Service B:**

-Receiving RPC request from Service A and handling<br />
-Storing and retrieval data to Redis and MySQL ( Redis as cache, MySQL as Database)<br />
-Return RPC response to Service A<br />

  
