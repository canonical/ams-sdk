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

	api "github.com/anbox-cloud/ams-sdk/api/ams"
	restapi "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/api"
	restclient "github.com/anbox-cloud/ams-sdk/pkg/ams/shared/rest/client"
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
	ListContainersWithFilters(filters []string) ([]api.Container, error)
	LaunchContainer(details *api.ContainersPost, noWait bool) (restclient.Operation, error)
	RetrieveContainerByID(id string) (*api.Container, string, error)
	UpdateContainerByID(id string, details *api.ContainerPatch, noWait bool) (restclient.Operation, error)
	DeleteContainerByID(id string, force bool) (restclient.Operation, error)
	DeleteContainers(ids []string, force bool) (restclient.Operation, error)
	RetrieveContainerLog(id, name string, downloader func(header *http.Header, body io.ReadCloser) error) error
	ExecuteContainer(id string, details *api.ContainerExecPost, args *ContainerExecArgs) (restclient.Operation, error)

	// Instances
	ListInstances() ([]api.Instance, error)
	ListInstancesWithFilters(filters []string) ([]api.Instance, error)
	LaunchInstance(details *api.InstancesPost, noWait bool) (restclient.Operation, error)
	RetrieveInstanceByID(id string) (*api.Instance, string, error)
	UpdateInstanceByID(id string, details *api.InstancePatch, noWait bool) (restclient.Operation, error)
	DeleteInstanceByID(id string, force bool) (restclient.Operation, error)
	DeleteInstances(ids []string, force bool) (restclient.Operation, error)
	RetrieveInstanceLog(id, name string, downloader func(header *http.Header, body io.ReadCloser) error) error
	ExecuteInstance(id string, details *api.InstanceExecPost, args *InstanceExecArgs) (restclient.Operation, error)

	// Shares
	CreateInstanceShare(id string, details *api.InstanceSharesPost) (*api.InstanceSharesPostResponse, error)
	UpdateInstanceShareByID(instanceID, shareID string, details *api.InstanceSharePatch) error
	DeleteInstanceShareByID(instanceID, shareID string) (restclient.Operation, error)

	// Config
	SetConfigItem(name, value string) error
	RetrieveConfigItems() (map[string]interface{}, error)

	// Applications
	CreateApplication(packagePath string, sentBytes chan float64) (restclient.Operation, error)
	CreateApplicationWithArgs(args *ApplicationCreateArgs) (restclient.Operation, error)
	UpdateApplicationWithPackage(id, packagePath string, sentBytes chan float64) (restclient.Operation, error)
	UpdateApplicationWithDetails(id string, details api.ApplicationPatch) error
	UpdateApplication(id string) (restclient.Operation, error)
	ListApplications() ([]api.Application, error)
	ListApplicationsWithFilters(filters []string) ([]api.Application, error)
	FindApplicationsByName(pattern string) ([]api.Application, error)
	RetrieveApplicationByID(id string) (*api.Application, string, error)
	DeleteApplicationByID(id string, force bool) (restclient.Operation, error)
	DeleteApplications(ids []string, force bool) (restclient.Operation, error)
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
	ImportImage(name, path string, isDefault bool) (restclient.Operation, error)
	ImportImageByType(name, path string, imgType api.ImageType, isDefault bool) (restclient.Operation, error)
	SetDefaultImage(id string) error
	DeleteImageByIDOrName(id string, force bool, imgType api.ImageType) (restclient.Operation, error)
	DeleteImageVersion(id string, version int) (restclient.Operation, error)
	RetrieveImageByIDOrName(id string, imgType api.ImageType) (*api.Image, string, error)
	RetrieveDefaultImage() (*api.Image, string, error)
	TriggerImageSync(id string) error

	// Services
	RetrieveServiceStatus() (*api.ServiceStatus, string, error)
	HasExtension(name string) (bool, error)
	ListTasks() ([]api.Task, error)
	GetVersion() (string, error)

	// Registry
	ListApplicationsFromRegistry() ([]api.RegistryApplication, error)
	PushApplicationToRegistry(id string) (restclient.Operation, error)
	PullApplicationFromRegistry(id string) (restclient.Operation, error)
	DeleteApplicationFromRegistry(id string) (restclient.Operation, error)

	GetEvents() (*restclient.EventListener, error)

	// Operations
	ListOperations() (map[string][]*restapi.Operation, error)
	ShowOperation(id string) (*restapi.Operation, error)
	CancelOperation(id string) error

	// Auth
	GetOIDCConfig(grantType string) (*restapi.OIDCResponse, string, error)
	CreateIdentity(details *api.IdentityPost) (restclient.Operation, error)
	ListIdentitiesWithFilters(filters []string) ([]api.Identity, error)
	RetrieveIdentityByID(id string) (*api.Identity, string, error)
	DeleteIdentity(id string, force bool) (restclient.Operation, error)
	SetGroupsForIdentity(id string, groups []string) (restclient.Operation, error)

	CreateAuthGroup(details *api.AuthGroupPost) (restclient.Operation, error)
	RetrieveAuthGroup(name string) (*api.AuthGroup, string, error)
	ListAuthGroupsWithFilters(filters []string) ([]api.AuthGroup, error)
	DeleteAuthGroup(name string, force bool) (restclient.Operation, error)
	UpdateAuthGroupDescription(name, description string) (restclient.Operation, error)
	SetPermissionsForGroup(name string, permissions []api.Permission) (restclient.Operation, error)
}

// clientImpl encapsulates a client to the AMS service and allows performing
// various operations with the service
type clientImpl struct {
	restclient.Client
	serviceStatus      *api.ServiceStatus
	hasInstanceSupport bool
}

// New creates a new client talking to the AMS service at the specified URL or unix.socket path
// and with the specified tls config, if provided
func New(addr any, tlsConfig *tls.Config) (Client, error) {
	c, err := restclient.NewTLSClient(addr, tlsConfig)
	if err != nil {
		return nil, err
	}

	client := clientImpl{Client: c}
	client.hasInstanceSupport, err = client.HasExtension("instance_support")
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// NewOIDCClient creates a new client OIDC talking to the AMS service at the specified URL
// and with the specified OIDC token
func NewOIDCClient(addr any, tlsConfig *tls.Config, tokenProvider restclient.TokenProvider) (Client, error) {
	c, err := restclient.NewOIDCClient(addr, tlsConfig, tokenProvider)
	if err != nil {
		return nil, err
	}
	client := clientImpl{Client: c}
	client.hasInstanceSupport, err = client.HasExtension("instance_support")
	if err != nil {
		return nil, err
	}
	return &client, nil
}
