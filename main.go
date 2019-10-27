package main

import (
	"flag"
	"fmt"
	"github.com/pingcap/tidb/util/codec"
	"log"
	"net/http"
	"os"
)

var keyFormat = flag.String("format", "proto", "output format (go/hex/base64/proto)")
var enableServer = flag.Bool("enable-server", false, "enable server")

func main() {
	flag.Parse()

	if *enableServer {
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":4321", nil))
		return
	}

	if flag.NArg() != 1 {
		fmt.Println("usage:\nmok {flags} {key}")
		flag.PrintDefaults()
		os.Exit(1)
	}

	N("key", []byte(flag.Arg(0))).Print()
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	keys, ok := r.URL.Query()["decode"]
	if !ok || len(keys[0]) < 1 {

		log.Println("Url Param 'key' is missing")
		return
	}
	key := keys[0]
	N("key", []byte(key)).Expand()

	fmt.Fprintf(w, `{"test": "xzx"}`)
}

func findTableID(node *Node) string {
	if node.typ == "table_id" {
		_, tableId, _ := codec.DecodeInt(node.val)
		return fmt.Sprintf("%d",tableId)
	}
	for _, v := range node.variants {
		for _, c := range v.children {
			return findTableID(c)
		}
	}
	return ""
}

//func main() {
//	http.HandleFunc("/", handler)
//	log.Fatal(http.ListenAndServe(":4321", nil))
//}