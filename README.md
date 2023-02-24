# Fusion360Updater
Updates a Fusion360 install automatically

This solution is slightly specific but it will automatically deploy and upgrade Fusion360.

  
## Notes
I use this along with Task Scheduler in Windows to automatically check daily for a new version. I work IT at a school and this has been a nightmare to stay on top of manually. 

Below is the batch script I use in my Intune deployment

```
@echo off

MD C:\FusionUpdater
xcopy /s "%~dp0\fusionUpdater.exe" "C:\FusionUpdater"

SCHTASKS /CREATE /SC DAILY /TN "MyTasks\Fusion Update Task" /TR "C:\FusionUpdater\fusionUpdater.exe" /ST 14:30 /RU System /F /RL HIGHEST

start /wait C:\FusionUpdater\fusionUpdater.exe

```
