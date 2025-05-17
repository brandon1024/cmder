/*
cmder is a simple and flexible library for building command-line interfaces in Go. If you're coming from Cobra and have
used it for any length of time, you have surely had your fair share of difficulties with the library. cmder will feel
quite a bit more comfortable and easy to use, and the wide range of examples throughout the project should help you get
started.

cmder takes a very opinionated approach to building command-line interfaces. The library will help you define, structure
and execute your commands, but that's about it. cmder embraces simplicity because sometimes, less is better.

To define a new commands, simply define a type that implements the [Command] interface. If you want your command to have
additional behaviour like flags or subcommands, simply implement the appropriate interfaces.

  - Bring your own types. cmder doens't force you to use special command structs. As long as you implement our narrow
    interfaces, you're good to go!
  - cmder is unobtrustive. Define your command and execute it. Simplicity above all else!
  - cmder is totally stateless making it super easy to unit test your commands. This isn't the case in other libraries.
  - We take great pride in our documentation. If you find anything unclear, please let us know so we can fix it.

To get started, see [Command] and [Execute].
*/
package cmder
