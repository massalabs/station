*** Settings ***
Library     String
Library     RequestsLibrary
Library     ./libs/ArgLibrary.py
Library     Collections
Resource    variables.resource


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

Get Plugin ID From Author and Name
    [Arguments]    ${author}    ${name}
    ${response}=    GET    ${API_URL}/plugin-manager

    ${pluginId}=    Set Variable    ${EMPTY}
    ${expectedURL}=    Set Variable    /plugin/${author}/${name}/

    FOR    ${element}    IN    @{response.json()}
        ${pluginURL}=    Get Regexp Matches    ${element['home']}    (\/.*){3}\/
        ${pluginURL}=    Evaluate    urllib.parse.unquote("${pluginURL[0]}")    modules=urllib.parse

        IF    "${expectedURL}" == "${pluginURL}"
            ${pluginId}=    Set Variable    ${element['id']}
            BREAK
        END
    END
    RETURN    ${pluginId}

String To Arg
    [Documentation]    Converts a string into an Arg object
    ...    It returns an U8 array composed of the length of the string
    ...    followed by the string itself
    [Arguments]    ${string}
    ${byteArray}=    ArgLibrary.String To Arg    ${string}
    RETURN    ${byteArray}


Switch To Buildnet
    ${headers}=    Create Dictionary    Origin=https://station.massa
    ${response}=    POST
    ...    ${API_URL}/network/buildnet
    ...    headers=${headers}
    ...    expected_status=${STATUS_OK}
    Should Be Equal    ${response.json()['currentNetwork']}    buildnet
    Log To Console    Switched to buildnet
