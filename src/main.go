package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"rackKV/pkg"
	"strconv"
)


func clearScreen() {
	fmt.Print("\033[H\033[2J") // Moves cursor home and clears the screen
}



func main(){
	handler := &pkg.RackHandle{
	}
	fmt.Println("Generating Hint Files")
	if err := pkg.Generate_hintFiles(); err != nil{
		 fmt.Println(err);
	}
	fmt.Println("Generating KeyDir");
	if err := pkg.GenerateKeyDir(handler); err != nil{
		 fmt.Println(err);
	}
	clearScreen()
	fmt.Println("Done");
	clearScreen()
		banner := ` ________  ________  ________  ___  __    ___  __    ___      ___ 
|\   __  \|\   __  \|\   ____\|\  \|\  \ |\  \|\  \ |\  \    /  /|
\ \  \|\  \ \  \|\  \ \  \___|\ \  \/  /|\ \  \/  /|\ \  \  /  / /
 \ \   _  _\ \   __  \ \  \    \ \   ___  \ \   ___  \ \  \/  / / 
  \ \  \\  \\ \  \ \  \ \  \____\ \  \\ \  \ \  \\ \  \ \    / /  
   \ \__\\ _\\ \__\ \__\ \_______\ \__\\ \__\ \__\\ \__\ \__/ /   
    \|__|\|__|\|__|\|__|\|_______|\|__| \|__|\|__| \|__|\|__|/    
                                                                  
                                                                  
                                                                  `
	fmt.Println(banner)
	fmt.Println("SERVER LISTENING")

	in_mem_map := handler.KeyDir;
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go pkg.MergerListener(ctx, in_mem_map, handler)


	defer handler.Close()
	defer handler.ActiveFile.Close();
	http.HandleFunc("/open" , func(w http.ResponseWriter , r *http.Request){
		fmt.Println("connected")
		rw := r.URL.Query().Get("rw")
		readwrite , _ := strconv.ParseBool(rw);
		syn := r.URL.Query().Get("syn")
		sync, _ := strconv.ParseBool(syn);
		if(handler.Mode.IsUp){
			w.Write([]byte("Db is already opened"))
			return
		}
		if err := pkg.Open(".", pkg.Mode{ReadWrite : readwrite ,SyncOnWrite : sync}, handler);(err != nil){
			w.Write([]byte(err.Error()))
		}else{
			w.Write([]byte("OK"));
		}
	})	

	
	http.HandleFunc("/put",func(w http.ResponseWriter , r *http.Request){
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")

		// fmt.Print("in put: ",handler)
		if(!handler.Mode.IsUp || !handler.Mode.ReadWrite ){
			w.Write([]byte("permission denied: Db is in read-only mode"))
			return
		}
		_ , err := pkg.PUT(handler,key,value)
		if(err != nil){
			fmt.Println("Error: ",err);
		}else{
			w.Write([]byte("OK"));
		}

	})

	http.HandleFunc("/get" , func(w http.ResponseWriter , r *http.Request){
		key := r.URL.Query().Get("key");
		val , err := pkg.GET(handler,key);
		if(!handler.Mode.IsUp){
			fmt.Println("Db is not up")
			w.Write([]byte("Db is not up"))
			return
		}
		if(err != nil) {
			fmt.Println("Error: ",err);
		}else{
			w.Write([]byte(val));
		}
	})
	http.HandleFunc("/delete" , func(w http.ResponseWriter , r *http.Request){
		key := r.URL.Query().Get("key");
		val , err := pkg.PUT(handler,key , "");
		if(!handler.Mode.IsUp){
			fmt.Println("Db is not up")
			w.Write([]byte("Db is not up"))
			return
		}
		if(err != nil) {
			fmt.Println("Error: ",err);
		}else{
			w.Write([]byte(val));
		}
	})


	http.ListenAndServe(":8080" , nil);
}