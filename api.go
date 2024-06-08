package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"projects/gobank/storage"
	"projects/gobank/types"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
  listenAddr string
  store storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
  return &APIServer{
    listenAddr: listenAddr,
    store: store,
  }
}

type ApiError struct {
  Error string `json:"error"`
}

func (s *APIServer) Run() {
  router := mux.NewRouter()

  router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
  router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountById))
  router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransferAccount))

  log.Println("JSON API server running on port", s.listenAddr)
  
  http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
  if r.Method == "GET" {
    return s.handleGetAccount(w, r)
  }

  if r.Method == "POST" {
    return s.handleCreateAccount(w, r)
  }

  return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
  if r.Method == "GET" {
    return s.handleGetAccountById(w, r)
  }

  if r.Method == "DELETE" {
    return s.handleDeleteAccount(w, r)
  }

  return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
  accounts, err := s.store.GetAccounts()
  if err != nil {
    return nil
  }

  return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
  id, err := getRequestId(r)
  if err != nil {
    return err
  }

  account, err := s.store.GetAccountById(id)
  if err != nil {
    return err
  }

  return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
  createAccountRequest := new(types.CreateAccountRequest)
  if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
    return err
  }

  account := types.NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
  if err := s.store.CreateAccount(account); err != nil {
    return err
  }

  return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
  id, err := getRequestId(r)
  if err != nil {
    return err
  }

  err = s.store.DeleteAccount(id)
  if err != nil {
    return err
  }
  
  return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
  if r.Method != "POST" {
    return fmt.Errorf("method not allowed %s", r.Method)
  }

  transferReq := new(types.TransferRequest)

  err := json.NewDecoder(r.Body).Decode(transferReq)
  if err != nil {
    return err
  }
  defer r.Body.Close()

  return WriteJSON(w, http.StatusOK, transferReq)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if err := f(w, r); err != nil {
      WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
    }
  }
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(status)
  return json.NewEncoder(w).Encode(v)
}

func getRequestId(r *http.Request) (int, error) {
  idStr := mux.Vars(r)["id"]

  id, err := strconv.Atoi(idStr)
  if err != nil {
    return id, fmt.Errorf("invalid id %s", idStr)
  }

  return id, nil
}

