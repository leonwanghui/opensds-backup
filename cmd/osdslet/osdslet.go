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
This module implements a entry into the OpenSDS REST service.

*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/opensds/opensds/pkg/api"
)

var (
	edp            string
	osdsletLogFile string
)

func init() {
	flag.StringVar(&edp, "endpoint", "127.0.0.1:50040", "Listen endpoint of controller service")
	flag.StringVar(&osdsletLogFile, "osdsletlog-file", "/var/log/opensds/osdslet.log", "Location of osdslet log file")
	flag.Parse()
}

func main() {
	// Open OpenSDS orchestrator service log file
	f, err := os.OpenFile(osdsletLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Assign it to the standard logger
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Start OpenSDS northbound REST service.
	api.Run(edp)
}
