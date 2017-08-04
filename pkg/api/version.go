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
This module implements a entry into the OpenSDS northbound REST service.

*/

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/opensds/opensds/pkg/model"

	"github.com/astaxie/beego"
)

const (
	KnownVersions = `[
		{
			"name": "v1alpha",
			"description": "v1alpha version",
			"status": "DEPRECATED",
			"updatedAt": "2017-04-10T14:36:58.014Z"
		},
		{
			"name": "v1beta1",
			"description": "first v1beta version",
			"status": "CURRENT",
			"updatedAt": "2017-06-10T14:36:58.014Z"
		},
		{
			"name": "v1beta2",
			"description": "second v1beta version",
			"status": "SUPPORTED",
			"updatedAt": "2017-07-10T14:36:58.014Z"
		},
	]`
)

type VersionPortal struct {
	beego.Controller
}

func (this *VersionPortal) ListVersions() {
	result, err := SearchVersions()
	if err != nil {
		rbody, _ := json.Marshal("List versions failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(rbody)
	} else {
		rbody, _ := json.Marshal(&result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(rbody)
	}
}

type SpecifiedVersionPortal struct {
	beego.Controller
}

func (this *SpecifiedVersionPortal) GetVersion() {
	version := this.Ctx.Input.Param(":apiVersion")

	result, err := SearchVersion(version)
	if err != nil {
		rbody, _ := json.Marshal("Get version failed:" + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(rbody)
	} else {
		rbody, _ := json.Marshal(&result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(rbody)
	}
}

func SearchVersions() ([]model.VersionSpec, error) {
	var versions []model.VersionSpec

	err := json.Unmarshal([]byte(KnownVersions), &versions)
	if err != nil {
		log.Println(err)
		return []model.VersionSpec{}, err
	}
	return versions, nil
}

func SearchVersion(versionName string) (model.VersionSpec, error) {
	versions, err := SearchVersions()
	if err != nil {
		log.Println(err)
		return model.VersionSpec{}, err
	}

	for _, version := range versions {
		if version.GetName() == versionName {
			return version, nil
		}
	}

	log.Println(errors.New("Can't find v1 in available versions!"))
	return model.VersionSpec{}, err
}
