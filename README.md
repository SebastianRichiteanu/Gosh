# TODO

## Major features

- pipes
- job control: fg, bg & jobs
- alias? unalias?
- persistent history? (history command + HISTFILE + HISTSIZE)
- variable interpolation? (${})
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
- readme
- color scheme? :D

## Bugs

### Command not found after vi

To reporduce:
   - open shell
   - vi random_file
   - write + save & exit
   - try "ls"
   - "ls: not found"

### ^C on shell exit sometimes

The number is arbitray, not sure where is it coming from for now
