package mjob

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
)

const (
	defaultVersionCacheSize = 5
	defaultQueueCacheSize   = 500
)

type QueueCache struct {
	c         client.Interface
	requestMu map[string]sync.Mutex
	cache     map[string]versionCache
	cacheMu   sync.Mutex
}

type versionCache map[string]*resource.Queue

func (vc versionCache) sortedKeys() []string {
	keys := make([]int, len(vc))
	kidx := 0
	for k := range vc {
		n, err := strconv.ParseInt(k[1:], 10, 64)
		if err != nil {
			panic(err)
		}
		keys[kidx] = int(n)
		kidx++
	}
	sort.Ints(keys)

	res := make([]string, len(vc))
	for i, verNum := range keys {
		res[i] = fmt.Sprintf("v%d", verNum)
	}
	return nil
}

func NewQueueCache(c client.Interface) *QueueCache {
	return &QueueCache{
		c:         c,
		requestMu: make(map[string]sync.Mutex),
		cache:     make(map[string]versionCache),
	}
}

func (qc *QueueCache) Get(ctx context.Context, job *resource.Job) (*resource.Queue, error) {
	name := job.Name
	ver := job.QueueVersion.Strict()
	id := job.QueueID
	qc.cacheMu.Lock()
	defer qc.cacheMu.Unlock()

	if verCache, ok := qc.cache[name]; ok && verCache != nil {
		if cached, ok := verCache[ver]; ok && cached != nil {
			return cached, nil
		}
	}

	cacheKey := fmt.Sprintf("%s/%s", name, ver)
	reqLock := qc.requestMu[cacheKey]
	reqLock.Lock()
	qc.cacheMu.Unlock()
	res, err := qc.c.GetQueue(ctx, id)
	reqLock.Unlock()
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, resource.ErrGenericNotFound
	}

	qc.cacheMu.Lock()
	for len(qc.cache) > defaultQueueCacheSize {
		for k := range qc.cache {
			delete(qc.cache, k)
			break
		}
	}

	verCache, ok := qc.cache[name]
	if !ok {
		verCache = make(versionCache)
		qc.cache[name] = verCache
	} else {
		if len(verCache) > defaultVersionCacheSize {
			keys := verCache.sortedKeys()
			for i := 0; i < defaultVersionCacheSize-len(verCache); i++ {
				key := keys[i]
				delete(verCache, key)
			}
		}
	}
	verCache[ver] = res
	return res, nil
}

func (qc *QueueCache) Reset() {
	qc.cacheMu.Lock()
	defer qc.cacheMu.Unlock()
	for k := range qc.cache {
		delete(qc.cache, k)
	}
	for k := range qc.requestMu {
		delete(qc.requestMu, k)
	}
}
