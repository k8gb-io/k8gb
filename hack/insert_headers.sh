for f in $(find . -name '*.go'); do
  if grep -q "Apache License, Version 2.0" $f; then
      echo "License header is already present. Skipping $f"
      continue
  fi
  echo "Adding header to $f"
  cat ./hack/boilerplate.go.txt $f > $f.new
  mv $f.new $f
done
