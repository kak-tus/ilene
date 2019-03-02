package model

import (
	"strings"
	"sync"

	"github.com/go-redis/redis"
)

// Stream object
type Stream struct {
	Key    string `json:"key"`
	Length int64  `json:"length"`
	Type   string `json:"type"`
}

// StreamsList object
type StreamsList struct {
	List []Stream `json:"list"`
}

// ListStreams list streams
func (m *Type) ListStreams() StreamsList {
	var mtx sync.Mutex
	sharded := make(map[string]int64)
	full := make(map[string]int64)

	err := m.rDB.ForEachMaster(func(c *redis.Client) error {
		iter := c.Scan(0, "qu*", 100).Iterator()
		if iter.Err() != nil {
			return iter.Err()
		}

		for iter.Next() {
			key := iter.Val()

			ctype := c.Type(key)
			if ctype.Err() != nil {
				return ctype.Err()
			}
			if ctype.Val() != "stream" {
				continue
			}

			clen := c.XLen(key)
			if clen.Err() != nil {
				return clen.Err()
			}

			mtx.Lock()

			pos1 := strings.Index(key, "{")
			pos2 := strings.LastIndex(key, "}")

			if pos1 >= 0 && pos2 >= 0 && pos1 < pos2 {
				fKey := key[0:pos1] + key[pos2+1:]
				full[fKey] += clen.Val()
			}

			sharded[key] += clen.Val()

			mtx.Unlock()
		}

		return nil
	})

	var list StreamsList

	if err != nil {
		m.log.Error(err)
		return list
	}

	for k, v := range sharded {
		list.List = append(list.List, Stream{
			Key:    k,
			Length: v,
			Type:   "sharded",
		})
	}
	for k, v := range full {
		list.List = append(list.List, Stream{
			Key:    k,
			Length: v,
			Type:   "full",
		})
	}

	return list
}
