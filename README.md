# rackKV
🚀 RackKV is a lightweight, log-structured key-value storage engine written in Go.
It is written in Go and designed for simplicity, durability, and fast read/write performance.
It supports both a simple HTTP API and a CLI interface, making it easy to experiment with database internals and performance.


✨ Features

> ⚡ High-throughput writes with append-only log files
> 🗂️ In-memory key directory for O(1) lookups
>  🪦 Tombstone support for deletes
> 📝 Hint files for faster recovery on restart
> ♻️ Compaction & merging to reclaim space and remove stale entries
> 🌐 HTTP-based API for easy integration
> 💻 CLI interface for direct terminal interaction


🏗️ Architecture

Write Path (PUT/DELETE)
  Append entry to active data file
  Update in-memory keydir
  Mark deletes with tombstones

Read Path (GET)
  Lookup keydir → get file & offset
  Read value directly from disk

Recovery
  Use hint files to rebuild keydir quickly
  Fallback to scanning data files if needed

Compaction & Merging
  Periodically rewrite only live keys into new data files
  Remove obsolete/tombstoned entries
  Generate fresh hint files


🔑 Interfaces
🌐 HTTP API 
• PUT
    
```/put?key=<key>&val=<value>```


