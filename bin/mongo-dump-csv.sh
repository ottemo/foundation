#!/bin/sh

# Use this script to dump all mongodb collections to CSV format for import
# into mysql or sqlite tables.  Tables will be dumped into a local folder with
# the same name as DBNAME and dropped into files named after their collection.

OIFS=$IFS;
IFS=",";

# fill in your details here
dbname=kg-dev
user=USER
pass=PASSWD
host=candidate.42.mongolayer.com:10243

# first get all collections in the database
collections=`mongo "$host/$dbname" --username $user --password $pass --eval "rs.slaveOk();db.getCollectionNames();" --quiet`;
collectionArray=($collections);

# for each collection
for ((i=0; i<${#collectionArray[@]}; ++i));
do
    echo 'exporting collection' ${collectionArray[$i]}
    # get comma separated list of keys. do this by peeking into the first document in the collection and get his set of keys
    keys=`mongo "$host/$dbname" --username $user --password $pass --eval "rs.slaveOk();var keys = []; for(var key in db.${collectionArray[$i]}.find().sort({_id: -1}).limit(1)[0]) { keys.push(key);  }; keys;" --quiet`;
    # now use mongoexport with the set of keys to export the collection to csv
    mongoexport --host $host --username $user --password $pass --quiet --db $dbname --collection ${collectionArray[$i]} --type=csv --fields "$keys" --out $dbname.${collectionArray[$i]}.csv;
done

IFS=$OIFS;
