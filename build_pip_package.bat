@echo off

REM Cleanup build folder
if exist tmp rd tmp /s /q
del "hipache_cli-*" /q

mkdir tmp
xcopy hipache_cli\* tmp /Y
xcopy hic tmp /Y
xcopy LICENSE tmp /Y
xcopy MANIFEST.in tmp /Y
xcopy README.rst tmp /Y
xcopy setup.py tmp /Y

pushd tmp
python setup.py sdist
popd

xcopy tmp\dist\* .

REM Cleanup build folder
if exist tmp rd tmp /s /q