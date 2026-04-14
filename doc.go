/*
Package cmder is a simple and flexible library for building command-line interfaces in Go. cmder takes a very
opinionated approach to building command-line interfaces. The library will help you define, structure and execute your
commands, but that's about it. cmder embraces simplicity because sometimes, less is better. The wide range of examples
throughout the project should help you get started.

To define a new command, simply define a type that implements the [Command] interface. If you want your command to have
additional behavior like flags or subcommands, simply implement the appropriate interfaces.

  - cmder doesn't force you to use special command structs. As long as you implement our narrow interfaces, you're good
    to go.
  - cmder is unobtrusive. Define your command and execute it.
  - cmder is totally stateless making it super easy to unit test your commands. This isn't the case in other libraries.
  - cmder natively supports getopt-style flag parsing, but you can use the standard flag library instead if you prefer.
  - We take great pride in our documentation. If you find anything unclear, please let us know so we can fix it.

To get started, see [Command] and [Execute].

For POSIX/GNU flag parsing, see package [getopt].
*/
package cmder
