// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements the database operation of data structure
defined in api module.

*/

package db

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/clientv3"
)

var (
	c          = &Client{}
	URL_PREFIX = "/api/v1alpha1/block/"
	TIME_OUT   = 3 * time.Second
)

func init() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:2380"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		cli.Close()
		fmt.Errorf(err.Error())
	}

	c.client = cli
}

type Request struct {
	Url        string `json:"url"`
	Content    string `json:"content"`
	NewContent string `json:"newContent"`
}

type Response struct {
	Status  string   `json:"status"`
	Message []string `json:"message"`
	Error   string   `json:"error"`
}

type Client struct {
	client *clientv3.Client
	lock   sync.Mutex
}

func Create(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.client.Put(ctx, req.Url, req.Content)
	if err != nil {
		log.Println("[Error] When create db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &Response{
		Status:  "Success",
		Message: []string{req.Content},
	}
}

func Get(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	resp, err := c.client.Get(ctx, req.Url)
	if err != nil {
		log.Println("[Error] When get db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	if len(resp.Kvs) == 0 {
		return &Response{
			Status: "Failure",
			Error:  "Wrong volume_id or attachment_id provided!",
		}
	}
	return &Response{
		Status:  "Success",
		Message: []string{string(resp.Kvs[0].Value)},
	}
}

func List(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	resp, err := c.client.Get(ctx, req.Url, clientv3.WithPrefix())
	if err != nil {
		log.Println("[Error] When get db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	var message = []string{}
	for _, v := range resp.Kvs {
		message = append(message, string(v.Value))
	}
	return &Response{
		Status:  "Success",
		Message: message,
	}
}

func Update(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.client.Put(ctx, req.Url, req.NewContent)
	if err != nil {
		log.Println("[Error] When update db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &Response{
		Status:  "Success",
		Message: []string{req.NewContent},
	}
}

func Delete(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.client.Delete(ctx, req.Url)
	if err != nil {
		log.Println("[Error] When delete db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &Response{
		Status: "Success",
	}
}
