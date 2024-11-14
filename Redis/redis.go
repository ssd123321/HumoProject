package Redis

import (
	"Tasks/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var (
	NotFound = errors.New("redis: nil")
)

type Cache struct {
	client *redis.Client
}

type ServerRedis struct {
	Addr     string
	Password string
	DB       int
}

var ServerRed = ServerRedis{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func NewCache(addr string, password string, db int) Cache {
	cache := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return Cache{cache}
}
func (c *Cache) GetPerson(ctx context.Context) (*model.Person, error) {
	var personJSON string
	var person model.Person
	err := c.client.Get(ctx, fmt.Sprint(ctx.Value("id"))).Scan(&personJSON)
	if err != nil {
		log.Printf("GetPersonError - %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal([]byte(personJSON), &person)
	if err != nil {
		return nil, err
	}
	person.Cache = true
	return &person, nil
}
func (c *Cache) SetPerson(ctx context.Context, person *model.Person) error {
	payload, err := json.Marshal(person)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, fmt.Sprint(ctx.Value("id")), string(payload), time.Second*20).Err()
	if err != nil {
		log.Printf("SetPersonError - %v", err.Error())
		return err
	}
	return nil
}
func (c *Cache) DeletePerson(ctx context.Context) error {
	err := c.client.Del(ctx, fmt.Sprint(ctx.Value("id"))).Err()
	if err != nil {
		log.Printf("DeletePersonError - %v", err)
		return err
	}
	return nil
}
func (c *Cache) SetSlice(ctx context.Context, people []model.Person) error {
	for _, person := range people {
		payload, err := json.Marshal(person)
		if err != nil {
			log.Printf("SetSlice - json.Marshal: %v", err)
			return err
		}
		err = c.client.Set(ctx, fmt.Sprint(person.ID), string(payload), time.Minute*2).Err()
		if err != nil {
			log.Printf("SetSlice - c.client.Set: %v", err)
			return err
		}
	}
	return nil
}
