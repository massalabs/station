*** Settings ***
Documentation       This is a test suite for Thyra Plugin Manager endpoints.

Library             RequestsLibrary
Library             String
Resource            keywords.resource
Resource            ../variables.resource

Suite Setup         Suite Setup


*** Test Cases ***
GET /plugin-manager with no plugins
    ${response}=    GET    ${API_URL}/plugin-manager
    Status Should Be    ${STATUS_OK}
    ${listLength}=    Get Length    ${response.json()}
    Should Be Equal As Integers    ${listLength}    0

POST /plugin-manager?source={{pluginSource}}
    ${source}=    Set Variable
    ...    https://github.com/massalabs/thyra-plugin-hello-world/releases/download/0.0.3/thyra-plugin-hello-world_${OS}-${ARCH}.zip
    ${response}=    POST
    ...    ${API_URL}/plugin-manager
    ...    params=source=${source}
    ...    expected_status=${STATUS_NO_CONTENT}
    Sleep    1 seconds    # Wait for the plugin to be registered

GET /plugin-manager with one plugin
    ${response}=    GET    ${API_URL}/plugin-manager
    Status Should Be    ${STATUS_OK}
    ${listLength}=    Get Length    ${response.json()}
    Should Be Equal As Integers    ${listLength}    1

GET /plugin-manager/{id}
    ${response}=    GET    ${API_URL}/plugin-manager
    ${id}=    Set Variable    ${response.json()[0]['id']}
    ${response}=    GET    ${API_URL}/plugin-manager/${id}
    Status Should Be    ${STATUS_OK}
    Should Be Equal As Strings    ${response.json()['status']}    Up

POST /plugin-manager/{id}/execute with stop command
    ${id}=    Get Plugin ID From Author and Name    massalabs    hello-world
    ${data}=    Create Dictionary    command=stop
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/${id}/execute
    ...    expected_status=${STATUS_NO_CONTENT}
    ...    json=${data}

POST /plugin-manager/{id}/execute with start command
    ${id}=    Get Plugin ID From Author and Name    massalabs    hello-world
    ${data}=    Create Dictionary    command=start
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/${id}/execute
    ...    expected_status=${STATUS_NO_CONTENT}
    ...    json=${data}
    Sleep    1 seconds    # Wait for the plugin to be started

POST /plugin-manager/{id}/execute with restart command
    ${id}=    Get Plugin ID From Author and Name    massalabs    hello-world
    ${data}=    Create Dictionary    command=restart
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/${id}/execute
    ...    expected_status=${STATUS_NO_CONTENT}
    ...    json=${data}
    Sleep    1 seconds    # Wait for the plugin to be restarted

GET /thyra/plugin/{author}/{name}/
    ${response}=    GET
    ...    ${API_URL}/thyra/plugin/massalabs/hello-world/
    ...    expected_status=${STATUS_OK}
    Should Contain    ${response.text}    Hello, world!

# We can register multiple times the same plugin, but the aliases list isn't updated correctly.
# This causes the previous plugin alias to be considered as still valid.
# TODO: Uncomment this test and make sure it passes once https://github.com/massalabs/thyra/issues/574 is fixed.
# POST /plugin-manager/register
#    ${id}=    Get Plugin ID From Author and Name    massalabs    hello-world
#    ${data}=    Create Dictionary
#    ...    id=${id}
#    ...    name=aliqua
#    ...    author=adipisicing
#    ...    description=minim consectetur dolore,
#    ...    logo=id et sunt irure,
#    ...    home=sunt
#    ...    api_spec=culpa enim sint aliqua
#    ...    url=oluptate
#    ${response}=    POST
#    ...    ${API_URL}/plugin-manager/register
#    ...    expected_status=${STATUS_NO_CONTENT}
#    ...    json=${data}
#    ${newid}=    Get Plugin ID From Author and Name    adipisicing    aliqua
#    Should Be Equal As Strings    ${newid}    ${id}

# Error cases

