package cmder

// collectSubcommands collects the immediate subcommands of the given [Command] into a map keyed by the command
// [Command] Name(). Returns an empty map if the command is not a [RootCommand].
func collectSubcommands(cmd Command) map[string]Command {
	subcommands := map[string]Command{}

	c, ok := cmd.(RootCommand)
	if !ok {
		return subcommands
	}

	for _, subcommand := range c.Subcommands() {
		subcommands[subcommand.Name()] = subcommand
	}

	return subcommands
}
