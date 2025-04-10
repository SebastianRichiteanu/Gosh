# TODO

## Major features

- pipes
- job control: fg, bg & jobs
- tests:
   - autocompletion:
      Builtin completion
      Completion with arguments
      Missing completions
      Executable completion
      Multiple completions
      Partial completions
   - history
   - env vars
   - 

## Nice to Have

- support for ctrl+left/right arrow?
- support for temporary var assignment, eg: var TEST='123' echo $TEST
- support autocompletion for env vars
- readme
- color scheme? :D
- main.go is a bit of a mess
- what if I only write the history/alias file on exit?

## Bugs

### Command not found after vi

To reproduce:
   - open shell
   - vi random_file
   - write something + save & exit
   - prompt "ls + <ENTER>"
   - "ls: not found"

Also happens with "nano", probably something to do with the shell

### ^C on shell exit sometimes

The number is arbitray, not sure where is it coming from for now
- DisableRawMode()/EnableRawMode()

### autocompletion doesn't work again :((

$ /home/zee|/.gosh/goshrc
$ if i hit tab with the above prompt, it's bell not autocompletion?

### another issue with autocomplete

if the current token is a dir, it will directly add / instead of checking the path for multiplem matches

$ /mnt/d/Programming/C|
$ /mnt/d/Programming/C/ (after tab)
$ C C++ (should be)

### Closer not called if init fails

If for example the NewPrompt fct returns an error, the closer will not execute and mess the tty

