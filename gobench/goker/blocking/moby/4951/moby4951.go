/*
 * Project: moby
 * Issue or PR  : https://github.com/moby/moby/pull/4951
 * Buggy version: 81f148be566ab2b17810ad4be61a5d8beac8330f
 * fix commit-id: 2ffef1b7eb618162673c6ffabccb9ca57c7dfce3
 * Flaky: 100/100
 * Description:
 *   The root cause and patch is clearly explained in the commit
 * description. The global lock is devices.Lock(), and the device
 * lock is baseInfo.lock.Lock(). It is very likely that this bug
 * can be reproduced.
 */
package main

import (
	"sync"

	"time"
)

type DeviceSet struct {
	mu               sync.Mutex
	infos            *DevInfo
	nrDeletedDevices int
}

func (devices *DeviceSet) DeleteDevice(hash string) {
	devices.mu.Lock()
	defer devices.mu.Unlock()

	devices.infos.lock.Lock()
	defer devices.infos.lock.Unlock()

	devices.deleteDevice(devices.infos)
}

func (devices *DeviceSet) deleteDevice(info *DevInfo) {
	devices.removeDeviceAndWait(info.Name())
}

func (devices *DeviceSet) removeDeviceAndWait(devname string) {
	/// remove devices by devname
	devices.mu.Unlock()
	time.Sleep(300 * time.Nanosecond)
	devices.mu.Lock()
}

type DevInfo struct {
	lock sync.Mutex
	name string
}

func (info *DevInfo) Name() string {
	return info.name
}

func NewDeviceSet() *DeviceSet {
	devices := &DeviceSet{
		infos: &DevInfo{
			name: "info1",
		},
	}

	return devices
}

func main() {

	ds := NewDeviceSet()
	/// Delete devices by the same info
	go ds.DeleteDevice("info1")
	go ds.DeleteDevice("info1")
}



