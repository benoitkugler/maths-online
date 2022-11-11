echo "Grouping SQL statements in create_all_gen.sql..." &&
echo "" > create_all_gen.sql && 
echo "-- sql/teacher/gen_create.sql" >> create_all_gen.sql &&
cat ../src/sql/teacher/gen_create.sql >> create_all_gen.sql 
echo "-- sql/editor/gen_create.sql" >> create_all_gen.sql &&
cat ../src/sql/editor/gen_create.sql >> create_all_gen.sql 
echo "-- sql/trivial/gen_create.sql" >> create_all_gen.sql &&
cat ../src/sql/trivial/gen_create.sql >> create_all_gen.sql 
echo "-- sql/tasks/gen_create.sql" >> create_all_gen.sql &&
cat ../src/sql/tasks/gen_create.sql >> create_all_gen.sql 
echo "-- sql/homework/gen_create.sql" >> create_all_gen.sql &&
cat ../src/sql/homework/gen_create.sql >> create_all_gen.sql && 
echo "-- sql/reviews/gen_create.sql" >> create_all_gen.sql &&
cat ../src/sql/reviews/gen_create.sql >> create_all_gen.sql && 
echo "Splitting tables, constraints and json functions..."
cd sql_statements && 
go run sql.go &&
echo "Done."
 