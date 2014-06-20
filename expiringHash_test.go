// expiringHash_test.go
package expiringHash

import (
	"fmt"
	"testing"
	"time"
)

func TestPutGetExpire(t *testing.T) {
	testKey := "testKey"
	testValue := "testValue"
	eh := New()
	eh.Put(testKey, testValue, time.Second)
	value, found := eh.Get(testKey)
	if value.(string) != testValue {
		t.Errorf("expected %s received %s\n", testValue, value)
	}
	if !found {
		t.Error("found should have been true, but was false")
	}
	if eh.Len() != 1 {
		t.Errorf("Len() did not work correctly, received len of %d but expected 1\n", eh.Len())
	}
	t.Log("len worked correctly")
	t.Log("found value correctly before expiration")
	time.Sleep(2 * time.Second)
	value, found = eh.Get(testKey)
	if found {
		t.Errorf("should not have found %s after ttl\n", testKey)
	}
	t.Log("found value is correct for missing item")
	t.Log("expiration worked")
	if eh.Len() != 0 {
		t.Errorf("Len() did not work correctly, received len of %d but expected 0\n", eh.Len())
	}
}

func TestItemNotFound(t *testing.T) {
	testKey := "testKey"
	eh := New()
	if _, found := eh.Get(testKey); found {
		t.Errorf("Get() returned true for key %s when it did not exist\n", testKey)
	}
}

func TestMultipleUpdates(t *testing.T) {
	testKey := "testKey"
	values := []string{"1", "2", "3"}
	eh := New()
	for _, value := range values {
		eh.Put(testKey, value, time.Second)
		retrieved, found := eh.Get(testKey)
		if !found {
			t.Error("found not true")
		}
		if retrieved.(string) != value {
			t.Errorf("expected %s received %s\n", value, retrieved)
		}
	}
	time.Sleep(2 * time.Second)
	_, found := eh.Get(testKey)
	if found {
		t.Error("expire should have happened")
	}
}

func TestStats(t *testing.T) {
	eh := New()
	testKey := "testKey"
	for i := 0; i < 100; i++ {
		value := fmt.Sprintf("%d", i)
		eh.Put(value, value, 100*time.Millisecond)
		eh.Get(value)
		eh.Get(value)
		eh.Get("notfound")
		eh.Del(value)
	}
	eh.Put(testKey, testKey, time.Nanosecond)
	time.Sleep(time.Millisecond)
	stats := eh.Stats()
	if stats.Puts != 101 {
		t.Errorf("puts should be 101 but was %d\n", stats.Puts)
	}
	if stats.GetHits != 200 {
		t.Errorf("getHits should be 200 but was %d\n", stats.GetHits)
	}
	if stats.GetMisses != 100 {
		t.Errorf("getMisses should be 100 but was %d\n", stats.GetMisses)
	}
	if stats.Expired != 1 {
		t.Errorf("expired should be 1 but was %d\n", stats.Expired)
	}
}
