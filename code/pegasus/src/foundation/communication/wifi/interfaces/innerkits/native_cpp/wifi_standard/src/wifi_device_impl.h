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
#ifndef OHOS_WIFI_DEVICE_IMPL_H
#define OHOS_WIFI_DEVICE_IMPL_H

#include "wifi_device.h"
#include "wifi_device_proxy.h"

namespace OHOS {
namespace Wifi {
class WifiDeviceImpl : public WifiDevice {
public:
    explicit WifiDeviceImpl(int systemAbilityId);
    ~WifiDeviceImpl();

    bool Init();

    /**
     * @Description Turn on Wi-Fi
     *
     * @return ErrCode - operation result
     */
    ErrCode EnableWifi() override;

    /**
     * @Description Turn off Wi-Fi
     *
     * @return ErrCode - operation result
     */
    ErrCode DisableWifi() override;

    /**
     * @Description Add a wifi device configuration.
     *
     * @param config - WifiDeviceConfig object
     * @param result - the device configuration's network id
     * @return ErrCode - operation result
     */
    ErrCode AddDeviceConfig(const WifiDeviceConfig &config, int &result) override;

    /**
     * @Description Remove the wifi device config equals to input network id
     *
     * @param networkId - want to remove device config's network id
     * @return ErrCode - operation result
     */
    ErrCode RemoveDevice(int networkId) override;

    /**
     * @Description Remove all device configs.
     *
     * @return ErrCode - operation result
     */
    ErrCode RemoveAllDevice() override;

    /**
     * @Description Get all the device configs
     *
     * @param result - Get result vector of WifiDeviceConfig
     * @return ErrCode - operation result
     */
    ErrCode GetDeviceConfigs(std::vector<WifiDeviceConfig> &result) override;

    /**
     * @Description Connecting to a Specified Network
     *
     * @param networkId - network id
     * @return ErrCode - operation result
     */
    ErrCode ConnectToNetwork(int networkId) override;

    /**
     * @Description Connect To a network base WifiDeviceConfig object
     *
     * @param config - WifiDeviceConfig object
     * @return ErrCode - operation result
     */
    ErrCode ConnectToDevice(const WifiDeviceConfig &config) override;

    /**
     * @Description Disconnect
     *
     * @return ErrCode - operation result
     */
    ErrCode Disconnect(void) override;

    /**
     * @Description Check whether Wi-Fi is active
     *
     * @param bActive - active / inactive
     * @return ErrCode - operation result
     */
    ErrCode IsWifiActive(bool &bActive) override;

    /**
     * @Description Get the Wifi State
     *
     * @param state - return current wifi state
     * @return ErrCode - operation result
     */
    ErrCode GetWifiState(int &state) override;

    /**
     * @Description Obtains the current Wi-Fi connection information
     *
     * @param info - WifiLinkedInfo object
     * @return ErrCode - operation result
     */
    ErrCode GetLinkedInfo(WifiLinkedInfo &info) override;

    /**
     * @Description Set the Country Code
     *
     * @param countryCode - country code
     * @return ErrCode - operation result
     */
    ErrCode SetCountryCode(const std::string &countryCode) override;

    /**
     * @Description Obtains the country code
     *
     * @param countryCode - output the country code
     * @return ErrCode - operation result
     */
    ErrCode GetCountryCode(std::string &countryCode) override;

    /**
     * @Description Register callback function.
     *
     * @param callback - IWifiDeviceCallBack object
     * @return ErrCode - operation result
     */
#ifdef OHOS_ARCH_LITE
    ErrCode RegisterCallBack(const std::shared_ptr<IWifiDeviceCallBack> &callback) override;
#else
    ErrCode RegisterCallBack(const sptr<IWifiDeviceCallBack> &callback) override;
#endif

    /**
     * @Description Get the Signal Level object
     *
     * @param rssi - rssi
     * @param band - band
     * @param level - return the level
     * @return ErrCode - operation result
     */
    ErrCode GetSignalLevel(const int &rssi, const int &band, int &level) override;

    /**
     * @Description Get supported features
     *
     * @param features - return supported features
     * @return ErrCode - operation result
     */
    ErrCode GetSupportedFeatures(long &features) override;

    /**
     * @Description Check if supported input feature
     *
     * @param feature - input feature
     * @return true - supported
     * @return false - unsupported
     */
    bool IsFeatureSupported(long feature) override;

    /**
     * @Description Enable device config, when set attemptEnable, disable other device config
     *
     * @param networkId - need enable device config's network id
     * @param attemptEnable - if set true, disable other device config
     * @return ErrCode - operation result
     */
    ErrCode EnableDeviceConfig(int networkId, bool attemptEnable);

    /**
     * @Description Disable Wi-Fi device configuration.
     *
     * @param networkId - device config's network id
     * @return ErrCode - operation result
     */
    ErrCode DisableDeviceConfig(int networkId);

    /**
     * @Description Obtaining ip Request Information
     *
     * @param info - IpInfo object
     * @return ErrCode - operation result
     */
    ErrCode GetIpInfo(IpInfo &info);

    /**
     * @Description Reconnect to the currently active network
     *
     * @return ErrCode - operation result
     */
    ErrCode ReConnect();

    /**
     * @Description ReAssociate network
     *
     * @return ErrCode - operation result
     */
    ErrCode ReAssociate(void);

    /**
     * @Description Enable WPS connection
     *
     * @param config - WpsConfig object
     * @return ErrCode - operation result
     */
    ErrCode StartWps(const WpsConfig &config);

    /**
     * @Description Close the WPS connection
     *
     * @return ErrCode - operation result
     */
    ErrCode CancelWps(void);

    /**
     * @Description  Get the device MAC address.
     *
     * @param result - Get device mac String
     * @return ErrCode - operation result
     */
    ErrCode GetDeviceMacAddress(std::string &result) override;

private:
    int systemAbilityId_;
#ifdef OHOS_ARCH_LITE
    IWifiDevice *client_;
#else
    sptr<IWifiDevice> client_;
#endif
};
}  // namespace Wifi
}  // namespace OHOS
#endif