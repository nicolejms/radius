// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	"github.com/project-radius/radius/pkg/azure/azresources"
	"github.com/project-radius/radius/pkg/radlogger"
	"github.com/project-radius/radius/pkg/rp/armerrors"
	"github.com/project-radius/radius/pkg/ucp/resources"
)

const (
	LinkedResourceUpdateErrorFormat = "Attempted to deploy existing resource '%s' which has a different application and/or environment. Options to resolve the conflict are: change the name of the '%s' resource in %s to create a new resource, or use '%s' application and '%s' environment to update the existing resource '%s'."
)

// Translation of internal representation of health state to user facing values
var InternalToUserHealthStateTranslation = map[string]string{
	HealthStateUnknown:       HealthStateUnhealthy,
	HealthStateHealthy:       HealthStateHealthy,
	HealthStateUnhealthy:     HealthStateUnhealthy,
	HealthStateDegraded:      HealthStateDegraded,
	HealthStateNotSupported:  "",
	HealthStateNotApplicable: HealthStateHealthy,
	HealthStateError:         HealthStateUnhealthy,
}

// Response represents a category of HTTP response (eg. OK with payload).
type Response interface {
	// Apply modifies the ResponseWriter to send the desired details back to the client.
	Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error
}

// OKResponse represents an HTTP 200 with a JSON payload.
//
// This is used when modification to an existing resource is processed synchronously.
type OKResponse struct {
	Body    interface{}
	Headers map[string]string
}

// NewOKResponse creates an OKResponse that will write a 200 OK with the provided body as JSON.
// Set the body to nil to write an empty 200 OK.
func NewOKResponse(body interface{}) Response {
	return &OKResponse{Body: body}
}

// NewOKResponseWithHeaders creates an OKResponse that will write a 200 OK with the provided body as JSON.
// Set the body to nil to write an empty 200 OK.
func NewOKResponseWithHeaders(body interface{}, headers map[string]string) Response {
	return &OKResponse{
		Body:    body,
		Headers: headers,
	}
}

func (r *OKResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusOK), radlogger.LogHTTPStatusCode, http.StatusOK)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	for key, element := range r.Headers {
		w.Header().Add(key, element)
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}
	return nil
}

// CreatedResponse represents an HTTP 201 with a JSON payload.
//
// This is used when a request to create a new resource is processed synchronously.
type CreatedResponse struct {
	Body interface{}
}

func NewCreatedResponse(body interface{}) Response {
	response := &CreatedResponse{Body: body}
	return response
}

func (r *CreatedResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusCreated), radlogger.LogHTTPStatusCode, http.StatusCreated)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// CreatedAsyncResponse represents an HTTP 201 with a JSON payload and location header.
//
// This is used when a request to create a new resource is processed asynchronously.
type CreatedAsyncResponse struct {
	Body     interface{}
	Location string
	Scheme   string
}

func NewCreatedAsyncResponse(body interface{}, location string, scheme string) Response {
	return &CreatedAsyncResponse{Body: body, Location: location, Scheme: scheme}
}

func (r *CreatedAsyncResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusCreated), radlogger.LogHTTPStatusCode, http.StatusCreated)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	location := url.URL{
		Host:   req.Host,
		Scheme: req.URL.Scheme,
		Path:   r.Location,
	}

	// In production this is the header we get from app service for the 'real' protocol
	protocol := req.Header.Get(textproto.CanonicalMIMEHeaderKey("X-Forwarded-Proto"))
	if protocol != "" {
		location.Scheme = protocol
	}

	if location.Scheme == "" {
		location.Scheme = r.Scheme
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", location.String())
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// AcceptedAsyncResponse represents an HTTP 202 with a JSON payload and location header.
//
// This is used when a request to create an existing resource is processed asynchronously.
type AcceptedAsyncResponse struct {
	Body     interface{}
	Location string
	Scheme   string
}

// NewAcceptedAsyncResponse creates an AcceptedAsyncResponse
func NewAcceptedAsyncResponse(body interface{}, location string, scheme string) Response {
	return &AcceptedAsyncResponse{Body: body, Location: location, Scheme: scheme}
}

func (r *AcceptedAsyncResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusAccepted), radlogger.LogHTTPStatusCode, http.StatusAccepted)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	location := url.URL{
		Host:   req.Host,
		Scheme: req.URL.Scheme,
		Path:   r.Location,
	}

	// In production this is the header we get from app service for the 'real' protocol
	protocol := req.Header.Get(textproto.CanonicalMIMEHeaderKey("X-Forwarded-Proto"))
	if protocol != "" {
		location.Scheme = protocol
	}

	if location.Scheme == "" {
		location.Scheme = r.Scheme
	}

	logger.Info(fmt.Sprintf("Returning location: %s", location.String()))

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", location.String())
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// AsyncOperationResponse represents the response for an async operation request.
type AsyncOperationResponse struct {
	Body        interface{}
	Location    string
	Code        int
	ResourceID  resources.ID
	OperationID uuid.UUID
	APIVersion  string
}

