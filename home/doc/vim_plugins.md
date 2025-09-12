# NAME

**vim-plugins** - install instructions for Vim plugins.


# DESCRIPTION

The instructions use Vim packages to install plugins directly from GitHub,
`:h packages`.

## Why manual labor with Vim packages?

There is a number of Vim package managers available on the market: Pathogen,
VimPlug, etc. These are pretty much Vim scripts manipulating Vim environment in
any way.

Luckily Vim 8 comes with built-in plugin manager. The nice features of this
manager, is that plugins reside at isolated paths, grouped under a package and
can be distributed together:

```
~/.vim/pack/<package>/start/<plugin>/{autoload,plugin,doc,...}
```

One can even change the location of default `~/.vim/pack` or use multiple
locations via `packpath`.

## Install a plugin

We'll use github.com/noahfrederick/vim-skeleton as an example. The username goes
under `pack/`. The repository name is the package name:

```console
% mkdir -p ~/.vim/pack/noahfrederick/start
% git -C ~/.vim/pack/noahfrederick/start clone https://github.com/noahfrederick/vim-skeleton.git
```

Generate help tags for the newly installed plugin (see `:help package-doc`):

```console
% vim -c 'helptags ~/.vim/pack/noahfrederick/start/vim-skeleton/doc'
```

## Update the plugin

```console
% git -C ~/.vim/pack/noahfrederick/start/vim-skeleton pull --rebase
```

Regenerate help tags.


# CONFIGURATION

## github.com/noahfrederick/vim-skeleton

The [vimrc](../vimrc) configures the plugins with custom replacement functions
to inject the year, Go-package, etc. into the templates.

The templates are included in [vim](../vim).
