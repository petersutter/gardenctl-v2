/*
SPDX-FileCopyrightText: Contributors to the Gardener project

SPDX-License-Identifier: Apache-2.0
*/

package tools_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hack Scripts Test Suite")
}
