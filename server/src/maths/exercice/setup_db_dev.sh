echo "Resetting DB..."
dropdb --if-exists maths_dev &&
createdb maths_dev &&
echo "Creating tables" && 
psql maths_dev < create_gen.sql &&
echo "Populating..." && 
psql maths_dev < setup_db_dev_populate.sql &&
echo "Done."