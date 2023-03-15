*** Settings ***
Library         REST   http://localhost
Resource        variables.resource
Documentation   This is a test suite for Thyra Plugin Manager endpoints.

*** Test Cases ***
GET /plugin-manager
    GET             /plugin-manager
    Integer         response status            ${OK}
    Array           response body              minItems=1  maxItems=1  uniqueItems=true
    [Teardown]      Output   response body


GET /plugin-manager/invalid
    GET             /plugin-manager/invalid
    Integer         response status            ${NOT_FOUND}
    String          response body code         "Plugin-0001"
    [Teardown]      Output   response body

GET /thyra/plugin/invalid/invalid
    GET             /thyra/plugin/invalid/invalid
    Integer         response status            ${NOT_IMPLEMENTED}
    [Teardown]      Output   response body

Test get Operating System
    ${system}=    Evaluate    platform.system()    platform
    log to console    \nI am running on ${system}

Test get Architecture
    ${arch}=    Evaluate    platform.machine()    platform
    log to console    \nI am running on ${arch}