// NewAsyncOperationResponse creates an AsyncOperationResponse
func NewAsyncOperationResponse(body interface{}, location string, code int, resourceID resources.ID, operationID uuid.UUID, apiVersion string) Response {
	return &AsyncOperationResponse{
		Body:        body,
		Location:    location,
		Code:        code,
		ResourceID:  resourceID,
		OperationID: operationID,
		APIVersion:  apiVersion,
	}
}

func (r *AsyncOperationResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	// Write Body
	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	locationHeader := r.getAsyncLocationPath(req, "operationResults")
	azureAsyncOpHeader := r.getAsyncLocationPath(req, "operationStatuses")

	// Write Headers
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", locationHeader)
	w.Header().Add("Azure-AsyncOperation", azureAsyncOpHeader)
	w.Header().Add("Retry-After", v1.DefaultRetryAfter)

	w.WriteHeader(r.Code)

	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// getAsyncLocationPath returns the async operation location path for the given resource type.
func (r *AsyncOperationResponse) getAsyncLocationPath(req *http.Request, resourceType string) string {
	dest := url.URL{
		Host:   req.Host,
		Scheme: req.URL.Scheme,
		Path: fmt.Sprintf("%s/providers/%s/locations/%s/%s/%s", r.ResourceID.PlaneScope(),
			r.ResourceID.ProviderNamespace(), r.Location, resourceType, r.OperationID.String()),
	}

	query := url.Values{}
	query.Add("api-version", r.APIVersion)
	dest.RawQuery = query.Encode()

	// In production this is the header we get from app service for the 'real' protocol
	protocol := req.Header.Get("X-Forwarded-Proto")
	if protocol != "" {
		dest.Scheme = protocol
	}

	if dest.Scheme == "" {
		dest.Scheme = "http"
	}

	return dest.String()
}

// NoContentResponse represents an HTTP 204.
//
// This is used for delete operations.
type NoContentResponse struct {
}

func NewNoContentResponse() Response {
	return &NoContentResponse{}
}

func (r *NoContentResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	w.WriteHeader(204)
	return nil
}

// BadRequestResponse represents an HTTP 400 with an error message in ARM error format.
//
// This is used for any operation that fails due to bad data with a simple error message.
type BadRequestResponse struct {
	Body armerrors.ErrorResponse
}

// NewLinkedResourceUpdateErrorResponse represents a HTTP 400 with an error message when user updates environment id and application id.
func NewLinkedResourceUpdateErrorResponse(resourceID resources.ID, oldProp *v1.BasicResourceProperties, newProp *v1.BasicResourceProperties) Response {
	newAppEnv := ""
	if newProp.Application != "" {
		name := newProp.Application
		if rid, err := resources.Parse(newProp.Application); err == nil {
			name = rid.Name()
		}
		newAppEnv += fmt.Sprintf("'%s' application", name)
	}
	if newProp.Environment != "" {
		if newAppEnv != "" {
			newAppEnv += " and "
		}
		name := newProp.Environment
		if rid, err := resources.Parse(newProp.Environment); err == nil {
			name = rid.Name()
		}
		newAppEnv += fmt.Sprintf("'%s' environment", name)
	}

	message := fmt.Sprintf(LinkedResourceUpdateErrorFormat, resourceID.Name(), resourceID.Name(), newAppEnv, oldProp.Application, oldProp.Environment, resourceID.Name())
	return &BadRequestResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.Invalid,
				Message: message,
				Target:  resourceID.String(),
			},
		},
	}
}

