#!/bin/zsh

names="$1"
weights="$2"
place="$3"

IFS=',' read -r -A names_array <<< "$names"
IFS=',' read -r -A weights_array <<< "$weights"

typeset -A name_sums

for ((i = 1; i <= ${#names_array[@]}; i++)); do
  name="${names_array[$i]}"
  weight="${weights_array[$i]}"
  sum=0  
  for ((j = 0; j < ${#name}; j++)); do
    letter="${name:$j:1}"
    letter=$(echo "$letter" | tr '[:upper:]' '[:lower:]')
    ascii_value=$(printf "%d" "'$letter")
    index=$((ascii_value - 96))
    sum=$((sum + index + 1))
  done
  name_sums["$name"]=$((sum * weight))
done

sorted_names=(
  $(for value in "${(@k)name_sums}"; do
    echo "$value ${name_sums[$value]}"
  done | sort -k2,2nr | awk '{print $1}')
)

echo "${sorted_names[(($place))]}: ${name_sums[${sorted_names[(($place))]}]}"