POST /plugins-manager/{id}/execute already started
    ${id}=    Get Plugin ID From Author and Name    massalabs    hello-world
    ${data}=    Create Dictionary    command=start
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/${id}/execute
    ...    expected_status=${STATUS_BAD_REQUEST}
    ...    json=${data}

    ${expectedError}=    Set Variable
    ...    "[start]${SPACE}${SPACE}(Error while starting plugin thyra-plugin-hello-world_${OS}-${ARCH}: plugin is not ready to start.\n). Current plugin status is Up."
    Should Be Equal As Strings    ${response.json()['code']}    Plugin-0030
    Should Be Equal As Strings    "${response.json()['message']}"    ${expectedError}

GET /plugin-manager/{id} with invalid id
    ${response}=    GET    ${API_URL}/plugin-manager/invalid    expected_status=${STATUS_NOT_FOUND}
    Should Be Equal As Strings    ${response.json()['code']}    Plugin-0001
    Should Be Equal As Strings    ${response.json()['message']}    get plugin error: no plugin matching correlationID invalid

GET /thyra/plugin/${author}/${name} with invalid author and name
    ${response}=    GET    ${API_URL}/thyra/plugin/invalid/invalid    expected_status=${STATUS_NOT_FOUND}
    Should Be Equal As Strings    ${response.text}    plugin not found for alias invalid/invalid

POST /plugin-manager/{id}/execute with invalid id
    ${data}=    Create Dictionary    command=start
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/invalid/execute
    ...    expected_status=${STATUS_NOT_FOUND}
    ...    json=${data}
    Should Be Equal As Strings    ${response.json()['code']}    Plugin-0001
    Should Be Equal As Strings    ${response.json()['message']}    get plugin error: no plugin matching correlationID invalid

DELETE /plugin-manager/{id} with invalid id
    ${response}=    DELETE    ${API_URL}/plugin-manager/3829029    expected_status=${STATUS_INTERNAL_SERVER_ERROR}
    Should Be Equal As Strings
    ...    ${response.json()['message']}
    ...    deleting plugin 3829029: no plugin matching correlationID 3829029

POST /plugin-manager/{id}/execute with invalid body
    ${data}=    Create Dictionary    command=test
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/invalid/execute
    ...    expected_status=${STATUS_UNPROCESSABLE_ENTITY}
    ...    json=${data}
    Should Be Equal As Strings    ${response.json()['code']}    606
    Should Be Equal As Strings
    ...    ${response.json()['message']}
    ...    body.command in body should be one of [update stop start restart]

POST /plugin-manager/register with invalid id
    ${data}=    Create Dictionary
    ...    id=1
    ...    name=ut aliqua non
    ...    author=adipisicing
    ...    description=minim consectetur dolore,
    ...    logo=id et sunt irure,
    ...    home=sunt
    ...    api_spec=culpa enim sint aliqua
    ...    url=oluptate
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/register
    ...    expected_status=${STATUS_NOT_FOUND}
    ...    json=${data}

    Should Be Equal As Strings    ${response.json()['code']}    Plugin-0001
    Should Be Equal As Strings    ${response.json()['message']}    get plugin error: no plugin matching correlationID 1

POST /plugin-manager/register with invalid body
    ${data}=    Create Dictionary
    ...    id=-65217
    ...    name=ut aliqua non
    ...    author=adipisicing
    ...    description=minim consectetur dolore,
    ...    logo=id et sunt irure,
    ...    home=sunt
    ...    api_spec=culpa enim sint aliqua
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/register
    ...    expected_status=${STATUS_UNPROCESSABLE_ENTITY}
    ...    json=${data}
    Should Be Equal As Strings    ${response.json()['code']}    602
    Should Be Equal As Strings    ${response.json()['message']}    body.url in body is required

POST /plugin-manager/{id}/execute with NotImplemented update command
    ${id}=    Get Plugin ID From Author and Name    massalabs    hello-world
    ${data}=    Create Dictionary    command=update
    ${response}=    POST
    ...    ${API_URL}/plugin-manager/${id}/execute
    ...    expected_status=${STATUS_NOT_IMPLEMENTED}
    ...    json=${data}
