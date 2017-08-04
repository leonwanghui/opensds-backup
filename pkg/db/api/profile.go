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

func CreateProfile(prf *api.ProfileSpec) (*api.ProfileSpec, error) {
	prfBody, err := json.Marshal(prf)
	if err != nil {
		return &api.ProfileSpec{}, err
	}

	dbReq := &db.Request{
		Url:     db.URL_PREFIX + "profiles/" + prf.Id,
		Content: string(prfBody),
	}
	dbRes := db.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}

	return prf, nil
}

func GetProfile(prfID string) (*api.ProfileSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "profiles/" + prfID,
	}
	dbRes := db.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}

	var prf = &api.ProfileSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), prf); err != nil {
		log.Println("[Error] When parsing profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}
	return prf, nil
}

func ListProfiles() (*[]api.ProfileSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "profiles",
	}
	dbRes := db.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list profiles in db:", dbRes.Error)
		return &[]api.ProfileSpec{}, errors.New(dbRes.Error)
	}

	var prfs = []api.ProfileSpec{}
	if len(dbRes.Message) == 0 {
		return &prfs, nil
	}
	for _, msg := range dbRes.Message {
		var prf = api.ProfileSpec{}
		if err := json.Unmarshal([]byte(msg), &prf); err != nil {
			log.Println("[Error] When parsing profile in db:", dbRes.Error)
			return &[]api.ProfileSpec{}, errors.New(dbRes.Error)
		}
		prfs = append(prfs, prf)
	}
	return &prfs, nil
}

func UpdateProfile(prfID string, input *api.ProfileSpec) (*api.ProfileSpec, error) {
	prf, err := GetProfile(prfID)
	if err != nil {
		return &api.ProfileSpec{}, err
	}
	if name := input.GetName(); name != "" {
		prf.Name = name
	}
	if desp := input.GetDescription(); desp != "" {
		prf.Description = desp
	}
	if props := input.Extra; len(props) != 0 {
		return &api.ProfileSpec{}, errors.New("Failed to update extra properties!")
	}

	prfBody, err := json.Marshal(prf)
	if err != nil {
		return &api.ProfileSpec{}, err
	}

	dbReq := &db.Request{
		Url:        db.URL_PREFIX + "profiles/" + prfID,
		NewContent: string(prfBody),
	}
	dbRes := db.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}
	return prf, nil
}

func DeleteProfile(prfID string) error {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "profiles/" + prfID,
	}
	dbRes := db.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete profile in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func AddExtraProperty(prfID string, ext api.ExtraSpec) (*api.ExtraSpec, error) {
	prf, err := GetProfile(prfID)
	if err != nil {
		return &api.ExtraSpec{}, err
	}

	for k, v := range ext {
		prf.Extra[k] = v
	}

	prf, err = CreateProfile(prf)
	if err != nil {
		return &api.ExtraSpec{}, err
	}
	return &prf.Extra, nil
}

func ListExtraProperties(prfID string) (*api.ExtraSpec, error) {
	prf, err := GetProfile(prfID)
	if err != nil {
		return &api.ExtraSpec{}, err
	}
	return &prf.Extra, nil
}

func RemoveExtraProperty(prfID, extraKey string) error {
	prf, err := GetProfile(prfID)
	if err != nil {
		return err
	}

	delete(prf.Extra, extraKey)
	prf, err = CreateProfile(prf)
	if err != nil {
		return err
	}
	return nil
}
