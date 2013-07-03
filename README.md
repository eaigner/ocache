**ocache** is a wrapper for [gomemcache](https://github.com/bradfitz/gomemcache) that adds support for namespacing and object en-/decoding.

### Example

```go
oc := ocache.New("127.0.0.1:11211")

// Simple set
oc.Set(&dataIn, 3600, "simpleKey")

// Namespaced set
oc.Set(&dataIn2, 3600, "ns1", "namespacedKey")

// Simple get
oc.Get(&dataOut, "simpleKey")

// Namespaced get
oc.Get(&dataOut2, "ns1", "namespacedKey")

// Simple delete
oc.Delete("simpleKey")

// Delete entire namespace
oc.DeleteNamespace("ns1")

// Delete specific item in namespace
oc.Delete("ns1", "namespacedKey")
```