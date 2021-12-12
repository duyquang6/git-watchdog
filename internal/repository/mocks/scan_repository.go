// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "github.com/duyquang6/git-watchdog/internal/model"
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"

	null "github.com/duyquang6/git-watchdog/pkg/null"
)

// ScanRepository is an autogenerated mock type for the ScanRepository type
type ScanRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: tx, scan
func (_m *ScanRepository) Create(tx *gorm.DB, scan *model.Scan) error {
	ret := _m.Called(tx, scan)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.Scan) error); ok {
		r0 = rf(tx, scan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: tx, id
func (_m *ScanRepository) Delete(tx *gorm.DB, id uint) error {
	ret := _m.Called(tx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) error); ok {
		r0 = rf(tx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: tx, id
func (_m *ScanRepository) GetByID(tx *gorm.DB, id uint) (model.Scan, error) {
	ret := _m.Called(tx, id)

	var r0 model.Scan
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) model.Scan); ok {
		r0 = rf(tx, id)
	} else {
		r0 = ret.Get(0).(model.Scan)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: tx, repoID, offset, limit
func (_m *ScanRepository) List(tx *gorm.DB, repoID null.Uint, offset uint, limit uint) ([]model.Scan, uint, error) {
	ret := _m.Called(tx, repoID, offset, limit)

	var r0 []model.Scan
	if rf, ok := ret.Get(0).(func(*gorm.DB, null.Uint, uint, uint) []model.Scan); ok {
		r0 = rf(tx, repoID, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Scan)
		}
	}

	var r1 uint
	if rf, ok := ret.Get(1).(func(*gorm.DB, null.Uint, uint, uint) uint); ok {
		r1 = rf(tx, repoID, offset, limit)
	} else {
		r1 = ret.Get(1).(uint)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*gorm.DB, null.Uint, uint, uint) error); ok {
		r2 = rf(tx, repoID, offset, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Update provides a mock function with given fields: tx, scan
func (_m *ScanRepository) Update(tx *gorm.DB, scan *model.Scan) error {
	ret := _m.Called(tx, scan)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.Scan) error); ok {
		r0 = rf(tx, scan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}