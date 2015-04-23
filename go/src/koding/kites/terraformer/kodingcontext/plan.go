package kodingcontext

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/hashicorp/terraform/command"
	"github.com/hashicorp/terraform/terraform"
)

func (c *Context) Plan(content io.Reader, destroy bool) (*terraform.Plan, error) {
	cmd := command.PlanCommand{
		Meta: command.Meta{
			ContextOpts: c.TerraformContextOpts(),
			Ui:          c.ui,
		},
	}

	// copy all contents from remote to local for operating
	if err := c.RemoteStorage.Clone(c.Location, c.LocalStorage); err != nil {
		return nil, err
	}

	basePath, err := c.LocalStorage.BasePath()
	if err != nil {
		return nil, err
	}

	outputDir := path.Join(basePath, c.Location)
	mainFileRelativePath := path.Join(c.Location, mainFileName+terraformFileExt)
	planFilePath := path.Join(outputDir, planFileName+terraformPlanFileExt)
	stateFilePath := path.Join(outputDir, stateFileName+terraformStateFileExt)

	// override the current main file
	if err := c.LocalStorage.Write(mainFileRelativePath, content); err != nil {
		return nil, err
	}

	// TODO: doesn't work because Module is not initializde and it panics.
	// Module is initialized inside cmd.Run, but then it will override
	// anything we set here. Seems the only way to pass the variable is to
	// pass with the file it self - arslan
	//
	// variables := []*config.Variable{}
	// for k, v := range c.Variables {
	// 	variables = append(variables, &config.Variable{
	// 		Name:    k,
	// 		Default: v,
	// 	})
	// }
	// cmd.ContextOpts.Module.Config().Variables = variables

	args := []string{
		"-no-color", // dont write with color
		// "-detailed-exitcode", // give more info on exit
		"-out", planFilePath, // save plan to a file
		"-state", stateFilePath,
		outputDir,
	}

	if destroy {
		args = append([]string{"-destroy"}, args...)
	}

	exitCode := cmd.Run(args)

	fmt.Printf("Debug output: %+v\n", c.Buffer.String())

	if exitCode != 0 {
		return nil, fmt.Errorf("plan failed with code: %d", exitCode)
	}

	planFile, err := os.Open(planFilePath)
	if err != nil {
		return nil, err
	}
	defer planFile.Close()

	// copy all contents from local to remote for later operating
	if err := c.LocalStorage.Clone(c.Location, c.RemoteStorage); err != nil {
		return nil, err
	}

	return terraform.ReadPlan(planFile)
}
