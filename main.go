package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/joho/godotenv"
)

const envFile = ".env"
const dataFile = "data/forms.json"

var loadEnv = env.Load

type formInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (f formInput) validate() error {
	if f.FirstName == "" || f.LastName == "" || f.Email == "" || f.PhoneNumber == "" {
		return errors.New("invalid input")
	}
	return nil
}

func (f formInput) save() error {
	file, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return err
	}
	var forms []formInput
	err = json.Unmarshal(file, &forms)
	if err != nil {
		return err
	}
	//FIXED:appends new data instead of replacing it
	forms = append(forms, f)
	toSave, err := json.Marshal(forms)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dataFile, toSave, os.ModeAppend)
	return err
}

func handleFunc(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := req.ParseForm()
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(resp, err.Error())
			return
		}
		//FIXED:map values from form to model
		f := formInput{
			FirstName:   req.FormValue("first_name"),
			LastName:    req.FormValue("last_name"),
			PhoneNumber: req.FormValue("phone_number"),
			Email:       req.FormValue("email"),
		}
		err = f.validate()
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(resp, err.Error())
			return
		}
		err = f.save()
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(resp, err.Error())
			return
		}
		//FIXED:changed to 201 to follow REST principles
		resp.WriteHeader(http.StatusCreated)
		fmt.Fprint(resp, "form saved")
	case http.MethodGet:
		//FIXED:returns all form data
		file, err := ioutil.ReadFile(dataFile)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(resp, err.Error())
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(file)
		log.Println("data returned")
	default:
		log.Println("error no 404")
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprint(resp, "not found")
	}
}

func run() (s *http.Server) {
	err := loadEnv(envFile)
	if err != nil {
		log.Fatal(err)
	}
	port, exist := os.LookupEnv("PORT")
	//FIXED:check whether port is empty or not
	if !exist || port == "" {
		log.Fatal("no port specified")
	}
	port = fmt.Sprintf(":%s", port)

	//FIXED:create new file to prevent data from not being written
	//if the file was not present
	ioutil.WriteFile(dataFile, []byte("[]"), os.ModeAppend)
	log.Println("created new empty data file")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/inputs", handleFunc)

	s = &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux, //FIXED: attach handler to server
	}

	go func() {
		log.Printf("starting server on port %s", port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return
}

func main() {
	s := run()
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown")
	}
	log.Println("Server exiting")
}
