It is necessary to write an in-memory cache that will return a profile and a list of orders based on a key (user UUID).

1. The cache must have a TTL (2 seconds)
2. The cache can be used by function(s) that work with orders (add/update/delete). If the TTL has expired, it returns nil. Upon update, the TTL is reset to 2 seconds. The methods must be thread-safe
3. Test scenarios for using this cache must be written (do not modify the basic structures)

Additional task: automatic cleanup of expired entries