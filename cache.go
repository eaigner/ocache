package ocache

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
	"sync"
	"time"
)

type Ocache struct {
	c         *memcache.Client
	ensureMtx sync.Mutex
}

func New(servers ...string) *Ocache {
	return &Ocache{c: memcache.New(servers...)}
}

func (o *Ocache) Get(v interface{}, compositeKey ...string) error {
	switch len(compositeKey) {
	case 1:
		return o.getSimple(v, compositeKey[0])
	case 2:
		return o.getNamespaced(v, compositeKey[0], compositeKey[1])
	}
	panic("invalid key")
}

func (o *Ocache) getSimple(v interface{}, key string) error {
	item, err := o.c.Get(key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(item.Value)).Decode(v)
}

func (o *Ocache) getNamespaced(v interface{}, ns, key string) error {
	nsItem, err := o.c.Get(ns)
	if err != nil {
		return err
	}
	nsk := string(nsItem.Value)
	return o.getSimple(v, makeKey(nsk, key))
}

func (o *Ocache) Set(v interface{}, expire time.Duration, compositeKey ...string) error {
	switch len(compositeKey) {
	case 1:
		return o.setSimple(v, expire, compositeKey[0])
	case 2:
		return o.setNamespaced(v, expire, compositeKey[0], compositeKey[1])
	}
	panic("invalid key")
}

func (o *Ocache) setSimple(v interface{}, expire time.Duration, key string) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return err
	}
	return o.c.Set(&memcache.Item{
		Key:        key,
		Value:      buf.Bytes(),
		Expiration: int32(expire.Seconds()),
	})
}

func (o *Ocache) setNamespaced(v interface{}, expire time.Duration, ns, key string) error {
	nsk := o.makeNamespaceKey(ns)
	return o.setSimple(v, expire, makeKey(nsk, key))
}

func (o *Ocache) Delete(compositeKey ...string) error {
	switch len(compositeKey) {
	case 1:
		return o.deleteSimple(compositeKey[0])
	case 2:
		return o.deleteNamespaced(compositeKey[0], compositeKey[1])
	}
	panic("invalid key")
}

func (o *Ocache) DeleteNamespace(ns string) error {
	return o.c.Set(&memcache.Item{
		Key:   ns,
		Value: []byte(generateKey()),
	})
}

func (o *Ocache) deleteSimple(key string) error {
	return o.c.Delete(key)
}

func (o *Ocache) deleteNamespaced(ns, key string) error {
	return o.deleteSimple(makeKey(o.makeNamespaceKey(ns), key))
}

func makeKey(ns, key string) string {
	return ns + `:` + key
}

func (o *Ocache) makeNamespaceKey(ns string) string {
	o.ensureMtx.Lock()
	defer o.ensureMtx.Unlock()

	item, err := o.c.Get(ns)
	switch err {
	case nil:
		return string(item.Value)
	case memcache.ErrCacheMiss:
		newKey := generateKey()
		err2 := o.c.Add(&memcache.Item{
			Key:   ns,
			Value: []byte(newKey),
		})
		switch err2 {
		case nil:
			return newKey
		case memcache.ErrNotStored:
			item, err3 := o.c.Get(ns)
			if err3 != nil {
				panic(err3) // should be unreachable
			}
			return string(item.Value)
		default:
			panic(err2)
		}
	default:
		panic(err)
	}
	panic("unreachable")
}

// generateKey generates a unique 128 bit value base64 encoded for use as key
func generateKey() string {
	b := make([]byte, 8)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, time.Now().UnixNano())
	buf.Write(b)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
