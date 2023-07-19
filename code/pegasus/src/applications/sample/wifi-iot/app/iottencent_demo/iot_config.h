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

#ifndef IOT_CONFIG_H
#define IOT_CONFIG_H

// <CONFIG THE LOG
/* if you need the iot log for the development , please enable it, else please comment it */
#define CONFIG_LINKLOG_ENABLE   1

// < CONFIG THE WIFI
/* Please modify the ssid and pwd for the own */
#define CONFIG_AP_SSID  "iQOO Neo5" // WIFI SSID
#define CONFIG_AP_PWD   "1348564418" // WIFI PWD
/* Tencent iot Cloud user ID , password */
#define CONFIG_USER_ID    "EFUHAVJCTYHi3861_text;12010126;a1ef6;1696003200"
#define CONFIG_USER_PWD   "ab0d0e52883a9d2d6059bbffbd41596bea44ab5a94832234a69e5938a0bb9cd5;hmacsha256"
#define CN_CLIENTID     "EFUHAVJCTYHi3861_text" // Tencent cloud ClientID format: Product ID + device name
#endif