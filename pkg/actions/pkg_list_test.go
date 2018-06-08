// Copyright 2018 The ksonnet authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package actions

import (
	"bytes"
	"testing"

	"github.com/ksonnet/ksonnet/pkg/app"
	amocks "github.com/ksonnet/ksonnet/pkg/app/mocks"
	"github.com/ksonnet/ksonnet/pkg/registry"
	"github.com/stretchr/testify/require"
)

func TestPkgList(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		libaries := app.LibraryRefSpecs{
			"lib1": &app.LibraryRefSpec{},
		}

		appMock.On("Libraries").Return(libaries, nil)

		spec := &registry.Spec{
			Libraries: registry.LibraryRefSpecs{
				"lib1": &registry.LibraryRef{},
				"lib2": &registry.LibraryRef{},
			},
		}

		incubator := mockRegistry("incubator", false)
		incubator.On("FetchRegistrySpec").Return(spec, nil)

		cases := []struct {
			name          string
			onlyInstalled bool
			outputName    string
		}{
			{
				name:       "list all packages",
				outputName: "pkg/list/output.txt",
			},
			{
				name:          "installed packages",
				onlyInstalled: true,
				outputName:    "pkg/list/installed.txt",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				in := map[string]interface{}{
					OptionApp:       appMock,
					OptionInstalled: tc.onlyInstalled,
				}

				a, err := NewPkgList(in)
				require.NoError(t, err)

				a.registryListFn = func(app.App) ([]registry.Registry, error) {
					registries := []registry.Registry{incubator}
					return registries, nil
				}

				var buf bytes.Buffer
				a.out = &buf

				err = a.Run()
				require.NoError(t, err)

				assertOutput(t, tc.outputName, buf.String())
			})
		}

	})
}

func TestPkgList_requires_app(t *testing.T) {
	in := make(map[string]interface{})
	_, err := NewPkgList(in)
	require.Error(t, err)
}
