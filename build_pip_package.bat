@echo off

REM Cleanup build folder
if exist tmp rd tmp /s /q
del "hic-*" /q

mkdir tmp
xcopy *.py tmp /Y
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