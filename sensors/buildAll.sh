source_dir=./cmd
build_dir=./exe

if [ ! -d "$build_dir" ]; then
    mkdir "$build_dir"
fi
rm -r "$build_dir"/*

sensor_files=("pressureSensor" "tempSensor" "windSensor")

for file in "${sensor_files[@]}"; do
    echo "Building $file"
    go build -o "$build_dir/$file" "$source_dir/$file"
done

echo "Build process completed."

for file in "${sensor_files[@]}"; do
    echo "Running $file"
    ."/$build_dir/$file" -id 44 -airport A234C -frequency 10 &
done

# ."/$build_dir/tempSensor" -id 44 -airport A234C -frequency 10 &
# ."/$build_dir/tempSensor" -id 55 -airport A333C -frequency 10 &
# ."/$build_dir/tempSensor" -id 66 -airport A555C -frequency 10 &

cleanup() {
    echo "Sensors stopping..."
    
    for file in "${sensor_files[@]}"; do
        pkill -TERM -f "$build_dir/$file"
    done
    
    wait
    
    echo "Stopped"
}

trap 'cleanup' INT

wait

