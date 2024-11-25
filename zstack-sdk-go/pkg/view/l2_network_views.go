// Copyright (c) ZStack.io, Inc.

package view

type L2NetworkInventoryView struct {
	BaseInfoView
	BaseTimeView

	Vlan                 int      `json:"vlan"`
	ZoneUuid             string   `json:"zoneUuid"` // Zone UUID. If specified, the cloud host will be created in the specified zone.
	PhysicalInterface    string   `json:"physicalInterface"`
	Type                 string   `json:"type"`
	AttachedClusterUuids []string `json:"attachedClusterUuids"`
}

type ClusterHostNetworkFactsView struct {
	Bondings []HostBondings `json:"bondings"`
	Nics     []HostNics     `json:"nics"`
	Success  bool           `json:"success"`
}

type HostNics struct {
	CarrierActive    bool   `json:"carrierActive"`
	CreateDate       string `json:"createDate"`
	HostUuid         string `json:"hostUuid"`
	InterfaceName    string `json:"interfaceName"`
	InterfaceType    string `json:"interfaceType"`
	LastOpDate       string `json:"lastOpDate"`
	Mac              string `json:"mac"`
	PciDeviceAddress string `json:"pciDeviceAddress"`
	SlaveActive      bool   `json:"slaveActive"`
	Speed            int    `json:"speed"`
	Uuid             string `json:"uuid"`
}

type HostBondings struct {
	AllSlavesActive bool   `json:"allSlavesActive"`
	BondingName     string `json:"bondingName"`
	CreateDate      string `json:"createDate"`
	HostUuid        string `json:"hostUuid"`
	LastOpDate      string `json:"lastOpDate"`
	Mac             string `json:"mac"`
	MiiStatus       string `json:"miiStatus"`
	Miimon          int    `json:"miimon"`
	Mode            string `json:"mode"`
	Slaves          []struct {
		BondingUuid      string `json:"bondingUuid"`
		CarrierActive    bool   `json:"carrierActive"`
		CreateDate       string `json:"createDate"`
		HostUuid         string `json:"hostUuid"`
		InterfaceName    string `json:"interfaceName"`
		InterfaceType    string `json:"interfaceType"`
		LastOpDate       string `json:"lastOpDate"`
		Mac              string `json:"mac"`
		PciDeviceAddress string `json:"pciDeviceAddress"`
		SlaveActive      bool   `json:"slaveActive"`
		Speed            int    `json:"speed"`
		Uuid             string `json:"uuid"`
	} `json:"slaves"`
	Uuid           string `json:"uuid"`
	XmitHashPolicy string `json:"xmitHashPolicy"`
}
