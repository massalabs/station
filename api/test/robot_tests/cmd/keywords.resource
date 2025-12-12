*** Settings ***
Documentation       Keywords specific to cmd test suite

Library             RequestsLibrary
Library             OperatingSystem
Resource            ../keywords.resource
Resource            ../variables.resource

*** Variables ***
${WALLET_PLUGIN_VERSION}       v0.5.6

*** Keywords ***
Suite Setup
    Basic Suite Setup
    Switch To Buildnet
    Init SC build environment
    Install Wallet Plugin

Delete Wallet Plugin
    Log To Console    Attempt to delete Massa Wallet Plugin if exist
    ${pluginId}=    Get Plugin ID From Author and Name    massalabs    wallet
    IF    "${pluginId}" != "${EMPTY}"
        Log To Console    Deleting plugin Massa Wallet with ID ${pluginId}
        ${headers}=    Create Dictionary    Origin=http://localhost
        ${response}=    DELETE    ${API_URL}/plugin-manager/${pluginId}    headers=${headers}
        Sleep    1 seconds    # Wait for the plugin to be deleted
    END

Install Wallet Plugin
    Log To Console    Attempt to install Massa Wallet Plugin if not already installed
    ${pluginId}=    Get Plugin ID From Author and Name    massalabs    wallet
    IF      "${pluginId}" == "${EMPTY}"
        Log To Console    Installing Massa Wallet Plugin because it does not exist yet
        ${source}=    Set Variable
        ...    https://massa-station-assets.s3.eu-west-3.amazonaws.com/plugins/wallet/${WALLET_PLUGIN_VERSION}/wallet-plugin_${OS}-${ARCH}.zip
        ${headers}=    Create Dictionary    Origin=http://localhost
        ${response}=    POST
        ...    ${API_URL}/plugin-manager
        ...    params=source=${source}
        ...    headers=${headers}
        ...    expected_status=${STATUS_NO_CONTENT}
        Sleep    1 seconds    # Wait for the plugin to be registered
    END

Init SC build environment
    Log To Console    Initializing SC build environment
    Run    cd testSC; npm install
    Run    cd testSC; npm run build
