*** Settings ***
Documentation       Test suite for Thyra front ends.
...                 This test suite mainly focuses on making sure the front end endpoints are working as expected.

Library             RequestsLibrary
Library             SeleniumLibrary
Resource            ../variables.resource
Resource            keywords.resource
Resource            variables.resource

Suite Teardown      Close All Browsers


*** Test Cases ***
GET /home/
    Open Thyra Page    ${API_URL}/home/
    Page Should Contain    Which Plugins
    Page Should Contain    Registry
    Page Should Contain    Web On Chain
    Page Should Contain    Manage plugin

GET /home/index.html
    Open Thyra Page    ${API_URL}/home/index.html
    Page Should Contain    Which Plugins
    Page Should Contain    Registry
    Page Should Contain    Web On Chain
    Page Should Contain    Manage plugin

GET /thyra/plugin-manager/
    Open Thyra Page    ${API_URL}/thyra/plugin-manager/
    Page Should Contain    Plugin Manager
    Page Should Contain    Install a plugin
    Page Should Contain    Install a plugin using .zip URL

GET /thyra/plugin-manager/index.html
    Open Thyra Page    ${API_URL}/thyra/plugin-manager/index.html
    Page Should Contain    Plugin Manager
    Page Should Contain    Install a plugin
    Page Should Contain    Install a plugin using .zip URL

GET /thyra/registry/
    Open Thyra Page    ${API_URL}/thyra/registry/
    Page Should Contain    Registry
    Page Should Contain    Browse decentralized websites
    Page Should Contain    Website name
    Page Should Contain    Address
    Page Should Contain    URL

GET /thyra/registry/index.html
    Open Thyra Page    ${API_URL}/thyra/registry/index.html
    Page Should Contain    Registry
    Page Should Contain    Browse decentralized websites
    Page Should Contain    Website name
    Page Should Contain    Address
    Page Should Contain    URL

GET /thyra/websiteCreator/
    Open Thyra Page    ${API_URL}/thyra/websiteCreator/
    Page Should Contain    Decentralized website storage
    Page Should Contain    Upload a website
    Page Should Contain    On wallet
    Page Should Contain    Website Name
    Page Should Contain    Use alphanumerical characters and lowercase
    Page Should Contain    Website File
    Page Should Contain Button    website-upload
    Page Should Contain Button    file-select-button

GET /thyra/websiteCreator/index.html
    Open Thyra Page    ${API_URL}/thyra/websiteCreator/index.html
    Page Should Contain    Decentralized website storage
    Page Should Contain    Upload a website
    Page Should Contain    On wallet
    Page Should Contain    Website Name
    Page Should Contain    Use alphanumerical characters and lowercase
    Page Should Contain    Website File
    Page Should Contain Button    website-upload
    Page Should Contain Button    file-select-button

# Error cases

GET /home/{resource} with invalid resource
    ${response}=    GET    ${API_URL}/home/invalid    expected_status=${STATUS_NOT_FOUND}

GET /thyra/plugin-manager/{resource} with invalid resource
    ${response}=    GET    ${API_URL}/thyra/plugin-manager/invalid    expected_status=${STATUS_NOT_FOUND}

GET /thyra/registry/{resource} with invalid resource
    ${response}=    GET    ${API_URL}/thyra/registry/invalid    expected_status=${STATUS_NOT_FOUND}

GET /thyra/websiteCreator/{resource} with invalid resource
    ${response}=    GET    ${API_URL}/thyra/websiteCreator/invalid    expected_status=${STATUS_NOT_FOUND}
