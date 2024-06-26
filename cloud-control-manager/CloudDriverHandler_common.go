// Cloud Driver Manager of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2020.12.

package clouddriverhandler

import (
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	icon "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/connect"
	icdrs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	im "github.com/cloud-barista/cb-spider/cloud-info-manager"
	ccim "github.com/cloud-barista/cb-spider/cloud-info-manager/connection-config-info-manager"
	cim "github.com/cloud-barista/cb-spider/cloud-info-manager/credential-info-manager"
	dim "github.com/cloud-barista/cb-spider/cloud-info-manager/driver-info-manager"
	rim "github.com/cloud-barista/cb-spider/cloud-info-manager/region-info-manager"

	"fmt"
	"strings"
)

/*
func ListCloudDriver() []string {
        var cloudDriverList []string
        // @todo get list from storage
        return cloudDriverList
}
*/

// 1. get the ConnectionConfig Info
// 2. get the driver info
// 3. get CloudDriver
func GetCloudDriver(cloudConnectName string) (idrv.CloudDriver, error) {
	cccInfo, err := ccim.GetConnectionConfig(cloudConnectName)
	if err != nil {
		return nil, err
	}

	cldDrvInfo, err := dim.GetCloudDriver(cccInfo.DriverName)
	if err != nil {
		return nil, err
	}

	return getCloudDriver(*cldDrvInfo)
}

// 1. get the driver info
// 2. get CloudDriver
func getCloudDriverByDriverName(driverName string) (idrv.CloudDriver, error) {
	cldDrvInfo, err := dim.GetCloudDriver(driverName)
	if err != nil {
		return nil, err
	}

	return getCloudDriver(*cldDrvInfo)
}

func getProviderNameByDriverName(driverName string) (string, error) {
	cldDrvInfo, err := dim.GetCloudDriver(driverName)
	if err != nil {
		return "", err
	}

	return cldDrvInfo.ProviderName, nil
}

