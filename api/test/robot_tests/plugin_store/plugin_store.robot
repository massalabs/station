*** Settings ***
Documentation       This is a test suite for Massa Station Plugin Store endpoints.

Library             RequestsLibrary
Resource            ../variables.resource
Resource            ../keywords.resource

Suite Setup         Basic Suite Setup


*** Test Cases ***
GET /plugin-store
    ${response}=    GET    ${API_URL}/plugin-store
    ${response}=    Set Variable    ${response.json()}
    Should Be Equal As Strings    ${response[0]['name']}    Node Manager
    Should Be Equal As Strings    ${response[0]['description']}    Join Massa network in a single click! Install, configure and manage Massa nodes.
    #TODO : Update link to after upgrading thyra-node-manager-plugin to station-massa-node-manager
    Should Be Equal As Strings    ${response[0]['url']}    https://github.com/massalabs/station-massa-node-manager
