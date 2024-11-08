package dashboards

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/perses/perses/go-sdk/dashboard"
	"gopkg.in/yaml.v3"
)

const (
	JSONOutput = "json"
	YAMLOutput = "yaml"
)

func init() {
	flag.String("output", YAMLOutput, "output format of the exec")
	flag.String("output-dir", "./dist", "output directory of the exec")
}

func executeDashboardBuilder(builder dashboard.Builder, outputFormat string, outputDir string, errWriter io.Writer) {
	var err error
	var output []byte
	if outputFormat == YAMLOutput {
		output, err = yaml.Marshal(builder.Dashboard)
	} else if outputFormat == JSONOutput {
		output, err = json.Marshal(builder.Dashboard)
	} else {
		err = fmt.Errorf("--output must be %q or %q", JSONOutput, YAMLOutput)
	}

	if err != nil {
		fmt.Fprint(errWriter, err)
		os.Exit(-1)
	}

	// create output directory if not exists
	_, err = os.Stat(outputDir)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprint(errWriter, err)
		os.Exit(-1)
	}

	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(outputDir, os.ModePerm)
	}

	_ = os.WriteFile(fmt.Sprintf("%s/%s.%s", outputDir, builder.Dashboard.Metadata.Name, outputFormat), output, os.ModePerm)
}

func NewExec() Exec {
	output := flag.Lookup("output").Value.String()
	outputDir := flag.Lookup("output-dir").Value.String()

	return Exec{
		outputFormat: output,
		outputDir:    outputDir,
	}
}

type Exec struct {
	outputFormat string
	outputDir    string
}

// BuildDashboard is a helper to print the result of a dashboard builder in stdout and errors to stderr
func (b *Exec) BuildDashboard(builder dashboard.Builder, err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
	executeDashboardBuilder(builder, b.outputFormat, b.outputDir, os.Stdout)
}
