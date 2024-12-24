package main

import (
	"Tasks/Redis"
	"Tasks/Service"
	"Tasks/TelegramBotAPI"
	"Tasks/handlers/grpc"
	"Tasks/handlers/grpc/gprc_api"
	http2 "Tasks/handlers/http"
	"Tasks/handlers/http/middleware"
	"Tasks/repository"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	grpc2 "google.golang.org/grpc"
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
	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	rep := repository.Repository{DB: db}
	cache := Redis.NewCache(Redis.ServerRed.Addr, Redis.ServerRed.Password, Redis.ServerRed.DB)
	service := Service.CreateService(&rep, &cache)
	serverRegistr := grpc2.NewServer()
	myGrpc := grpc.NewGrpcBaningServer(&gprc_api.UnimplementedGrpcBankingServer{}, service)
	net, err := grpc.NewGrpcServer(":8089")
	if err != nil {
		log.Fatal(err)
	}
	gprc_api.RegisterGrpcBankingServer(serverRegistr, myGrpc)
	go func() {
		fmt.Println("Starting grpc server on port :8089")
		err = serverRegistr.Serve(net)
		if err != nil {
			log.Fatal(err)
		}
	}()
	router := mux.NewRouter()
	mux.CORSMethodMiddleware(router)
	handler := http2.NewHandler(service, router)
	handler.Router.Use(middleware.Recovery)
	handler.Router.HandleFunc("/login", handler.Login).Methods("POST")
	handler.Router.HandleFunc("/AddPerson", handler.SignUp).Methods("POST")
	handler.Router.Handle("/GetPeople", middleware.Authentication(http.HandlerFunc(handler.GetPeople))).Methods("GET")
	handler.Router.Handle("/DeletePerson/{id}", middleware.Authentication(http.HandlerFunc(handler.DeletePerson))).Methods("DELETE")
	handler.Router.Handle("/GetPerson/{id}", middleware.Authentication(http.HandlerFunc(handler.GetPatientByID))).Methods("GET")
	handler.Router.Handle("/UpdatePerson", middleware.Authentication(http.HandlerFunc(handler.UpdatePerson))).Methods("PUT")
	handler.Router.Handle("/GetPersonInFile/{id}", middleware.Authentication(http.HandlerFunc(handler.GetPersonInFile))).Methods("GET")
	handler.Router.Handle("/AddPersonInFile", middleware.Authentication(http.HandlerFunc(handler.AddPeopleFromFile))).Methods("POST")
	handler.Router.Handle("/AddCard", middleware.Authentication(http.HandlerFunc(handler.AddCard))).Methods("POST")
	handler.Router.Handle("/Transfer", middleware.Authentication(http.HandlerFunc(handler.TransferMoney))).Methods("POST")
	handler.Router.Handle("/change", middleware.Authentication(http.HandlerFunc(handler.ChangePassword))).Methods("POST")
	b, err := TelegramBotAPI.InitializeBot("8193972383:AAF470KuDnyRSi6EZUpaNgnSFjhumE480YY", service, make(map[int64]int), make(map[int64]int), make(map[int64]string))
	if err != nil {
		log.Fatal(err)
	}
	updates, err := b.GetUpdates()
	if err != nil {
		log.Fatal(err)
	}
	go b.HandlerPollingData(updates)
	fmt.Printf("Server start on port %d\n", 8080)
	err = http.ListenAndServe("localhost:8080", handler.Router)
	if err != nil {
		log.Fatal(err)
	}
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
