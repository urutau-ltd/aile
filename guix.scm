;; guix.scm --- Guix package file for this library -*- mode: scheme; -*-
;;
;; SPDX-License-Identifier: AGPL-3.0-or-later or LGPL-3.0-or-later
;; Copyright © 2026 Urutau-Ltd <softwarelibre@urutau-ltd.org>
;;
;;   , _ ,      _    _            _                     _ _      _
;;  ( o o )    | |  | |          | |                   | | |    | |
;; /'` ' `'\   | |  | |_ __ _   _| |_ __ _ _   _ ______| | |_ __| |
;; |'''''''|   | |  | | '__| | | | __/ _` | | | |______| | __/ _` |
;; |\\'''//|   | |__| | |  | |_| | || (_| | |_| |      | | || (_| |
;;    """       \____/|_|   \__,_|\__\__,_|\__,_|      |_|\__\__,_|
;;
;; This program is free software: you can redistribute it and/or modify
;; it under the terms of the GNU General Public License as published by
;; the Free Software Foundation, either version 3 of the License, or (at
;; your option) any later version.
;;
;; This program is distributed in the hope that it will be useful, but
;; WITHOUT ANY WARRANTY; without even the implied warranty of
;; MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
;; General Public License for more details.
;;
;; You should have received a copy of the GNU General Public License
;; along with this program. If not, see <https://www.gnu.org/licenses/>.
(use-modules (gnu packages golang)
             (guix packages)
             (guix git-download)
             (guix build-system go)
             (guix utils)
             ((guix licenses)
              #:prefix license:)
             (guix gexp))

(define-public go-codeberg-org-urutau-ltd-aile
  (package
    (name "go-codeberg-org-urutau-ltd-aile")
    (version "1.0.0-1")
    (source
     (origin
       (method git-fetch)
       (uri (git-reference
             (url "https://codeberg.org/urutau-ltd/aile.git")
             (commit (string-append "v" version))))
       (file-name (git-file-name name version))
       (sha256
        (base32 "1jmym2q9p9b36vhic31vf42p8gcc6g9gxapds5v127g0hwbpkpz1"))))
    (build-system go-build-system)
    (arguments
     (list
      #:go go-1.25
      ;; Tests disabled until the go-build-system stops breaking
      ;; unit tests with 404 errors. They work fine locally, idk why
      ;; the tests phase is dumb and likes to break things.
      ;#:tests? #f
      #:import-path "codeberg.org/urutau-ltd/aile"))
    (home-page "https://codeberg.org/urutau-ltd/aile")
    (synopsis "Small http runtime for Go")
    (description
     "Package aile provides a small stdlib-first HTTP runtime for Go.")
    (license license:agpl3)))

go-codeberg-org-urutau-ltd-aile
