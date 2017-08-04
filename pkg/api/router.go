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
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

const (
	StatusOK       = http.StatusOK
	StatusCreated  = http.StatusCreated
	StatusAccepted = http.StatusAccepted

	StatusBadRequest   = http.StatusBadRequest
	StatusUnauthorized = http.StatusUnauthorized
	StatusForbidden    = http.StatusForbidden
	StatusNotFound     = http.StatusNotFound

	StatusInternalServerError = http.StatusInternalServerError
	StatusNotImplemented      = http.StatusNotImplemented
)

func Run(host string) {
	ns :=
		beego.NewNamespace("/api",
			beego.NSCond(func(ctx *context.Context) bool {
				if ctx.Input.Scheme() == "http" {
					return true
				}
				return false
			}),
			beego.NSRouter("/", &VersionPortal{}, "get:ListVersions"),
			beego.NSRouter("/:apiVersion", &SpecifiedVersionPortal{}, "get:GetVersion"),
			beego.NSNamespace("/v1alpha1",
				beego.NSNamespace("/block",
					beego.NSNamespace("/docks",
						beego.NSRouter("/", &DockPortal{}, "get:ListDocks"),
						beego.NSRouter("/:dockId", &SpecifiedDockPortal{}, "get:GetDock"),
					),
					beego.NSNamespace("/pools",
						beego.NSRouter("/", &PoolPortal{}, "get:ListPools"),
						beego.NSRouter("/:poolId", &SpecifiedPoolPortal{}, "get:GetPool"),
					),
					beego.NSNamespace("/profiles",
						beego.NSRouter("/", &ProfilePortal{}, "post:CreateProfile;get:ListProfiles"),
						beego.NSRouter("/:profileId", &SpecifiedProfilePortal{}, "get:GetProfile;put:UpdateProfile;delete:DeleteProfile"),
						beego.NSRouter("/:profileId/extras", &ProfileExtrasPortal{}, "post:AddExtraProperty;get:ListExtraProperties"),
						beego.NSRouter("/:profileId/extras/:extraKey", &ProfileExtrasPortal{}, "delete:RemoveExtraProperty"),
					),
					beego.NSNamespace("/volumes",
						beego.NSRouter("/", &VolumePortal{}, "post:CreateVolume;get:ListVolumes"),
						beego.NSRouter("/:volumeId", &SpecifiedVolumePortal{}, "get:GetVolume;put:UpdateVolume;delete:DeleteVolume"),
					),
					beego.NSNamespace("/attachments",
						beego.NSRouter("/", &VolumeAttachmentPortal{}, "post:CreateVolumeAttachment;get:ListVolumeAttachments"),
						beego.NSRouter("/:attachmentId", &SpecifiedVolumeAttachmentPortal{}, "get:GetVolumeAttachment;put:UpdateVolumeAttachment;delete:DeleteVolumeAttachment"),
					),
					beego.NSNamespace("/snapshots",
						beego.NSRouter("/", &VolumeSnapshotPortal{}, "post:CreateVolumeSnapshot;get:ListVolumeSnapshots"),
						beego.NSRouter("/:snapshotId", &SpecifiedVolumeSnapshotPortal{}, "get:GetVolumeSnapshot;put:UpdateVolumeSnapshot;delete:DeleteVolumeSnapshot"),
					),
				),
			),
		)

	beego.AddNamespace(ns)
	beego.Run(host)
}
