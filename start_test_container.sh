#
# Copyright © 2022 Paulo Villela. All rights reserved.
# Use of this source code is governed by the Apache 2.0 license
# that can be found in the LICENSE file.
#

docker run -d --rm -p 9999:5432 --name realworld-pg -v realworldVol:/var/lib/postgresql/data -e \
    POSTGRES_PASSWORD=realworld-pg postgres
