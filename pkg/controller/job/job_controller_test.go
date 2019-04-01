/*
Copyright 2019 The Upbound Authors.

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

package job

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	batchv1alpha1 "github.com/saravanakumar-periyasamy/kubebuilder-demo/pkg/apis/batch/v1alpha1"
	"golang.org/x/net/context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "job-a", Namespace: "default"}}
var jobKey = types.NamespacedName{Name: "job-a", Namespace: "default"}

const timeout = time.Minute * 5

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	job_c := &batchv1alpha1.Job{ObjectMeta: metav1.ObjectMeta{Name: "job-c", Namespace: "default"}}
	job_b := &batchv1alpha1.Job{ObjectMeta: metav1.ObjectMeta{Name: "job-b", Namespace: "default"}, Spec: batchv1alpha1.JobSpec{DependOnJobs: []string{"job-c"}}}
	job_a := &batchv1alpha1.Job{ObjectMeta: metav1.ObjectMeta{Name: "job-a", Namespace: "default"}, Spec: batchv1alpha1.JobSpec{DependOnJobs: []string{"job-b"}}}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the Job object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), job_c)
	err = c.Create(context.TODO(), job_b)
	err = c.Create(context.TODO(), job_a)

	// The job_a object may not be a valid object because it might be missing some required fields.
	// Please modify the job_a object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), job_a)
	defer c.Delete(context.TODO(), job_b)
	defer c.Delete(context.TODO(), job_c)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	//Eventually the job should get processed
	g.Eventually(func() batchv1alpha1.State {
		instance := &batchv1alpha1.Job{}
		c.Get(context.TODO(), jobKey, instance)
		return instance.Status.State
	}, timeout).Should(gomega.Equal(batchv1alpha1.Succeeded))

}
