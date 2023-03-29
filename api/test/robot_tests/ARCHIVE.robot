*** Settings ***
Documentation       This is an archive that shows how to get the operating system and architecture
...                 of the machine you are running on.
...
...                 TO BE DELETED


*** Test Cases ***
Test get Operating System
    ${system}=    Evaluate    platform.system()    platform
    log to console    \nI am running on ${system}

Test get Architecture
    ${arch}=    Evaluate    platform.machine()    platform
    log to console    \nI am running on ${arch}
