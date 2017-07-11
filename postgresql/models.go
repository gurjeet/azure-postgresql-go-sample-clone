package postgresql

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// ServerProperties hold the parameters for creating a server instance
// Note: these are flattened and simplified for this sample
type ServerProperties struct {
	Location                   string
	AdministratorLogin         string
	AdministratorLoginPassword string
	Version                    ServerVersion
	Tier                       SkuTier // Basic, Standard
	ComputeUnits               int32   // if empty use default for tier
	StorageMB                  int32   // if empty use default for tier
}

// ServerVersion enumerates the values for server version.
type ServerVersion string

const (
	ServerVersionNineFullStopFive ServerVersion = "9.5"
	ServerVersionNineFullStopSix  ServerVersion = "9.6"
)

// SkuTier enumerates skus
type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierStandard SkuTier = "Standard"
)

// CreateMode enumerates create modes
type CreateMode string

const (
	CreateModeDefault            CreateMode = "Default"
	CreateModePointInTimeRestore CreateMode = "PointInTimeRestore"
)
