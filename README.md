### Bitcask

This is a golang implementation of [Bitcask by Riak](https://riak.com/assets/bitcask-intro.pdf) paper.

### TODO
- [ ] Transactions
- [ ] BTree/Radix Tree based Indexing
- [ ] Benchmarking
- [ ] RAFT implementation

### Benefits of this approach

- Low Latency: Write queries are handled with a single O(1) disk seek. Keys lookup happen in memory using a hash table lookup. This makes it possible to achieve low latency even with a lot of keys/values in the database. Bitcask also relies on the filesystem read-ahead cache for a faster reads.
- High Throughput: Since the file is opened in "append only" mode, it can handle large volumes of write operations with ease. 
- Predictable performance: The DB has a consistent performance even with growing number of records. This can be seen in benchmarks as well.
- Crash friendly: Bitcask commits each record to the disk and also generates a "hints" file which makes it easy to recover in case of a crash.
- Elegant design: Bitcask achieves a lot just by keeping the architecture simple and relying on filesystem primitives for handling complex scenarios (for eg: backup/recovery, cache etc).
- Ability to handle datasets larger than RAM.

### Limitations

- The main limitation is that all the keys must fit in RAM since they're held inside as an in-memory hash table. A potential workaround for this could be to shard the keys in multiple buckets. Incoming records can be hashed into different buckets based on the key. A shard based approach allows each bucket to have limited RAM usage.
