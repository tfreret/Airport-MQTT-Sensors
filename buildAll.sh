source_dir=./cmd
build_dir=./exe

if [ ! -d "$build_dir" ]; then
    mkdir "$build_dir"
fi

if [ "$(ls -A "$build_dir")" ]; then
    rm -r "$build_dir"/*
fi

files=("sensors/pressureSensor" "sensors/tempSensor" "sensors/windSensor")
args=("-id 44 -airport KJFK -frequency 10" "-id 55 -airport KJFK -frequency 10" "-id 66 -airport KJFK -frequency 10")

for file in "${files[@]}"; do
    echo "Building $file"
    go build -o "$build_dir/$file" "$source_dir/$file"
done

echo "Build completed."

# echo "lunching mosquitto serveur"
# mosquitto &

for i in "${!files[@]}"; do
    file="${files[i]}"
    arg="${args[i]}"

    echo "Running $file with arguments: $arg"
    ."/$build_dir/$file" $arg &
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