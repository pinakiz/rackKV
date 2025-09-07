package main

import (
	"fmt"
	"net/http"
	"rackKV/pkg"
	"strconv"
)


func main(){
	handler := &pkg.RackHandle{
		Mode: pkg.Mode{IsUp: true},
	}

	handler.Mode.IsUp = false
	defer handler.Close()
	http.HandleFunc("/open" , func(w http.ResponseWriter , r *http.Request){
		fmt.Println("hehe: ",handler)
		rw := r.URL.Query().Get("rw")
		readwrite , _ := strconv.ParseBool(rw);
		syn := r.URL.Query().Get("syn")
		sync, _ := strconv.ParseBool(syn);
		if(handler.Mode.IsUp){
			w.Write([]byte("Db is already opened"))
			return
		}
		temp , err := pkg.Open(".", pkg.Mode{ReadWrite : readwrite ,SyncOnWrite : sync})
		handler = temp;
		if(err != nil){
			w.Write([]byte(err.Error()))
		}else{
			w.Write([]byte("OK"));
		}
	})	
	
	http.HandleFunc("/put",func(w http.ResponseWriter , r *http.Request){
		// key := r.URL.Query().Get("key");
		// value := r.URL.Query().Get("Value");
		fmt.Println(handler.Mode.IsUp)
		fmt.Println(handler.Mode.ReadWrite)

		if(!handler.Mode.IsUp || !handler.Mode.ReadWrite ){
			w.Write([]byte("permission denied: Db is in read-only mode"))
			return
		}
		pkg.PUT(handler,"hi","bye")
		w.Write([]byte("OK"));
	})

	http.HandleFunc("/get" , func(w http.ResponseWriter , r *http.Request){
		// key := r.URL.Query().Get("key");
		w.Write([]byte("OK"));
	})

	http.ListenAndServe(":8080" , nil);
}




























































