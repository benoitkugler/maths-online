cd migrations &&
echo "Bundling create SQL commands..." && 
bash make_sql.sh &&
echo "Resetting DB..."
dropdb --if-exists maths_dev &&
createdb maths_dev &&
echo "Creating tables" && 
psql maths_dev < create_all_gen.sql &&
echo "" &&
echo "Populating..." && 
psql maths_dev < setup_db_dev_populate.sql &&
echo "Done."