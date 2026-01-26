package util

import (
	"errors"
	"net/http"
)

type DeviceInfo struct {
	DeviceID        string
	DeviceType      string
	DeviceModel     string
	DeviceOS        string
	DeviceOSVersion string
	UserAgent       string
	IPAddress       string
}

func ExtractDeviceInfo(r *http.Request) (*DeviceInfo, error) {
	deviceID := r.Header.Get("X-Device-ID")
	if deviceID == "" {
		return nil, errors.New("X-Device-ID header is required")
	}

	deviceType := r.Header.Get("X-Device-Type")
	deviceModel := r.Header.Get("X-Device-Model")
	deviceOS := r.Header.Get("X-Device-OS")

	if deviceType == "" || deviceModel == "" || deviceOS == "" {
		return nil, errors.New("device info headers (X-Device-Type, X-Device-Model, X-Device-OS) are required")
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	return &DeviceInfo{
		DeviceID:        deviceID,
		DeviceType:      deviceType,
		DeviceModel:     deviceModel,
		DeviceOS:        deviceOS,
		DeviceOSVersion: r.Header.Get("X-Device-OS-Version"),
		UserAgent:       r.Header.Get("User-Agent"),
		IPAddress:       ip,
	}, nil
}
