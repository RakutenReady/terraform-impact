package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RakutenReady/terraform-impact/impact"
)

func TestCreateCommandLineImpacter(t *testing.T) {
	opts := validImpactOptions()
	opts.Files = []string{"File_1", "File_2", "File_3"}
	result := createImpacter(opts)

	assert := assert.New(t)
	assert.IsType(impact.ImpacterImpl{}, result, "Result should be of ImpacterImpl type")

	impacter := result.(impact.ImpacterImpl)
	assert.IsType(impact.CommandLineImpacter{}, impacter.Inner, "Impacter.Inner should be of CommandLineImpacter type")

	inner := impacter.Inner.(impact.CommandLineImpacter)
	assert.ElementsMatch(opts.Files, inner.Files)
}

func TestCreateGitHubPullRequestImpacter(t *testing.T) {
	filesWithGitHubUrl := []string{"https://github.com/bob", "whatever-comes-second", "or-third-doesnt-matter"}
	testCases := []struct {
		Credentials  string
		WantUsername string
		WantPassword string
	}{
		{"a_user@hotmail.com:nice_gh_pwd123!", "a_user@hotmail.com", "nice_gh_pwd123!"},
		{"simple_user:simple_password", "simple_user", "simple_password"},
		{"weird_one:because:password:is:split", "weird_one", "because:password:is:split"},
		{"", "", ""},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		opts := validImpactOptions()
		opts.Files = filesWithGitHubUrl
		opts.Credentials = testCase.Credentials

		result := createImpacter(opts)

		assert.IsType(impact.ImpacterImpl{}, result, "Result should be of ImpacterImpl type")

		impacter := result.(impact.ImpacterImpl)
		assert.IsType(impact.GitHubPullRequestImpacter{}, impacter.Inner, "Impacter.Inner should be of GitHubPullRequestImpacter type")

		inner := impacter.Inner.(impact.GitHubPullRequestImpacter)
		assert.Equal(opts.Files[0], inner.Url, "url should match")
		assert.Equal(testCase.WantUsername, inner.Username, "username should match")
		assert.Equal(testCase.WantPassword, inner.Password, "password should match")
	}
}

func TestCreateGitHubPullRequestImpacterPanicsOnInvalidCreds(t *testing.T) {
	wantPanicMsg := "Invalid credentials format, use <username:password>"

	opts := validImpactOptions()
	opts.Files = []string{"https://github.com/bob", "whatever-comes-second", "or-third-doesnt-matter"}
	opts.Credentials = "user-password-without-separator"

	shouldPanicFn := func() {
		createImpacter(opts)
	}

	assert.PanicsWithErrorf(t, wantPanicMsg, shouldPanicFn, "createImpacter should panic because of credentials [%v]", opts.Credentials)
}
