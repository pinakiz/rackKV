# rackKV
🚀 RackKV is a lightweight, log-structured key-value storage engine written in Go.
It is written in Go and designed for simplicity, durability, and fast read/write performance.
It supports both a simple HTTP API and a CLI interface, making it easy to experiment with database internals and performance.


## ✨ Features  

- ⚡ **High-throughput writes** with append-only log files  
- 🗂️ **In-memory key directory** for O(1) lookups  
- 🪦 **Tombstone support** for deletes  
- 📝 **Hint files** for faster recovery on restart  
- ♻️ **Compaction & merging** to reclaim space and remove stale entries  
- 🌐 **HTTP-based API** for easy integration  
- 💻 **CLI interface** for direct terminal interaction  



## 🏗️ Architecture

- Write Path (PUT/DELETE)
  - Append entry to active data file
  - Update in-memory keydir
  - Mark deletes with tombstones

- Read Path (GET)
  - Lookup keydir → get file & offset
  - Read value directly from disk

- Recovery
  - Use hint files to rebuild keydir quickly
  - Fallback to scanning data files if needed

- Compaction & Merging
  - Periodically rewrite only live keys into new data files
  - Remove obsolete/tombstoned entries
  - Generate fresh hint files

## 🔑 Interfaces
🌐 HTTP API 
- PUT

   ```/put?key=<key>&val=<value>```
- GET
 
   ```/get?key```
- DELETE
 
   ```/delete?key```

## 📊 Benchmarks
We benchmarked using wrk with 8 threads and 1000 connections for 30s:
```
wrk -t8 -c1000 -d30s -s put.lua http://localhost:8080

```
### Results:
- Requests/sec: 35,668
- Latency (avg): 28.25 ms
- Throughput: 4.01 MB/s
- Total Requests: 107,285 in 30s

💻 CLI
```rackkv open           # start server  
rackkv put key value  # insert key-value pair  
rackkv get key        # fetch value for key  
rackkv delete key     # delete key  
```

## 🚀 Running RackKV
Clone the repo and run:
```
go run main.go
```
To use the CLI:
```
go run client/client.go
```
## 📚 Learnings
RackKV is not just a project, but also a learning journey into database internals — from log-structured storage to compaction strategies, performance benchmarking and much more.

 ________  ________  ________  ___  __    ___  __    ___      ___ 
|\   __  \|\   __  \|\   ____\|\  \|\  \ |\  \|\  \ |\  \    /  /|
\ \  \|\  \ \  \|\  \ \  \___|\ \  \/  /|\ \  \/  /|\ \  \  /  / /
 \ \   _  _\ \   __  \ \  \    \ \   ___  \ \   ___  \ \  \/  / / 
  \ \  \\  \\ \  \ \  \ \  \____\ \  \\ \  \ \  \\ \  \ \    / /  
   \ \__\\ _\\ \__\ \__\ \_______\ \__\\ \__\ \__\\ \__\ \__/ /   
    \|__|\|__|\|__|\|__|\|_______|\|__| \|__|\|__| \|__|\|__|/    
                                                                  
