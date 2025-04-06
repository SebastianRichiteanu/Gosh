[![progress-banner](https://backend.codecrafters.io/progress/shell/6165853b-9273-400c-ae13-04c0b8d1fc81)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview).

In this challenge, you'll build your own POSIX compliant shell that's capable of
interpreting shell commands, running external programs and builtin commands like
cd, pwd, echo and more. Along the way, you'll learn about shell command parsing,
REPLs, builtin commands, and more.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# Passing the first stage

The entry point for your `shell` implementation is in `cmd/myshell/main.go`.
Study and uncomment the relevant code, and push your changes to pass the first
stage:

```sh
git commit -am "pass 1st stage" # any msg
git push origin master
```

Time to move on to the next stage!

# Stage 2 & beyond

Note: This section is for stages 2 and beyond.

1. Ensure you have `go (1.19)` installed locally
1. Run `./your_program.sh` to run your program, which is implemented in
   `cmd/myshell/main.go`.
1. Commit your changes and run `git push origin master` to submit your solution
   to CodeCrafters. Test output will be streamed to your terminal.


# TODO

- autocomplete still broken for some stuff, really need to implement it better:
   $ ls
   test.sh
   $ ./test.<TAB>
   test..sh  test.s

probably split then input into tokens, get the current token(where cursor is) and autocomplete that
   - check if it's ~
   - check if it's a file or a dir (append / if dir)








- for ctrl+c print ^C and skip line

- -bash: cd: too many arguments
- todos
- pipes
- job control: fg, bg & jobs
- alias? unalias?
- persistent history?
- variable interpolation? 
- readme
- run linter :) and maybe add it as make and add a conifg for it?
- color scheme? :D
- license?

- for history: if i do up arrow a couple of times and then delete a char, should I reset the index or just go up? :hm:
   - i think I should keep the index but also can edit history
   - so for example if I delete a couple of chars from an above entry 
   - if i go up/down and again to that index the change should be there 
   - also maybe MaxHistorySize?


- tests for autocompletion!!!
   Builtin completion
   Completion with arguments
   Missing completions
   Executable completion
   Multiple completions
   Partial completions

