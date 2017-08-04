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

func CreateVolume(vol *api.VolumeSpec) (*api.VolumeSpec, error) {
	volBody, err := json.Marshal(vol)
	if err != nil {
		return &api.VolumeSpec{}, err
	}

	dbReq := &db.Request{
		Url:     db.URL_PREFIX + "volumes/" + vol.Id,
		Content: string(volBody),
	}
	dbRes := db.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume in db:", dbRes.Error)
		return &api.VolumeSpec{}, errors.New(dbRes.Error)
	}

	return vol, nil
}

func GetVolume(volID string) (*api.VolumeSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volumes/" + volID,
	}
	dbRes := db.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume in db:", dbRes.Error)
		return &api.VolumeSpec{}, errors.New(dbRes.Error)
	}

	var vol = &api.VolumeSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vol); err != nil {
		log.Println("[Error] When parsing volume in db:", dbRes.Error)
		return &api.VolumeSpec{}, errors.New(dbRes.Error)
	}
	return vol, nil
}

func ListVolumes() (*[]api.VolumeSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volumes",
	}
	dbRes := db.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volumes in db:", dbRes.Error)
		return &[]api.VolumeSpec{}, errors.New(dbRes.Error)
	}

	var vols = []api.VolumeSpec{}
	if len(dbRes.Message) == 0 {
		return &vols, nil
	}
	for _, msg := range dbRes.Message {
		var vol = api.VolumeSpec{}
		if err := json.Unmarshal([]byte(msg), &vol); err != nil {
			log.Println("[Error] When parsing volume in db:", dbRes.Error)
			return &[]api.VolumeSpec{}, errors.New(dbRes.Error)
		}
		vols = append(vols, vol)
	}
	return &vols, nil
}

func DeleteVolume(volID string) error {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volumes/" + volID,
	}
	dbRes := db.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func CreateVolumeAttachment(volID string, atc *api.VolumeAttachmentSpec) (*api.VolumeAttachmentSpec, error) {
	atcBody, err := json.Marshal(atc)
	if err != nil {
		return &api.VolumeAttachmentSpec{}, err
	}

	dbReq := &db.Request{
		Url:     db.URL_PREFIX + "volume/" + volID + "/attachments/" + atc.Id,
		Content: string(atcBody),
	}
	dbRes := db.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}

	return atc, nil
}

func GetVolumeAttachment(volID, attachmentID string) (*api.VolumeAttachmentSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volume/" + volID + "/attachments/" + attachmentID,
	}
	dbRes := db.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}

	var atc = &api.VolumeAttachmentSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), atc); err != nil {
		log.Println("[Error] When parsing volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}
	return atc, nil
}

func ListVolumeAttachments(volID string) (*[]api.VolumeAttachmentSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volume/" + volID + "/attachments",
	}
	dbRes := db.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volume attachments in db:", dbRes.Error)
		return &[]api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}

	var atcs = []api.VolumeAttachmentSpec{}
	if len(dbRes.Message) == 0 {
		return &atcs, nil
	}
	for _, msg := range dbRes.Message {
		var atc = api.VolumeAttachmentSpec{}
		if err := json.Unmarshal([]byte(msg), &atc); err != nil {
			log.Println("[Error] When parsing volume attachment in db:", dbRes.Error)
			return &[]api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
		}
		atcs = append(atcs, atc)
	}
	return &atcs, nil
}

func UpdateVolumeAttachment(volID, attachmentID, mountpoint string, hostInfo *api.HostInfo) (*api.VolumeAttachmentSpec, error) {
	atc, err := GetVolumeAttachment(volID, attachmentID)
	if err != nil {
		return &api.VolumeAttachmentSpec{}, err
	}

	atc.HostInfo = hostInfo
	atc.Mountpoint = mountpoint
	atcBody, err := json.Marshal(atc)
	if err != nil {
		return &api.VolumeAttachmentSpec{}, err
	}

	dbReq := &db.Request{
		Url:        db.URL_PREFIX + "volume/" + volID + "/attachments/" + attachmentID,
		NewContent: string(atcBody),
	}
	dbRes := db.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}
	return atc, nil
}

func DeleteVolumeAttachment(volID, attachmentID string) error {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volume/" + volID + "/attachments/" + attachmentID,
	}
	dbRes := db.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume attachment in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func CreateVolumeSnapshot(vs *api.VolumeSnapshotSpec) (*api.VolumeSnapshotSpec, error) {
	vsBody, err := json.Marshal(vs)
	if err != nil {
		return &api.VolumeSnapshotSpec{}, err
	}

	dbReq := &db.Request{
		Url:     db.URL_PREFIX + "volume/snapshots/" + vs.Id,
		Content: string(vsBody),
	}
	dbRes := db.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume snapshot in db:", dbRes.Error)
		return &api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}

	return vs, nil
}

func GetVolumeSnapshot(snapshotID string) (*api.VolumeSnapshotSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volume/snapshots/" + snapshotID,
	}
	dbRes := db.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume attachment in db:", dbRes.Error)
		return &api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}

	var vs = &api.VolumeSnapshotSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vs); err != nil {
		log.Println("[Error] When parsing volume snapshot in db:", dbRes.Error)
		return &api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}
	return vs, nil
}

func ListVolumeSnapshots() (*[]api.VolumeSnapshotSpec, error) {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volume/snapshots",
	}
	dbRes := db.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volume snapshots in db:", dbRes.Error)
		return &[]api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}

	var vss = []api.VolumeSnapshotSpec{}
	if len(dbRes.Message) == 0 {
		return &vss, nil
	}
	for _, msg := range dbRes.Message {
		var vs = api.VolumeSnapshotSpec{}
		if err := json.Unmarshal([]byte(msg), &vs); err != nil {
			log.Println("[Error] When parsing volume snapshot in db:", dbRes.Error)
			return &[]api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
		}
		vss = append(vss, vs)
	}
	return &vss, nil
}

func DeleteVolumeSnapshot(snapshotID string) error {
	dbReq := &db.Request{
		Url: db.URL_PREFIX + "volume/snapshots/" + snapshotID,
	}
	dbRes := db.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
