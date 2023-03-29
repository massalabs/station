*** Settings ***
Library     RequestsLibrary
Resource    ../keywords.resource
Resource    ../variables.resource


*** Keywords ***
Suite Setup
    [Documentation]    Suite Setup specific to this test suite
    Basic Suite Setup
    Delete all plugins

Delete all plugins
    Log To Console    Deleting all plugins
    ${response}=    GET    ${API_URL}/plugin-manager
    FOR    ${element}    IN    @{response.json()}
        ${response}=    DELETE    ${API_URL}/plugin-manager/${element['id']}
        Status Should Be    ${STATUS_NO_CONTENT}
    END

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
