# mysql-operator
Simple MySQL Kubernetes operator that allows create databases with users in external cluster via CRDs.

We use this operator to manage databases in our dev environments.

## build & deploy

1. Download Waypoint https://www.waypointproject.io/downloads
2. Run `waypoint up -var=mysql_host=... -var=mysql_user=... -var=mysql_password=... -var=mysql_default_db=mysql -var=dockerconfigjson=...`
3. ???
4. PROFIT
