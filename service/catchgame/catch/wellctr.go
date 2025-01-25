package catch

import (
	"sync"
)

type wellKey struct {
	UserId int64
	ObjId  int64
}

// 吃井计算器
type wellCtr struct {
	data map[wellKey]int64
	lock sync.Mutex

	wellNum int64 // 吃井数量
	maxCtr  int64 // 计数器个数
}

func NewWellCtr(wellNum, maxCtr int64) *wellCtr {
	return &wellCtr{
		wellNum: wellNum,
		maxCtr:  maxCtr,
		data:    make(map[wellKey]int64),
	}
}

// 吃井检测
func (ctr *wellCtr) Check(key wellKey, result bool) bool {
	ctr.lock.Lock()
	defer ctr.lock.Unlock()

	if ctr.maxCtr > 0 && len(ctr.data) > int(ctr.maxCtr) { // 超限自动清理
		ctr.data = make(map[wellKey]int64)
	}
	// 成功，清理失败计数器
	if result {
		ctr.data[key] = 0
		return true
	} else { // 失败，计算是否吃井，改换结果
		if ctr.data[key] < ctr.wellNum {
			ctr.data[key]++
			return false
		} else {
			ctr.data[key] = 0
			return true
		}
	}
}
