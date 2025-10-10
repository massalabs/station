*** Settings ***
Documentation       Keywords specific to cmd test suite

Library             RequestsLibrary
Library             OperatingSystem
Resource            ../keywords.resource
Resource            ../variables.resource

*** Variables ***
${WALLET_PLUGIN_VERSION}       v0.5.0

*** Keywords ***
Suite Setup
    Basic Suite Setup
    Switch To Buildnet
    Init SC build environment
    ${GITHUB_ACTIONS}    Get Environment Variable    GITHUB_ACTIONS    false
    IF    "${GITHUB_ACTIONS}" == "true"
        Log To Console    We are in the CI
        Delete Wallet Plugin
        Install Wallet Plugin
    END

Delete Wallet Plugin
    Log To Console    Deleting Massa Wallet Plugin
    ${pluginId}=    Get Plugin ID From Author and Name    massalabs    wallet
    IF    "${pluginId}" != "${EMPTY}"
        Log To Console    Deleting plugin with ID ${pluginId}
        ${headers}=    Create Dictionary    Origin=http://localhost
        ${response}=    DELETE    ${API_URL}/plugin-manager/${pluginId}    headers=${headers}
        Sleep    1 seconds    # Wait for the plugin to be deleted
    END

Install Wallet Plugin
    Log To Console    Installing Massa Wallet Plugin
    ${source}=    Set Variable
    ...    https://massa-station-assets.s3.eu-west-3.amazonaws.com/plugins/wallet/${WALLET_PLUGIN_VERSION}/wallet-plugin_${OS}-${ARCH}.zip
    ${headers}=    Create Dictionary    Origin=http://localhost
    ${response}=    POST
    ...    ${API_URL}/plugin-manager
    ...    params=source=${source}
    ...    headers=${headers}
    ...    expected_status=${STATUS_NO_CONTENT}
    Sleep    1 seconds    # Wait for the plugin to be registered

Init SC build environment
    Log To Console    Initializing SC build environment
    Log To Console    Installing SC dependencies
    Run    cd testSC; npm install
    Log To Console    Building SC
    Run    cd testSC; npm run build
