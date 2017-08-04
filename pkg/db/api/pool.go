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

func CreatePool(pol *api.StoragePoolSpec) (*api.StoragePoolSpec, error) {
	polBody, err := json.Marshal(pol)
	if err != nil {
		return &api.StoragePoolSpec{}, err
	}

	dbReq := &db.Request{
		Url:     db.URL_PREFIX + "pools/" + pol.Id,
		Content: string(polBody),
	}
	dbRes := db.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create pol in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}

	return pol, nil
}

func GetPool(polID string) (*api.StoragePoolSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "pools/" + polID,
	}
	dbRes := db.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get pool in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}

	var pol = &api.StoragePoolSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), pol); err != nil {
		log.Println("[Error] When parsing pool in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}
	return pol, nil
}

func ListPools() (*[]api.StoragePoolSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "pools",
	}
	dbRes := db.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list pools in db:", dbRes.Error)
		return &[]api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}

	var pols = []api.StoragePoolSpec{}
	if len(dbRes.Message) == 0 {
		return &pols, nil
	}
	for _, msg := range dbRes.Message {
		var pol = api.StoragePoolSpec{}
		if err := json.Unmarshal([]byte(msg), &pol); err != nil {
			log.Println("[Error] When parsing pool in db:", dbRes.Error)
			return &[]api.StoragePoolSpec{}, errors.New(dbRes.Error)
		}
		pols = append(pols, pol)
	}
	return &pols, nil
}

func UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*api.StoragePoolSpec, error) {
	pol, err := GetPool(polID)
	if err != nil {
		return &api.StoragePoolSpec{}, err
	}
	if name != "" {
		pol.Name = name
	}
	if desp != "" {
		pol.Description = desp
	}
	polBody, err := json.Marshal(pol)
	if err != nil {
		return &api.StoragePoolSpec{}, err
	}

	dbReq := &db.Request{
		Url:        db.URL_PREFIX + "pools/" + polID,
		NewContent: string(polBody),
	}
	dbRes := db.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update pool in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}
	return pol, nil
}

func DeletePool(polID string) error {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "pools/" + polID,
	}
	dbRes := db.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete pool in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
