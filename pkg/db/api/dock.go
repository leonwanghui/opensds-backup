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
This module implements the database operation interface of data structure
defined in api module.

*/

package api

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/opensds/opensds/pkg/db"
	api "github.com/opensds/opensds/pkg/model"
)

func CreateDock(dck *api.DockSpec) (*api.DockSpec, error) {
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return &api.DockSpec{}, err
	}

	dbReq := &db.Request{
		Url:     db.URL_PREFIX + "docks/" + dck.Id,
		Content: string(dckBody),
	}
	dbRes := db.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}

	return dck, nil
}

func GetDock(dckID string) (*api.DockSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "docks/" + dckID,
	}
	dbRes := db.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), dck); err != nil {
		log.Println("[Error] When parsing dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}
	return dck, nil
}

func ListDocks() (*[]api.DockSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "docks",
	}
	dbRes := db.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list docks in db:", dbRes.Error)
		return &[]api.DockSpec{}, errors.New(dbRes.Error)
	}

	var dcks = []api.DockSpec{}
	if len(dbRes.Message) == 0 {
		return &dcks, nil
	}
	for _, msg := range dbRes.Message {
		var dck = api.DockSpec{}
		if err := json.Unmarshal([]byte(msg), &dck); err != nil {
			log.Println("[Error] When parsing dock in db:", dbRes.Error)
			return &[]api.DockSpec{}, errors.New(dbRes.Error)
		}
		dcks = append(dcks, dck)
	}
	return &dcks, nil
}

func UpdateDock(dckID, name, desp string) (*api.DockSpec, error) {
	dck, err := GetDock(dckID)
	if err != nil {
		return &api.DockSpec{}, err
	}
	if name != "" {
		dck.Name = name
	}
	if desp != "" {
		dck.Description = desp
	}
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return &api.DockSpec{}, err
	}

	dbReq := &db.Request{
		Url:        db.URL_PREFIX + "docks/" + dckID,
		NewContent: string(dckBody),
	}
	dbRes := db.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}
	return dck, nil
}

func DeleteDock(dckID string) error {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "docks/" + dckID,
	}
	dbRes := db.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
