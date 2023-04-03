*** Settings ***
Library     String


*** Keywords ***
Get Operating System
    ${OS}=    Evaluate    platform.system()    platform
    RETURN    ${OS}

Get Architecture
    ${ARCH}=    Evaluate    platform.machine()    platform
    RETURN    ${ARCH}

Basic Suite Setup
    [Documentation]    This is the default suite setup
    [Tags]    suite_setup

    ${OS}=    Get Operating System
    ${OS}=    Convert To Lower Case    ${OS}
    ${ARCH}=    Get Architecture

    IF    "${ARCH}" == "x86_64"
        ${ARCH}=    Set Variable    amd64
    END

    Set Global Variable    ${OS}
    Set Global Variable    ${ARCH}