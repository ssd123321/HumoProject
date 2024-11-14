package main

import (
	"Tasks/Redis"
	"Tasks/Service"
	"Tasks/handlers"
	"Tasks/handlers/middleware"
	"Tasks/repository"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	/*
		var (
			resultCh    = make(chan string)
			ctx, cancel = context.WithCancel(context.Background())
			services    = []string{"Google", "Yandex", "Yahoo", "Baidu"}
			wg          sync.WaitGroup
			winner      string
		)
		defer cancel()
		for i := range services {
			svc := services[i]
			wg.Add(1)
			go func() {
				requestRide(ctx, svc, resultCh)
				wg.Done()
			}()
			go func() {
				winner = <-resultCh
				cancel()
			}()
		}
		wg.Wait()
		log.Printf("found car in %q", winner)
	*/

	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	rep := repository.Repository{DB: db}
	cache := Redis.NewCache(Redis.ServerRed.Addr, Redis.ServerRed.Password, Redis.ServerRed.DB)
	service := Service.CreateService(&rep, &cache)
	router := mux.NewRouter()
	mux.CORSMethodMiddleware(router)
	handler := handlers.NewHandler(service, router)
	handler.Router.Use(middleware.Recovery)
	handler.Router.Use(middleware.SetLimit)
	handler.Router.HandleFunc("/AddPerson", handler.AddPerson).Methods("POST")
	handler.Router.HandleFunc("/DeletePerson/{id}", handler.DeletePerson).Methods("DELETE")
	handler.Router.HandleFunc("/GetPerson/{id}", handler.GetPatientByID).Methods("GET")
	handler.Router.HandleFunc("/UpdatePerson", handler.UpdatePerson).Methods("PUT")
	handler.Router.HandleFunc("/GetPersonInFile/{id}", handler.GetPersonInFile).Methods("GET")
	handler.Router.HandleFunc("/AddPersonInFile", handler.AddPeopleFromFile).Methods("POST")
	fmt.Printf("Server start on port %d\n", 8080)
	err = http.ListenAndServe("localhost:8080", handler.Router)
	if err != nil {
		log.Fatal(err)
	}

	/*
		ctx, cancel := context.WithCancel(context.Background())

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("ceded") // exit properly on cancellation
				default:
					// do work
				}
			}
		}(ctx)

		time.Sleep(1 * time.Second)

		cancel()
	*/
}
func cancel() {
	_, cancel := context.WithCancel(context.Background())
	cancel() // cancel the context
}
func OpenDB() (*gorm.DB, error) {
	dsn := "host=localhost user=humo password=humo dbname=postgres port=5432 sslmode=disable"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func requestRide(ctx context.Context, serviceName string, resultCh chan string) {
	time.Sleep(time.Second * 3)
	for {
		select {
		case <-ctx.Done():
			log.Printf("stopped the search in %q (%v)", serviceName, ctx.Err())
			return
		default:
			if rand.Float64() > 0.75 {
				fmt.Println("1")
				resultCh <- serviceName
				return
			}
			continue
		}
	}
}
func ProcessData(ctx context.Context) string {
	id := ctx.Value(0)
	select {
	case <-time.After(time.Second * 4):
		return fmt.Sprint("done processing id: ", id)
	case <-ctx.Done():
		return "Canceled"
	}
}
