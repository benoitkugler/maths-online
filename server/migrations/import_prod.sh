# This script is a workaround for a bug related to user defined functions
# It imports an SQL dump in two steps : functions and schema then data
# The first argument is the name of the file dump
#
# The typical commands to run previoulsy are 
# pg_dump -U <user> -h <host> -d education_prod > isyro.dump
# scp <user>@<host>:isyro.dump isyro.dump

echo "Spliting schema and data..." && 
line=$(grep -n "COPY" $1 | cut -d: -f1 | head -1) &&  
line="$((line-1))" && 
(head -$line > schema.sql; cat > data.sql) < $1 && 
echo "Resetting DB..." && 
dropdb --if-exists --force isyro_prod && createdb isyro_prod && 
echo "Importing..." && 
psql isyro_prod < schema.sql && 
psql isyro_prod < data.sql && 
echo "Cleaning up" && 
rm schema.sql data.sql && 
echo "Setting passwords to 1234..." && 
psql isyro_prod -c  "UPDATE teachers SET passwordcrypted = '\xcd9c6d2bb0ce633cbc55d4a138f6945f89a5351db810c556a5a36f140e6116c1'"