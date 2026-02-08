package cmder_test

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/brandon1024/cmder"
	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates how to register command flags. The flags and documentation were borrowed from 'kubectl
// apply' to showcase what a real-world example woud look like.
func ExampleFlagInitializer() {
	var cmd KubectlApply

	args := []string{"-h"}

	if err := cmder.Execute(context.Background(), &cmd, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
}

const KubectlApplyDesc = `
Apply a configuration to a resource by file name or stdin. The resource name must be specified. This resource will be
created if it doesn't exist yet. To use 'apply', always create the resource initially with either 'apply' or
'create --save-config'.

JSON and YAML formats are accepted.
`

const KubectlApplyExamples = `
# Apply the configuration in pod.json to a pod
kubectl apply -f ./pod.json

# Apply resources from a directory containing kustomization.yaml - e.g. dir/kustomization.yaml
kubectl apply -k dir/

# Apply the JSON passed into stdin to a pod
cat pod.json | kubectl apply -f -

# Apply the configuration from all files that end with '.json'
kubectl apply -f '*.json'

# Note: --prune is still in Alpha
# Apply the configuration in manifest.yaml that matches label app=nginx and delete all other resources that are not in the file and match label app=nginx
kubectl apply --prune -f manifest.yaml -l app=nginx

# Apply the configuration in manifest.yaml and delete all the other config maps that are not in the file
kubectl apply --prune -f manifest.yaml --all --prune-allowlist=core/v1/ConfigMap
`

type KubectlApply struct {
	all                      bool
	allowMissingTemplateKeys bool
	cascade                  string
	dryRun                   string
	fieldManager             string
	filename                 getopt.StringsVar
	force                    bool
	forceConflicts           bool
	gracePeriod              int
	kustomize                string
	openapiPatch             bool
	output                   string
	overwrite                bool
	prune                    bool
	pruneAllowlist           getopt.StringsVar
	recursive                bool
	selector                 string
	serverSide               bool
	showManagedFields        bool
	subresource              string
	template                 string
	timeout                  time.Duration
	validate                 string
	wait                     bool
}

func (a *KubectlApply) InitializeFlags(fs *flag.FlagSet) {
	fs.BoolVar(&a.all, "all", false,
		"Select all resources in the namespace of the specified resource types.")
	fs.BoolVar(&a.allowMissingTemplateKeys, "allow-missing-template-keys", true,
		"If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats.")
	fs.StringVar(&a.cascade, "cascade", "background",
		"Must be \"background\", \"orphan\", or \"foreground\". Selects the deletion cascading strategy for the dependents (e.g. Pods created by a ReplicationController). Defaults to background.")
	fs.StringVar(&a.dryRun, "dry-run", "none",
		"Must be \"none\", \"server\", or \"client\". If client strategy, only print the object that would be sent, without sending it. If server strategy, submit server-side request without persisting the resource.")
	fs.StringVar(&a.fieldManager, "field-manager", "kubectl-client-side-apply",
		"Name of the manager used to track field ownership.")
	fs.Var(&a.filename, "filename",
		"The files that contain the configurations to apply.")
	fs.Var(&a.filename, "f",
		"The files that contain the configurations to apply.")
	fs.BoolVar(&a.force, "force", false,
		"If true, immediately remove resources from API and bypass graceful deletion. Note that immediate deletion of some resources may result in inconsistency or data loss and requires confirmation.")
	fs.BoolVar(&a.forceConflicts, "force-conflicts", false,
		"If true, server-side apply will force the changes against conflicts.")
	fs.IntVar(&a.gracePeriod, "grace-period", -1,
		"Period of time in seconds given to the resource to terminate gracefully. Ignored if negative. Set to 1 for immediate shutdown. Can only be set to 0 when --force is true (force deletion).")
	fs.StringVar(&a.kustomize, "kustomize", "",
		"Process a kustomization directory. This flag can't be used together with -f or -R.")
	fs.BoolVar(&a.openapiPatch, "openapi-patch", true,
		"If true, use openapi to calculate diff when the openapi presents and the resource can be found in the openapi spec. Otherwise, fall back to use baked-in types.")
	fs.StringVar(&a.output, "output", "",
		"Output format. One of: (json, yaml, kyaml, name, go-template, go-template-file, template, templatefile, jsonpath, jsonpath-as-json, jsonpath-file).")
	fs.StringVar(&a.output, "o", "",
		"Output format. One of: (json, yaml, kyaml, name, go-template, go-template-file, template, templatefile, jsonpath, jsonpath-as-json, jsonpath-file).")
	fs.BoolVar(&a.overwrite, "overwrite", true,
		"Automatically resolve conflicts between the modified and live configuration by using values from the modified configuration.")
	fs.BoolVar(&a.prune, "prune", false,
		"Automatically delete resource objects, that do not appear in the configs and are created by either apply or create --save-config. Should be used with either -l or --all.")
	fs.Var(&a.pruneAllowlist, "prune-allowlist",
		"Overwrite the default allowlist with <group/version/kind> for --prune.")
	fs.BoolVar(&a.recursive, "recursive", false,
		"Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	fs.BoolVar(&a.recursive, "R", false,
		"Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	fs.StringVar(&a.selector, "selector", "",
		"Selector (label query) to filter on, supports '=', '==', '!=', 'in', 'notin'.(e.g. -l key1=value1,key2=value2,key3 in (value3)). Matching objects must satisfy all of the specified label constraints.")
	fs.StringVar(&a.selector, "l", "",
		"Selector (label query) to filter on, supports '=', '==', '!=', 'in', 'notin'.(e.g. -l key1=value1,key2=value2,key3 in (value3)). Matching objects must satisfy all of the specified label constraints.")
	fs.BoolVar(&a.serverSide, "server-side", false,
		"If true, apply runs in the server instead of the client.")
	fs.BoolVar(&a.showManagedFields, "show-managed-fields", false,
		"If true, keep the managedFields when printing objects in JSON or YAML format.")
	fs.StringVar(&a.subresource, "subresource", "",
		"If specified, apply will operate on the subresource of the requested object. Only allowed when using --server-side.")
	fs.StringVar(&a.template, "template", "",
		"Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")
	fs.DurationVar(&a.timeout, "timeout", time.Duration(0),
		"The length of time to wait before giving up on a delete, zero means determine a timeout from the size of the object.")
	fs.StringVar(&a.validate, "validate", "strict",
		"Must be one of: strict (or true), warn, ignore (or false). \"true\" or \"strict\" will use a schema to validate the input and fail the request if invalid. It will perform server side validation if ServerSideFieldValidation is enabled on the api-server, but will fall back to less reliable client-side validation if not. \"warn\" will warn about unknown or duplicate fields without blocking the request if server-side field validation is enabled on the API server, and behave as \"ignore\" otherwise. \"false\" or \"ignore\" will not perform any schema validation, silently dropping any unknown or duplicate fields.")
	fs.BoolVar(&a.wait, "wait", false, "If true, wait for resources to be gone before returning. This waits for finalizers.")
}

func (a *KubectlApply) Run(ctx context.Context, args []string) error {
	// left as an exercise for the reader...
	return nil
}

func (a *KubectlApply) Name() string {
	return "apply"
}

func (a *KubectlApply) UsageLine() string {
	return "kubectl apply (-f FILENAME | -k DIRECTORY)"
}

func (a *KubectlApply) ShortHelpText() string {
	return "Apply a configuration change to a resource from a file or stdin."
}

func (a *KubectlApply) HelpText() string {
	return KubectlApplyDesc
}

func (a *KubectlApply) ExampleText() string {
	return KubectlApplyExamples
}
