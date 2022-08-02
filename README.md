# Fusion360Updater
Updates a Fusion360 install automatically

This solution is slightly specific but it will automatically upgrade a Fusion360 install.

## Requirements
### Local file in C:\FusionUpdater named currentVersion.ini
  I prefill this to 0.0.0 but you can insert your current installed version
### Fusion installer found [HERE](https://dl.appstreaming.autodesk.com/production/installers/Fusion%20360%20Admin%20Install.exe) renamed to Fusion360Install.exe (without spaces) in the same directory (C:\FusionUpdater\)
  This is the installer file but also takes arguements for updating and uninstalling
  
## Notes
I use this along with Task Scheduler in Windows to automatically check daily for a new version. I work IT at a school and this has been a nightmare to stay on top of manually. 
