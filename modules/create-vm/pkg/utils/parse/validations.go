package parse

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"strings"
	"unicode"
)

func (c *CLIOptions) assertValidTypes() error {
	if !output.IsOutputType(string(c.Output)) {
		return zerrors.NewMissingRequiredError("%v is not a valid output type", c.Output)
	}
	return nil
}

func (c *CLIOptions) resolveTemplateParams() error {
	var paramsError zerrors.MultiError

	for i, param := range c.TemplateParams {
		trimmedParam := strings.TrimLeftFunc(param, unicode.IsSpace)
		c.TemplateParams[i] = trimmedParam
		split := strings.SplitN(trimmedParam, templateParamSep, 2)
		if len(split) < 2 || split[0] == "" {
			paramsError.Add(fmt.Sprintf("param %d \"%v\"", i, param), zerrors.NewMissingRequiredError("param %v has incorrect format: should be KEY:VAL", param))
		}
	}

	return paramsError.
		ShortPrint("following parameters have incorrect format: should be KEY:VAL :").
		AsOptional()
}

func (c *CLIOptions) trimSpacesAndReduceCount() {
	c.TemplateName = strings.TrimSpace(c.TemplateName)
	c.setTemplateNamespace(strings.TrimSpace(c.GetTemplateNamespace()))             // reduce count to 1
	c.setVirtualMachineNamespace(strings.TrimSpace(c.GetVirtualMachineNamespace())) // reduce count to 1

	for i, v := range c.TemplateParams {
		c.TemplateParams[i] = strings.TrimLeftFunc(v, unicode.IsSpace)
	}
	for i, v := range c.DataVolumes {
		c.DataVolumes[i] = strings.TrimSpace(v)
	}
	for i, v := range c.OwnDataVolumes {
		c.OwnDataVolumes[i] = strings.TrimSpace(v)
	}
	for i, v := range c.PersistentVolumeClaims {
		c.PersistentVolumeClaims[i] = strings.TrimSpace(v)
	}
	for i, v := range c.OwnPersistentVolumeClaims {
		c.OwnPersistentVolumeClaims[i] = strings.TrimSpace(v)
	}
}

func (c *CLIOptions) resolveDefaultNamespaces() error {
	tempNamespace := c.GetTemplateNamespace()
	vmNamespace := c.GetVirtualMachineNamespace()

	if tempNamespace == "" || vmNamespace == "" {
		activeNamespace, err := env.GetActiveNamespace()
		if err != nil {
			return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), c.getMissingNamespaceOptionNames())
		}
		if tempNamespace == "" {
			c.setTemplateNamespace(activeNamespace)
		}
		if vmNamespace == "" {
			c.setVirtualMachineNamespace(activeNamespace)
		}
	}
	return nil
}
