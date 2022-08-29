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
cat ../src/sql/homework/gen_create.sql >> create_all_gen.sql 
 