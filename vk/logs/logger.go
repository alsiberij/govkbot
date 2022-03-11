package logs

import (
	"log"
	"os"
	"time"
)

func Log(rs []byte) {
	f, err := os.Create("vk/logs/response" + time.Now().Format("2006-01-02--15-04-05") + ".json")
	if err != nil {
		log.Print("Error creating log file")
	}
	_, err = f.Write(rs)
	if err != nil {
		log.Print("Error writing log file")
	}
	err = f.Close()
	if err != nil {
		log.Print("Error closing log file")
	}
}
