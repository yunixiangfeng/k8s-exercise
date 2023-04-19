/*
Copyright 2023 wu123.

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

package controllers

import (
	"context"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	traefikv1 "k8s-operator/api/v1"
)

// TraefikServiceReconciler reconciles a TraefikService object
type TraefikServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=traefik.wu123.com,resources=traefikservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=traefik.wu123.com,resources=traefikservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=traefik.wu123.com,resources=traefikservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TraefikService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
// 每当traefikv1.TraefikService对象有变化时，我们就会收到一个请求
func (r *TraefikServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// 获取日志对象
	l := log.FromContext(ctx, "namespace", req.Namespace)
	// TODO(user): your logic here

	// 1.通过名称获取TraefikService对象，并打印
	var obj traefikv1.TraefikService
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		// 如果Not Found则表示该资源已经删除，需要做删除处理
		if apierrors.IsNotFound(err) {
			l.Info("delete service ...")
			err = nil
		} else {
			l.Error(err, "unable to fetch TraefikService")
		}
		return ctrl.Result{}, err
	}
	l.Info("get TraefikService object",
		"name", obj.Name,
		"ur1", obj.Spec.URL,
		"entrypoint", obj.Spec.Entrypoint)

	// 2.更新状态 排除已经同步完成的对象
	if obj.Status.Active {
		l.Info("traefik service is active, skip ...")
		return ctrl.Result{}, nil
	}

	// 3.注册服务到Traefik service中，比如写入到etcd provider中
	l.Info("set traefik service ...")

	// 4，修改成功调整对象的状态
	obj.Status.Active = true
	obj.Status.LastUpdateTime = &metav1.Time{Time: time.Now()}
	if err := r.Status().Update(ctx, &obj); err != nil {
		l.Info("unable to update status", "reason", err.Error())
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TraefikServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&traefikv1.TraefikService{}).
		Complete(r)
}
