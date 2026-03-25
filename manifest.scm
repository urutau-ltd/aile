;; Development manifest for aile.
;;
;; Most development in this repository is expected to happen through this
;; manifest. The repo uses Go 1.26.

(specifications->manifest
 (list "emacs"
       "gcc-toolchain"
       "git"
       "go@1.26"
       "go-golang-org-x-tools-godoc"
       "gopls"
       "make"
       "podman"
       "podman-compose"
       "ripgrep"))
