package ocache

import (
	"testing"
)

var oc = New("127.0.0.1:11211")

func TestSimple(t *testing.T) {
	v := "asdf"
	k := "key1"

	if err := oc.Set(&v, 0, k); err != nil {
		t.Fatal(err)
	}

	var v2 string
	if err := oc.Get(&v2, k); err != nil {
		t.Fatal(err)
	}
	if v != v2 {
		t.Fatal(v2)
	}

	if err := oc.Delete(k); err != nil {
		t.Fatal(err)
	}

	var v3 string
	if err := oc.Get(&v3, k); err == nil {
		t.Fatal("should report error")
	}
	if v3 != "" {
		t.Fatal(v3)
	}
}

func TestNamespaced(t *testing.T) {
	v, v2, v3 := "a", "b", "c"
	ns, ns2 := "ns1", "ns2"
	k, k2 := "key1", "key2"

	if err := oc.Set(&v, 0, ns, k); err != nil {
		t.Fatal(err)
	}
	if err := oc.Set(&v2, 0, ns, k2); err != nil {
		t.Fatal(err)
	}
	if err := oc.Set(&v3, 0, ns2, k2); err != nil {
		t.Fatal(err)
	}

	var vv string
	if err := oc.Get(&vv, ns, k); err != nil {
		t.Fatal(err)
	}
	if vv != v {
		t.Fatal(vv)
	}
	var vv2 string
	if err := oc.Get(&vv2, ns, k2); err != nil {
		t.Fatal(err)
	}
	if vv2 != v2 {
		t.Fatal(vv2)
	}
	var vv3 string
	if err := oc.Get(&vv3, ns2, k2); err != nil {
		t.Fatal(err)
	}
	if vv3 != v3 {
		t.Fatal(vv3)
	}

	// Remove first namespace
	if err := oc.DeleteNamespace(ns); err != nil {
		t.Fatal(err)
	}

	vv = ""
	if err := oc.Get(&vv, ns, k); err == nil {
		t.Fatal()
	}
	if vv != "" {
		t.Fatal(vv)
	}
	vv2 = ""
	if err := oc.Get(&vv2, ns, k2); err == nil {
		t.Fatal()
	}
	if vv2 != "" {
		t.Fatal(vv2)
	}
	vv3 = ""
	if err := oc.Get(&vv3, ns2, k2); err != nil {
		t.Fatal(err)
	}
	if vv3 != v3 {
		t.Fatal(vv3)
	}

	// Delete object directly
	if err := oc.Delete(ns2, k2); err != nil {
		t.Fatal(err)
	}

	vv3 = ""
	if err := oc.Get(&vv3, ns2, k2); err == nil {
		t.Fatal(err)
	}
	if vv3 != "" {
		t.Fatal(vv3)
	}

}
