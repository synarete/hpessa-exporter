// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"os"
	"os/user"
	"path"
	"runtime"

	"github.com/go-logr/logr"
)

var (
	version = "devel"
	gitrev  = "unknown"
)

const (
	HostnameEnvKey     = "HOSTNAME"
	PodNameEnvKey      = "POD_NAME"
	PodNamespaceEnvKey = "POD_NAMESPACE"
	PodIPEnvKey        = "POD_IP"
	HostIPEnvKey       = "HOST_IP"
)

type Ident struct {
	Progname  string     `json:"program"`
	Version   string     `json:"version"`
	Hostname  string     `json:"hostname"`
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	HostIP    string     `json:"hostip"`
	PodIP     string     `json:"podip"`
	User      *user.User `json:"user"`
}

func SelfIdent() *Ident {
	return &Ident{
		Progname:  Progname(),
		Version:   Version(),
		Hostname:  os.Getenv(HostnameEnvKey),
		Name:      os.Getenv(PodNameEnvKey),
		Namespace: os.Getenv(PodNamespaceEnvKey),
		PodIP:     os.Getenv(PodIPEnvKey),
		HostIP:    os.Getenv(HostIPEnvKey),
		User:      currentUser(),
	}
}

func currentUser() *user.User {
	currUser, err := user.Current()
	if err != nil || currUser == nil {
		return &user.User{
			Username: "unknown",
			Name:     "unknown",
		}
	}
	return currUser
}

func Progname() string {
	return path.Base(os.Args[0])
}

func Version() string {
	return version + "-" + gitrev
}

func LogSelfIdent(log logr.Logger) {
	ident := SelfIdent()
	progname := ident.Progname
	log.Info(progname, "Version", ident.Version)
	log.Info(progname, "Name", ident.Name)
	log.Info(progname, "Namespace", ident.Namespace)
	log.Info(progname, "PodIP", ident.PodIP)
	log.Info(progname, "HostIP", ident.HostIP)
	log.Info(progname, "GOOS", runtime.GOOS)
	log.Info(progname, "GOARCH", runtime.GOARCH)
}
