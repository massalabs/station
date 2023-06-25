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
    Should Be Equal As Strings    ${response[0]['name']}    Massa Node-manager
    Should Be Equal As Strings    ${response[0]['description']}    Join Massa network now. Automatically install, configure and manage Massa nodes.
