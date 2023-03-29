*** Settings ***
Library     RequestsLibrary
Resource    ../keywords.resource
Resource    ../variables.resource


*** Keywords ***
Suite Setup
    [Documentation]    Suite Setup specific to this test suite
    Basic Suite Setup
    Delete Hello World Plugin

# Unused for now
Delete all plugins
    Log To Console    Deleting all plugins
    ${response}=    GET    ${API_URL}/plugin-manager
    FOR    ${element}    IN    @{response.json()}
        ${response}=    DELETE    ${API_URL}/plugin-manager/${element['id']}
        Status Should Be    ${STATUS_NO_CONTENT}
    END

Delete Hello World Plugin
    Log To Console    Deleting Hello World Plugin
    ${pluginId}=    Get Plugin ID From Author and Name    massalabs    hello world
    ${response}=    DELETE    ${API_URL}/plugin-manager/${pluginId}

Get Plugin ID From Author and Name
    [Arguments]    ${author}    ${name}
    ${response}=    GET    ${API_URL}/plugin-manager

    ${pluginId}=    Set Variable    ${EMPTY}
    ${expectedURL}=    Set Variable    /thyra/plugin/${author}/${name}/

    FOR    ${element}    IN    @{response.json()}
        ${pluginURL}=    Get Regexp Matches    ${element['home']}    (\/.*){3}\/
        ${pluginURL}=    Evaluate    urllib.parse.unquote("${pluginURL[0]}")    modules=urllib.parse

        IF    "${expectedURL}" == "${pluginURL}"
            ${pluginId}=    Set Variable    ${element['id']}
            BREAK
        END
    END
    RETURN    ${pluginId}
