# Ebox
*Manage all your Emacs distributions with ease.*


Ebox is [emacs-distribution](https://www.emacswiki.org/emacs/emacs-distribution) on steroids. It is a single binary, independent of operating system and shell, and offers a complete Emacs sandbox.
All your Emacs distributions live in `~/emacs` and can be executed independently from one another. So you can e. g. use [spacemacs](https://github.com/syl20bnr/spacemacs) alongside [prelude](https://github.com/bbatsov/prelude), or [Emacs live](https://github.com/overtone/emacs-live) alongside a custom configuration, without them interfering with each other.

## Building

Run `make`. Ebox is installed as `$GOPATH/bin/ebox`.

## Usage

Running just `ebox` lists all distributions in `~/emacs`.
```
$ ebox
custom
spacemacs
```

Running Ebox with the name of an existing distribution spawns an instance of that distribution.
```
$ ebox spacemacs
```
In this example the environment variable `HOME` is set to `~/emacs/spacemacs` inside that Emacs instance.

Ebox can also download distributions. For some well-known distributions you don't even have to specify a URL:
```
$ #This example clones https://github.com/bbatsov/prelude
$ ebox prelude
Cloning into '/home/laerling/emacs/prelude/.emacs.d'...
remote: Counting objects: 5783, done.
remote: Compressing objects: 100% (10/10), done.
remote: Total 5783 (delta 5), reused 7 (delta 3), pack-reused 5770
Receiving objects: 100% (5783/5783), 4.54 MiB | 874.00 KiB/s, done.
Resolving deltas: 100% (3435/3435), done.
```

If the distribution is hosted on github, you don't have to write the complete URL either:
```
$ #This example clones https://github.com/githubusername/distroname
$ ebox githubusername/distroname
```

For all other distributions, simply specify the URL:
```
$ #This example clones https://domain.tld/gitlabusername/distroname.git
$ ebox domain.tld/gitlabusername/distroname.git
```
