package luckid_test

import (
	"go-nuxt-blogs/pkg/luckid"
	"log"
	"sync"
	"testing"
)

func TestLuckID(t *testing.T) {
	err := luckid.New("./luck.txt", 3, 33)
	if err != nil {
		t.Errorf("got %s want nil", err)
	}

	if id := luckid.Get(); id != 1 {
		t.Errorf("got %d want 1", id)
	}

	testNum := 10000
	total := 0
	data := make(map[int64]struct{})

	defer func() {
		log.Printf("测试 %d 次，重复 %d 次：", testNum, total)
		if err := luckid.Save(); err != nil {
			t.Errorf("got %v want nil", err)
		}
	}()

	var dataMutex sync.Mutex
	for i := 0; i < testNum; i++ {
		newID := luckid.Next()

		dataMutex.Lock()
		_, exists := data[newID]
		data[newID] = struct{}{}
		dataMutex.Unlock()

		if exists {
			t.Errorf("got %t want %t", false, true)
			total++
		}

		log.Println(newID)
	}

}
