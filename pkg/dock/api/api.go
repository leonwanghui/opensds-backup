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

package api

import (
	"encoding/json"
	"log"
	"strings"

	db "github.com/opensds/opensds/pkg/db/api"
	"github.com/opensds/opensds/pkg/dock"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	api "github.com/opensds/opensds/pkg/model"
)

func CreateVolume(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	vol, err := dock.CreateVolume(dck.GetDriverName(), req.GetVolumeName(), req.GetVolumeSize())
	if err != nil {
		log.Println("[Error] When create volume in dock module:", err)
		return &pb.DockResponse{}, err
	}
	vol.ProfileId = req.GetProfileId()
	vol.PoolId = req.GetPoolId()

	result, err := db.CreateVolume(vol)
	if err != nil {
		log.Println("[Error] When create volume in db module:", err)
		return &pb.DockResponse{}, err
	}

	volBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

func GetVolume(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	result, err := dock.GetVolume(dck.DriverName, req.GetVolumeId())
	if err != nil {
		log.Println("[Error] When get volume in dock module:", err)
		return &pb.DockResponse{}, err
	}

	volBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

func DeleteVolume(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	if err := dock.DeleteVolume(dck.DriverName, req.GetVolumeId()); err != nil {
		log.Println("Error occured in dock module when delete volume:", err)
		return &pb.DockResponse{}, err
	}

	if err := db.DeleteVolume(req.GetVolumeId()); err != nil {
		log.Println("Error occured in dock module when delete volume in db:", err)
		return &pb.DockResponse{}, err
	}

	return &pb.DockResponse{
		Status:  "Success",
		Message: "Delete volume success",
	}, nil
}

func CreateVolumeAttachment(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck, hostInfo = &api.DockSpec{}, &api.HostInfo{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}
	if err := json.Unmarshal([]byte(req.GetHostInfo()), hostInfo); err != nil {
		log.Println("Error occured in dock module when parsing host info:", err)
		return &pb.DockResponse{}, err
	}

	atc, err := dock.CreateVolumeAttachment(dck.DriverName, req.GetVolumeId(), req.GetDoLocalAttach(), req.GetMultiPath(), hostInfo)
	if err != nil {
		log.Println("Error occured in dock module when create volume attachment:", err)
		return &pb.DockResponse{}, err
	}

	result, err := db.CreateVolumeAttachment(req.GetVolumeId(), atc)
	if err != nil {
		log.Println("Error occured in dock module when create volume attachment in db:", err)
		return &pb.DockResponse{}, err
	}

	atcBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(atcBody),
	}, nil
}

func UpdateVolumeAttachment(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck, hostInfo = &api.DockSpec{}, &api.HostInfo{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}
	if err := json.Unmarshal([]byte(req.GetHostInfo()), hostInfo); err != nil {
		log.Println("Error occured in dock module when parsing host info:", err)
		return &pb.DockResponse{}, err
	}

	if err := dock.UpdateVolumeAttachment(dck.DriverName, req.GetVolumeId(), hostInfo.Host, req.GetMountpoint()); err != nil {
		log.Println("Error occured in dock module when update volume attachment:", err)
		return &pb.DockResponse{}, err
	}

	result, err := db.UpdateVolumeAttachment(req.GetVolumeId(), req.GetAttachmentId(), req.GetMountpoint(), hostInfo)
	if err != nil {
		log.Println("Error occured in dock module when update volume attachment in db:", err)
		return &pb.DockResponse{}, err
	}

	atcBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(atcBody),
	}, nil
}

func DeleteVolumeAttachment(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	if err := dock.DeleteVolumeAttachment(dck.DriverName, req.GetVolumeId()); err != nil {
		log.Println("Error occured in dock module when delete volume attachment:", err)
		if strings.Contains(err.Error(), "The status of volume is not in-use") {
			if err = db.DeleteVolumeAttachment(req.GetVolumeId(), req.GetAttachmentId()); err != nil {
				log.Println("Error occured in dock module when delete volume attachment in db:", err)
				return &pb.DockResponse{}, err
			}
		}
		return &pb.DockResponse{}, nil
	}

	if err := db.DeleteVolumeAttachment(req.GetVolumeId(), req.GetAttachmentId()); err != nil {
		log.Println("Error occured in dock module when delete volume attachment in db:", err)
		return &pb.DockResponse{}, err
	}

	return &pb.DockResponse{
		Status:  "Success",
		Message: "Delete volume attachment success",
	}, nil
}

func CreateVolumeSnapshot(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	snp, err := dock.CreateSnapshot(dck.DriverName,
		req.GetSnapshotName(),
		req.GetVolumeId(),
		req.GetSnapshotDescription())
	if err != nil {
		log.Println("Error occured in dock module when create snapshot:", err)
		return &pb.DockResponse{}, err
	}

	result, err := db.CreateVolumeSnapshot(snp)
	if err != nil {
		log.Println("Error occured in dock module when create volume snapshot in db:", err)
		return &pb.DockResponse{}, err
	}

	snpBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(snpBody),
	}, nil
}

func GetVolumeSnapshot(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	result, err := dock.GetSnapshot(dck.DriverName, req.GetSnapshotId())
	if err != nil {
		log.Println("Error occured in dock module when get snapshot:", err)
		return &pb.DockResponse{}, err
	}

	snpBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(snpBody),
	}, nil
}

func DeleteVolumeSnapshot(req *pb.DockRequest) (*pb.DockResponse, error) {
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	if err := dock.DeleteSnapshot(dck.DriverName, req.GetSnapshotId()); err != nil {
		log.Println("Error occured in dock module when delete snapshot:", err)
		return &pb.DockResponse{}, err
	}

	if err := db.DeleteVolumeSnapshot(req.GetSnapshotId()); err != nil {
		log.Println("Error occured in dock module when delete volume snapshot in db:", err)
		return &pb.DockResponse{}, err
	}

	return &pb.DockResponse{
		Status:  "Success",
		Message: "Delete snapshot success",
	}, nil
}
