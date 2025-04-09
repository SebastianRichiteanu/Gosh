# TODO

## Major features

- pipes
- job control: fg, bg & jobs
- alias? unalias?
- persistent history? (history command + HISTFILE + HISTSIZE)
- tests:
   - autocompletion:
      Builtin completion
      Completion with arguments
      Missing completions
      Executable completion
      Multiple completions
      Partial completions
   - history
   - ?

## Nice to Have

- support for ctrl+left/right arrow?
- support for temporary var assignment, eg: var TEST='123' echo $TEST
- support autocompletion for env vars
- readme
- color scheme? :D

## Bugs

### Command not found after vi

To reproduce:
   - open shell
   - vi random_file
   - write something + save & exit
   - prompt "ls + <ENTER>"
   - "ls: not found"

### ^C on shell exit sometimes

The number is arbitray, not sure where is it coming from for now
