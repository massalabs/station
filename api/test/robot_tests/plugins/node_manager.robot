*** Settings ***
Documentation       Test node manager plugin

Library             RequestsLibrary
Resource            ../keywords.resource
Resource            ../variables.resource

Suite Setup         Basic Suite Setup

*** Variables ***
${NODE_MANAGER_PLUGIN_VERSION}       v0.4.3

*** Test Cases ***
Install node manager plugin
    Log To Console    Installing Node Manager Plugin
    ${source}=    Set Variable
    ...    https://massa-station-assets.s3.eu-west-3.amazonaws.com/plugins/node-manager/${NODE_MANAGER_PLUGIN_VERSION}/node-manager-plugin_${OS}-${ARCH}.zip
    ${headers}=    Create Dictionary    Origin=http://localhost
    ${response}=    POST
    ...    ${API_URL}/plugin-manager
    ...    params=source=${source}
    ...    headers=${headers}
    ...    expected_status=${STATUS_NO_CONTENT}
    Sleep    1 seconds    # Wait for the plugin to be registered
    ${pluginId}=    Get Plugin ID From Author and Name    massa-labs    node-manager
    Should Not Be Equal    ${pluginId}    ${EMPTY}    Plugin ID should not be empty

Close and restart node manager plugin
    Log To Console      Close node manager plugin
    ${id}=    Get Plugin ID From Author and Name    massa-labs    node-manager
    ${data}=    Create Dictionary    command=stop
    ${headers}=    Create Dictionary    Origin=http://localhost
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/${id}/execute
    ...    headers=${headers}
    ...    expected_status=${STATUS_NO_CONTENT}
    ...    json=${data}
    ${data}=    Create Dictionary    command=start
    ${headers}=    Create Dictionary    Origin=http://localhost
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/${id}/execute
    ...    headers=${headers}
    ...    expected_status=${STATUS_NO_CONTENT}
    ...    json=${data}
    Sleep    1 seconds

Uninstall node manager plugin
    ${pluginId}=    Get Plugin ID From Author and Name    massa-labs    node-manager
    Log To Console      Uninstall node manager plugin with ID ${pluginId}
    ${headers}=    Create Dictionary    Origin=http://localhost
    ${response}=    DELETE    ${API_URL}/plugin-manager/${pluginId}    headers=${headers}
    Sleep    1 seconds    # Wait for the plugin to be deleted
    ${pluginId}=    Get Plugin ID From Author and Name    massa-labs    node-manager
    Should Be Equal    ${pluginId}    ${EMPTY}    Plugin ID should be empty

