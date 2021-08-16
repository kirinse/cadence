// Copyright (c) 2021 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package frontend

import (
	"context"

	"github.com/uber/cadence/common/authorization"
	"github.com/uber/cadence/common/config"
	"github.com/uber/cadence/common/resource"
	"github.com/uber/cadence/common/types"
)

// AccessControlledWorkflowAdminHandler frontend handler wrapper for authentication and authorization
type AccessControlledWorkflowAdminHandler struct {
	resource.Resource

	adminHandler AdminHandler
	authorizer      authorization.Authorizer
}

var _ AdminHandler = (*AccessControlledWorkflowAdminHandler)(nil)

// NewAccessControlledHandlerImpl creates frontend handler with authentication support
func NewAccessControlledAdminHandlerImpl(adminHandler AdminHandler, resource resource.Resource, authorizer authorization.Authorizer, cfg config.Authorization) *AccessControlledWorkflowAdminHandler {
	if authorizer == nil {
		authorizer = authorization.NewAuthorizer(cfg, resource.GetLogger(), resource.GetDomainCache())
	}
	return &AccessControlledWorkflowAdminHandler{
		Resource:        resource,
		adminHandler: adminHandler,
		authorizer:      authorizer,
	}
}

func (a *AccessControlledWorkflowAdminHandler) AddSearchAttribute(ctx context.Context, request *types.AddSearchAttributeRequest) error {
	attr := &authorization.Attributes{
		APIName:    "AddSearchAttribute",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.AddSearchAttribute(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) CloseShard(ctx context.Context, request *types.CloseShardRequest) error {
	attr := &authorization.Attributes{
		APIName:    "CloseShard",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.CloseShard(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) DescribeCluster(ctx context.Context) (*types.DescribeClusterResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "DescribeCluster",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.DescribeCluster(ctx)
}

func (a *AccessControlledWorkflowAdminHandler) DescribeShardDistribution(ctx context.Context, request *types.DescribeShardDistributionRequest) (*types.DescribeShardDistributionResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "DescribeShardDistribution",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.DescribeShardDistribution(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) DescribeHistoryHost(ctx context.Context, request *types.DescribeHistoryHostRequest) (*types.DescribeHistoryHostResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "DescribeHistoryHost",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.DescribeHistoryHost(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) DescribeQueue(ctx context.Context, request *types.DescribeQueueRequest) (*types.DescribeQueueResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "DescribeQueue",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.DescribeQueue(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) DescribeWorkflowExecution(ctx context.Context, request *types.AdminDescribeWorkflowExecutionRequest) (*types.AdminDescribeWorkflowExecutionResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "DescribeWorkflowExecution",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.DescribeWorkflowExecution(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) GetDLQReplicationMessages(ctx context.Context, request *types.GetDLQReplicationMessagesRequest) (*types.GetDLQReplicationMessagesResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "GetDLQReplicationMessages",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.GetDLQReplicationMessages(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) GetDomainReplicationMessages(ctx context.Context, request *types.GetDomainReplicationMessagesRequest) (*types.GetDomainReplicationMessagesResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "GetDomainReplicationMessages",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.GetDomainReplicationMessages(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) GetReplicationMessages(ctx context.Context, request *types.GetReplicationMessagesRequest) (*types.GetReplicationMessagesResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "GetReplicationMessages",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.GetReplicationMessages(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) GetWorkflowExecutionRawHistoryV2(ctx context.Context, request *types.GetWorkflowExecutionRawHistoryV2Request) (*types.GetWorkflowExecutionRawHistoryV2Response, error) {
	attr := &authorization.Attributes{
		APIName:    "GetWorkflowExecutionRawHistoryV2",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.GetWorkflowExecutionRawHistoryV2(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) MergeDLQMessages(ctx context.Context, request *types.MergeDLQMessagesRequest) (*types.MergeDLQMessagesResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "MergeDLQMessages",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errUnauthorized
	}

	return a.adminHandler.MergeDLQMessages(ctx,request)
}

func (a *AccessControlledWorkflowAdminHandler) PurgeDLQMessages(ctx context.Context, request *types.PurgeDLQMessagesRequest) error {
	attr := &authorization.Attributes{
		APIName:    "PurgeDLQMessages",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.PurgeDLQMessages(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) ReadDLQMessages(ctx context.Context, request *types.ReadDLQMessagesRequest) (*types.ReadDLQMessagesResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "ReadDLQMessages",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil,err
	}
	if !isAuthorized {
		return nil,errUnauthorized
	}

	return a.adminHandler.ReadDLQMessages(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) ReapplyEvents(ctx context.Context, request *types.ReapplyEventsRequest) error {
	attr := &authorization.Attributes{
		APIName:    "ReapplyEvents",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.ReapplyEvents(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) RefreshWorkflowTasks(ctx context.Context, request *types.RefreshWorkflowTasksRequest) error {
	attr := &authorization.Attributes{
		APIName:    "RefreshWorkflowTasks",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.RefreshWorkflowTasks(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) RemoveTask(ctx context.Context, request *types.RemoveTaskRequest) error {
	attr := &authorization.Attributes{
		APIName:    "RemoveTask",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.RemoveTask(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) ResendReplicationTasks(ctx context.Context, request *types.ResendReplicationTasksRequest) error {
	attr := &authorization.Attributes{
		APIName:    "ResendReplicationTasks",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.ResendReplicationTasks(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) ResetQueue(ctx context.Context, request *types.ResetQueueRequest) error {
	attr := &authorization.Attributes{
		APIName:    "ResetQueue",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errUnauthorized
	}

	return a.adminHandler.ResetQueue(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) GetCrossClusterTasks(ctx context.Context, request *types.GetCrossClusterTasksRequest) (*types.GetCrossClusterTasksResponse, error) {
	attr := &authorization.Attributes{
		APIName:    "GetCrossClusterTasks",
		Permission: authorization.PermissionAdmin,
	}
	isAuthorized, err := a.isAuthorized(ctx, attr)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil,errUnauthorized
	}

	return a.adminHandler.GetCrossClusterTasks(ctx, request)
}

func (a *AccessControlledWorkflowAdminHandler) isAuthorized(
	ctx context.Context,
	attr *authorization.Attributes,
) (bool, error) {
	result, err := a.authorizer.Authorize(ctx, attr)
	if err != nil {
		return false, err
	}
	isAuth := result.Decision == authorization.DecisionAllow
	return isAuth, nil
}