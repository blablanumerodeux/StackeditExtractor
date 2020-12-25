# StackEditExtractor

## How to run it 

1. Put your json extraction in the source folder  
2. Change the following const in the code accordingly:  
    - PATH_OUTPUT_FOLDER (default "target/")
    - JSON_INPUT_FILE (default "source/StackEdit workspace.json")
3. remove BOM sequence  
   Stackedit inserts a BOM sequence on their extract, remove it with dos2unix:
```bash
file -i out.file
sudo pacman -Sy dos2unix
dos2unix source/StackEdit\ workspace.json
```

[removing a BOM](https://unix.stackexchange.com/questions/381230/how-can-i-remove-the-bom-from-a-utf-8-file)
