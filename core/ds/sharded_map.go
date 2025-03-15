package ds

import (
	"hash/fnv"
	"sync"
)

type ShardedMap struct {
	shards []sync.Map
}

func InitShardedMap(shardCount int) *ShardedMap {
	shards := make([]sync.Map, shardCount)
	return &ShardedMap{shards: shards}
}

func (s *ShardedMap) getShards(key string) *sync.Map {
	hash := fnv.New64a()
	hash.Write([]byte(key))
	index := hash.Sum64() % uint64(len(s.shards))
	return &s.shards[index]
}

func (s *ShardedMap) Set(key string, value interface{}) {
	shard := s.getShards(key)
	shard.Store(key, value)
}

func (s *ShardedMap) Get(key string) (interface{}, bool) {
	shard := s.getShards(key)
	return shard.Load(key)
}
