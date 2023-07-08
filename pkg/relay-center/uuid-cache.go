package relay_center

import (
	"sync"
	"time"
)

const MaxUUIDCacheSize = 51200

type uuidCache struct {
	locker sync.Mutex
	uuids  map[uint32]*uuidNode
}

type uuidNode struct {
	uuid   string
	expire int64
}

var UUIDCache = &uuidCache{
	uuids: make(map[uint32]*uuidNode),
}

func init() {
	go UUIDCache.runCleanUUID()
}

// 存储UUID到缓存中, 如果已经存在则返回false
func (self *uuidCache) Store(uuid string, expire int64) bool {
	// 将UUID转换为hash值, 并限定在最大缓存大小内
	// 这里有一个潜在问题, 如果UUID的数量超过了最大缓存大小, 那么就会出现hash冲突
	// 为了性能考虑, 这里不做处理, 但是这种情况出现的概率非常小
	id := self.calcUUIDHash(uuid) % MaxUUIDCacheSize

	self.locker.Lock()
	defer self.locker.Unlock()

	node, ok := self.uuids[id]
	// 查uuid是否已经存在，并使用过
	if ok && node.uuid == uuid && node.expire > time.Now().Unix() {
		return false
	}

	// 存储新节点
	self.uuids[id] = &uuidNode{
		uuid:   uuid,
		expire: expire,
	}

	return true
}

func (self *uuidCache) calcUUIDHash(uuid string) uint32 {
	const seed = 131 // 31, 131, 1313, 13131, 131313, etc..
	var hash uint32 = 0

	for _, c := range []byte(uuid) {
		hash = hash*seed + uint32(c)
	}

	return hash
}

// 清理过期的UUID
func (self *uuidCache) runCleanUUID() {
	tick := time.NewTicker(time.Second * 60)
	defer tick.Stop()

	for {
		<-tick.C

		self.doCleanUUID()
	}
}

func (self *uuidCache) doCleanUUID() {
	self.locker.Lock()
	defer self.locker.Unlock()

	now := time.Now().Unix()
	remove := make([]uint32, 0, len(self.uuids))
	// 查找过期的UUID
	for id, node := range self.uuids {
		if node.expire < now {
			remove = append(remove, id)
		}
	}

	for _, id := range remove {
		delete(self.uuids, id)
	}
}
