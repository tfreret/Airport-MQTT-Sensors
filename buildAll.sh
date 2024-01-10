source_dir=./cmd
build_dir=./exe

if [ ! -d "$build_dir" ]; then
    mkdir "$build_dir"
fi

if [ "$(ls -A "$build_dir")" ]; then
    rm -r "${build_dir:?}/"*
fi

files=("sensors/pressureSensor" "sensors/tempSensor" "sensors/windSensor" "file-recorder")

for file in "${files[@]}"; do
    echo "Building $file"
    go build -o "$build_dir/$file" "$source_dir/$file"
done

echo "Build completed."

# echo "lunching mosquitto serveur"
# mosquitto &

for i in "${!files[@]}"; do
    file="${files[i]}"

    echo "Running $file"
    ."/$build_dir/$file" &
done


cleanup() {
    echo "Process stopping..."
    
    for file in "${files[@]}"; do
        pkill -TERM -f "$build_dir/$file"
    done
    
    wait
}

trap 'cleanup' INT

wait

echo "All process stopped"