package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/felipefbs/testcontainers/connections/queue"
	"github.com/felipefbs/testcontainers/message"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	handler := http.NewServeMux()

	handler.HandleFunc("/save", saveHandler)

	server := http.Server{
		Addr:    ":" + os.Getenv("HTTP_SERVER_PORT"),
		Handler: handler,
	}

	slog.Info("http server listening on", "port", os.Getenv("HTTP_SERVER_PORT"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("invalid method request")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("this route only accepts post requests"))

		return
	}

	defer r.Body.Close()

	var msg message.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		slog.Error("failed to read message", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("malformed body message"))

		return
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to read message", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("malformed body message"))

		return
	}

	messageSender := queue.NewMessageSender(os.Getenv("KAFKA_SERVER_URL"))
	deliveryChan := make(chan kafka.Event)
	err = messageSender.SendMessage(r.Context(), os.Getenv("KAFKA_TOPIC"), msgBytes, nil, deliveryChan)
	if err != nil {
		slog.Error("failed to send message", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}
