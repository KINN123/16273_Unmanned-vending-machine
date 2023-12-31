/*
 * Copyright (c) 2022 HiSilicon (Shanghai) Technologies CO., LIMITED.
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

#ifndef STATUSCODE_H
#define STATUSCODE_H

#include "status_code.h"
typedef enum StatusCode {
    StatusFailure = -1,
    StatusOk = 0,
    StatusTimeOut,
    StatusURIError,
    StatusIOError,
    StatusIOReady,
    StatusNoInitialize,
    StatusDeviceOffline,
    StatusNoDevice,
    StatusNoVideoDevices,
    StatusNewFramed,
    StatusFramedEnd,
    StatusComposedData,
    StatusNoSession,
    StatusTargetExist,
    StatusTargetNoExist,
    StatusQEmpty,
    StatusLostFrame,
} StatusCode;

#endif /* __STATUSCODE__H__ */
