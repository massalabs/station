*** Settings ***
Library     RequestsLibrary
Resource    ../keywords.resource
Resource    ../variables.resource


*** Keywords ***
Suite Setup
    [Documentation]    Suite Setup specific to this test suite
    Basic Suite Setup
    # To remove plugins that may have been installed at station's first run setup
    Delete all plugins


Delete all plugins
    Log To Console    Deleting all plugins
    ${response}=    GET    ${API_URL}/plugin-manager
    FOR    ${element}    IN    @{response.json()}
        ${headers}=    Create Dictionary    Origin=http://localhost
        ${response}=    DELETE    ${API_URL}/plugin-manager/${element['id']}    headers=${headers}
        Status Should Be    ${STATUS_NO_CONTENT}
    END

# Unused for now
Delete Hello World Plugin
    Log To Console    Deleting Hello World Plugin
    ${pluginId}=    Get Plugin ID From Author and Name    massalabs    hello-world
    IF    "${pluginId}" != "${EMPTY}"
        ${headers}=    Create Dictionary    Origin=http://localhost
        ${response}=    DELETE    ${API_URL}/plugin-manager/${pluginId}    headers=${headers}
    END
