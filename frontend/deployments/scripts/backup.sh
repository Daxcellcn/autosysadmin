#!/bin/bash
# frontend/deployments/scripts/backup.sh

#!/bin/bash
set -e

DATE=$(date +%Y-%m-%d_%H-%M-%S)
BACKUP_DIR="/backups/autosysadmin_$DATE"

echo "Creating backup directory $BACKUP_DIR..."
mkdir -p $BACKUP_DIR

echo "Backing up PostgreSQL database..."
PGPASSWORD=secret pg_dump -h localhost -U autosysadmin -d autosysadmin > $BACKUP_DIR/db_backup.sql

echo "Backing up Redis data..."
redis-cli SAVE
cp /var/lib/redis/dump.rdb $BACKUP_DIR/redis_backup.rdb

echo "Compressing backup..."
tar -czvf $BACKUP_DIR.tar.gz $BACKUP_DIR

echo "Backup completed: $BACKUP_DIR.tar.gz"