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
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	db "github.com/opensds/opensds/pkg/db/api"
	"github.com/opensds/opensds/pkg/model"
	"github.com/satori/go.uuid"
)

type ProfileRequest struct {
	Spec *model.ProfileSpec `json:"spec"`
}

func NewProfileRequest() *ProfileRequest {
	return &ProfileRequest{
		Spec: &model.ProfileSpec{
			BaseModel: &model.BaseModel{},
		},
	}
}

type ProfilePortal struct {
	beego.Controller
}

func (this *ProfilePortal) CreateProfile() {
	req := NewProfileRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse profile request body failed:", err)
		resBody, _ := json.Marshal("Parse profile request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}
	req.Spec.Id = uuid.NewV4().String()

	result, err := db.CreateProfile(req.Spec)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create storage profile failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusCreated)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *ProfilePortal) ListProfiles() {
	result, err := db.ListProfiles()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List storage profiles failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

type SpecifiedProfilePortal struct {
	beego.Controller
}

func (this *SpecifiedProfilePortal) GetProfile() {
	id := this.Ctx.Input.Param(":profileId")

	result, err := db.GetProfile(id)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get storage profile failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedProfilePortal) UpdateProfile() {
	id := this.Ctx.Input.Param(":profileId")

	req := NewProfileRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse profile request body failed:", err)
		resBody, _ := json.Marshal("Parse profile request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result, err := db.UpdateProfile(id, req.Spec)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get storage profile failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusAccepted)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedProfilePortal) DeleteProfile() {
	id := this.Ctx.Input.Param(":profileId")

	if err := db.DeleteProfile(id); err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Delete storage profile failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal("Delete profile success!")
		this.Ctx.Output.SetStatus(StatusAccepted)
		this.Ctx.Output.Body(resBody)
	}
}

type ExtraRequest struct {
	Spec model.ExtraSpec `json:"spec"`
}

func NewExtraRequest() *ExtraRequest {
	return &ExtraRequest{
		Spec: model.ExtraSpec{},
	}
}

type ProfileExtrasPortal struct {
	beego.Controller
}

func (this *ProfileExtrasPortal) AddExtraProperty() {
	id := this.Ctx.Input.Param(":profileId")

	req := NewExtraRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse extra request body failed:", err)
		resBody, _ := json.Marshal("Parse extra request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result, err := db.AddExtraProperty(id, req.Spec)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create extra property failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusCreated)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *ProfileExtrasPortal) ListExtraProperties() {
	id := this.Ctx.Input.Param(":profileId")

	result, err := db.ListExtraProperties(id)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List extra properties failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *ProfileExtrasPortal) RemoveExtraProperty() {
	id := this.Ctx.Input.Param(":profileId")
	extraKey := this.Ctx.Input.Param(":extraKey")

	if err := db.RemoveExtraProperty(id, extraKey); err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Remove profile extra property failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal("Remove extra property success!")
		this.Ctx.Output.SetStatus(StatusAccepted)
		this.Ctx.Output.Body(resBody)
	}
}
