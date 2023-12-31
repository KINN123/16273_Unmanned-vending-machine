/*
 * Copyright (C) 2021-2022 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include "wifi_permission_utils.h"
#ifndef OHOS_ARCH_LITE
#include "ipc_skeleton.h"
#endif
#include "wifi_auth_center.h"

namespace OHOS {
namespace Wifi {
#ifdef OHOS_ARCH_LITE
int WifiPermissionUtils::VerifySetWifiInfoPermission()
{
    return PERMISSION_GRANTED;
}

int WifiPermissionUtils::VerifyGetWifiInfoPermission()
{
    return PERMISSION_GRANTED;
}

int WifiPermissionUtils::VerifyWifiConnectionPermission()
{
    return PERMISSION_GRANTED;
}

int WifiPermissionUtils::VerifyGetScanInfosPermission()
{
    return PERMISSION_GRANTED;
}

int WifiPermissionUtils::VerifyGetWifiLocalMacPermission()
{
    return PERMISSION_GRANTED;
}

int WifiPermissionUtils::VerifySetWifiConfigPermission()
{
    return PERMISSION_GRANTED;
}

int WifiPermissionUtils::VerifyGetWifiDirectDevicePermission()
{
    return PERMISSION_GRANTED;
}

#else
int WifiPermissionUtils::VerifySetWifiInfoPermission()
{
    return WifiAuthCenter::GetInstance().VerifySetWifiInfoPermission(
        IPCSkeleton::GetCallingPid(), IPCSkeleton::GetCallingUid());
}

int WifiPermissionUtils::VerifyGetWifiInfoPermission()
{
    return WifiAuthCenter::GetInstance().VerifyGetWifiInfoPermission(
        IPCSkeleton::GetCallingPid(), IPCSkeleton::GetCallingUid());
}

int WifiPermissionUtils::VerifyWifiConnectionPermission()
{
    return WifiAuthCenter::GetInstance().VerifyWifiConnectionPermission(
        IPCSkeleton::GetCallingPid(), IPCSkeleton::GetCallingUid());
}

int WifiPermissionUtils::VerifyGetScanInfosPermission()
{
    return WifiAuthCenter::GetInstance().VerifyGetScanInfosPermission(
        IPCSkeleton::GetCallingPid(), IPCSkeleton::GetCallingUid());
}

int WifiPermissionUtils::VerifyGetWifiLocalMacPermission()
{
    return WifiAuthCenter::GetInstance().VerifyGetWifiLocalMacPermission(
        IPCSkeleton::GetCallingPid(), IPCSkeleton::GetCallingUid());
}

int WifiPermissionUtils::VerifySetWifiConfigPermission()
{
    return WifiAuthCenter::GetInstance().VerifySetWifiConfigPermission(
        IPCSkeleton::GetCallingPid(), IPCSkeleton::GetCallingUid());
}

int WifiPermissionUtils::VerifyGetWifiDirectDevicePermission()
{
    return WifiAuthCenter::GetInstance().VerifyGetWifiDirectDevicePermission(
        IPCSkeleton::GetCallingPid(), IPCSkeleton::GetCallingUid());
}
#endif
}  // namespace Wifi
}  // namespace OHOS