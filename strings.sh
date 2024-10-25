grep -rho 't\.Errorf\(.*\)' ./src --exclude='*_test.go' | sed -E 's/t\.Errorf\(\s*"(([^"\\]|\\.)*)".*/\1/' | sed -E 's/\\\"/"/g' | sort | uniq > keys.txt
if [ ! -f "keys.txt" ]; then
  echo "Error: keys.txt not found!"
  exit 1
fi

if [ ! -f "strings_de.ini" ]; then
  awk '{print $0"="}' "keys.txt" | sort > "strings_de.ini"
  echo "Created strings_de.ini with all keys."
else
  # Read existing keys from strings_de.ini
  existing_keys=$(awk -F= '{print $1}' "strings_de.ini")

  # Filter out existing keys from keys.txt
  new_keys=$(grep -Fxv -f <(echo "$existing_keys") "keys.txt")

  # Add new keys and sort the file
  {
    echo "$new_keys" | awk 'NF' | awk '{print $0"="}'
    cat "strings_de.ini"
  } | sort -u > "strings_de.ini.tmp" && mv "strings_de.ini.tmp" "strings_de.ini"

  echo "Updated strings_de.ini with new keys, if any."
fi