*** Settings ***
Documentation       Keywords specific to cmd test suite

Library             RequestsLibrary
Resource            ../keywords.resource
Resource            ../variables.resource


*** Keywords ***
Suite Setup
    Basic Suite Setup
    Delete Wallet Plugin
    Install Wallet Plugin

Delete Wallet Plugin
    Log To Console    Deleting Massa Wallet Plugin
    ${pluginId}=    Get Plugin ID From Author and Name    Massalabs    Massa${SPACE}Wallet
    IF    "${pluginId}" != "${EMPTY}"
        Log To Console    Deleting plugin with ID ${pluginId}
        ${response}=    DELETE    ${API_URL}/plugin-manager/${pluginId}
        Sleep    1 seconds    # Wait for the plugin to be deleted
    END

Install Wallet Plugin
    ${source}=    Set Variable
    ...    https://github.com/massalabs/thyra-plugin-wallet/releases/download/v0.0.4/thyra-plugin-wallet_${OS}-${ARCH}.zip
    ${response}=    POST
    ...    ${API_URL}/plugin-manager
    ...    params=source=${source}
    ...    expected_status=${STATUS_NO_CONTENT}
    Sleep    1 seconds    # Wait for the plugin to be registered
