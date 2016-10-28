// Copyright 2016 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package invoke_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rancher/cni/pkg/invoke"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindInPath", func() {
	var (
		multiplePaths  []string
		pluginName     string
		pluginDir      string
		anotherTempDir string
	)

	BeforeEach(func() {
		tempDir, err := ioutil.TempDir("", "cni-find")
		Expect(err).NotTo(HaveOccurred())

		plugin, err := ioutil.TempFile(tempDir, "a-cni-plugin")

		anotherTempDir, err = ioutil.TempDir("", "nothing-here")

		multiplePaths = []string{anotherTempDir, tempDir}
		pluginDir, pluginName = filepath.Split(plugin.Name())
	})

	Context("when multiple paths are provided", func() {
		It("returns only the path to the plugin", func() {
			pluginPath, err := invoke.FindInPath(pluginName, multiplePaths)
			Expect(err).NotTo(HaveOccurred())
			Expect(pluginPath).To(Equal(filepath.Join(pluginDir, pluginName)))
		})
	})

	Context("when an error occurs", func() {
		Context("when no paths are provided", func() {
			It("returns an error noting no paths were provided", func() {
				_, err := invoke.FindInPath(pluginName, []string{})
				Expect(err).To(MatchError("no paths provided"))
			})
		})

		Context("when no plugin is provided", func() {
			It("returns an error noting the plugin name wasn't found", func() {
				_, err := invoke.FindInPath("", multiplePaths)
				Expect(err).To(MatchError("no plugin name provided"))
			})
		})

		Context("when the plugin cannot be found", func() {
			It("returns an error noting the path", func() {
				pathsWithNothing := []string{anotherTempDir}
				_, err := invoke.FindInPath(pluginName, pathsWithNothing)
				Expect(err).To(MatchError(fmt.Sprintf("failed to find plugin %q in path %s", pluginName, pathsWithNothing)))
			})
		})
	})
})
