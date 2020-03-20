package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/andy1219111/influx-client/res"

	influx "github.com/influxdata/influxdb1-client"
)

//InfluxClient influxDB客户端
type InfluxClient struct {
	host     string
	port     string
	user     string
	password string
	client   *influx.Client
}

// InfluxPoint 插入的记录结构体
type InfluxPoint = influx.Point

//NewInfluxClient 得到influxDB客户端对象
func NewInfluxClient(host, port, user, password string) (*InfluxClient, error) {

	client := &InfluxClient{}
	client.host = host
	client.port = port
	client.user = user
	client.password = password

	url, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))
	if err != nil {
		return client, nil
	}
	conf := influx.Config{
		URL:      *url,
		Username: user,
		Password: password,
	}

	//创建influx client对象
	influxDB, err := influx.NewClient(conf)
	if err != nil {
		return client, err
	}
	client.client = influxDB

	return client, nil
}

//Ping ping方法
func (c *InfluxClient) Ping() (dur time.Duration, version string, err error) {
	if c.client == nil {
		err = errors.New("the influx client is nil")
		return
	}
	dur, version, err = c.client.Ping()
	if err != nil {
		return
	}
	return
}

//Query 执行influxDB查询
func (c *InfluxClient) Query(influxQL, database string) (*influx.Response, error) {
	if c.client == nil {
		err := errors.New("the influx client is nil")
		return nil, err
	}

	q := influx.Query{
		Command:  influxQL,
		Database: database,
	}
	response, err := c.client.Query(q)
	if err != nil {
		return response, err
	}

	if err = response.Error(); err != nil {
		return response, err
	}

	return response, nil
}

// QueryMap 将查询结果以map方式返回
func (c *InfluxClient) QueryMap(influxQL, database string) ([]map[string]interface{}, error) {
	//初始化存储查询结果的map数组
	results := make([]map[string]interface{}, 10)
	results = results[0:0]

	response, err := c.Query(influxQL, database)
	if err != nil {
		return results, err
	}
	resJSON, err := response.MarshalJSON()
	if err != nil {
		return results, err
	}

	res := &res.Res{}
	err = json.Unmarshal([]byte(resJSON), res)
	if err != nil {
		return results, err
	}
	//未查询到记录
	if len(res.Results) == 0 {
		results, nil
	}
	for _, row := range res.Results[0].Series[0].Values {
		rowMap := make(map[string]interface{})
		for key, column := range res.Results[0].Series[0].Columns {
			rowMap[column] = row[key]
		}
		results = append(results, rowMap)
	}

	return results, nil
}

//Insert  插入数据
func (c *InfluxClient) Insert(data []InfluxPoint, database string, retentionPolicy string) error {
	bps := influx.BatchPoints{
		Points:          data,
		Database:        database,
		RetentionPolicy: retentionPolicy,
	}
	_, err := c.client.Write(bps)
	if err != nil {
		return err
	}
	return nil
}

/*
func main() {

	influxClient, err := NewInfluxClient("192.168.1.150", "8086", "work", "123456")

	dur, version, err := influxClient.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("ping success %v,%s \n", dur, version)

	influxQl := " select * from device_flow"
	res, err := influxClient.QueryMap(influxQl, "dt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("record: %+v,%d", res, len(res))
}
*/
