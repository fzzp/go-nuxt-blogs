package luckid

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var once sync.Once
var lk *luck

// luck 幸运id结构体，生成唯一ID
type luck struct {
	next     int64
	min      int64 // 增量最小值
	max      int64 // 增量最大值
	mutex    sync.RWMutex
	rnd      *rand.Rand
	filepath string // 存储id的文件路径
}

// New 初始化幸运id实例
func New(filepath string, min, max int64) error {
	var err error
	once.Do(func() {
		lk, err = newLuckID(filepath, min, max)
	})
	return err
}

// Get 获取当前的id
func Get() int64 {
	return lk.get()
}

// Next 生成下一个幸运id
func Next() int64 {
	return lk._next()
}

// Save 最后生成的幸运id保存到指定文件
func Save() error {
	return lk.save()
}

// new 创建一个幸运id对象
func newLuckID(filepath string, min, max int64) (*luck, error) {
	if min <= 0 {
		min = 1
	}
	if max <= min {
		max = min + 3
	}
	var id = 1
	// 从文件读取幸运id
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		// 文件不存在则创建并写入1，幸运id从1开始
		if errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(filepath)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			_, err = f.Write([]byte("1"))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		tmp, err := strconv.Atoi(string(bytes))
		if err != nil {
			return nil, err
		}
		id = tmp
	}

	// 不符合，则默认为1，
	if id <= 0 {
		id = 1
	}

	luck := &luck{
		next:     int64(id),
		mutex:    sync.RWMutex{},
		rnd:      rand.New(rand.NewSource(time.Now().UnixNano())),
		min:      min,
		max:      max,
		filepath: filepath,
	}

	return luck, nil
}

// Get 获取当前的id
func (l *luck) get() int64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.next
}

// Next 生成下一个幸运id
func (l *luck) _next() int64 {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	num := l.min + l.rnd.Int63n(l.max-l.min+1)
	l.next = l.next + num
	return l.next
}

// Save 最后生成的幸运id保存到指定文件
func (l *luck) save() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	f, err := os.Create(l.filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d", l.next)
	if err != nil {
		return err
	}
	return err
}
