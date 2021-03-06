// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/duyquang6/git-watchdog/pkg/dto"
	mock "github.com/stretchr/testify/mock"

	null "github.com/duyquang6/git-watchdog/pkg/null"
)

// RepositoryService is an autogenerated mock type for the RepositoryService type
type RepositoryService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, req
func (_m *RepositoryService) Create(ctx context.Context, req dto.CreateRepositoryRequest) (*dto.CreateRepositoryResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *dto.CreateRepositoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateRepositoryRequest) *dto.CreateRepositoryResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.CreateRepositoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.CreateRepositoryRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *RepositoryService) Delete(ctx context.Context, id uint) (*dto.DeleteRepositoryResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *dto.DeleteRepositoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, uint) *dto.DeleteRepositoryResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.DeleteRepositoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOne provides a mock function with given fields: ctx, id
func (_m *RepositoryService) GetOne(ctx context.Context, id uint) (*dto.GetOneRepositoryResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *dto.GetOneRepositoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, uint) *dto.GetOneRepositoryResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.GetOneRepositoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IssueScanRepo provides a mock function with given fields: ctx, repoID
func (_m *RepositoryService) IssueScanRepo(ctx context.Context, repoID uint) (*dto.IssueScanResponse, error) {
	ret := _m.Called(ctx, repoID)

	var r0 *dto.IssueScanResponse
	if rf, ok := ret.Get(0).(func(context.Context, uint) *dto.IssueScanResponse); ok {
		r0 = rf(ctx, repoID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.IssueScanResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, repoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListScan provides a mock function with given fields: ctx, repoID, page, limit
func (_m *RepositoryService) ListScan(ctx context.Context, repoID null.Uint, page uint, limit uint) (*dto.ListScanResponse, error) {
	ret := _m.Called(ctx, repoID, page, limit)

	var r0 *dto.ListScanResponse
	if rf, ok := ret.Get(0).(func(context.Context, null.Uint, uint, uint) *dto.ListScanResponse); ok {
		r0 = rf(ctx, repoID, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.ListScanResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, null.Uint, uint, uint) error); ok {
		r1 = rf(ctx, repoID, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, req
func (_m *RepositoryService) Update(ctx context.Context, id uint, req dto.UpdateRepositoryRequest) (*dto.UpdateRepositoryResponse, error) {
	ret := _m.Called(ctx, id, req)

	var r0 *dto.UpdateRepositoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, uint, dto.UpdateRepositoryRequest) *dto.UpdateRepositoryResponse); ok {
		r0 = rf(ctx, id, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.UpdateRepositoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint, dto.UpdateRepositoryRequest) error); ok {
		r1 = rf(ctx, id, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
