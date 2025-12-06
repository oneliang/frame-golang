#!/bin/bash
project_array=("test" "transfer")
dependencies_array=("github.com/oneliang/util-golang/constants@main" "github.com/oneliang/util-golang/base@main" "github.com/oneliang/util-golang/common@main")
for ((i=0;i<${#project_array[@]};i++))
do
  echo "updating ${project_array[$i]} dependencies"
  if [ $i -eq 0 ]; then
    cd ${project_array[$i]}
  else
    cd "../${project_array[$i]}"
  fi
  for ((j=0;j<${#dependencies_array[@]};j++))
  do
    echo "exec go get, getting ${dependencies_array[$j]}"
    go get -u ${dependencies_array[$j]}
  done
done