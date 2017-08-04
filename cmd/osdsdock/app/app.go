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
This module implements the entry into operations of storageDock module.

*/

package app

import (
	"errors"
	"log"

	"github.com/opensds/opensds/cmd/osdsdock/app/dock"
	"github.com/opensds/opensds/cmd/osdsdock/app/pool"
	db "github.com/opensds/opensds/pkg/db/api"
)

func ResourceDiscovery() error {
	if err := dockResourceDiscovery(); err != nil {
		log.Println("[Error] When discover dock resource:", err)
		return err
	}
	if err := poolResourceDiscovery(); err != nil {
		log.Println("[Error] When discover pool resource:", err)
		return err
	}
	return nil
}

func dockResourceDiscovery() error {
	dcks, err := dock.ListDocks()
	if err != nil {
		log.Println("[Error] When list docks:", err)
		return err
	}
	if len(*dcks) == 0 {
		return errors.New("The dock resource is empty")
	}

	for _, dck := range *dcks {
		if _, err = db.CreateDock(&dck); err != nil {
			log.Printf("[Error] When create dock %s in db: %v\n", dck.Id, err)
			return err
		}
	}
	return nil
}

func poolResourceDiscovery() error {
	pols, err := pool.ListPools()
	if err != nil {
		log.Println("[Error] When list pools:", err)
		return err
	}
	if len(*pols) == 0 {
		return errors.New("The pool resource is empty")
	}

	for _, pol := range *pols {
		if _, err = db.CreatePool(&pol); err != nil {
			log.Printf("[Error] When create pool %s in db: %v\n", pol.Id, err)
			return err
		}
	}
	return nil
}
