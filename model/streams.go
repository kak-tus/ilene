package model

import (
	"sync"

	"github.com/go-redis/redis"
)

// Stream object
type Stream struct {
	Key    string `json:"key"`
	Length int64  `json:"length"`
}

// StreamsList object
type StreamsList struct {
	List []Stream `json:"list"`
}

// ListStreams list streams
func (m *Type) ListStreams() StreamsList {
	var list []Stream
	var mtx sync.Mutex

	err := m.rDB.ForEachMaster(func(c *redis.Client) error {
		iter := c.Scan(0, "qu*", 100).Iterator()
		if iter.Err() != nil {
			return iter.Err()
		}

		for iter.Next() {
			ctype := c.Type(iter.Val())
			if ctype.Err() != nil {
				return ctype.Err()
			}
			if ctype.Val() != "stream" {
				continue
			}

			clen := c.XLen(iter.Val())
			if clen.Err() != nil {
				return clen.Err()
			}

			mtx.Lock()
			list = append(list, Stream{
				Key:    iter.Val(),
				Length: clen.Val(),
			})
			mtx.Unlock()
		}

		return nil
	})

	if err != nil {
		m.log.Error(err)
		return StreamsList{}
	}

	return StreamsList{List: list}
}
