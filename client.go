package main

/*
 * Client - implement Command line and send key==>value to server
 */


import (
	"fmt"
	"net/http"
	"bufio"
	"os"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"sort"
)


type InputArgs struct {
	Key  int
	Value   string
}

type ReplyMessage struct {
	Message string
}

type Data struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

func main() {

    var rec InputArgs
	var n int
	stdin := bufio.NewReader(os.Stdin)
	n = 0
	initConsistentHash();
	for n!=3 {

		switch n {
			case 1:{
				fmt.Println("PUT (format: Key Value)")
				_,err := fmt.Fscanln(stdin,&rec.Key)
				if err != nil {
					fmt.Println("Format Error: Do not input the character");
					stdin.ReadString('\n')
					break;
				}
				_,err =fmt.Fscanln(stdin,&rec.Value)
				if err != nil {
					fmt.Println("Format Error: Do not input the space");
					stdin.ReadString('\n')
					break;
				}
				fmt.Println("value"+rec.Value);
				mes,err:= PUTdata(rec)
				if err != nil {
					fmt.Println(err);
					break
				}
				fmt.Println("\n***************************************")
				fmt.Println(mes.Message)
				fmt.Println("***************************************")
				break;
			}
			case 2:{
				fmt.Println("GET (format: key)")
				_,err := fmt.Fscanln(stdin,&rec.Key)
				if err != nil {
					fmt.Println("Format Error: Do not input the character");
					stdin.ReadString('\n')
					break;
				}
				mes,err:= GET(rec.Key)
				if err != nil {
					fmt.Println(err);
					break
				}
				fmt.Println("\n***************************************")
				fmt.Println(mes.Message)
				fmt.Println("***************************************")
				break;
			}
			case 3:{
				fmt.Println("Remove (format: serverPort(3000,3001,3002))")
				_,err := fmt.Fscanln(stdin,&rec.Key)
				if err != nil {
					fmt.Println("Format Error: Do not input the character");
					stdin.ReadString('\n')
					break;
				}
				mes,err:= GET(rec.Key)
				if err != nil {
					fmt.Println(err);
					break
				}
				fmt.Println("\n***************************************")
				fmt.Println(mes.Message)
				fmt.Println("***************************************")
				break;
			}
			default:{

			}
		}

		fmt.Println("\nPlease input the number of function that you want to choose:" )
		fmt.Println("1. PUT key and value ")
		fmt.Println("2. GET value by key ")
		fmt.Println("3. Remove one serverNode ")
		fmt.Println("4. Exit ")
		n=0;
		rec.Key = 0;
		rec.Value=""
		_,err := fmt.Fscanln(stdin,&n)
		if (err != nil) || (n>3) || (n<1) {
			fmt.Println("Input Error: please enter 1, 2, 3");
			stdin.ReadString('\n')
			n = 0;
		}
	}
}

var serverNode = []int{3000,3001,3002}
var circle = make(map[int]string)
var circleKey []int

func PUTdata(args InputArgs) (ReplyMessage, error) {

	var reply ReplyMessage

	serverPort := GetServerNode(args.Key)

	urlPath :=  "http://127.0.0.1:"+serverPort+"/keys/"+strconv.Itoa(args.Key)+"/"+args.Value

	client := &http.Client{}
	req, err := http.NewRequest("PUT", urlPath, nil)
	if err != nil {
		fmt.Println(err)
		reply.Message = "call : Http.NewRequest error" 
		return reply, err
	}
	req.Header.Add("Content-Type", "application/json")
	
	res, err := client.Do(req)
	if err!=nil {
		fmt.Println("call : client.Do error",err)
		reply.Message = "call : client.Do error" 
		return reply, err
	}
	
	defer res.Body.Close()

	reply.Message = "PUT "+strconv.Itoa(args.Key)+" ==> "+args.Value+" to serverNode["+serverPort+"] successful!"

	return reply, nil
}


func GET(key int) (ReplyMessage, error) {

	var reply ReplyMessage

	serverPort := GetServerNode(key)

	urlPath :=  "http://127.0.0.1:"+serverPort+"/keys/"+strconv.Itoa(key)

	res, err := http.Get(urlPath)
	if err!=nil {
		fmt.Println("GET: http.Get",err)
		reply.Message = "GET: http.Get error"
		return reply, err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println("GET: ioutil.ReadAll",err)
		reply.Message = "GET: ioutil.ReadAll error"
		return reply, err
	}
	var data Data
	err = json.Unmarshal(body, &data)
	//fmt.Println(prices)
	reply.Message = "GET "+strconv.Itoa(data.Key)+" ==> "+data.Value+" from serverNode["+serverPort+"]"
	//fmt.Println("result: ",result)
	return reply, nil
}

func initConsistentHash() {
	add(3000)
	fmt.Println(circle)
	//fmt.Println(circleKey)
	add(3001)
	fmt.Println(circle)
	//fmt.Println(circleKey)
	add(3002)
	circleKey = sortMap(circle)
	fmt.Println("initConsistentHash")
	fmt.Println(circle)
	fmt.Println(circleKey)

}


func add(servernode int) {
  	key := hashFunction(servernode)
  	circle[key] = strconv.Itoa(servernode)
}

func remove(servernode int) {
  key := hashFunction(servernode)
  delete(circle, key)  
}

func GetServerNode(key int) string {
    if len(circle) == 0 {
      return "";
    }
	fmt.Println("GetServerNode")
    hash := hashFunction(key)
	fmt.Println(hash)
	//fmt.Println("circle[key]:" + circle[hash])
    _, ok := circle[hash]
    fmt.Println(ok)
    if (ok == false) {
    	i := 0
    	var res string
		for k := range circleKey {
			//fmt.Println(k)
        	if i == 0 {
        		res = circle[circleKey[k]]
        		i = -1
        	}
			if circleKey[k] >= hash {
				fmt.Println("circle[k]:" + circle[circleKey[k]])
				return circle[circleKey[k]]
			} 
		}
		return res
    }

    return circle[hash]
}

func hashFunction(key int) int {
	hash := 7;
	hash = 31*hash + key;
    return 31*hash%13 
}

func sortMap(circle map[int]string) []int {
	mk := make([]int, len(circle))
	i := 0
    for k, _ := range circle {
        mk[i] = k
        i++
    }
    sort.Ints(mk)
    fmt.Println(mk)
    return mk
}


