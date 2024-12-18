package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

// DataServer is the main struct for the WebSocket server
type DataServer struct {
	listenAddress string
	store         Store
	logger        *log.Logger
	conns         map[net.Conn]bool
}

func NewDataServer(listenAddr string, store Store) *DataServer {
	return &DataServer{
		listenAddress: listenAddr,
		store:         store,
		logger:        log.New(os.Stdout, "[websocket-server] ", log.LstdFlags),
		conns:         make(map[net.Conn]bool),
	}
}

func (s *DataServer) Run() error {
	router := http.NewServeMux()
	router.Handle("/", websocket.Handler(s.handleWSConn))

	go s.farmDataFetcher()
	s.logger.Printf("DataStream server is live on: %s\n", s.listenAddress)

	if err := http.ListenAndServe(s.listenAddress, router); err != nil {
		return err
	}

	return nil
}

func (s *DataServer) farmDataFetcher() {
	// Get latest data from the store
	latestData, err := s.store.FetchData(os.Getenv("FARM_DATA_ID"))
	dataLen := len(latestData.data)
	if err != nil {
		s.logger.Printf("failed to fetch data: %s %v\n", err, dataLen)
		return
	}

	for {
		// TODO:  Make api request to fetch data
		//  Check if the data is new
		//  If the data is new
		//  Store the data
		s.store.StoreData(os.Getenv("FARM_DATA_ID"), latestData.data)
		s.logger.Printf("new form submitted: %s\n", latestData.data)
		//  Broadcast the data
		s.broadcast(latestData.data)
		//  Sleep for 1min
		time.Sleep(1 * time.Minute)
	}
}

func (s *DataServer) handleWSConn(conn *websocket.Conn) {
	s.logger.Printf("incomming connection from: %+v\n", conn.RemoteAddr())

	// TODO: Validate the connection

	s.conns[conn] = true
}

func (s *DataServer) broadcast(msg []byte) {
	for conn := range s.conns {
		go func(conn net.Conn, msg []byte) {
			_, err := conn.Write(msg)
			if err != nil {
				s.logger.Printf("broadcast write error: %s\n", err)
				return
			}
		}(conn, msg)
	}
	s.logger.Printf("broadcasted farm data: %s\n", msg)
}

func (s *DataServer) Close() {
	for conn := range s.conns {
		conn.Close()
	}
}
