#https://zsh.sourceforge.io/Doc/Release/Functions.html
autoload -U colors && colors

# https://zsh.sourceforge.io/Doc/Release/Prompt-Expansion.html
PS1="%{$fg[green]%}%n%{$reset_color%}@%m:%{$fg[cyan]%}%3~%{$reset_color%} %# "
