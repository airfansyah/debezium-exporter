package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "strconv"
  "time"
  "os"
  "strings"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
  debeziumConnectorState = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "debezium_connector_state",
      Help: "Debezium connector state",
    },
    []string{"app", "connector_name"},
  )
)

var (
  debeziumTaskState = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "debezium_task_state",
      Help: "Debezium task state",
    },
    []string{"app", "connector_name", "task_id"},
  )
)

func init(){
  prometheus.MustRegister(debeziumConnectorState)
  prometheus.MustRegister(debeziumTaskState)
}

type Status struct {
  Name string `json:"name"`
  Connector struct {
    State string `json:"state"`
    WorkerID string `json:"worker_id"`
  } `json:"connector"`
  Tasks []struct {
    ID int `json:"id"`
    State string `json:"state"`
    WorkerID string `json:"worker_id"`
  } `json:"tasks"`
  Type string `json:"type"`
}

func getMetrics() {
  urls := os.Getenv("DEBEZIUM_URL")
  urlList := strings.Split(urls, ",")

  for _, url:= range urlList {
      app := strings.Split(url, ".")[0]

      
      if url == "" {
        fmt.Println("no url")
        return
      }
      response, err := http.Get("http://" + url + "/connectors?expand=status")
      if err != nil {
        fmt.Println("can not connect to",url)
        return
      }
      defer response.Body.Close()
      body, _ := ioutil.ReadAll(response.Body)
      
      var data map[string]struct {
        Status `json:"status"`
      }
      
      json.Unmarshal(body,&data)
      
      for key, value := range data {
        conStatNum := 0.0
        if value.Connector.State == "RUNNING" {
          conStatNum = 1.0
        }
        debeziumConnectorState.WithLabelValues(app, key).Set(conStatNum)
        for _, task := range value.Tasks {
          taskID := strconv.Itoa(task.ID)
          taskStatNum := 0.0
          if task.State == "RUNNING" {
            taskStatNum = 1.0
          }
          debeziumTaskState.WithLabelValues(app, key, taskID).Set(taskStatNum)
          
        }
      }
  }

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "ok")
}


func main(){
  timer := time.NewTicker(5* time.Second)
  defer timer.Stop()
  
  go func() {
    for {
      select {
      case <-timer.C:
        getMetrics()
      }
    }
  }()
  http.HandleFunc("/",rootHandler)
  http.Handle("/metrics", promhttp.Handler())
  port := 9100
  addr := fmt.Sprintf(":%d", port)
  fmt.Printf("Exporter is listening on port %d...\n", port)
  http.ListenAndServe(addr, nil)
}
