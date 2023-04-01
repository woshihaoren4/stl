package vector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestVector(t *testing.T) {
	vec := NewVector(1, 2, 3)
	assert.Equal(t, 3, vec.Len(), "初始化长度为3")

	vec.Push(4, 5, 6)
	buf, _ := vec.UnsafeRangeDef(2, 6)
	assert.Equal(t, []int{3, 4, 5, 6}, buf, "向末尾插入三个值")

	vec.Insert(0, 0)
	val, _ := vec.UnsafeDef(0)
	assert.Equal(t, 0, *val, "在开始位置插入")

	vec.Insert(7, 7)
	v, _ := vec.Pop()
	assert.Equal(t, 7, v, "在末尾位置插入")

	vec.Insert(1, 666)
	v, _ = vec.Remove(1)
	assert.Equal(t, 666, v, "在位置[1]插入")

	v, _ = vec.Remove(0)
	assert.Equal(t, 0, v, "在开始位置移除")

	v, _ = vec.Remove(vec.Len() - 1)
	assert.Equal(t, 6, v, "在末尾位置移除")
	assert.Equal(t, 5, vec.Len(), "当前剩余长度应该为5")

	//1,2,3,4,5
	list, _ := vec.UnsafeRangeDef(0, 5)
	for i, _ := range list {
		list[i] = i
	}
	list, _ = vec.UnsafeRangeDef(1, 2)
	assert.Equal(t, []int{1}, list, "不安全的引用可以变更内容")

	v, _ = vec.Load(0)
	assert.Equal(t, 0, v, "load第一个位置为0")

	v, _ = vec.Load(4)
	assert.Equal(t, 4, v, "load最后一个位置为4")

	list, _ = vec.RangeLoad(0, vec.Len())
	list[0] = 666
	list, _ = vec.RangeLoad(0, vec.Len())
	assert.Equal(t, []int{0, 1, 2, 3, 4}, list, "load range 不改变原来变量")

	vec.Insert(1, 666)
	exist := vec.Contain(666)
	assert.Equal(t, true, exist, "当前的结果集里存在666")

	//0,1,2,3,4,666
	vec.Filter(func(t *int) bool {
		return *t < 10
	})
	list = vec.UnsafeIntoSlice()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, list, "过滤掉大于10的变量")

	vec.Foreach(func(i int, v int) {
		fmt.Printf("index[%d]--->%d\n", i, v)
	})
	vec.Reverse()
	vec.Update(0, 5)
	vec.Update(1, 4)
	vec.Update(2, 3)
	vec.Update(3, 2)
	vec.Update(4, 1)
	vec.Update(5, 0)
	list = vec.UnsafeIntoSlice()
	assert.Equal(t, []int{5, 4, 3, 2, 1, 0}, list, "翻转并且+1并且+1")

	vec.Push(666)
	assert.Equal(t, 7, vec.Len(), "长度又变成7个")
	vec.Clean()
	assert.Equal(t, 0, vec.Len(), "清理完毕长度为0")
}

func TestConcurrency(t *testing.T) {
	vec := NewVector[int]()
	wg := sync.WaitGroup{}
	wg.Add(20)
	//起10线程向里面加数据
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 10000; i++ {
				vec.Push(i)
			}
			wg.Done()
			fmt.Println("--->数据添加完毕")
		}()
	}
	//再起十个线程拿数据
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 10000; i++ {
				_, ok := vec.Pop()
				if !ok {
					i--
				}
			}
			wg.Done()
			fmt.Println("--->数据获取完毕")
		}()
	}
	wg.Wait()
	assert.Equal(t, 0, vec.Len(), "拿去完毕，剩余0")

}
