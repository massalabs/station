*** Settings ***
Documentation       A test suite for every endpoints related to the Web on Chain.

Library             RequestsLibrary
Library             SeleniumLibrary
Resource            keywords.resource
Resource            variables.resource
Resource            ../keywords.resource
Resource            ../variables.resource
Resource            ../front_ends/variables.resource
Resource            ../front_ends/keywords.resource

Suite Setup         Suite Setup
Suite Teardown      Close All Browsers


*** Test Cases ***
PUT /websiteCreator/prepare
    [Documentation]    Prepare a website for creation.

    ${website_url}=    Generate Random String    8    [LOWER][NUMBERS]
    ${website_url}=    Set Variable    test${website_url}
    ${zip}=    Get File For Streaming Upload    ${CURDIR}/${TEST_ZIP_FILE}
    ${data}=    Create Dictionary    nickname=${WALLET_NICKNAME}    url=${website_url}
    ${file}=    Create Dictionary    zipfile=${zip}
    ${response}=    PUT
    ...    ${API_URL}/websiteCreator/prepare
    ...    data=${data}
    ...    files=${file}
    ...    expected_status=500

    Should Be Equal As Strings    ${response.json()['name']}    ${website_url}
    Set Global Variable    ${WEBSITE_NAME}    ${response.json()['name']}
    Set Global Variable    ${WEBSITE_ADDRESS}    ${response.json()['address']}

POST /websiteCreator/upload
    [Documentation]    Upload the content of the website to the blockchain.

    ${zip}=    Get File For Streaming Upload    ${CURDIR}/${TEST_ZIP_FILE}
    ${data}=    Create Dictionary    nickname=${WALLET_NICKNAME}    address=${WEBSITE_ADDRESS}
    ${file}=    Create Dictionary    zipfile=${zip}
    ${response}=    POST
    ...    ${API_URL}/websiteCreator/upload
    ...    data=${data}
    ...    files=${file}
    ...    expected_status=${STATUS_OK}

    Should Be Equal As Strings    ${response.json()['address']}    ${WEBSITE_ADDRESS}

Check content of the uploaded website
    Open Thyra Page    ${API_URL}/browse/${WEBSITE_ADDRESS}/index.html
    Page Should Contain    My test website!
    Page Should Contain    Decentralization is non-negotiable
