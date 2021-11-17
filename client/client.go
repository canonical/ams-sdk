// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the Lesser GNU General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the Lesser GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package client

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/anbox-cloud/ams-sdk/api"
	restapi "github.com/anbox-cloud/ams-sdk/shared/rest/api"
	"github.com/anbox-cloud/ams-sdk/shared/rest/client"
	restclient "github.com/anbox-cloud/ams-sdk/shared/rest/client"
)

const (
	// By default we're creating always applications of type 'game'
	defaultAppType           = "game"
	extendedTransportTimeout = 300 * time.Second
)

// Client is the interface used to communicate with an AMS server
type Client interface {
	// Nodes
	ListNodes() ([]api.Node, error)
	AddNode(node *api.NodesPost) (restclient.Operation, error)
	RemoveNode(name string, force, keepInCluster bool) (restclient.Operation, error)
	RetrieveNodeByName(name string) (*api.Node, string, error)
	UpdateNode(name string, details *api.NodePatch) (restclient.Operation, error)

	// Certificates
	ListCertificates() ([]restapi.Certificate, error)
	AddCertificate(details *restapi.CertificatesPost) (*restapi.Response, error)
	DeleteCertificate(fingerprint string) error

	// Containers
	ListContainers() ([]api.Container, error)
	LaunchContainer(details *api.ContainersPost) (restclient.Operation, error)
	RetrieveContainerByID(id string) (*api.Container, string, error)
	DeleteContainerByID(id string, force bool) (restclient.Operation, error)
	RetrieveContainerLog(id, name string, downloader func(header *http.Header, body io.ReadCloser) error) error
	ExecuteContainer(id string, details *api.ContainerExecPost, args *ContainerExecArgs) (restclient.Operation, error)

	// Config
	SetConfigItem(name, value string) error
	RetrieveConfigItems() (map[string]interface{}, error)

	// Applications
	CreateApplication(packagePath string, sentBytes chan float64) (restclient.Operation, error)
	UpdateApplicationWithPackage(id, packagePath string, sentBytes chan float64) (restclient.Operation, error)
	UpdateApplicationWithDetails(id string, details api.ApplicationPatch) error
	UpdateApplication(id string) (restclient.Operation, error)
	ListApplications() ([]api.Application, error)
	FindApplicationsByName(pattern string) ([]api.Application, error)
	RetrieveApplicationByID(id string) (*api.Application, string, error)
	DeleteApplicationByID(id string, force bool) (restclient.Operation, error)
	ExportApplicationByVersion(id string, version int, downloader func(header *http.Header, body io.ReadCloser) error) error
	PublishApplicationVersion(id string, version int) (restclient.Operation, error)
	RevokeApplicationVersion(id string, version int) (restclient.Operation, error)
	DeleteApplicationVersion(id string, version int, force bool) (restclient.Operation, error)

	// Addons
	AddAddon(name string, packagePath string, sentBytes chan float64) (restclient.Operation, error)
	UpdateAddon(name, packagePath string, sentBytes chan float64) (restclient.Operation, error)
	RetrieveAddon(name string) (*api.Addon, string, error)
	DeleteAddon(name string) (restclient.Operation, error)
	DeleteAddonVersion(name string, version int) (restclient.Operation, error)
	ListAddons() ([]api.Addon, error)

	// Images
	ListImages() ([]api.Image, error)
	AddImage(name, packagePath string, isDefault bool, sentBytes chan float64) (restclient.Operation, error)
	UpdateImage(id, packagePath string, sentBytes chan float64) (restclient.Operation, error)
	ImportImage(name, path string, isDefault bool) (client.Operation, error)
	SetDefaultImage(id string) error
	DeleteImageByIDOrName(id string, force bool) (restclient.Operation, error)
	DeleteImageVersion(id string, version int) (restclient.Operation, error)
	RetrieveImageByIDOrName(id string) (*api.Image, string, error)
	RetrieveDefaultImage() (*api.Image, string, error)

	// Services
	RetrieveServiceStatus() (*api.ServiceStatus, string, error)
	HasExtension(name string) bool
	ListTasks() ([]api.Task, error)
	GetVersion() (string, error)

	// Registry
	ListApplicationsFromRegistry() ([]api.RegistryApplication, error)
	PushApplicationToRegistry(id string) (client.Operation, error)
	PullApplicationFromRegistry(id string) (client.Operation, error)
	DeleteApplicationFromRegistry(id string) (client.Operation, error)

	GetEvents() (*restclient.EventListener, error)

	// Operations
	ListOperations() (map[string][]*restapi.Operation, error)
	ShowOperation(id string) (*restapi.Operation, error)
	CancelOperation(id string) error
}

// clientImpl encapsulates a client to the AMS service and allows performing
// various operations with the service
type clientImpl struct {
	restclient.Client
	serviceStatus *api.ServiceStatus
}

// New creates a new client talking to the AMS service at the specified URL or unix.socket path
// and with the specified tls config, if provided
func New(addr interface{}, tlsConfig *tls.Config) (Client, error) {
	c, err := restclient.New(addr, tlsConfig)
	if err != nil {
		return nil, err
	}
	return &clientImpl{Client: c}, nil
}
