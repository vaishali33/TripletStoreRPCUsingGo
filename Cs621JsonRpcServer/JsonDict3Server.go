package main 
//concurrency map - https://blog.golang.org/go-maps-in-action
import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"encoding/json"
	"os"
	"net"
	"bufio"
    "log"
    "strings"
    "io/ioutil"
)

type DKey struct{
	KeyA, RelA string
}
type Request struct{
	KeyRel DKey
	Val map[string]interface{}
	Id int
}
type Response struct{
	Tripair Triplet
	Done bool
	ID int
	Err error
}

type Triplet struct{
	Key, Rel string
	Val map[string]interface{}
}

type ListResponse struct{
	List interface{}
	Id int
	Err error
}

type DIC3 int

var dict3 map[DKey]map[string]interface{}
var configMap map[string]interface{}
var protocol string
var ipAdd string
var dict3File string
var methods []interface{}
var port string
var confLoc string


//server shutdown 
func (t *DIC3) Shutdown(dKey *DKey, reply *int) error{
	fmt.Println("shuting server down!!!!!!")
	*reply = 9
	persistDict3()
	os.Exit(0)
	return nil
}
//Lookup returns value stored for given key and relation
func (t *DIC3) Lookup(req *Request, reply *Response) error {
	fmt.Println("in Lookup : ", req.KeyRel.KeyA, req.KeyRel.RelA, req.Id)
	val := dict3[req.KeyRel]
	
	reply.ID = req.Id
	reply.Tripair.Key = req.KeyRel.KeyA
	reply.Tripair.Rel = req.KeyRel.RelA
	reply.Tripair.Val = val
	reply.Done = true
	reply.Err = nil
	return nil
}
//Insert a given triplet in DICT3
func (t *DIC3) Insert(triplet *Request, reply *Response) error {
	fmt.Println("in Insert : ", triplet.KeyRel.KeyA, triplet.KeyRel.RelA, triplet.Val, triplet.Id)
	
	dict3[DKey{triplet.KeyRel.KeyA, triplet.KeyRel.RelA}] = triplet.Val
	_, ok := dict3[DKey{triplet.KeyRel.KeyA, triplet.KeyRel.RelA}]
	reply.ID = triplet.Id
	reply.Done = ok
	reply.Err = nil
	return nil
}

//InsertOrUpdate given triplet in DICT3
func (t *DIC3) InsertOrUpdate(triplet *Request, reply *Response) error {
	fmt.Println("in InsertOrUpdate : ", triplet.KeyRel.KeyA, triplet.KeyRel.RelA, triplet.Val, triplet.Id)
	keyRel := DKey{triplet.KeyRel.KeyA, triplet.KeyRel.RelA}
	_, ok := dict3[keyRel]
	if !ok {
		//Insert
		fmt.Println("Inserting.....")
		dict3[keyRel] = triplet.Val
	}else{
		//Update
		fmt.Println("Updating.....")
		delete(dict3, keyRel)
		dict3[DKey{triplet.KeyRel.KeyA, triplet.KeyRel.RelA}] = triplet.Val
	}
	reply.ID = triplet.Id
	reply.Done = ok
	reply.Err = nil
	return nil
}

//Delete from DICT3
func (t *DIC3) Delete(req *Request, reply *Response) error {
	fmt.Println("in Delete : ", req.Id, req.KeyRel.KeyA, req.KeyRel.RelA)
	delete(dict3, req.KeyRel)
	fmt.Println("after delete DICT3 ", dict3)
	reply.ID = req.Id
	reply.Done = true
	reply.Err = nil
	return nil
}

//list keys in DICT3
func (t *DIC3) Listkeys(req *Request, reply *ListResponse) error {
	fmt.Println("in ListKeys:: Request ID : ", req.Id)
	keys := make([]string, 0, len(dict3))
    for k := range dict3 {
    	found := checkIfPresent(k.KeyA, keys)
    	if !found{
        	keys = append(keys, k.KeyA)
        	}
    }
	reply.Id = req.Id
	reply.List = keys
	reply.Err = nil
	return nil
}

//list key-relation pairs in DICT3
func (t *DIC3) ListIDs(req *Request, reply *ListResponse) error {
	fmt.Println("in ListIDs:: Request ID : ", req.Id)
	keys := make([]string, 0, len(dict3))
    for k := range dict3 {
        keys = append(keys, "[",k.KeyA,",", k.RelA, "]")
    }
	reply.Id = req.Id
	reply.List = keys
	reply.Err = nil
	return nil
}

func main() {
	
	loadConfig()
	
	loadDict3()
	
	dic3 := new(DIC3)
	rpc.Register(dic3)
	
	tcpAddr, err := net.ResolveTCPAddr(protocol, port)
	checkError(err)
	fmt.Println("Server started........")
	listener, err := net.ListenTCP(protocol, tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}

}
//load config in configMap
func loadConfig() error{
	configMap = make(map[string]interface{})
	fmt.Println("Reading ", os.Args[1])
	dat, err := ioutil.ReadFile(os.Args[1])
    checkError(err)
    //fmt.Print(dat)
	
	if err := json.Unmarshal(dat, &configMap); err != nil {
        log.Fatal("Error in loading config ", err)
    }
	protocol = configMap["protocol"].(string)
	ipAdd = configMap["ipAddress"].(string)
	port = configMap["port"].(string)
	persiStorage := configMap["persistentStorageContainer"].(map[string]interface{})
	dict3File = persiStorage["file"].(string)
	methods = configMap["methods"].([]interface{})
	fmt.Println("Methods exposed by server: ", methods)
	return nil
}
//load DICT3 in memory from persistent storage
func loadDict3() error{
	dict3 = make(map[DKey]map[string]interface{})
	
	file, err := os.Open(dict3File)
	if err != nil {
    	log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
    	arr := strings.Split(scanner.Text(), "=")
    	if len(arr) == 3{
    		key_rel := DKey{arr[0], arr[1]}
    		b :=[]byte(arr[2])
			var f map[string]interface{}
			err := json.Unmarshal(b, &f)
			if err != nil{
				log.Fatal(err)
			}
			dict3[key_rel] = f
		}
	}
	fmt.Println(dict3)
	if err := scanner.Err(); err != nil {
    	log.Fatal(err)
	}
	return nil
}

func persistDict3()error{
	 // For more granular writes, open a file for writing.
    f, err := os.Create(dict3File)
    checkError(err)
    defer f.Close()
	
	for k, v := range dict3{
		b, err := json.Marshal(v)
		val := string(b[:])
		s := []string{k.KeyA, "=", k.RelA,"=", val,"\n"}
		_, err = f.WriteString(strings.Join(s, ""))
		checkError(err)
		f.Sync()
	}
	return nil
}
func checkIfPresent(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