func NewBadRequestResponse(message string) Response {
	return &BadRequestResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.Invalid,
				Message: message,
			},
		},
	}
}

func NewBadRequestARMResponse(body armerrors.ErrorResponse) Response {
	return &BadRequestResponse{
		Body: body,
	}
}

func (r *BadRequestResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusBadRequest), radlogger.LogHTTPStatusCode, http.StatusBadRequest)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// ValidationErrorResponse represents an HTTP 400 with validation errors in ARM error format.
type ValidationErrorResponse struct {
	Body armerrors.ErrorResponse
}

func NewValidationErrorResponse(errors validator.ValidationErrors) Response {
	body := armerrors.ErrorResponse{
		Error: armerrors.ErrorDetails{
			Code:    armerrors.Invalid,
			Message: errors.Error(),
		},
	}

	for _, fe := range errors {
		if err, ok := fe.(error); ok {
			detail := armerrors.ErrorDetails{
				Target:  fe.Field(),
				Message: err.Error(),
			}
			body.Error.Details = append(body.Error.Details, detail)
		}
	}

	return &ValidationErrorResponse{Body: body}
}

func (r *ValidationErrorResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusBadRequest), radlogger.LogHTTPStatusCode, http.StatusBadRequest)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// NotFoundResponse represents an HTTP 404 with an ARM error payload.
//
// This is used for GET operations when the response does not exist.
type NotFoundResponse struct {
	Body armerrors.ErrorResponse
}

func NewLegacyNotFoundResponse(id azresources.ResourceID) Response {
	return &NotFoundResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.NotFound,
				Message: fmt.Sprintf("the resource with id '%s' was not found", id.ID),
				Target:  id.ID,
			},
		},
	}
}

// NewNotFoundMessageResponse represents an HTTP 404 with string message.
func NewNotFoundMessageResponse(message string) Response {
	return &NotFoundResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.NotFound,
				Message: message,
			},
		},
	}
}

func NewNotFoundResponse(id resources.ID) Response {
	return &NotFoundResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.NotFound,
				Message: fmt.Sprintf("the resource with id '%s' was not found", id.String()),
				Target:  id.String(),
			},
		},
	}
}

// NewNotFoundAPIVersionResponse creates Response for unsupported api version. (message is consistent with ARM).
func NewNotFoundAPIVersionResponse(resourceType string, namespace string, apiVersion string) Response {
	return &NotFoundResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.InvalidResourceType, // ARM uses "InvalidResourceType" code with 404 http code.
				Message: fmt.Sprintf("The resource type '%s' could not be found in the namespace '%s' for api version '%s'.", resourceType, namespace, apiVersion),
			},
		},
	}
}

func (r *NotFoundResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusNotFound), radlogger.LogHTTPStatusCode, http.StatusNotFound)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// ConflictResponse represents an HTTP 409 with an ARM error payload.
//
// This is used for delete operations.
type ConflictResponse struct {
	Body armerrors.ErrorResponse
}

func NewConflictResponse(message string) Response {
	return &ConflictResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.Conflict,
				Message: message,
			},
		},
	}
}

func (r *ConflictResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusConflict), radlogger.LogHTTPStatusCode, http.StatusConflict)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

type InternalServerErrorResponse struct {
	Body armerrors.ErrorResponse
}

func NewInternalServerErrorARMResponse(body armerrors.ErrorResponse) Response {
	return &InternalServerErrorResponse{
		Body: body,
	}
}

