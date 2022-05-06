# This script is a workaround a bug related to used defined functions
# It imports an SQL dump in two steps : functions and schema then data
line=$(grep -n "COPY" $1 | cut -d: -f1 | head -1) && 
line="$((line-1))" && 
(head -$line > schema.sql; cat > data.sql) < $1 && 
dropdb --if-exists isyro_prod && createdb isyro_prod && 
psql isyro_prod < schema.sql && 
psql isyro_prod < data.sql && 
echo "Cleaning up" && 
rm schema.sql data.sql