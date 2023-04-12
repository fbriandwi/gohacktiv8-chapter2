package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "strings"
    "time"
)

const (
    url = "https://jsonplaceholder.typicode.com/posts"
)

type data struct {
    Water  int    `json:"water"`
    Wind   int    `json:"wind"`
    Status string `json:"status"`
}

func main() {
    rand.Seed(time.Now().UnixNano())

    for {
        // generate random values for water and wind
        water := rand.Intn(100) + 1
        wind := rand.Intn(100) + 1
        // set status based on water and wind values
        var status string
        if water > 50 && wind > 50 {
            status = "Hujan dan Angin Kencang"
        } else if water > 50 {
            status = "Hujan"
        } else if wind > 50 {
            status = "Angin Kencang"
        } else {
            status = "Normal"
        }
        // create data struct
        data := &data{
            Water:  water,
            Wind:   wind,
            Status: status,
        }
        // convert data to JSON format
        jsonData, err := json.Marshal(data)
        if err != nil {
            fmt.Println(err)
            continue
        }
        // send POST request with JSON data
        resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
        if err != nil {
            fmt.Println(err)
            continue
        }
        // print response
        fmt.Println(resp.Status)
        // wait for 15 seconds before sending next request
        time.Sleep(15 * time.Second)
    }
}