func (r *InternalServerErrorResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusInternalServerError), radlogger.LogHTTPStatusCode, http.StatusInternalServerError)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// GetUserFacingResourceHealthState computes the aggregate health state to be shown to the user
// It also modifies the state of the individual output resources to the user facing value as needed
func GetUserFacingResourceHealthState(restOutputResources []OutputResource) (string, string) {
	aggregateHealthState := HealthStateHealthy
	aggregateHealthStateErrorDetails := ""
	foundNotSupported := false
	foundHealthyOrUnhealthy := false

	for i, or := range restOutputResources {
		userHealthState := InternalToUserHealthStateTranslation[or.Status.HealthState]
		if userHealthState != or.Status.HealthState {
			// Set the individual output resource to the user facing value
			restOutputResources[i].Status.HealthState = userHealthState
		}

		switch or.Status.HealthState {
		case HealthStateUnknown:
			aggregateHealthState = userHealthState
			aggregateHealthStateErrorDetails = "Health state unknown"
		case HealthStateHealthy:
			foundHealthyOrUnhealthy = true
		case HealthStateUnhealthy:
			// If any one of the resources is unhealthy, the aggregate is unhealthy
			aggregateHealthState = userHealthState
			foundHealthyOrUnhealthy = true
		case HealthStateNotSupported:
			// If any one of the resources is not supported, the user facing aggregate is ""
			aggregateHealthState = userHealthState
			foundNotSupported = true
		case HealthStateNotApplicable:
			// This case is ignored and has no effect on the aggregate state
		default:
			// Unexpected state
			or.Status.HealthState = InternalToUserHealthStateTranslation[HealthStateUnhealthy]
			aggregateHealthStateErrorDetails = fmt.Sprintf("output resource found in unexpected state: %s", or.Status.HealthState)
		}
	}

	if foundNotSupported && foundHealthyOrUnhealthy {
		// We do not expect a combination of not supported and supported health reporting for output resources
		// This will result in an aggregation logic error
		aggregateHealthState = InternalToUserHealthStateTranslation[HealthStateError]
		aggregateHealthStateErrorDetails = "Health aggregation error"
	}

	return aggregateHealthState, aggregateHealthStateErrorDetails
}

func GetUserFacingResourceProvisioningState(restOutputResources []OutputResource) string {
	var aggregateProvisiongState = ProvisioningStateProvisioned
forLoop:
	for _, or := range restOutputResources {
		switch or.Status.ProvisioningState {
		case ProvisioningStateFailed:
			// If any of the output resources is Failed, then the aggregate is Failed
			aggregateProvisiongState = ProvisioningStateFailed
			break forLoop
		case ProvisioningStateProvisioning, ProvisioningStateNotProvisioned:
			// If any of the output resources is not in Provisioned state, the aggregate is Provisioning
			aggregateProvisiongState = ProvisioningStateProvisioning
		}
	}
	return aggregateProvisiongState
}

// GetUserFacingAppHealthState computes the aggregate application health based on the input child resource status
// It accepts a map with key as resource name and status as the resource status and returns the aggregate health
// state and health state error details
func GetUserFacingAppHealthState(resourceStatuses map[string]ResourceStatus) (string, string) {
	aggregateHealthState := HealthStateHealthy
	aggregateHealthStateErrorDetails := ""

forloop:
	for r, rs := range resourceStatuses {
		userHealthState := InternalToUserHealthStateTranslation[rs.HealthState]

		switch rs.HealthState {
		case HealthStateUnknown:
			aggregateHealthState = userHealthState
			aggregateHealthStateErrorDetails = fmt.Sprintf("Resource %s has unknown health state", r)
		case HealthStateHealthy:
			// No change since default aggregated value is Healthy
		case HealthStateUnhealthy:
			// If any one of the resources is unhealthy, the aggregate is unhealthy
			aggregateHealthState = userHealthState
			aggregateHealthStateErrorDetails = fmt.Sprintf("Resource %s is unhealthy", r)
		case HealthStateNotSupported:
			// Will ignore NotSupported state for aggregation at application level
		default:
			// Unexpected state
			rs.HealthState = InternalToUserHealthStateTranslation[HealthStateUnhealthy]
			aggregateHealthStateErrorDetails = fmt.Sprintf("Resource %s found in unexpected state: %s", r, rs.HealthState)
		}

		if userHealthState == HealthStateUnhealthy {
			break forloop
		}
	}

	return aggregateHealthState, aggregateHealthStateErrorDetails
}

