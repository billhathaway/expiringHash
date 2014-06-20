expiringHash
==
This was an experiment to see how easy it would be to write a library for an expiring hash table.

It is not optimized, but seems to work correctly.  

[![GoDoc](https://godoc.org/github.com/billhathaway/expiringHash?status.png)](http://godoc.org/github.com/billhathaway/expiringHash)

Example:
```
eh := expiringHash.New()

// insert an item with a one second TTL
eh.Put("key","value",time.Second)

// value will be the value put in, found will be true
value,found := eh.Get("key")

// sleep to let the item expire
time.Sleep(2 * time.Second)

// value will be nil, found will be false
value,found = eh.Get("key")
```
