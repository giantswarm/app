package validation

import (
	"context"
	"fmt"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/app/v6/pkg/key"
)

const (
	catalogNotFoundTemplate           = "catalog %#q not found"
	nameTooLongTemplate               = "name %#q is %d chars and exceeds max length of %d chars"
	nameNotFoundReasonTemplate        = "name is not specified for %s"
	targetNamespaceNotAllowedTemplate = "target namespace %s is not allowed for in-cluster apps"
	namespaceNotFoundReasonTemplate   = "namespace is not specified for %s %#q"
	labelInvalidValueTemplate         = "label %#q has invalid value %#q"
	labelNotFoundTemplate             = "label %#q not found"
	resourceNotFoundTemplate          = "%s %#q in namespace %#q not found"

	defaultCatalogName            = "default"
	nginxIngressControllerAppName = "nginx-ingress-controller-app"

	// nameMaxLength is 53 characters as this is the maximum allowed for Helm
	// release names.
	nameMaxLength = 53
)

func (v *Validator) ValidateApp(ctx context.Context, app v1alpha1.App) (bool, error) {
	var err error

	err = v.validateCatalog(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateKubeConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateLabels(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateMetadataConstraints(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateName(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateNamespaceConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateTargetNamespace(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateUserConfig(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func (v *Validator) ValidateAppUpdate(ctx context.Context, app, currentApp v1alpha1.App) (bool, error) {
	err := v.validateNamespaceUpdate(ctx, app, currentApp)
	if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func (v *Validator) validateCatalog(ctx context.Context, cr v1alpha1.App) error {
	var err error

	if key.CatalogName(cr) == "" {
		return nil
	}

	var namespaces []string
	{
		if key.CatalogNamespace(cr) != "" {
			namespaces = []string{key.CatalogNamespace(cr)}
		} else {
			namespaces = []string{metav1.NamespaceDefault, "giantswarm"}
		}
	}

	var matchedCatalog *v1alpha1.Catalog

	for _, ns := range namespaces {
		var catalog v1alpha1.Catalog
		err = v.g8sClient.Get(ctx, client.ObjectKey{
			Namespace: ns,
			Name:      key.CatalogName(cr),
		}, &catalog)
		if apierrors.IsNotFound(err) {
			// no-op
			continue
		} else if err != nil {
			return microerror.Mask(err)
		}
		matchedCatalog = &catalog
		break
	}

	if matchedCatalog == nil || matchedCatalog.Name == "" {
		return microerror.Maskf(validationError, catalogNotFoundTemplate, key.CatalogName(cr))
	}

	return nil
}

func (v *Validator) validateConfig(ctx context.Context, cr v1alpha1.App) error {
	if key.AppConfigMapName(cr) != "" {
		err := v.validateConfigMapExists(ctx, key.AppConfigMapName(cr), key.AppConfigMapNamespace(cr), "configmap", cr)
		if apierrors.IsNotFound(err) {
			// appConfigMapNotFoundError is used rather than a validation error because
			// during cluster creation there is a short delay while it is generated.
			return microerror.Maskf(appConfigMapNotFoundError, resourceNotFoundTemplate, "configmap", key.AppConfigMapName(cr), key.AppConfigMapNamespace(cr))
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	if key.AppSecretName(cr) != "" {
		err := v.validateSecretExists(ctx, key.AppSecretName(cr), key.AppSecretNamespace(cr), "secret", cr)
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(validationError, resourceNotFoundTemplate, "secret", key.AppSecretName(cr), key.AppSecretNamespace(cr))
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (v *Validator) validateName(ctx context.Context, cr v1alpha1.App) error {
	if len(cr.Name) > nameMaxLength {
		return microerror.Maskf(validationError, nameTooLongTemplate, cr.Name, len(cr.Name), nameMaxLength)
	}

	return nil
}

// We make sure users cannot create in-cluster Apps outside their organization
// or WC namespaces. Otherwise `.spec.namespace` could be exploited to override permissions.
func (v *Validator) validateTargetNamespace(ctx context.Context, cr v1alpha1.App) error {
	isInCluster := key.InCluster(cr)
	isNotGs := cr.Namespace != "giantswarm"
	isOutsideOrg := cr.Namespace != cr.Spec.Namespace

	if isInCluster && isNotGs && isOutsideOrg {
		return microerror.Maskf(validationError, targetNamespaceNotAllowedTemplate, cr.Spec.Namespace)
	}

	return nil
}

func (v *Validator) validateNamespaceConfig(ctx context.Context, cr v1alpha1.App) error {
	annotations := key.AppNamespaceAnnotations(cr)
	labels := key.AppNamespaceLabels(cr)

	if annotations == nil && labels == nil {
		// no-op
		return nil
	}

	var apps []v1alpha1.App
	{
		fieldSelector, err := fields.ParseSelector(fmt.Sprintf("metadata.name!=%s", cr.Name))
		if err != nil {
			return microerror.Mask(err)
		}

		lo := client.ListOptions{
			Namespace:     cr.Namespace,
			FieldSelector: fieldSelector,
		}
		var appList v1alpha1.AppList
		err = v.g8sClient.List(ctx, &appList, &lo)
		if err != nil {
			return microerror.Mask(err)
		}

		apps = appList.Items
	}

	for _, app := range apps {
		if key.AppNamespace(cr) != key.AppNamespace(app) {
			continue
		}

		targetAnnotations := key.AppNamespaceAnnotations(app)
		if targetAnnotations != nil && annotations != nil {
			for k, v := range targetAnnotations {
				originalValue, ok := annotations[k]
				if ok && originalValue != v {
					return microerror.Maskf(validationError, "app %#q annotation %#q for target namespace %#q collides with value %#q for app %#q",
						key.AppName(cr), k, key.AppNamespace(cr), v, app.Name)
				}
			}
		}

		targetLabels := key.AppNamespaceLabels(app)
		if targetLabels != nil && labels != nil {
			for k, v := range targetLabels {
				originalValue, ok := labels[k]
				if ok && originalValue != v {
					return microerror.Maskf(validationError, "app %#q label %#q for target namespace %#q collides with value %#q for app %#q",
						key.AppName(cr), k, key.AppNamespace(cr), v, app.Name)
				}
			}
		}
	}

	return nil
}

func (v *Validator) validateKubeConfig(ctx context.Context, cr v1alpha1.App) error {
	if !key.InCluster(cr) {
		err := v.validateSecretExists(ctx, key.KubeConfigSecretName(cr), key.KubeConfigSecretNamespace(cr), "kubeconfig secret", cr)
		if apierrors.IsNotFound(err) {
			// kubeConfigNotFoundError is used rather than a validation error because
			// during cluster creation there is a short delay while it is generated.
			return microerror.Maskf(kubeConfigNotFoundError, resourceNotFoundTemplate, "kubeconfig secret", key.KubeConfigSecretName(cr), key.KubeConfigSecretNamespace(cr))
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (v *Validator) validateLabels(ctx context.Context, cr v1alpha1.App) error {
	// This is for migration towards managing App CRs in the org namespace.
	// For the transition time being we must retain backward compatibility for the
	// App CRs in cluster namespaces.
	isManagedInOrg := !key.InCluster(cr) && key.IsInOrgNamespace(cr)

	if isManagedInOrg {
		return v.validateOrgLabels(ctx, cr)
	}

	return v.validateClusterLabels(ctx, cr)
}

func (v *Validator) validateClusterLabels(ctx context.Context, cr v1alpha1.App) error {
	if key.VersionLabel(cr) == "" {
		return microerror.Maskf(validationError, labelNotFoundTemplate, label.AppOperatorVersion)
	}
	if key.VersionLabel(cr) == key.LegacyAppVersionLabel {
		return microerror.Maskf(validationError, labelInvalidValueTemplate, label.AppOperatorVersion, key.VersionLabel(cr))
	}

	return nil
}

func (v *Validator) validateOrgLabels(ctx context.Context, cr v1alpha1.App) error {
	if key.ClusterLabel(cr) == "" {
		return microerror.Maskf(validationError, labelNotFoundTemplate, label.Cluster)
	}

	return nil
}

func (v *Validator) validateMetadataConstraints(ctx context.Context, cr v1alpha1.App) error {
	name := key.AppCatalogEntryName(key.CatalogName(cr), key.AppName(cr), key.Version(cr))

	var entry v1alpha1.AppCatalogEntry
	err := v.g8sClient.Get(ctx, client.ObjectKey{
		Namespace: metav1.NamespaceDefault,
		Name:      name,
	}, &entry)
	if apierrors.IsNotFound(err) {
		v.logger.Debugf(ctx, "appcatalogentry %#q not found, skipping metadata validation", name)
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	if entry.Spec.Restrictions == nil {
		// no-op
		return nil
	}

	if len(entry.Spec.Restrictions.CompatibleProviders) > 0 {
		if !contains(entry.Spec.Restrictions.CompatibleProviders, v1alpha1.Provider(v.provider)) {
			return microerror.Maskf(validationError, "app %#q can only be installed for providers %#q not %#q",
				cr.Spec.Name, entry.Spec.Restrictions.CompatibleProviders, v.provider)
		}
	}

	if entry.Spec.Restrictions.FixedNamespace != "" {
		if entry.Spec.Restrictions.FixedNamespace != cr.Spec.Namespace {
			return microerror.Maskf(validationError, "app %#q can only be installed in namespace %#q only, not %#q",
				cr.Spec.Name, entry.Spec.Restrictions.FixedNamespace, cr.Spec.Namespace)
		}
	}

	var apps []v1alpha1.App
	if entry.Spec.Restrictions.ClusterSingleton || entry.Spec.Restrictions.NamespaceSingleton {
		fieldSelector, err := fields.ParseSelector(fmt.Sprintf("metadata.name!=%s", cr.Name))
		if err != nil {
			return microerror.Mask(err)
		}

		lo := client.ListOptions{
			FieldSelector: fieldSelector,
			Namespace:     cr.Namespace,
		}
		var appList v1alpha1.AppList
		err = v.g8sClient.List(ctx, &appList, &lo)
		if err != nil {
			return microerror.Mask(err)
		}

		apps = appList.Items
	}

	for _, app := range apps {
		if app.Spec.Name == cr.Spec.Name {
			if entry.Spec.Restrictions.ClusterSingleton {
				return microerror.Maskf(validationError, "app %#q can only be installed once in cluster %#q",
					cr.Spec.Name, cr.Namespace)
			}
			if entry.Spec.Restrictions.NamespaceSingleton {
				if app.Spec.Namespace == cr.Spec.Namespace {
					return microerror.Maskf(validationError, "app %#q can only be installed only once in namespace %#q",
						cr.Spec.Name, key.Namespace(cr))
				}
			}
		}
	}

	return nil
}

func (v *Validator) validateNamespaceUpdate(ctx context.Context, app, currentApp v1alpha1.App) error {
	if key.Namespace(app) != key.Namespace(currentApp) {
		return microerror.Maskf(validationError, "target namespace for app %#q cannot be changed from %#q to %#q", app.Name,
			key.Namespace(currentApp), key.Namespace(app))
	}

	return nil
}

func (v *Validator) validateUserConfig(ctx context.Context, cr v1alpha1.App) error {
	if key.UserConfigMapName(cr) != "" {
		// NGINX Ingress Controller is no longer a pre-installed app
		// managed by cluster-operator. So we don't need to restrict
		// the name.
		if key.CatalogName(cr) == defaultCatalogName && key.AppName(cr) != nginxIngressControllerAppName {
			configMapName := fmt.Sprintf("%s-user-values", cr.Name)
			if key.UserConfigMapName(cr) != configMapName {
				return microerror.Maskf(validationError, "user configmap must be named %#q for app in default catalog", configMapName)
			}
		}

		err := v.validateConfigMapExists(ctx, key.UserConfigMapName(cr), key.UserConfigMapNamespace(cr), "configmap", cr)
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(validationError, resourceNotFoundTemplate, "configmap", key.UserConfigMapName(cr), key.UserConfigMapNamespace(cr))
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	if key.UserSecretName(cr) != "" {
		if key.CatalogName(cr) == defaultCatalogName {
			secretName := fmt.Sprintf("%s-user-secrets", cr.Name)
			if key.UserSecretName(cr) != secretName {
				return microerror.Maskf(validationError, "user secret must be named %#q for app in default catalog", secretName)
			}
		}

		err := v.validateSecretExists(ctx, key.UserSecretName(cr), key.UserSecretNamespace(cr), "secret", cr)
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(validationError, resourceNotFoundTemplate, "secret", key.UserSecretName(cr), key.UserSecretNamespace(cr))
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (v *Validator) validateConfigMapExists(ctx context.Context, name, namespace, kind string, cr v1alpha1.App) error {
	if namespace == "" {
		return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, kind, name)
	}

	if name == "" {
		return microerror.Maskf(validationError, nameNotFoundReasonTemplate, kind)
	}

	if v.enableManagedByLabel && key.IsManagedByFlux(cr, v.projectName) {
		v.logger.Debugf(ctx, "skipping validation of app '%s/%s' dependencies due to '%s=%s' label", cr.Namespace, cr.Name, label.ManagedBy, key.ManagedByLabel(cr))
		return nil
	}

	_, err := v.k8sClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (v *Validator) validateSecretExists(ctx context.Context, name, namespace, kind string, cr v1alpha1.App) error {
	if namespace == "" {
		return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, kind, name)
	}

	if name == "" {
		return microerror.Maskf(validationError, nameNotFoundReasonTemplate, kind)
	}

	// Check basic things, like empty name or namespace, but skip
	// existance validation when managed by Fux.
	if v.enableManagedByLabel && key.IsManagedByFlux(cr, v.projectName) {
		v.logger.Debugf(ctx, "skipping validation of app '%s/%s' dependencies due to '%s=%s' label", cr.Namespace, cr.Name, label.ManagedBy, key.ManagedByLabel(cr))
		return nil
	}

	_, err := v.k8sClient.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func contains(s []v1alpha1.Provider, e v1alpha1.Provider) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
