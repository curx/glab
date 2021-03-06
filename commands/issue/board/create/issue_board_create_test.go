package create

import (
	"fmt"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/commands/cmdutils"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func TestNewCmdCreate(t *testing.T) {
	api.CreateIssueBoard = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateIssueBoardOptions) (*gitlab.IssueBoard, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "NS/WRONG_REPO" {
			return nil, fmt.Errorf("error expected")
		}

		return &gitlab.IssueBoard{
			ID:        11,
			Name:      *opts.Name,
			Project:   &gitlab.Project{PathWithNamespace: projectID.(string)},
			Milestone: nil,
			Lists:     nil,
		}, nil
	}
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr bool
	}{
		{
			name: "Name passed as arg",
			arg:  `"Test"`,
			want: `✓ Board created: "Test"`,
		},
		{
			name: "Name passed in name flag",
			arg:  `--name "Test"`,
			want: `✓ Board created: "Test"`,
		},
		{
			name:    "WRONG_REPO",
			arg:     `"Test" -R NS/WRONG_REPO`,
			wantErr: true,
		},
	}
	io, _, stdout, stderr := iostreams.Test()

	f := cmdtest.StubFactory("https://gitlab.com/glab-cli/test")
	f.IO = io
	f.IO.IsaTTY = true
	f.IO.IsErrTTY = true

	cmd := NewCmdCreate(f)
	cmdutils.EnableRepoOverride(cmd, f)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cmdtest.RunCommand(cmd, tc.arg)
			if tc.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			out := stripansi.Strip(stdout.String())

			assert.Contains(t, out, tc.want)
			assert.Contains(t, stderr.String(), "")

		})
	}
}
