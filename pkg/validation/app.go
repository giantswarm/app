package validation

import (
	"context"
	"fmt"
	"strings"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/app/v7/pkg/key"
)

const (
	catalogNotFoundTemplate           = "catalog %#q not found"
	nameTooLongTemplate               = "name %#q is %d chars and exceeds max length of %d chars"
	nameNotFoundReasonTemplate        = "name is not specified for %s"
	targetNamespaceNotAllowedTemplate = "target namespace %s is not allowed for in-cluster apps"
	namespaceMismatchTemplate         = "wrong %#q namespace for the `chart-operator.giantswarm.io/app-namespace` annotation"
	namespaceNotFoundReasonTemplate   = "namespace is not specified for %s %#q"
	labelInvalidValueTemplate         = "label %#q has invalid value %#q"
	labelNotFoundTemplate             = "label %#q not found"
	labelInClusterAppTemplate         = "label %#q must be set to `0.0.0` for in-cluster app"
	resourceNotFoundTemplate          = "%s %#q in namespace %#q not found"

	defaultCatalogName = "default"

	// nameMaxLength is 53 characters as this is the maximum allowed for Helm
	// release names.
	nameMaxLength = 53
)

func (v *Validator) ValidateApp(ctx context.Context, app v1alpha1.App) (bool, error) {
	var err error

	err = v.validateAnnotations(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateCatalog(ctx, app)
	if err != nil {
		return false, microerror.Mask(err)
	}

	err = v.validateLabels(ctx, app)
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

	err = v.validateUniqueInClusterAppName(ctx, app)
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

// This is for preventing chart-operator to select elevated
// client for a given app, by making an impression it comes from
// a different namespace, see explanation:
// https://github.com/giantswarm/giantswarm/issues/22100#issuecomment-1131723221
func (v *Validator) validateAnnotations(ctx context.Context, cr v1alpha1.App) error {
	namespaceAnnotation := key.AppNamespaceAnnotation(cr)
	if namespaceAnnotation != "" && namespaceAnnotation != cr.ObjectMeta.Namespace {
		return microerror.Maskf(validationError, namespaceMismatchTemplate, namespaceAnnotation)
	}

	return nil
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
		if err := v.validateNameAndNamespaceAreSet(key.AppConfigMapName(cr), key.AppConfigMapNamespace(cr), "configmap"); err != nil {
			return microerror.Mask(err)
		}

		if v.isAdmissionController {
			v.logger.Debugf(ctx, "skipping '.spec.config.configMap' validation of app '%s/%s' in admission controllers", cr.Namespace, cr.Name)
		} else {
			err := v.validateConfigMapExists(ctx, key.AppConfigMapName(cr), key.AppConfigMapNamespace(cr), "configmap")
			if apierrors.IsNotFound(err) {
				// appConfigMapNotFoundError is used rather than a validation error because
				// during cluster creation there is a short delay while it is generated.
				return microerror.Maskf(appConfigMapNotFoundError, resourceNotFoundTemplate, "configmap", key.AppConfigMapName(cr), key.AppConfigMapNamespace(cr))
			} else if err != nil {
				return microerror.Mask(err)
			}
		}

	}

	if key.AppSecretName(cr) != "" {
		if err := v.validateNameAndNamespaceAreSet(key.AppSecretName(cr), key.AppSecretNamespace(cr), "secret"); err != nil {
			return microerror.Mask(err)
		}

		if v.isAdmissionController {
			v.logger.Debugf(ctx, "skipping '.spec.config.secret' validation of app '%s/%s' in admission controllers", cr.Namespace, cr.Name)
		} else {
			err := v.validateSecretExists(ctx, key.AppSecretName(cr), key.AppSecretNamespace(cr), "secret")
			if apierrors.IsNotFound(err) {
				return microerror.Maskf(validationError, resourceNotFoundTemplate, "secret", key.AppSecretName(cr), key.AppSecretNamespace(cr))
			} else if err != nil {
				return microerror.Mask(err)
			}
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
		if err := v.validateNameAndNamespaceAreSet(key.KubeConfigSecretName(cr), key.KubeConfigSecretNamespace(cr), "kubeconfig secret"); err != nil {
			return microerror.Mask(err)
		}

		if v.isAdmissionController {
			v.logger.Debugf(ctx, "skipping '.spec.kubeConfig.secret' validation of remote cluster app '%s/%s' in admission controllers", cr.Namespace, cr.Name)
		} else {
			err := v.validateSecretExists(ctx, key.KubeConfigSecretName(cr), key.KubeConfigSecretNamespace(cr), "kubeconfig secret")
			if apierrors.IsNotFound(err) {
				// kubeConfigNotFoundError is used rather than a validation error because
				// during cluster creation there is a short delay while it is generated.
				return microerror.Maskf(kubeConfigNotFoundError, resourceNotFoundTemplate, "kubeconfig secret", key.KubeConfigSecretName(cr), key.KubeConfigSecretNamespace(cr))
			} else if err != nil {
				return microerror.Mask(err)
			}
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
	if key.InCluster(cr) && key.VersionLabel(cr) != key.UniqueAppVersionLabel {
		return microerror.Maskf(validationError, labelInClusterAppTemplate, label.AppOperatorVersion)
	}

	return nil
}

func (v *Validator) validateOrgLabels(ctx context.Context, cr v1alpha1.App) error {
	if key.ClusterLabel(cr) == "" {
		return microerror.Maskf(validationError, labelNotFoundTemplate, label.Cluster)
	}
	if key.InCluster(cr) && key.VersionLabel(cr) != key.UniqueAppVersionLabel {
		return microerror.Maskf(validationError, labelInClusterAppTemplate, label.AppOperatorVersion)
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
		if !contains(entry.Spec.Restrictions.CompatibleProviders, v.provider) {
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

		var labelSelector labels.Selector
		if key.IsInOrgNamespace(cr) {
			labelSelector, err = labels.Parse(fmt.Sprintf("%s=%s", label.Cluster, key.ClusterID(cr)))
			if err != nil {
				return microerror.Mask(err)
			}
		}

		fieldSelector, err := fields.ParseSelector(fmt.Sprintf("metadata.name!=%s", cr.Name))
		if err != nil {
			return microerror.Mask(err)
		}

		lo := client.ListOptions{
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
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
		if app.Spec.Name != cr.Spec.Name {
			continue
		}

		if entry.Spec.Restrictions.ClusterSingleton {
			clusterId := key.ClusterID(cr)

			if clusterId == "" {
				clusterId = cr.Namespace
			}
			return microerror.Maskf(validationError, "app %#q can only be installed once in cluster %#q",
				cr.Spec.Name, clusterId)
		}

		if !entry.Spec.Restrictions.NamespaceSingleton {
			continue
		}

		if app.Spec.Namespace == cr.Spec.Namespace {
			return microerror.Maskf(validationError, "app %#q can only be installed only once in namespace %#q",
				cr.Spec.Name, key.Namespace(cr))
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
		if key.CatalogName(cr) == defaultCatalogName {
			// This check is for `cluster-operator` only. For CAPI clusters, that does not rely on
			// `cluster-operator`, the names could be any, but since it hasn't been conditioned earlier,
			// making the whole function conditional now, when CAPI is well-established, may have some
			// hard-to-predict consequences. Hence it is safer to let the function operate as usual,
			// but in addition put some extra conditions like below.

			// User ConfigMap name must, in general, match the '<App_CR_name>-user-values',
			// with one exception below.
			configMapName := fmt.Sprintf("%s-user-values", cr.Name)

			nameMismatch := key.UserConfigMapName(cr) != configMapName

			// The Cluster Operator may be told to create in-cluster app (aka bundle) for the WC,
			// what requires special naming, like adding prefix or suffix, because otherwise a conflict
			// may arise in within the MC, see:
			// https://github.com/giantswarm/cluster-operator/blob/1681bd32d19e4fda44993d7629eed227ff3cdc59/service/controller/resource/app/desired.go#L216,
			// Some default apps may hence carry names prefixed with the cluster ID, which are different from
			// names of these apps as configured in the Release CR (without prefix). This is a problem, because
			// Cluster Operator uses Release CR-originated names, to define user ConfigMap names. This in turn,
			// fails the above condition, so another check is needed against name without the prefix.
			if !key.IsInOrgNamespace(cr) && key.ClusterLabel(cr) != "" {
				configMapName = strings.TrimPrefix(
					configMapName,
					fmt.Sprintf("%s-", key.ClusterLabel(cr)),
				)

				nameMismatch = nameMismatch && key.UserConfigMapName(cr) != configMapName
			}

			if nameMismatch {
				return microerror.Maskf(validationError, "user configmap must be named %#q for app in default catalog", configMapName)
			}
		}

		if err := v.validateNameAndNamespaceAreSet(key.UserConfigMapName(cr), key.UserConfigMapNamespace(cr), "configmap"); err != nil {
			return microerror.Mask(err)
		}

		if v.isAdmissionController {
			v.logger.Debugf(ctx, "skipping '.spec.userConfig.configMap' validation of app '%s/%s' in admission controllers", cr.Namespace, cr.Name)
		} else {
			err := v.validateConfigMapExists(ctx, key.UserConfigMapName(cr), key.UserConfigMapNamespace(cr), "configmap")
			if apierrors.IsNotFound(err) {
				return microerror.Maskf(validationError, resourceNotFoundTemplate, "configmap", key.UserConfigMapName(cr), key.UserConfigMapNamespace(cr))
			} else if err != nil {
				return microerror.Mask(err)
			}
		}

	}

	if key.UserSecretName(cr) != "" {
		if key.CatalogName(cr) == defaultCatalogName {
			// This check is for `cluster-operator` only. For CAPI clusters, that does not rely on
			// `cluster-operator`, the names could be any, but since it hasn't been conditioned earlier,
			// making the whole function conditional now, when CAPI is well-established, may have some
			// hard-to-predict consequences. Hence it is safer to let the function operate as usual,
			// but in addition put some extra conditions like below.

			// User Secret name must, in general, match the '<App_CR_name>-user-secrets',
			// with one exception below.
			secretName := fmt.Sprintf("%s-user-secrets", cr.Name)

			nameMismatch := key.UserSecretName(cr) != secretName

			// The Cluster Operator may be told to create in-cluster app (aka bundle) for the WC,
			// what requires special naming, like adding prefix or suffix, because otherwise a conflict
			// may arise in within the MC, see:
			// https://github.com/giantswarm/cluster-operator/blob/1681bd32d19e4fda44993d7629eed227ff3cdc59/service/controller/resource/app/desired.go#L216,
			// Some default apps may hence carry names prefixed with the cluster ID, which are different from
			// names of these apps as configured in the Release CR (without prefix). This is a problem, because
			// Cluster Operator uses Release CR-originated names, to define user Secret names. This in turn,
			// fails the above condition, so another check is needed against name without the prefix.
			if !key.IsInOrgNamespace(cr) && key.ClusterLabel(cr) != "" {
				secretName = strings.TrimPrefix(
					secretName,
					fmt.Sprintf("%s-", key.ClusterLabel(cr)),
				)

				nameMismatch = nameMismatch && key.UserSecretName(cr) != secretName
			}

			if nameMismatch {
				return microerror.Maskf(validationError, "user secret must be named %#q for app in default catalog", secretName)
			}
		}

		if err := v.validateNameAndNamespaceAreSet(key.UserSecretName(cr), key.UserSecretNamespace(cr), "secret"); err != nil {
			return microerror.Mask(err)
		}

		if v.isAdmissionController {
			v.logger.Debugf(ctx, "skipping '.spec.userConfig.secret' validation of app '%s/%s' in admission controllers", cr.Namespace, cr.Name)
		} else {
			err := v.validateSecretExists(ctx, key.UserSecretName(cr), key.UserSecretNamespace(cr), "secret")
			if apierrors.IsNotFound(err) {
				return microerror.Maskf(validationError, resourceNotFoundTemplate, "secret", key.UserSecretName(cr), key.UserSecretNamespace(cr))
			} else if err != nil {
				return microerror.Mask(err)
			}
		}

	}

	return nil
}

func (v *Validator) validateNameAndNamespaceAreSet(name, namespace, kind string) error {
	if namespace == "" {
		return microerror.Maskf(validationError, namespaceNotFoundReasonTemplate, kind, name)
	}

	if name == "" {
		return microerror.Maskf(validationError, nameNotFoundReasonTemplate, kind)
	}

	return nil
}

func (v *Validator) validateConfigMapExists(ctx context.Context, name, namespace, kind string) error {
	_, err := v.k8sClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (v *Validator) validateSecretExists(ctx context.Context, name, namespace, kind string) error {
	_, err := v.k8sClient.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (v *Validator) validateUniqueInClusterAppName(ctx context.Context, cr v1alpha1.App) error {
	// WARNING: This part assumes knowledge of the internal workings of app-operator and the App Platform
	specialNamespace := "giantswarm"

	if !key.InCluster(cr) && cr.Namespace != specialNamespace {
		return nil
	}

	apps := &v1alpha1.AppList{}

	fieldName := "metadata.name"
	fieldValue := cr.Name

	err := v.g8sClient.List(ctx, apps, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(fieldName, fieldValue),
	})
	if err != nil {
		return microerror.Maskf(validationError, "failed to list apps with %#q set to %#q to validate unique in-cluster app name rule, %#v", fieldName, fieldValue, err)
	}

	for _, inspectedApp := range apps.Items {
		// If it is the same app (we are handling an update event for example) then skip over
		if inspectedApp.Namespace == cr.Namespace {
			continue
		}

		// Found another app that is in-cluster and bears the same name, you shall not pass!
		//
		// The extra check for comparing the names is redundant, it is added because of a long-standing bug in
		// the fake kubernetes client used in tests.
		//
		// See: https://github.com/kubernetes-sigs/controller-runtime/issues/1376
		// See: https://github.com/kubernetes-sigs/controller-runtime/issues/866
		if inspectedApp.Name == cr.Name {
			if inspectedApp.Namespace == specialNamespace {
				return microerror.Maskf(validationError, "found another app named %#q installed into the %#q namespace", inspectedApp.Name, specialNamespace)
			}

			if key.InCluster(inspectedApp) {
				if cr.Namespace == specialNamespace {
					return microerror.Maskf(validationError, "there is in-cluster app named %#q already installed in the %#q namespace that would cause name collision with the currently submitted app named %#q in the %#q namespace", inspectedApp.Name, inspectedApp.Namespace, cr.Name, cr.Namespace)
				}

				return microerror.Maskf(validationError, "in-cluster apps must be given a unique name, found an app named %#q as well in the %#q namespace", inspectedApp.Name, inspectedApp.Namespace)
			}
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
