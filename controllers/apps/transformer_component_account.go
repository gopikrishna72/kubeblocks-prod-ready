/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package apps

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sethvargo/go-password/password"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/component"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
)

// componentAccountTransformer handles component system accounts.
type componentAccountTransformer struct{}

var _ graph.Transformer = &componentAccountTransformer{}

func (t *componentAccountTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*componentTransformContext)
	if model.IsObjectDeleting(transCtx.ComponentOrig) {
		return nil
	}

	synthesizeComp := transCtx.SynthesizeComponent
	graphCli, _ := transCtx.Client.(model.GraphClient)

	for _, account := range synthesizeComp.SystemAccounts {
		secret, err := t.buildSystemAccount(ctx, synthesizeComp, account)
		if err != nil {
			return err
		}
		if err = t.createOrUpdate(ctx, dag, graphCli, secret); err != nil {
			return err
		}
	}
	return nil
}

func (t *componentAccountTransformer) buildSystemAccount(ctx graph.TransformContext,
	synthesizeComp *component.SynthesizedComponent, account appsv1alpha1.ComponentSystemAccount) (*corev1.Secret, error) {
	var password []byte
	if account.SecretRef != nil {
		var err error
		if password, err = t.getPasswordFromSecret(ctx, account); err != nil {
			return nil, err
		}
	} else {
		password = t.generatePassword(account)
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: synthesizeComp.Namespace,
			Name:      constant.GenerateAccountSecretName(synthesizeComp.ClusterName, synthesizeComp.Name, account.Name),
			Labels:    constant.GetComponentWellKnownLabels(synthesizeComp.ClusterName, synthesizeComp.Name),
		},
		Data: map[string][]byte{
			constant.AccountNameForSecret:   []byte(account.Name),
			constant.AccountPasswdForSecret: password,
		},
	}
	return secret, nil
}

func (t *componentAccountTransformer) getPasswordFromSecret(ctx graph.TransformContext, account appsv1alpha1.ComponentSystemAccount) ([]byte, error) {
	secretKey := types.NamespacedName{
		Namespace: account.SecretRef.Namespace,
		Name:      account.SecretRef.Name,
	}
	secret := &corev1.Secret{}
	if err := ctx.GetClient().Get(ctx.GetContext(), secretKey, secret); err != nil {
		return nil, err
	}
	if len(secret.Data) == 0 || len(secret.Data[constant.AccountPasswdForSecret]) == 0 {
		return nil, fmt.Errorf("referenced account secret has no required credential field")
	}
	return secret.Data[constant.AccountPasswdForSecret], nil
}

func (t *componentAccountTransformer) generatePassword(account appsv1alpha1.ComponentSystemAccount) []byte {
	config := account.PasswordGenerationPolicy
	passwd, _ := password.Generate((int)(config.Length), (int)(config.NumDigits), (int)(config.NumSymbols), false, false)
	switch config.LetterCase {
	case appsv1alpha1.UpperCases:
		passwd = strings.ToUpper(passwd)
	case appsv1alpha1.LowerCases:
		passwd = strings.ToLower(passwd)
	}
	return []byte(passwd)
}

func (t *componentAccountTransformer) createOrUpdate(ctx graph.TransformContext,
	dag *graph.DAG, graphCli model.GraphClient, secret *corev1.Secret) error {
	key := types.NamespacedName{
		Namespace: secret.Namespace,
		Name:      secret.Name,
	}
	obj := &corev1.Secret{}
	if err := ctx.GetClient().Get(ctx.GetContext(), key, obj); err != nil {
		if apierrors.IsNotFound(err) {
			graphCli.Create(dag, secret)
			return nil
		}
		return err
	}
	objCopy := obj.DeepCopy()
	objCopy.Immutable = secret.Immutable
	objCopy.Data = secret.Data
	objCopy.StringData = secret.StringData
	objCopy.Type = secret.Type
	if !reflect.DeepEqual(obj, objCopy) {
		graphCli.Update(dag, obj, objCopy)
	}
	return nil
}