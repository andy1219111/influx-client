package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"influx-client/res"
	"net/url"
	"time"

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
func (c *InfluxClient) Query(influxQL, database string) (*res.Res, error) {
	if c.client == nil {
		err := errors.New("the influx client is nil")
		return nil, err
	}

	res := &res.Res{}
	q := influx.Query{
		Command:  influxQL,
		Database: database,
	}
	response, err := c.client.Query(q)
	if err != nil {
		return res, err
	}
	if err = response.Error(); err != nil {
		return res, err
	}

	resJSON, err := response.MarshalJSON()
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(resJSON), res)
	if err != nil {
		return res, err
	}
	return res, nil
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
