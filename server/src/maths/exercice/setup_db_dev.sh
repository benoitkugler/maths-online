echo "Resetting DB..."
dropdb --if-exists maths_dev &&
createdb maths_dev &&
echo "Creating tables" && 
psql maths_dev < create.sql &&
echo "Done."