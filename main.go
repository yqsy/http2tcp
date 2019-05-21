package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"strings"
	"os"
)

var (
	Version  = "1.0"
	HttpAddr = ":50003"
	DestAddr = "127.0.0.1:50001"
)

func handler(w http.ResponseWriter, r *http.Request) {
	remoteConn, err := net.Dial("tcp", DestAddr)
	if err != nil {
		log.Printf("dial error [%s] -> [%s]", r.RemoteAddr, DestAddr)
		return
	}
	defer remoteConn.Close()

	// remote tcp line + jsonrpc
	bufReader := bufio.NewReader(remoteConn)
	bufWriter := bufio.NewWriter(remoteConn)
	textReader := textproto.NewReader(bufReader)
	textWriter := textproto.NewWriter(bufWriter)

	// TODO 外置配置,可选
	// 握手请求
	_ = textWriter.PrintfLine(`{"jsonrpc":"2.0","method":"server.version","params":["http2tcp","1.4"],"id":1}`)

	// 握手应答
	// {"jsonrpc": "2.0", "result": ["ElectrumX 1.9.5", "1.2"], "id": 1}
	_, err = textReader.ReadLine()
	if err != nil {
		log.Printf("handle error [%s] <- [%s]", r.RemoteAddr, DestAddr)
		return
	}

	log.Printf("new proxy [%s] <-> [%s] \n", r.RemoteAddr, DestAddr)

	// [local] -> [remote]
	jsonReqBytes, _ := ioutil.ReadAll(r.Body)
	jsonReq := string(jsonReqBytes)

	// TODO 外置配置,可选
	// 将收来请求的换行符号去除,electrumx的不识别带有换行的json数据. 因为他是以\n为终止的.
	newJsonReq := strings.Replace(jsonReq, "\n", "", -1)
	_ = textWriter.PrintfLine(newJsonReq)

	// [local] <- [remote]
	line, err := textReader.ReadLine()
	if err != nil {
		log.Printf("get response error [%s] <- [%s]", r.RemoteAddr, DestAddr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(line))
	log.Printf("stop proxy [%s] <-> [%s] \n", r.RemoteAddr, DestAddr)
}

func main() {
	arg := os.Args
	
	if len(arg) < 3 {
		log.Printf("Usage:\n %v httpaddr destaddr", arg[0])
		return
	}

	HttpAddr = arg[1]
	DestAddr = arg[2]

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("http2tcp version: %s\n", Version)
	log.Printf("relay http(%s) -> tcp(%s)\n", HttpAddr, DestAddr)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(HttpAddr, nil)
	if err != nil {
		panic(err)
	}
}
