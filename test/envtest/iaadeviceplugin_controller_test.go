// Copyright 2021-2022 Intel Corporation. All Rights Reserved.
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

package envtest

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	devicepluginv1 "github.com/intel/intel-device-plugins-for-kubernetes/pkg/apis/deviceplugin/v1"
)

var _ = Describe("IaaDevicePlugin Controller", func() {

	const timeout = time.Second * 30
	const interval = time.Second * 1

	Context("Basic CRUD operations", func() {
		It("should handle IaaDevicePlugin objects correctly", func() {
			spec := devicepluginv1.IaaDevicePluginSpec{
				Image: "testimage",
			}

			key := types.NamespacedName{
				Name:      "iaadeviceplugin-test",
				Namespace: ns,
			}

			toCreate := &devicepluginv1.IaaDevicePlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: ns,
				},
				Spec: spec,
			}

			By("creating IaaDevicePlugin successfully")
			Expect(k8sClient.Create(context.Background(), toCreate)).Should(Succeed())
			time.Sleep(time.Second * 5)

			fetched := &devicepluginv1.IaaDevicePlugin{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), key, fetched)
				return len(fetched.Status.ControlledDaemonSet.UID) > 0
			}, timeout, interval).Should(BeTrue())

			By("updating image name successfully")
			updatedImage := "updated-iaa-testimage"
			fetched.Spec.Image = updatedImage

			Expect(k8sClient.Update(context.Background(), fetched)).Should(Succeed())
			fetchedUpdated := &devicepluginv1.IaaDevicePlugin{}
			Eventually(func() string {
				_ = k8sClient.Get(context.Background(), key, fetchedUpdated)
				return fetchedUpdated.Spec.Image
			}, timeout, interval).Should(Equal(updatedImage))

			By("deleting IaaDevicePlugin successfully")
			Eventually(func() error {
				f := &devicepluginv1.IaaDevicePlugin{}
				_ = k8sClient.Get(context.Background(), key, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			Eventually(func() error {
				f := &devicepluginv1.IaaDevicePlugin{}
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})

	It("upgrades", func() {
		dp := &devicepluginv1.IaaDevicePlugin{}

		var image string

		testUpgrade("iaa", dp, &image, nil)

		Expect(dp.Spec.Image == image).To(BeTrue())
	})
})
