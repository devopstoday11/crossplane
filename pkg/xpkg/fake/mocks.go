/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/crossplane/crossplane/pkg/xpkg"
)

var _ xpkg.Cache = &MockCache{}

// MockCache is a mock Cache.
type MockCache struct {
	MockGet    func() (v1.Image, error)
	MockStore  func() error
	MockDelete func() error
}

// NewMockCacheGetFn creates a new MockGet function for MockCache.
func NewMockCacheGetFn(img v1.Image, err error) func() (v1.Image, error) {
	return func() (v1.Image, error) { return img, err }
}

// NewMockCacheStoreFn creates a new MockStore function for MockCache.
func NewMockCacheStoreFn(err error) func() error {
	return func() error { return err }
}

// NewMockCacheDeleteFn creates a new MockDelete function for MockCache.
func NewMockCacheDeleteFn(err error) func() error {
	return func() error { return err }
}

// Get calls the underlying MockGet.
func (c *MockCache) Get(source, id string) (v1.Image, error) {
	return c.MockGet()
}

// Store calls the underlying MockStore.
func (c *MockCache) Store(source, id string, img v1.Image) error {
	return c.MockStore()
}

// Delete calls the underlying MockDelete.
func (c *MockCache) Delete(id string) error {
	return c.MockDelete()
}

var _ xpkg.Fetcher = &MockFetcher{}

// MockFetcher is a mock fetcher.
type MockFetcher struct {
	MockFetch func() (v1.Image, error)
}

// NewMockFetchFn creates a new MockFetch function for MockFetcher.
func NewMockFetchFn(img v1.Image, err error) func() (v1.Image, error) {
	return func() (v1.Image, error) { return img, err }
}

// Fetch calls the underlying MockFetch.
func (m *MockFetcher) Fetch(ctx context.Context, ref name.Reference, secrets []string) (v1.Image, error) {
	return m.MockFetch()
}
