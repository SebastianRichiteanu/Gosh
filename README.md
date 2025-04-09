# TODO

## Major features

- pipes
- job control: fg, bg & jobs
- alias? unalias?
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

### home dir not handled in source
 source ~/.gosh/goshrc

### autocompletion doesn't work again :((

$ /home/zee|/.gosh/goshrc
$ if i hit tab with the above prompt, it's bell not autocompletion?