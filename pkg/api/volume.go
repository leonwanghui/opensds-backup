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

	"github.com/opensds/opensds/pkg/controller"
	db "github.com/opensds/opensds/pkg/db/api"
	"github.com/opensds/opensds/pkg/model"

	"github.com/astaxie/beego"
)

type VolumeRequest struct {
	Spec *model.VolumeSpec `json:"spec"`
}

func NewVolumeRequest() *VolumeRequest {
	return &VolumeRequest{
		Spec: &model.VolumeSpec{
			BaseModel: &model.BaseModel{},
		},
	}
}

type VolumePortal struct {
	beego.Controller
}

func (this *VolumePortal) CreateVolume() {
	req := NewVolumeRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	c, err := controller.NewControllerWithVolumeConfig(req.Spec, nil, nil)
	if err != nil {
		log.Println("Set up controller failed:", err)
		resBody, _ := json.Marshal("Set up controller failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result, err := c.CreateVolume()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create volume failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusCreated)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumePortal) ListVolumes() {
	result, err := db.ListVolumes()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List volumes failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

type SpecifiedVolumePortal struct {
	beego.Controller
}

func (this *SpecifiedVolumePortal) GetVolume() {
	id := this.Ctx.Input.Param(":volumeId")

	result, err := db.GetVolume(id)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get volume failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumePortal) UpdateVolume() {
	this.Ctx.Output.SetStatus(StatusNotImplemented)
}

func (this *SpecifiedVolumePortal) DeleteVolume() {
	volId := this.Ctx.Input.Param(":volumeId")

	req := NewVolumeRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}
	req.Spec.Id = volId

	c, err := controller.NewControllerWithVolumeConfig(req.Spec, nil, nil)
	if err != nil {
		log.Println("Set up controller failed:", err)
		resBody, _ := json.Marshal("Set up controller failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result := c.DeleteVolume()
	if result.Status != "Success" {
		resBody, _ := json.Marshal("Delete volume failed: " + result.Error)
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusAccepted)
		this.Ctx.Output.Body(resBody)
	}
}

type VolumeAttachmentRequest struct {
	Spec *model.VolumeAttachmentSpec `json:"spec"`
}

func NewVolumeAttachmentRequest() *VolumeAttachmentRequest {
	return &VolumeAttachmentRequest{
		Spec: &model.VolumeAttachmentSpec{
			BaseModel: &model.BaseModel{},
		},
	}
}

type VolumeAttachmentPortal struct {
	beego.Controller
}

func (this *VolumeAttachmentPortal) CreateVolumeAttachment() {
	req := NewVolumeAttachmentRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse volume attachment request body failed:", err)
		resBody, _ := json.Marshal("Parse volume attachment request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	c, err := controller.NewControllerWithVolumeConfig(nil, req.Spec, nil)
	if err != nil {
		log.Println("Set up controller failed:", err)
		resBody, _ := json.Marshal("Set up controller failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result, err := c.CreateVolumeAttachment()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create volume attachment failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusCreated)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeAttachmentPortal) ListVolumeAttachments() {
	volId := this.GetString("volumeId")

	result, err := db.ListVolumeAttachments(volId)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List volume attachments failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

type SpecifiedVolumeAttachmentPortal struct {
	beego.Controller
}

func (this *SpecifiedVolumeAttachmentPortal) GetVolumeAttachment() {
	id := this.Ctx.Input.Param(":attachmentId")
	volId := this.GetString("volumeId")

	result, err := db.GetVolumeAttachment(volId, id)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get volume attachment failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumeAttachmentPortal) UpdateVolumeAttachment() {
	id := this.Ctx.Input.Param(":attachmentId")

	req := NewVolumeAttachmentRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse volume attachment request body failed:", err)
		resBody, _ := json.Marshal("Parse volume attachment request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}
	req.Spec.Id = id

	c, err := controller.NewControllerWithVolumeConfig(nil, req.Spec, nil)
	if err != nil {
		log.Println("Set up controller failed:", err)
		resBody, _ := json.Marshal("Set up controller failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result, err := c.UpdateVolumeAttachment()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Update volume attachment failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumeAttachmentPortal) DeleteVolumeAttachment() {
	id := this.Ctx.Input.Param(":attachmentId")

	req := NewVolumeAttachmentRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse volume attachment request body failed:", err)
		resBody, _ := json.Marshal("Parse volume attachment request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}
	req.Spec.Id = id

	c, err := controller.NewControllerWithVolumeConfig(nil, req.Spec, nil)
	if err != nil {
		log.Println("Set up controller failed:", err)
		resBody, _ := json.Marshal("Set up controller failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result := c.DeleteVolumeAttachment()
	if result.Status != "Success" {
		resBody, _ := json.Marshal("Delete volume attachment failed: " + result.Error)
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusAccepted)
		this.Ctx.Output.Body(resBody)
	}
}

type VolumeSnapshotRequest struct {
	Spec *model.VolumeSnapshotSpec `json:"spec"`
}

func NewVolumeSnapshotRequest() *VolumeSnapshotRequest {
	return &VolumeSnapshotRequest{
		Spec: &model.VolumeSnapshotSpec{
			BaseModel: &model.BaseModel{},
		},
	}
}

type VolumeSnapshotPortal struct {
	beego.Controller
}

func (this *VolumeSnapshotPortal) CreateVolumeSnapshot() {
	req := NewVolumeSnapshotRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse volume snapshot request body failed:", err)
		resBody, _ := json.Marshal("Parse volume snapshot request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	c, err := controller.NewControllerWithVolumeConfig(nil, nil, req.Spec)
	if err != nil {
		log.Println("Set up controller failed:", err)
		resBody, _ := json.Marshal("Set up controller failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result, err := c.CreateVolumeSnapshot()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create volume snapshot failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusCreated)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeSnapshotPortal) ListVolumeSnapshots() {
	result, err := db.ListVolumeSnapshots()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List volume snapshots failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

type SpecifiedVolumeSnapshotPortal struct {
	beego.Controller
}

func (this *SpecifiedVolumeSnapshotPortal) GetVolumeSnapshot() {
	id := this.GetString("volumeId")

	result, err := db.GetVolumeSnapshot(id)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get volume snapshot failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumeSnapshotPortal) UpdateVolumeSnapshot() {
	this.Ctx.Output.SetStatus(StatusNotImplemented)
}

func (this *SpecifiedVolumeSnapshotPortal) DeleteVolumeSnapshot() {
	id := this.Ctx.Input.Param(":snapshotId")

	req := NewVolumeSnapshotRequest()
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(req); err != nil {
		log.Println("Parse volume snapshot request body failed:", err)
		resBody, _ := json.Marshal("Parse volume snapshot request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}
	req.Spec.Id = id

	c, err := controller.NewControllerWithVolumeConfig(nil, nil, req.Spec)
	if err != nil {
		log.Println("Set up controller failed:", err)
		resBody, _ := json.Marshal("Set up controller failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
	}

	result := c.DeleteVolumeSnapshot()
	if result.Status != "Success" {
		resBody, _ := json.Marshal("Delete volume snapshot failed: " + result.Error)
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusAccepted)
		this.Ctx.Output.Body(resBody)
	}
}
