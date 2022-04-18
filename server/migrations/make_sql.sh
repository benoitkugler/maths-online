echo "" > create_all_gen.sql && 
echo "-- maths/exercice/create_gen.sql" >> create_all_gen.sql &&
cat ../src/maths/exercice/create_gen.sql >> create_all_gen.sql  &&
echo "-- prof/trivial-poursuit/create_gen.sql" >> create_all_gen.sql &&
cat ../src/prof/trivial-poursuit/create_gen.sql >> create_all_gen.sql 
