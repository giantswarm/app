package values

import (
	"context"
	"fmt"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/imdario/mergo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/app/v7/pkg/key"
)

// MergeSecretData merges the data from the catalog, app, user and extra config secrets
// and returns a single set of values.
func (v *Values) MergeSecretData(ctx context.Context, app v1alpha1.App, catalog v1alpha1.Catalog) (map[string]interface{}, error) {
	appSecretName := key.AppSecretName(app)
	catalogSecretName := key.CatalogSecretName(catalog)
	userSecretName := key.UserSecretName(app)

	extraConfigs := key.ExtraConfigs(app)

	if appSecretName == "" && catalogSecretName == "" && userSecretName == "" && len(extraConfigs) == 0 {
		// Return early as there is no secret.
		return nil, nil
	}

	// We get the catalog level secrets if configured.
	rawCatalogData, err := v.getSecretDataForCatalog(ctx, catalog)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	catalogData, err := extractData(secret, "catalog", toStringMap(rawCatalogData))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if catalogData == nil {
		// If there is no catalog data then treat it as an empty map otherwise `mergo.Merge` will silently
		// fail to merge the first layers: `dst = nil; mergo.Merge(dst, MAP_OF_DATA)` and `dst` is still nil
		catalogData = map[string]interface{}{}
	}

	// Pre cluster extra secrets
	err = v.fetchAndMergeExtraConfigs(ctx, getPreClusterExtraSecretEntries(extraConfigs), v.getSecretAsString, catalogData)
	if err != nil {
		return nil, err
	}

	// We get the app level secrets if configured.
	rawAppData, err := v.getSecretDataForApp(ctx, app)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	appData, err := extractData(secret, "app", toStringMap(rawAppData))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Secrets are merged and in case of intersecting values the app level
	// secrets are preferred.
	err = mergo.Merge(&catalogData, appData, mergo.WithOverride)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Post cluster / pre user extra secrets
	err = v.fetchAndMergeExtraConfigs(ctx, getPostClusterPreUserExtraSecretEntries(extraConfigs), v.getSecretAsString, catalogData)
	if err != nil {
		return nil, err
	}

	// We get the user level values if configured and merge them.
	if userSecretName != "" {
		rawUserData, err := v.getUserSecretDataForApp(ctx, app)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		// Secrets are merged again and in case of intersecting values the user
		// level secrets are preferred.
		userData, err := extractData(secret, "user", toStringMap(rawUserData))
		if err != nil {
			return nil, microerror.Mask(err)
		}

		err = mergo.Merge(&catalogData, userData, mergo.WithOverride)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// Post user extra secrets
	err = v.fetchAndMergeExtraConfigs(ctx, getPostUserExtraSecretEntries(extraConfigs), v.getSecretAsString, catalogData)
	if err != nil {
		return nil, err
	}

	return catalogData, nil
}

func (v *Values) getSecretAsString(ctx context.Context, secretName, secretNamespace string) (map[string]string, error) {
	data, err := v.getSecret(ctx, secretName, secretNamespace)

	if err != nil {
		return nil, err
	}

	return toStringMap(data), nil
}

func (v *Values) getSecret(ctx context.Context, secretName, secretNamespace string) (map[string][]byte, error) {
	if secretName == "" {
		// Return early as no secret has been specified.
		return nil, nil
	}

	v.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("looking for secret %#q in namespace %#q", secretName, secretNamespace))

	secret, err := v.k8sClient.CoreV1().Secrets(secretNamespace).Get(ctx, secretName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, microerror.Maskf(notFoundError, "secret %#q in namespace %#q not found", secretName, secretNamespace)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	v.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found secret %#q in namespace %#q", secretName, secretNamespace))

	return secret.Data, nil
}

func (v *Values) getSecretDataForApp(ctx context.Context, app v1alpha1.App) (map[string][]byte, error) {
	secret, err := v.getSecret(ctx, key.AppSecretName(app), key.AppSecretNamespace(app))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return secret, nil
}

func (v *Values) getSecretDataForCatalog(ctx context.Context, catalog v1alpha1.Catalog) (map[string][]byte, error) {
	secret, err := v.getSecret(ctx, key.CatalogSecretName(catalog), key.CatalogSecretNamespace(catalog))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return secret, nil
}

func (v *Values) getUserSecretDataForApp(ctx context.Context, app v1alpha1.App) (map[string][]byte, error) {
	secret, err := v.getSecret(ctx, key.UserSecretName(app), key.UserSecretNamespace(app))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return secret, nil
}
