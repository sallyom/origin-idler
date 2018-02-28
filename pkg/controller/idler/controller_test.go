


package idler_test

import (
    . "github.com/openshift/origin-idler/pkg/apis/idling/v1alpha2"
    . "github.com/openshift/origin-idler/pkg/client/clientset_generated/clientset/typed/idling/v1alpha2"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic tests

var _ = Describe("Idler controller", func() {
    var instance Idler
    var expectedKey string
    var client IdlerInterface

    BeforeEach(func() {
        instance = Idler{}
        instance.Name = "instance-1"
        expectedKey = "default/instance-1"
    })

    AfterEach(func() {
        client.Delete(instance.Name, &metav1.DeleteOptions{})
    })

    Describe("when creating a new object", func() {
        It("invoke the reconcile method", func() {
            after := make(chan struct{})
            controller.AfterReconcile = func(key string, err error) {
                defer func() {
                    // Recover in case the key is reconciled multiple times
                    defer func() { recover() }()
                    close(after)
                }()
                Expect(key).To(Equal(expectedKey))
                Expect(err).ToNot(HaveOccurred())
            }

            // Create the instance
            client = cs.IdlingV1alpha2().Idlers("default")
            _, err := client.Create(&instance)
            Expect(err).ShouldNot(HaveOccurred())

            // Wait for reconcile to happen
            Eventually(after).Should(BeClosed())

            // INSERT YOUR CODE HERE - test conditions post reconcile
        })
    })
})
