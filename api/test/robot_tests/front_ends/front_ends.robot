*** Settings ***
Documentation       Test suite for Thyra front ends.
...                 This test suite mainly focuses on making sure the front end endpoints are working as expected.

Library             RequestsLibrary
Library             SeleniumLibrary
Resource            ../variables.resource

Suite Teardown      Close All Browsers

*** Test Cases ***
GET /thyra/home/{resource}
    Open Browser    ${API_URL}/thyra/home/
    Title Should Be    Thyra
    Page Should Contain    Which Plugins
    Page Should Contain    Registry
    Page Should Contain    Web On Chain
    Page Should Contain    Manage plugin

