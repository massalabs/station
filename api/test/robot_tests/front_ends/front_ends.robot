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
GET /
    Open Thyra Page    ${API_URL}/home/
    Page Should Contain    Which Plugins
    Page Should Contain    Registry
    Page Should Contain    Web On Chain
    Page Should Contain    Manage plugin

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

GET /store/
    Open Thyra Page    ${API_URL}/store/
    Page Should Contain    Plugin Manager
    Page Should Contain    Install a plugin
    Page Should Contain    Install a plugin using .zip URL

GET /store/index.html
    Open Thyra Page    ${API_URL}/store/index.html
    Page Should Contain    Plugin Manager
    Page Should Contain    Install a plugin
    Page Should Contain    Install a plugin using .zip URL

GET /search/
    Open Thyra Page    ${API_URL}/search/
    Page Should Contain    Registry
    Page Should Contain    Browse decentralized websites
    Page Should Contain    Website name
    Page Should Contain    Address
    Page Should Contain    URL

GET /search/index.html
    Open Thyra Page    ${API_URL}/search/index.html
    Page Should Contain    Registry
    Page Should Contain    Browse decentralized websites
    Page Should Contain    Website name
    Page Should Contain    Address
    Page Should Contain    URL

GET /websiteUploader/
    Open Thyra Page    ${API_URL}/websiteUploader/
    Page Should Contain    Decentralized website storage
    Page Should Contain    Upload a website
    Page Should Contain    On wallet
    Page Should Contain    Website Name
    Page Should Contain    Use alphanumerical characters and lowercase
    Page Should Contain    Website File
    Page Should Contain Button    website-upload
    Page Should Contain Button    file-select-button

GET /websiteUploader/index.html
    Open Thyra Page    ${API_URL}/websiteUploader/index.html
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

GET /store/{resource} with invalid resource
    ${response}=    GET    ${API_URL}/store/invalid    expected_status=${STATUS_NOT_FOUND}

GET /search/{resource} with invalid resource
    ${response}=    GET    ${API_URL}/search/invalid    expected_status=${STATUS_NOT_FOUND}

GET /websiteUploader/{resource} with invalid resource
    ${response}=    GET    ${API_URL}/websiteUploader/invalid    expected_status=${STATUS_NOT_FOUND}
