package common

import (
	"FGW_WEB/pkg/common/msg"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// InfoPC содержит информацию о ПК.
type InfoPC struct {
	Domain string `json:"domain"`
	IPAddr string `json:"ipAddr"`
}

// NewInfoPC возвращает новый экземпляр InfoPC.
func NewInfoPC() (*InfoPC, error) {
	hostname, err := hostName()
	if err != nil {
		return nil, err
	}

	ipStr, err := ipAddr()
	if err != nil {
		return nil, err
	}

	return &InfoPC{
		Domain: hostname,
		IPAddr: ipStr,
	}, nil
}

// HostName возвращает имя хоста.
func (i *InfoPC) HostName() string {
	return i.Domain
}

// AddrIP возвращает IP адрес.
func (i *InfoPC) AddrIP() string {
	return i.IPAddr
}

// hostName возвращает имя текущего хоста.
func hostName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s: %s", msg.E3100, err.Error()))
	}

	return hostname, nil
}

// ipAddr возвращает IP адрес текущего хоста.
func ipAddr() (string, error) {
	hostname, err := hostName()
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s: %s", msg.E3101, err.Error()))
	}

	ipStr := make([]string, len(ips))
	for _, ip := range ips {
		ipStr = append(ipStr, ip.String())
	}

	return strings.Join(ipStr, " | "), nil
}
