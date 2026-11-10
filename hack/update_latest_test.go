/*
SPDX-FileCopyrightText: Contributors to the Gardener project

SPDX-License-Identifier: Apache-2.0
*/

package hack_test

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("update-latest.sh", func() {
	var (
		script     string
		latestFile string
	)

	BeforeEach(func() {
		if _, err := exec.LookPath("bash"); err != nil {
			Skip("bash is not available")
		}

		var err error

		script, err = filepath.Abs("update-latest.sh")
		Expect(err).NotTo(HaveOccurred())

		latestFile = filepath.Join(GinkgoT().TempDir(), "LATEST")
	})

	writeLatest := func(content string) {
		Expect(os.WriteFile(latestFile, []byte(content), 0o600)).To(Succeed())
	}

	readLatest := func() string {
		content, err := os.ReadFile(latestFile)
		Expect(err).NotTo(HaveOccurred())

		return string(content)
	}

	run := func(args ...string) (string, error) {
		cmd := exec.Command("bash", append([]string{script}, args...)...) // #nosec G204 -- test executes the script under test with fixed arguments
		out, err := cmd.CombinedOutput()

		return string(out), err
	}

	DescribeTable("should move LATEST forward to a newer stable release",
		func(current, released string) {
			writeLatest(current + "\n")

			out, err := run(released, latestFile)
			Expect(err).NotTo(HaveOccurred(), out)
			Expect(readLatest()).To(Equal(released + "\n"))
		},
		Entry("patch release from a hotfix branch", "v2.19.0", "v2.19.1"),
		Entry("minor release", "v2.19.1", "v2.20.0"),
		Entry("major release", "v2.19.1", "v3.0.0"),
		Entry("numeric (not lexicographic) comparison", "v2.9.9", "v2.10.0"),
		Entry("components with leading zeros are compared in base 10", "v2.08.0", "v2.09.0"),
	)

	It("should normalize a LATEST file without a trailing newline", func() {
		writeLatest("v2.19.0")

		out, err := run("v2.19.1", latestFile)
		Expect(err).NotTo(HaveOccurred(), out)
		Expect(readLatest()).To(Equal("v2.19.1\n"))
	})

	DescribeTable("should leave LATEST untouched",
		func(current, released string) {
			writeLatest(current + "\n")

			out, err := run(released, latestFile)
			Expect(err).NotTo(HaveOccurred(), out)
			Expect(readLatest()).To(Equal(current + "\n"))
		},
		Entry("for a patch release of an older minor version", "v2.20.0", "v2.19.2"),
		Entry("for an older patch release of the same minor version", "v2.19.2", "v2.19.1"),
		Entry("for the version LATEST already references", "v2.19.1", "v2.19.1"),
		Entry("for a release candidate", "v2.19.1", "v2.20.0-rc.1"),
		Entry("for a development version", "v2.19.1", "v2.20.0-dev"),
		Entry("for a version with build metadata", "v2.19.1", "v2.20.0+build.1"),
	)

	DescribeTable("should fail without touching LATEST",
		func(current, released string) {
			writeLatest(current)

			out, err := run(released, latestFile)
			Expect(err).To(HaveOccurred(), out)
			Expect(out).To(ContainSubstring("error:"))
			Expect(readLatest()).To(Equal(current))
		},
		Entry("for an invalid released version", "v2.19.0\n", "2.19.1"),
		Entry("for a released version that is not semver", "v2.19.0\n", "v2.19"),
		Entry("for a LATEST file that does not reference a stable version", "v2.20.0-dev\n", "v2.20.0"),
		Entry("for an empty LATEST file", "", "v2.20.0"),
	)

	It("should fail if the LATEST file does not exist", func() {
		out, err := run("v2.19.1", latestFile)
		Expect(err).To(HaveOccurred(), out)
		Expect(latestFile).NotTo(BeAnExistingFile())
	})
})
