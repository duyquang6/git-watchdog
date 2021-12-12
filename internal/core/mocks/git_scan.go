// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	core "github.com/duyquang6/git-watchdog/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// GitScan is an autogenerated mock type for the GitScan type
type GitScan struct {
	mock.Mock
}

// Scan provides a mock function with given fields: repoName, repoURL
func (_m *GitScan) Scan(repoName string, repoURL string) ([]core.Finding, error) {
	ret := _m.Called(repoName, repoURL)

	var r0 []core.Finding
	if rf, ok := ret.Get(0).(func(string, string) []core.Finding); ok {
		r0 = rf(repoName, repoURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Finding)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(repoName, repoURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}