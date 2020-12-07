package cli

type ImpactCommand struct {
	Factory impactFactory
}

func NewImpactCommand() ImpactCommand {
	return ImpactCommand{newImpactFactory()}
}

func (cmd ImpactCommand) Usage() string {
	return `Terraform Impact.

This tool takes a list of files as input and outputs a list of all the Terraform states
impacted by any of those files. An impact is described as a file creation, modification
or deletion.

Usage:
  impact <files>... [--rootdir <dir>] [--pattern <string>] [--user <credentials>]
  impact -h | --help
  impact -v | --version

Arguments:
  <files>                  List of files that could impact any of the Terraform states.
                           When <files> is the url to a GitHub pull request, uses files
                           from the pull request.

Options:
  -r --rootdir <dir>       The directory from where the state discovery begins.
  -p --pattern <string>    A string to filter states. Only states whose path contains the
                           string will be taken into account.
  -u --user <credentials>  Credentials to access GitHub pull requests. Follows the curl
                           format 'username:password'. You can also pass credentials as
                           environment variables through: GITHUB_USERNAME and
                           GITHUB_PASSWORD. Note that the option always takes precendence
                           over environment variables.
  -h --help                Show this screen.
  -v --version             Show version.`
}

func (cmd ImpactCommand) Run(opts ImpactOptions) error {
	impacter, service, outputer := cmd.Factory.Create(opts)

	impacterFiles, err := impacter.List()
	if err != nil {
		return err
	}

	result, err := service.Impact(impacterFiles)
	if err != nil {
		return err
	}

	outputer.Output(result)

	return nil
}
