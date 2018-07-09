# Ebox
*Manage all your Emacs distributions with ease.*


Ebox is [emacs-distribution](https://www.emacswiki.org/emacs/emacs-distribution) on steroids. It is a single binary, independent of operating system and shell, and offers a complete Emacs sandbox.
All your Emacs distributions live in `~/emacs` and can be executed independently from one another. So you can e. g. use [spacemacs](https://github.com/syl20bnr/spacemacs) alongside [prelude](https://github.com/bbatsov/prelude), or [Emacs live](https://github.com/overtone/emacs-live) alongside a custom configuration, without them interfering with each other.

## Building

Run `make`. Ebox is installed as `$GOPATH/bin/ebox`.
