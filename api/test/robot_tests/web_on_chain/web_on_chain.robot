*** Settings ***
Documentation       A test suite for every endpoints related to the Web on Chain.

Library             RequestsLibrary
Library             SeleniumLibrary
Resource            variables.resource
Resource            ../variables.resource
Resource            keywords.resource
Resource            ../front_ends/variables.resource
Resource            ../front_ends/keywords.resource

Suite Setup         Suite Setup
Suite Teardown      Close All Browsers

*** Test Cases ***
PUT /websiteUploader/prepare
    [Documentation]    Prepare a website for creation.

    ${zip}=    Get File For Streaming Upload    ${CURDIR}/${TEST_ZIP_FILE}
    ${data}=    Create Dictionary    url=${website_url}    description=${website_description}    nickname=${WALLET_NICKNAME}
    ${file}=    Create Dictionary    zipfile=${zip}
    ${response}=    PUT
    ...    ${API_URL}/websiteUploader/prepare
    ...    data=${data}
    ...    files=${file}
    ...    expected_status=any

    Log To Console    json response: ${response.json()}    # Print the response content to the test log for debugging

    Should Be Equal As Integers    ${response.status_code}    ${STATUS_OK}    # Assert the status code is 200 OK
    Should Be Equal As Strings    ${response.json()['name']}    ${website_url}
    Set Global Variable    ${WEBSITE_NAME}    ${response.json()['name']}
    Set Global Variable    ${WEBSITE_ADDRESS}    ${response.json()['address']}

POST /websiteUploader/upload
    [Documentation]    Upload the content of the website to the blockchain.

    ${zip}=    Get File For Streaming Upload    ${CURDIR}/${TEST_ZIP_FILE}
    ${data}=    Create Dictionary    nickname=${WALLET_NICKNAME}    address=${WEBSITE_ADDRESS}
    ${file}=    Create Dictionary    zipfile=${zip}
    ${response}=    POST
    ...    ${API_URL}/websiteUploader/upload
    ...    data=${data}
    ...    files=${file}
    ...    expected_status=${STATUS_OK}

    Should Be Equal As Strings    ${response.json()['address']}    ${WEBSITE_ADDRESS}

POST /cmd/executeFunction To remove the dns entry
    ${argument}=    keywords.String To Arg    ${website_url}
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=deleteEntryFromDNS
    ...    at=AS12tsy5TtDyGsWggDvqDLagxZqkuTiTPqWhqmrx7KSQNJwQSbnBA
    ...    args=${argument}
    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_OK}
    
GET /all/domains 
    ${response}=    GET
    ...    ${API_URL}/all/domains
    ...    expected_status=${STATUS_OK}
    Should Contain    ${response.text}    flappy

GET /my/domains/{nickname}
    ${response}=    GET
    ...    ${API_URL}/my/domains/${WALLET_NICKNAME}
    ...    expected_status=${STATUS_OK}
    Should Contain    ${response.text}    ${website_url}
