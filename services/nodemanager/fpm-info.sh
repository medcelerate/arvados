# Copyright (C) The Arvados Authors. All rights reserved.
#
# SPDX-License-Identifier: AGPL-3.0

case "$TARGET" in
    debian* | ubuntu*)
        fpm_depends+=(libcurl3-gnutls libpython2.7)
        ;;
esac
