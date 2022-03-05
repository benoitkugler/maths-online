# Build the project and copy the files into the static server folder 
npm run build &&
cd .. &&
rm -r server/static/prof &&
mkdir server/static/prof && 
cp -r prof/dist/* server/static/prof/ &&
echo "Fichier copiés dans server/static/prof/"