package values

import (
	"context"
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/imdario/mergo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/app/v4/pkg/key"
)

const (
	configmap = "configmap"
	secret    = "secret"
)

// MergeConfigMapData merges the data from the catalog, app and user configmaps
// and returns a single set of values.
func (v *Values) MergeConfigMapData(ctx context.Context, app v1alpha1.App, catalog v1alpha1.Catalog) (map[string]interface{}, error) {
	appConfigMapName := key.AppConfigMapName(app)
	catalogConfigMapName := key.CatalogConfigMapName(catalog)
	userConfigMapName := key.UserConfigMapName(app)

	if appConfigMapName == "" && catalogConfigMapName == "" && userConfigMapName == "" {
		// Return early as there is no config.
		return nil, nil
	}

	// We get the catalog level values if configured.
	rawCatalogData, err := v.getConfigMapForCatalog(ctx, catalog)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	catalogData, err := extractData(configmap, "catalog", rawCatalogData)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// We get the app level values if configured.
	rawAppData, err := v.getConfigMapForApp(ctx, app)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	appData, err := extractData(configmap, "app", rawAppData)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = mergo.Merge(&catalogData, appData, mergo.WithOverride)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// We get the user level values if configured and merge them.
	if key.UserConfigMapName(app) != "" {
		rawUserData, err := v.getUserConfigMapForApp(ctx, app)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		userData, err := extractData(configmap, "user", rawUserData)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		err = mergo.Merge(&catalogData, userData, mergo.WithOverride)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return catalogData, nil
}

func (v *Values) getConfigMap(ctx context.Context, configMapName, configMapNamespace string) (map[string]string, error) {
	if configMapName == "" {
		// Return early as no configmap has been specified.
		return nil, nil
	}

	v.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("looking for configmap %#q in namespace %#q", configMapName, configMapNamespace))

	configMap, err := v.k8sClient.CoreV1().ConfigMaps(configMapNamespace).Get(ctx, configMapName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, microerror.Maskf(notFoundError, "configmap %#q in namespace %#q not found", configMapName, configMapNamespace)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	v.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found configmap %#q in namespace %#q", configMapName, configMapNamespace))

	return configMap.Data, nil
}

func (v *Values) getConfigMapForApp(ctx context.Context, app v1alpha1.App) (map[string]string, error) {
	configMap, err := v.getConfigMap(ctx, key.AppConfigMapName(app), key.AppConfigMapNamespace(app))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return configMap, nil
}

func (v *Values) getConfigMapForCatalog(ctx context.Context, catalog v1alpha1.Catalog) (map[string]string, error) {
	configMap, err := v.getConfigMap(ctx, key.CatalogConfigMapName(catalog), key.CatalogConfigMapNamespace(catalog))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return configMap, nil
}

func (v *Values) getUserConfigMapForApp(ctx context.Context, app v1alpha1.App) (map[string]string, error) {
	configMap, err := v.getConfigMap(ctx, key.UserConfigMapName(app), key.UserConfigMapNamespace(app))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return configMap, nil
}
