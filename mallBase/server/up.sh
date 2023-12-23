#!/usr/bin/env bash

git pull
declare -a arrs
count=1
while read rows
do
if [[ $rows =~ "gitlab.mall.com" ]]; then
 array=(${rows// v/ })
 if [[ ${array[0]} != "module" ]]; then
   arrs[$count]=${array[0]}
   count=$[$count + 1]
 fi
fi
done < go.mod
for var in ${arrs[@]}
do
echo "go get ${var}"
go get ${var}
done
git add ./
git commit -m "up-${name}"
git push