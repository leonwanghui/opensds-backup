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
)

type DockPortal struct {
	beego.Controller
}

func (this *DockPortal) ListDocks() {
	result, err := db.ListDocks()
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List docks failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}

type SpecifiedDockPortal struct {
	beego.Controller
}

func (this *SpecifiedDockPortal) GetDock() {
	id := this.Ctx.Input.Param(":dockId")

	result, err := db.GetDock(id)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get dock failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(StatusOK)
		this.Ctx.Output.Body(resBody)
	}
}