// GetUserFacingAppProvisioningState computes the aggregate application provisioning state based on the input
// child resource status. It accepts a map with key as resource name and status as the resource status and
// returns the aggregate provisioning state and provisioning state error details
func GetUserFacingAppProvisioningState(statuses map[string]ResourceStatus) (string, string) {
	var aggregateProvisiongState = ProvisioningStateProvisioned
	var aggregateProvisiongStateErrorDetails string
forLoop:
	for r, rs := range statuses {
		switch rs.ProvisioningState {
		case ProvisioningStateFailed:
			// If any of the resources is Failed, then the aggregate is Failed
			aggregateProvisiongState = ProvisioningStateFailed
			aggregateProvisiongStateErrorDetails = fmt.Sprintf("Resource %s is in Failed state", r)
			break forLoop
		case ProvisioningStateProvisioning, ProvisioningStateNotProvisioned:
			// If any of the resources is not in Provisioned state, the aggregate is Provisioning
			aggregateProvisiongState = ProvisioningStateProvisioning
			aggregateProvisiongStateErrorDetails = fmt.Sprintf("Resource %s is in %s state", r, rs.ProvisioningState)
		}
	}
	return aggregateProvisiongState, aggregateProvisiongStateErrorDetails
}

// PreconditionFailedResponse represents an HTTP 412 with an ARM error payload.
type PreconditionFailedResponse struct {
	Body armerrors.ErrorResponse
}

func NewPreconditionFailedResponse(target string, message string) Response {
	return &PreconditionFailedResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.PreconditionFailed,
				Message: message,
				Target:  target,
			},
		},
	}
}

func (r *PreconditionFailedResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusPreconditionFailed), radlogger.LogHTTPStatusCode, http.StatusPreconditionFailed)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusPreconditionFailed)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}

// ClientAuthenticationFailed represents an HTTP 401 with an ARM error payload.
type ClientAuthenticationFailed struct {
	Body armerrors.ErrorResponse
}

func NewClientAuthenticationFailedARMResponse() Response {
	return &ClientAuthenticationFailed{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.InvalidAuthenticationInfo,
				Message: "Server failed to authenticate the request",
			},
		},
	}
}
func (r *ClientAuthenticationFailed) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusUnauthorized), radlogger.LogHTTPStatusCode, http.StatusUnauthorized)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}
	return nil
}

// AsyncOperationResultResponse
type AsyncOperationResultResponse struct {
	Headers map[string]string
}

func NewAsyncOperationResultResponse(headers map[string]string) Response {
	return &AsyncOperationResultResponse{
		Headers: headers,
	}
}

func (r *AsyncOperationResultResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusAccepted), radlogger.LogHTTPStatusCode, http.StatusAccepted)

	w.Header().Add("Content-Type", "application/json")

	for key, element := range r.Headers {
		w.Header().Add(key, element)
	}

	w.WriteHeader(http.StatusAccepted)

	return nil
}

// MethodNotAllowedResponse represents an HTTP 405 with an ARM error payload.
type MethodNotAllowedResponse struct {
	Body armerrors.ErrorResponse
}

// NewMethodNotAllowedResponse creates MethodNotAllowedResponse instance.
func NewMethodNotAllowedResponse(target string, message string) Response {
	return &MethodNotAllowedResponse{
		Body: armerrors.ErrorResponse{
			Error: armerrors.ErrorDetails{
				Code:    armerrors.Invalid,
				Message: message,
				Target:  target,
			},
		},
	}
}

func (r *MethodNotAllowedResponse) Apply(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("responding with status code: %d", http.StatusMethodNotAllowed), radlogger.LogHTTPStatusCode, http.StatusMethodNotAllowed)

	bytes, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling %T: %w", r.Body, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing marshaled %T bytes to output: %s", r.Body, err)
	}

	return nil
}