package helpers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type hostData struct {
	Counter      uint64
	RangeStart   uint64
	RangeEnd     uint64
	RangeCurrent uint64
}

type getResponse struct {
	Status string
	Data   string
}

var (
	ETCDConn *clientv3.Client
)

func establishConnection() {
	var err error

	ETCDConn, err = clientv3.New(clientv3.Config{
		Endpoints:   GetEnvListString("ETCD_DSN"),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic("failed to connect etcd")
	}
}

func deferConnection() {
	defer ETCDConn.Close()
}

func getCounter(timeout time.Duration) *getResponse {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	tx, err := ETCDConn.Get(ctx, "counter")
	cancel()
	if err != nil {
		return &getResponse{
			Status: "error",
			Data:   err.Error(),
		}
	}

	if tx.Count == 0 {
		return &getResponse{
			Status: "error",
			Data:   "counter not found",
		}
	}
	return &getResponse{
		Status: "ok",
		Data:   string(tx.Kvs[0].Value),
	}
}

func getHostname(timeout time.Duration) *getResponse {
	hostname := os.Getenv("HOSTNAME")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	tx, err := ETCDConn.Get(ctx, hostname)
	cancel()
	if err != nil {
		return &getResponse{
			Status: "error",
			Data:   err.Error(),
		}
	}

	if tx.Count == 0 {
		return &getResponse{
			Status: "error",
			Data:   "hostname not found",
		}
	}

	return &getResponse{
		Status: "ok",
		Data:   string(tx.Kvs[0].Value),
	}
}

func putCounter(timeout time.Duration, counter string) *getResponse {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_, err := ETCDConn.Put(ctx, "counter", counter)
	cancel()
	if err != nil {
		return &getResponse{
			Status: "error",
			Data:   err.Error(),
		}
	}
	return &getResponse{
		Status: "ok",
		Data:   "ok",
	}
}

func putHostname(timeout time.Duration, hostname string, Data hostData) *getResponse {
	stringData, errData := json.Marshal(Data)
	if errData != nil {
		return &getResponse{
			Status: "error",
			Data:   errData.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_, err := ETCDConn.Put(ctx, hostname, string(stringData))
	cancel()

	if err != nil {
		return &getResponse{
			Status: "error",
			Data:   err.Error(),
		}
	}

	return &getResponse{
		Status: "ok",
		Data:   "ok",
	}
}

func InitEtcd() *getResponse {
	establishConnection()
	checkCounter := getCounter(5 * time.Second)
	if checkCounter.Status == "error" {
		putCounter(5*time.Second, "0")
	}
	checkCounter = getCounter(5 * time.Second)
	checkHostname := getHostname(5 * time.Second)
	if checkHostname.Status == "error" {
		counter, _ := strconv.ParseUint(checkCounter.Data, 10, 64)
		etcdRange := GetEnvInt("ETCD_RANGE")
		putHostname(5*time.Second, os.Getenv("HOSTNAME"), hostData{
			Counter:      counter,
			RangeStart:   counter * uint64(etcdRange),
			RangeEnd:     counter*uint64(etcdRange) + uint64(etcdRange-1),
			RangeCurrent: counter * uint64(etcdRange),
		})
	}
	deferConnection()
	return &getResponse{
		Status: "ok",
		Data:   "ok",
	}
}

func deleteHostname() *getResponse {
	establishConnection()
	_, err := ETCDConn.Delete(context.Background(), os.Getenv("HOSTNAME"))
	if err != nil {
		return &getResponse{
			Status: "error",
			Data:   err.Error(),
		}
	}
	deferConnection()
	return &getResponse{
		Status: "ok",
		Data:   "ok",
	}
}

func deleteCounter() *getResponse {
	establishConnection()
	_, err := ETCDConn.Delete(context.Background(), "counter")
	if err != nil {
		return &getResponse{
			Status: "error",
			Data:   err.Error(),
		}
	}
	deferConnection()
	return &getResponse{
		Status: "ok",
		Data:   "ok",
	}
}

func checkHostnameRangeEnd() *getResponse {
	establishConnection()
	checkHostname := getHostname(5 * time.Second)
	if checkHostname.Status == "error" {
		return &getResponse{
			Status: "error",
			Data:   checkHostname.Data,
		}
	}
	var Data hostData
	json.Unmarshal([]byte(checkHostname.Data), &Data)
	if Data.RangeCurrent >= Data.RangeEnd {
		deleteHostname()
		checkCounter := getCounter(5 * time.Second)
		counter, _ := strconv.ParseUint(checkCounter.Data, 10, 64)
		counter++
		log.Println(counter)
		putCounter(5*time.Second, strconv.FormatUint(counter, 10))
		etcdRange := GetEnvInt("ETCD_RANGE")
		putHostname(5*time.Second, os.Getenv("HOSTNAME"), hostData{
			Counter:      counter,
			RangeStart:   counter * uint64(etcdRange),
			RangeEnd:     counter*uint64(etcdRange) + uint64(etcdRange-1),
			RangeCurrent: counter * uint64(etcdRange),
		})
	}

	deferConnection()
	return &getResponse{
		Status: "ok",
		Data:   checkHostname.Data,
	}
}

func GetCurrentCounter() *getResponse {
	establishConnection()
	// checkHostnameRangeEnd()
	checkHostname := getHostname(5 * time.Second)
	if checkHostname.Status == "error" {
		deferConnection()
		return &getResponse{
			Status: "error",
			Data:   "hostname not found",
		}
	}
	var Data hostData
	json.Unmarshal([]byte(checkHostname.Data), &Data)
	Data.RangeCurrent++
	if Data.RangeCurrent >= Data.RangeEnd {
		checkCounter := getCounter(5 * time.Second)
		counter, _ := strconv.ParseUint(checkCounter.Data, 10, 64)
		counter++
		putCounter(5*time.Second, strconv.FormatUint(counter, 10))
		etcdRange := GetEnvInt("ETCD_RANGE")
		var newData hostData
		newData.Counter = counter
		newData.RangeStart = counter * uint64(etcdRange)
		newData.RangeEnd = counter*uint64(etcdRange) + uint64(etcdRange-1)
		newData.RangeCurrent = counter * uint64(etcdRange)
		putHostname(5*time.Second, os.Getenv("HOSTNAME"), newData)
		Data = newData
	}

	putHostname(5*time.Second, os.Getenv("HOSTNAME"), Data)
	deferConnection()
	return &getResponse{
		Status: "ok",
		Data:   strconv.FormatUint(Data.RangeCurrent, 10),
	}
}
