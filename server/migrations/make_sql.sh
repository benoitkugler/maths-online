echo "" > create_all_gen.sql && 
echo "-- prof/teacher/gen_create.sql" >> create_all_gen.sql &&
cat ../src/prof/teacher/gen_create.sql >> create_all_gen.sql 
echo "-- prof/editor/gen_create.sql" >> create_all_gen.sql &&
cat ../src/prof/editor/gen_create.sql >> create_all_gen.sql 
echo "-- tasks/gen_create.sql" >> create_all_gen.sql &&
cat ../src/tasks/gen_create.sql >> create_all_gen.sql 
echo "-- prof/trivial/gen_create.sql" >> create_all_gen.sql &&
cat ../src/prof/trivial/gen_create.sql >> create_all_gen.sql 
echo "-- prof/homework/gen_create.sql" >> create_all_gen.sql &&
cat ../src/prof/homework/gen_create.sql >> create_all_gen.sql 
 