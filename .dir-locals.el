;;; Directory Local Variables
;;; For more information see (info "(emacs) Directory Variables")

;; Emacs is the main editor for this project.
;;
;; Go workflow:
;; - GNU Guix as the main development environment
;; - Eglot as the LSP client
;; - gopls as the language server
;; - gofmt for formatting
;; - `make guix-check` as the project-wide verification command
;;
;; Scheme workflow:
;; - edit Guix manifests and package definitions directly
;; - use `make pkg` / `guix build -f ./guix.scm` as the build feedback loop

((nil . ((fill-column . 80)
         (sentence-end-double-space . nil)
         (require-final-newline . t)
         (show-trailing-whitespace . t)
         (compile-command . "make guix-check")))

 (go-mode . ((indent-tabs-mode . t)
             (tab-width . 8)
             (compile-command . "make guix-test")
             (eglot-workspace-configuration
              . (:gopls
                 ((gofumpt . t)
                  (staticcheck . t)
                  (usePlaceholders . t)
                  (completeUnimported . t)
                  (semanticTokens . t)
                  (matcher . "CaseSensitive")
                  (analyses . ((fieldalignment . t)
                               (nilness . t)
                               (shadow . t)
                               (unusedparams . t)
                               (unusedwrite . t)
                               (useany . t))))))
             (eval . (setq-local eglot-server-programs
                                 '(((go-mode go-ts-mode go-mod-mode go-mod-ts-mode)
                                    . ("gopls")))))))

 (go-ts-mode . ((indent-tabs-mode . t)
                (tab-width . 8)
                (compile-command . "make guix-test")
                (eglot-workspace-configuration
                 . (:gopls
                    ((gofumpt . t)
                     (staticcheck . t)
                     (usePlaceholders . t)
                     (completeUnimported . t)
                     (semanticTokens . t)
                     (matcher . "CaseSensitive")
                     (analyses . ((fieldalignment . t)
                                  (nilness . t)
                                  (shadow . t)
                                  (unusedparams . t)
                                  (unusedwrite . t)
                                  (useany . t))))))
                (eval . (setq-local eglot-server-programs
                                    '(((go-mode go-ts-mode go-mod-mode go-mod-ts-mode)
                                       . ("gopls")))))))

 (go-mod-mode . ((indent-tabs-mode . t)
                 (tab-width . 8)
                 (compile-command . "make guix-test")))

 (go-mod-ts-mode . ((indent-tabs-mode . t)
                    (tab-width . 8)
                    (compile-command . "make guix-test")))

 (scheme-mode . ((indent-tabs-mode . nil)
                 (fill-column . 78)
                 (compile-command . "make pkg")))

 (emacs-lisp-mode . ((indent-tabs-mode . nil)
                     (fill-column . 80)))

 (yaml-mode . ((indent-tabs-mode . nil)))
 (yaml-ts-mode . ((indent-tabs-mode . nil)))
 (sh-mode . ((indent-tabs-mode . nil)))
 (conf-unix-mode . ((indent-tabs-mode . nil)))
 (makefile-mode . ((indent-tabs-mode . t))))
