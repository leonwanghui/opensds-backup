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
				// To judge whether the scheme is legal or not.
				if ctx.Input.Scheme() != "http" && ctx.Input.Scheme() != "https" {
					return false
				}
				return true
			}),

			// List all opensds api versions
			beego.NSRouter("/", &VersionPortal{}, "get:ListVersions"),
			// Show specified api version
			beego.NSRouter("/:apiVersion", &SpecifiedVersionPortal{}, "get:GetVersion"),

			beego.NSNamespace("/v1alpha",
				// List all dock services, including a list of dock object
				beego.NSRouter("/docks", &DockPortal{}, "get:ListDocks"),
				// Show one dock service, including endpoint, driverName and so on
				beego.NSRouter("/docks/:dockId", &SpecifiedDockPortal{}, "get:GetDock"),

				// Profile is a set of policies configured by admin and provided for users
				// CreateProfile, UpdateProfile and DeleteProfile are used for admin only
				// ListProfiles and GetProfile are used for both admin and users
				beego.NSRouter("/profiles", &ProfilePortal{}, "post:CreateProfile;get:ListProfiles"),
				beego.NSRouter("/profiles/:profileId", &SpecifiedProfilePortal{}, "get:GetProfile;put:UpdateProfile;delete:DeleteProfile"),

				// All operations of extras are used for Admin only
				beego.NSRouter("/profiles/:profileId/extras", &ProfileExtrasPortal{}, "post:AddExtraProperty;get:ListExtraProperties"),
				beego.NSRouter("/profiles/:profileId/extras/:extraKey", &ProfileExtrasPortal{}, "delete:RemoveExtraProperty"),

				beego.NSNamespace("/block",
					// Pool is the virtual description of backend storage, usually divided into block, file and object,
					// and every pool is atomic, which means every pool contains a specific set of features.
					// ListPools and GetPool are used for checking the status of backend pool, admin only
					beego.NSRouter("/pools", &PoolPortal{}, "get:ListPools"),
					beego.NSRouter("/pools/:poolId", &SpecifiedPoolPortal{}, "get:GetPool"),

					// Volume is the logical description of a piece of storage, which can be directly used by users.
					// All operations of volume can be used for both admin and users.
					beego.NSRouter("/volumes", &VolumePortal{}, "post:CreateVolume;get:ListVolumes"),
					beego.NSRouter("/volumes/:volumeId", &SpecifiedVolumePortal{}, "get:GetVolume;put:UpdateVolume;delete:DeleteVolume"),

					// Creates, shows, lists, unpdates and deletes attachment.
					beego.NSRouter("/attachments", &VolumeAttachmentPortal{}, "post:CreateVolumeAttachment;get:ListVolumeAttachments"),
					beego.NSRouter("/attachments/:attachmentId", &SpecifiedVolumeAttachmentPortal{}, "get:GetVolumeAttachment;put:UpdateVolumeAttachment;delete:DeleteVolumeAttachment"),

					// Snapshot is a point-in-time copy of the data that a volume contains.
					// Creates, shows, lists, unpdates and deletes snapshot.
					beego.NSRouter("/snapshots", &VolumeSnapshotPortal{}, "post:CreateVolumeSnapshot;get:ListVolumeSnapshots"),
					beego.NSRouter("/snapshots/:snapshotId", &SpecifiedVolumeSnapshotPortal{}, "get:GetVolumeSnapshot;put:UpdateVolumeSnapshot;delete:DeleteVolumeSnapshot"),
				),
			),
		)

	beego.AddNamespace(ns)
	beego.Run(host)
}
