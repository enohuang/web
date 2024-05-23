package redis

import (
	"context"
	"dengming20240317/web/session"
	"errors"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session: 找不到 session")
)

type StoreOption func(store *Store)

type Store struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	//分两个命令有风险
	/*	_, err := s.client.HSet(ctx, id, id, id).Result()
		if err != nil {
			return nil, err
		}
		_, err = s.client.Expire(ctx, id, s.expiration).Result()
		if err != nil {
			return nil, err
		}
		return &Session{client: s.client, id: id}, nil*/

	//
	const lua = `
redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
return redis.call("pexpire", KEYS[1], ARGV[3])
`

	key := id
	_, err := s.client.Eval(ctx, lua, []string{key}, "_sess_id", id, s.expiration.Milliseconds()).Result()
	if err != nil {
		return nil, err
	}
	return &Session{client: s.client, id: id}, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	ok, err := s.client.Expire(ctx, id, s.expiration).Result()
	if err != nil {
		return nil
	}
	if !ok {
		return errors.New("session 不存在")
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {

	_, err := s.client.Del(ctx, id).Result()
	if err != nil {
		return err
	}
	return err
	/*//id 对应的session 不存在， 你没有删除任何东西
	if  cnt == 0{

	}*/

}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {

	cnt, err := s.client.Exists(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, ErrSessionNotFound
	}
	return &Session{client: s.client}, nil
}

func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	store := &Store{
		expiration: time.Minute * 15,
		client:     client,
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

func WithExpirationOption(expiration time.Duration) StoreOption {
	return func(store *Store) {
		store.expiration = expiration
	}
}

type Session struct {
	id     string
	client redis.Cmdable
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	cnt, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, ErrSessionNotFound
	}
	return &Session{
		client: s.client,
	}, nil
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else 
   return -1
end 
`
	res, err := s.client.Eval(ctx, lua, []string{s.id}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (s *Session) ID() string {
	return s.id
}

func key(prefix, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}