// CloudConnection for Region-Level Control (Except. DiskHandler)
func GetCloudConnection(cloudConnectName string) (icon.CloudConnection, error) {
	conn, err := commonGetCloudConnection(cloudConnectName, "")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// CloudConnection for Zone-Level Control (Ex. DiskHandler)
func GetZoneLevelCloudConnection(cloudConnectName string, targetZoneName string) (icon.CloudConnection, error) {
	conn, err := commonGetCloudConnection(cloudConnectName, targetZoneName)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// 1. get credential info
// 2. get region info
// 3. get CloudConneciton
func commonGetCloudConnection(cloudConnectName string, targetZoneName string) (icon.CloudConnection, error) {
	cccInfo, err := ccim.GetConnectionConfig(cloudConnectName)
	if err != nil {
		return nil, err
	}

	cldDriver, err := GetCloudDriver(cloudConnectName)
	if err != nil {
		return nil, err
	}

	crdInfo, err := cim.GetCredentialDecrypt(cccInfo.CredentialName)
	if err != nil {
		return nil, err
	}

	rgnInfo, err := rim.GetRegion(cccInfo.RegionName)
	if err != nil {
		return nil, err
	}

	regionName, zoneName, err := getRegionNameByRegionInfo(rgnInfo)
	if err != nil {
		return nil, err
	}

	connectionInfo := idrv.ConnectionInfo{ // @todo powerkim
		CredentialInfo: idrv.CredentialInfo{
			ClientId:         getValue(crdInfo.KeyValueInfoList, "ClientId"),
			ClientSecret:     getValue(crdInfo.KeyValueInfoList, "ClientSecret"),
			TenantId:         getValue(crdInfo.KeyValueInfoList, "TenantId"),
			SubscriptionId:   getValue(crdInfo.KeyValueInfoList, "SubscriptionId"),
			IdentityEndpoint: getValue(crdInfo.KeyValueInfoList, "IdentityEndpoint"),
			Username:         getValue(crdInfo.KeyValueInfoList, "Username"),
			Password:         getValue(crdInfo.KeyValueInfoList, "Password"),
			DomainName:       getValue(crdInfo.KeyValueInfoList, "DomainName"),
			ProjectID:        getValue(crdInfo.KeyValueInfoList, "ProjectID"),
			AuthToken:        getValue(crdInfo.KeyValueInfoList, "AuthToken"),
			ClientEmail:      getValue(crdInfo.KeyValueInfoList, "ClientEmail"),
			PrivateKey:       getValue(crdInfo.KeyValueInfoList, "PrivateKey"),
			Host:             getValue(crdInfo.KeyValueInfoList, "Host"),
			APIVersion:       getValue(crdInfo.KeyValueInfoList, "APIVersion"),
			MockName:         getValue(crdInfo.KeyValueInfoList, "MockName"),
			ApiKey:           getValue(crdInfo.KeyValueInfoList, "ApiKey"),
			ClusterId:        getValue(crdInfo.KeyValueInfoList, "ClusterId"),
			ConnectionName:   cloudConnectName,
		},
		RegionInfo: idrv.RegionInfo{ // @todo powerkim
			Region:     regionName,
			Zone:       zoneName,       // default Zone
			TargetZone: targetZoneName, // Target Zone for Zone-Level Control(Ex. DiskHandler)
		},
	}

	cldConnection, err := cldDriver.ConnectCloud(connectionInfo)
	if err != nil {
		return nil, err
	}

	return cldConnection, nil
}

// 1. get credential info
// 2. get region info
// 3. get CloudConneciton
func GetCloudConnectionByDriverNameAndCredentialName(driverName string, credentialName string) (icon.CloudConnection, error) {
	cldDriver, err := getCloudDriverByDriverName(driverName)
	if err != nil {
		return nil, err
	}

	crdInfo, err := cim.GetCredentialDecrypt(credentialName)
	if err != nil {
		return nil, err
	}

	providerName, err := getProviderNameByDriverName(driverName)
	if err != nil {
		return nil, err
	}

	connectionInfo := idrv.ConnectionInfo{ // @todo powerkim
		CredentialInfo: idrv.CredentialInfo{
			ClientId:         getValue(crdInfo.KeyValueInfoList, "ClientId"),
			ClientSecret:     getValue(crdInfo.KeyValueInfoList, "ClientSecret"),
			TenantId:         getValue(crdInfo.KeyValueInfoList, "TenantId"),
			SubscriptionId:   getValue(crdInfo.KeyValueInfoList, "SubscriptionId"),
			IdentityEndpoint: getValue(crdInfo.KeyValueInfoList, "IdentityEndpoint"),
			Username:         getValue(crdInfo.KeyValueInfoList, "Username"),
			Password:         getValue(crdInfo.KeyValueInfoList, "Password"),
			DomainName:       getValue(crdInfo.KeyValueInfoList, "DomainName"),
			ProjectID:        getValue(crdInfo.KeyValueInfoList, "ProjectID"),
			AuthToken:        getValue(crdInfo.KeyValueInfoList, "AuthToken"),
			ClientEmail:      getValue(crdInfo.KeyValueInfoList, "ClientEmail"),
			PrivateKey:       getValue(crdInfo.KeyValueInfoList, "PrivateKey"),
			Host:             getValue(crdInfo.KeyValueInfoList, "Host"),
			APIVersion:       getValue(crdInfo.KeyValueInfoList, "APIVersion"),
			MockName:         getValue(crdInfo.KeyValueInfoList, "MockName"),
			ApiKey:           getValue(crdInfo.KeyValueInfoList, "ApiKey"),
			ClusterId:        getValue(crdInfo.KeyValueInfoList, "ClusterId"),
		},
	}

	// get Provider's Meta Info for default region
	cloudOSMetaInfo, err := im.GetCloudOSMetaInfo(providerName)
	if err != nil {
		cblog.Error(err)
		return nil, err
	}
	if cloudOSMetaInfo.DefaultRegionToQuery != nil {
		if Length := len(cloudOSMetaInfo.DefaultRegionToQuery); Length == 1 {
			connectionInfo.RegionInfo.Region = cloudOSMetaInfo.DefaultRegionToQuery[0]
		} else if Length == 2 {
			connectionInfo.RegionInfo.Region = cloudOSMetaInfo.DefaultRegionToQuery[0]
			connectionInfo.RegionInfo.Zone = cloudOSMetaInfo.DefaultRegionToQuery[1]
		}
	}

	cldConnection, err := cldDriver.ConnectCloud(connectionInfo)
	if err != nil {
		return nil, err
	}

	return cldConnection, nil
}

func GetProviderNameByConnectionName(cloudConnectName string) (string, error) {
	cccInfo, err := ccim.GetConnectionConfig(cloudConnectName)
	if err != nil {
		return "", err
	}

	rgnInfo, err := rim.GetRegion(cccInfo.RegionName)
	if err != nil {
		return "", err
	}

	return rgnInfo.ProviderName, nil
}

func GetRegionNameByConnectionName(cloudConnectName string) (string, string, error) {
	cccInfo, err := ccim.GetConnectionConfig(cloudConnectName)
	if err != nil {
		return "", "", err
	}

	rgnInfo, err := rim.GetRegion(cccInfo.RegionName)
	if err != nil {
		return "", "", err
	}

	return getRegionNameByRegionInfo(rgnInfo)
}

func getRegionNameByRegionInfo(rgnInfo *rim.RegionInfo) (string, string, error) {

	// @todo should move KeyValueList into XXXDriver.go, powerkim
	var regionName string
	var zoneName string
	switch strings.ToUpper(rgnInfo.ProviderName) {
	case "AWS", "AZURE", "ALIBABA", "GCP", "TENCENT", "IBM", "NCP", "NCPVPC", "KTCLOUD", "NHNCLOUD", "KTCLOUDVPC":
		regionName = getValue(rgnInfo.KeyValueInfoList, "Region")
		zoneName = getValue(rgnInfo.KeyValueInfoList, "Zone")
	case "OPENSTACK", "CLOUDIT", "DOCKER", "CLOUDTWIN", "MOCK":
		regionName = getValue(rgnInfo.KeyValueInfoList, "Region")
	default:
		errmsg := rgnInfo.ProviderName + " is not a valid ProviderName!!"
		return "", "", fmt.Errorf(errmsg)
	}

	return regionName, zoneName, nil
}

func getValue(keyValueInfoList []icdrs.KeyValue, key string) string {
	for _, kv := range keyValueInfoList {
		if strings.EqualFold(kv.Key, key) { // ignore case
			return kv.Value
		}
	}
	return "Not set"
}
