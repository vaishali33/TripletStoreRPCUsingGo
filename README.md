# TripletStoreRPCUsingGo
Implementing a Simple Client Server With Go and JSON-RPC

Introduction:
Designed and developed a GO client-server application that implements JSON Remote Procedure calls for operations such as lookup, insert, insertOrUpdate, delete, listKeys, listIDs on DICT3 data structure made up of key, relation, value pair, where each key and relation maps to unique value which is any valid JSON object, and key and relations are strings. Server responds to all properly­structured JSON­RPC messages for above operations. Client reads a JSON­RPC request message from the standard input and makes the appropriate request to the server, and shows the response (if any). Both your client and server code takes as their 1st command line argument the filename that contains the client/server configuration JSON object, file name is config.json.
Design:
 
Go Client server application
Implementation:
Structure of DICT3:
Key and relation are strings and value is any valid JSON object so its stored as map[string]interface{}. As key and relation form unique combination to retrieve value, I have used struct DKey as shown in below struct list to use as key for DICT3. 
So DICT3 is of type map[DKey]map[string]interface{}

User defined structs used at client and server side are as follows:
type DKey struct{ KeyA, RelA string }
type Request struct{ KeyRel DKey	Val map[string]interface{}	Id int }
type Response struct{ Tripair Triplet 	    Done bool 	ID int	 Err error }
type Triplet struct{ Key, Rel string     Val map[string]interface{} }
type ListResponse struct{ List interface{}     Id int     Err error }
type DIC3 struct{}

Below methods implemented:
Method name	Client side implementation	Server side implementation	Description
Lookup	Lookup(key, relation)	(t *DIC3) Lookup(req *Request, reply *Response) error	Returns value for given key and relation pair from DICT3
Insert	Insert(key, relation, value)	(t *DIC3) Insert(triplet *Request, reply *Response) error 	Inserts a key, relation, value triplet in DICT3
InsertOrUpdate	InsertOrUpdate(key, relation, value)	(t *DIC3) InsertOrUpdate(triplet *Request, reply *Response) error	Inserts or updates a key, relation, value triplet in DICT3 based on key and relation pair
Delete	Delete(key, relation)	(t *DIC3) Delete(req *Request, reply *Response) error	Deletes key, relation, value triplet based on key and relation pair
Listkeys	Listkeys()	(t *DIC3) Listkeys(req *Request, reply *ListResponse) error	Returns set of keys in DICT3
ListIDs	ListIDs()	(t *DIC3) ListIDs(req *Request, reply *ListResponse) error	Returns set of key, relation pairs in DICT3

Client reads instructions to execute by calling server from commands.txt file, path to commands.txt is provided to client as command line arguments.
Both client and server reads config.json for configuration, path to this file is sent as first command line argument to both client and server.
Components of server application:
How to run? 
-	Execute start.sh script for server.
Config.json : reads config.json for configuration such as port number, IP Address, persistent storage location etc.
Dict3.txt : server persists all triplets into this file before shutdown using persistDict3() method, and when server starts it loads it into memory using loadDict3() method.
JsonDict3Server.go :
Consists of methods exposed to clients given in above table and other helper methods such as loadDict3(), loadConfig(), persistDict3(), checkError(),checkIfPresent() etc. Main() method is used to start server. Also has all above user defined structs. Server uses many go packages such as log, os, strings, io/ioutil, net/rpc/jsonrpc, bufio, fmt, encoding/json etc.

Components of client application:
How to run? 
-	Execute start.sh script for server and start_client.sh for client.
Config.json file: reads config.json for configuration such as port number, IP Address, persistent storage location etc.
commands.txt file: this takes instruction to be executed using client to server such as lookup(keyA, relA), listKeys() etc.
JsonDict3Client.go :
Consists of main method and other helper methods such as loadConfig(), checkError () etc. Also has all above user defined structs. Client uses many go packages such as log, os, strings, io/ioutil, net/rpc/jsonrpc, bufio, fmt etc.

Server Console :
Reading  config.json
Methods exposed by server:  [lookup insert insertOrUpdate delete listKeys listIDs shutdown]
map[{key1  rel3}:map[v1:9 v2:[12 aabc 14] v3:3] {key1 rel1}:map[v1:901 v2:[12 3 14] v3:3] {key2 rel1}:map[t1:201 t2:[a 56 c] t3:29] {key3 rel1}:map[v1:901 v2:[12 3 14] v3:3] {key3 rel2}:map[v1:91 v2:[aa 3 14] v3:2] {key1 rel2}:map[v1:901 v2:[12 aa 14] v3:3] {key2 rel2}:map[t1:201 t2:[a 56 c] t3:20] {key1 rel3}:map[v1:9 v2:[12 aabc 14] v3:3] {key2 rel3}:map[t1:21 t2:[a 5 c] t3:200] {key2  rel2}:map[t1:201 t2:[a 56 c] t3:20]]
Server started........
in Lookup :  key1 rel1 96
in Lookup :  keyA relA 97
in Insert :  key1 rel3 map[v1:9 v2:[12 aabc 14] v3:3] 0
in InsertOrUpdate :  key2 rel2 map[t1:201 t2:[a 56 c] t3:20] 0
Updating.....
in Delete :  5 key1 rel1
after delete DICT3  map[{key1  rel3}:map[v1:9 v2:[12 aabc 14] v3:3] {key2 rel1}:map[t1:201 t2:[a 56 c] t3:29] {key3 rel1}:map[v1:901 v2:[12 3 14] v3:3] {key2 rel2}:map[t1:201 t2:[a 56 c] t3:20] {key3 rel2}:map[v1:91 v2:[aa 3 14] v3:2] {key1 rel2}:map[v1:901 v2:[12 aa 14] v3:3] {key1 rel3}:map[v1:9 v2:[12 aabc 14] v3:3] {key2 rel3}:map[t1:21 t2:[a 5 c] t3:200] {key2  rel2}:map[t1:201 t2:[a 56 c] t3:20]]
in ListKeys:: Request ID :  691
in ListIDs:: Request ID :  699
shuting server down!!!!!!
 
Client Console:
Reading  config.json
Methods available to client:  [lookup insert insertOrUpdate delete listKeys listIDs shutdown]
lookup::  {{key1 rel1 map[v1:901 v2:[12 3 14] v3:3]} true 96 <nil>}
lookup:: No such triplet for key relation :  {keyA relA}
insert successful::  true
insert or update executed
delete executed
list keys::  {[key1 key2 key3] 691 <nil>}
list ids::  {[[ key1 ,  rel3 ] [ key2 , rel1 ] [ key3 , rel1 ] [ key2 , rel2 ] [ key3 , rel2 ] [ key1 , rel2 ] [ key1 , rel3 ] [ key2 , rel3 ] [ key2 ,  rel2 ]] 699 <nil>}
No such method named -  xx
shutdown executed..
