# CnCNet map loader for RedAlert2 YR
This app download all maps RedAlert2 YR (and other) from CnCNet database.


### Download
Download last app from release page https://github.com/vmpartner/cncnet-map-loader/releases


### Usage
Just run app, then loader create 2 folder, tmp and maps. App will download maps to tmp and unzip to maps.
Folder tmp autoclean.
```
./cncnet-map-loader_windows_amd64.exe -timeout=5000 -game="yr"
```

### Config
```
game - Set game search (default "yr")
timeout - Set timeout between request to cncnet server in milliseconds (default 500)
```

### Install maps
Just copy all maps to RedAlert2 YR game folder in Maps/Custom/

### Thanks
If app was helpful for you, please add star, thank you!

### Warning
Game start long time if there is too much map files!
