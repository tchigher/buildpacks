// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package acceptance

import (
	"testing"

	"github.com/GoogleCloudPlatform/buildpacks/pkg/acceptance"
)

func init() {
	acceptance.DefineFlags()
}

func TestAcceptance(t *testing.T) {
	// TODO(b/163086885): Remove after publishing images.
	if acceptance.PullImages() {
		t.Skip("Tests are skipped until stack images are published to gcr.io")
	}
	builder, cleanup := acceptance.CreateBuilder(t)
	t.Cleanup(cleanup)

	testCases := []acceptance.Test{
		{
			App: "no_requirements_txt",
		},
		{
			App: "requirements_txt",
		},
		{
			App: "pip_dependency",
		},
		{
			App: "gunicorn_present",
		},
		{
			App: "gunicorn_outdated",
		},
		{
			App: "custom_entrypoint",
			Env: []string{"GOOGLE_ENTRYPOINT=uwsgi --http :$PORT --wsgi-file custom.py --callable app"},
		},
		{
			Name: "custom gunicorn entrypoint",
			App:  "gunicorn_present",
			Env:  []string{"GOOGLE_ENTRYPOINT=gunicorn main:app"},
		},
	}
	for _, tc := range testCases {
		tc := tc
		if tc.Name == "" {
			tc.Name = tc.App
		}
		tc.Env = append(tc.Env, "GOOGLE_RUNTIME=python39")

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			acceptance.TestApp(t, builder, tc)
		})
	}
}

func TestFailures(t *testing.T) {
	builder, cleanup := acceptance.CreateBuilder(t)
	t.Cleanup(cleanup)

	testCases := []acceptance.FailureTest{
		{
			App:       "pip_check",
			MustMatch: `sub-dependency-\w 1\.0\.0 has requirement sub-dependency-\w.*1\.0\.0, but you.* have sub-dependency-\w 1\.0\.0`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.App, func(t *testing.T) {
			t.Parallel()

			tc.Env = append(tc.Env, "GOOGLE_RUNTIME=python39")

			acceptance.TestBuildFailure(t, builder, tc)
		})
	}
}
