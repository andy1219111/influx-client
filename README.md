# influx-client
influxDB读写库

## 使用示例

```go
package main

import (
	"fmt"
	"log"
	influx "test_go/lib/influx-client"
	"test_go/params"
	"test_go/utils"
)

func main() {

	params.INIParser = &utils.IniParser{}
	err := params.INIParser.Load("./conf/conf.ini")
	if err != nil {
		log.Println("load the config file failed", err)
	}

	influxClient, err := influx.NewInfluxClient(
		"127.0.0.1",
		"8086",
		"username",
		"password")

	dur, version, err := influxClient.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("ping success %v,%s \n", dur, version)

	influxQl := "select * from device_flow"
	res, err := influxClient.Query(influxQl, "dt")
	if err != nil {
		fmt.Println(err)
	}
	resString, _ := res.MarshalJSON()
	fmt.Printf("record: %s", string(resString))

	points := []influx.InfluxPoint{
		influx.InfluxPoint{
			Raw:       "device_flow,device_id=88,port_id=12 down_speed=55.00,down_sum=863541.00,up_speed=60.00,up_sum=863000.00",
			Precision: "s",
		},
		influx.InfluxPoint{
			Raw:       "device_flow,device_id=88,port_id=11 down_speed=55.00,down_sum=863541.00,up_speed=60.00,up_sum=863000.00",
			Precision: "s",
		},
	}

	err = influxClient.Insert(points, "dt", "autogen")
	if err != nil {
		fmt.Println(err)
	}
}


```
