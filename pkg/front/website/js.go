// GENERATED BY textFileToGoConst
// GitHub:     github.com/logrusorgru/textFileToGoConst
// input file: html\front\website.js
// generated:  Fri Sep 30 17:40:04 CEST 2022

package website

const JS = "let gWallets = [];\r\nlet deployers = [];\r\nlet actualTxType = \"\";\r\nlet nextFileToUpload;\r\nlet uploadable = false;\r\n\r\n// INIT\r\ngetWallets();\r\ngetWebsiteDeployerSC();\r\ninitializeDefaultWallet();\r\nsetupModal();\r\n\r\nconst eventManager = new EventManager();\r\n\r\n// Write the default wallet text in wallet popover component\r\nasync function getWebsiteDeployerSC() {\r\n  let defaultWallet = getDefaultWallet();\r\n\r\n  $(\"#website-deployers-table tbody tr\").remove();\r\n\r\n  axios\r\n    .get(\"/my/domains/\" + defaultWallet)\r\n    .then((websites) => {\r\n      let count = 0;\r\n      for (const website of websites.data) {\r\n        tableInsert(website, count);\r\n        count++;\r\n      }\r\n      deployers = websites.data;\r\n    })\r\n    .catch((e) => {\r\n      console.error(e);\r\n      errorAlert(getErrorMessage(e.response.data.code));\r\n    });\r\n}\r\n\r\n// Write the default wallet text in wallet popover component\r\nfunction initializeDefaultWallet() {\r\n  let defaultWallet = getDefaultWallet();\r\n  if (defaultWallet === \"\") {\r\n    defaultWallet = \"Connect\";\r\n  }\r\n  $(\".popover__title\").html(defaultWallet);\r\n}\r\n\r\n// Retrieve the default wallet nickname in cookies\r\nfunction getDefaultWallet() {\r\n  let defaultWallet = \"\";\r\n  const cookies = document.cookie.split(\";\");\r\n  cookies.forEach((cookie) => {\r\n    const keyValue = cookie.split(\"=\");\r\n    if (keyValue[0] === \"defaultWallet\") {\r\n      defaultWallet = keyValue[1];\r\n    }\r\n  });\r\n  return defaultWallet;\r\n}\r\n\r\nfunction getWallet(nickname) {\r\n  return gWallets.find((w) => w.nickname === nickname);\r\n}\r\n\r\nfunction getDeployerByDns(dns) {\r\n  return deployers.find((c) => c.name === dns);\r\n}\r\n\r\nfunction setupModal() {\r\n  $(\"#passwordModal\").on(\"shown.bs.modal\", function () {\r\n    $(\"#passwordModal\").trigger(\"focus\");\r\n  });\r\n}\r\n\r\nfunction setTxType(txType) {\r\n  actualTxType = txType;\r\n}\r\n\r\nasync function callTx() {\r\n  const passwordValue = $(\"#walletPassword\").val();\r\n\r\n  if (actualTxType === \"deployWebsiteAndUpload\") {\r\n    deployWebsiteAndUpload(passwordValue);\r\n  }\r\n  if (actualTxType.includes(\"uploadWebsiteCreator\")) {\r\n    const websiteIndex = actualTxType.split(\"uploadWebsiteCreator\")[1];\r\n    uploadWebsite(nextFileToUpload, websiteIndex, passwordValue);\r\n  }\r\n}\r\n\r\n// open file upload\r\nfunction openDialog() {\r\n  document.getElementById(\"fileid0\").value = null;\r\n  document.getElementById(\"fileid0\").click();\r\n}\r\n\r\n// Handle event on file selecting\r\nfunction handleFileSelect(evt) {\r\n  let files = evt.target.files; // get files\r\n  let f = files[0];\r\n  const reader = new FileReader();\r\n  reader.onload = (event) => importWallet(JSON.parse(event.target.result)); // desired file content\r\n  reader.onerror = (error) => reject(error);\r\n  reader.readAsText(f);\r\n}\r\n\r\n// Append wallet accounts in popover component list\r\nasync function feedWallet(w) {\r\n  let counter = 0;\r\n  for (const wallet of w) {\r\n    $(\"#wallet-list\").append(\r\n      \"<li class='wallet-item'><a class='wallet-link' id='wallet-link-\" +\r\n        counter +\r\n        \"' onclick='changeDefaultWallet(event)' href='#'>\" +\r\n        wallet.nickname +\r\n        \"</a></li>\"\r\n    );\r\n    counter++;\r\n  }\r\n}\r\n\r\n// Handle popover click & update default wallet in cookies\r\nfunction changeDefaultWallet(event) {\r\n  const idElementClicked = event.target.id;\r\n  const newDefaultWalletId = idElementClicked.split(\"-\")[2];\r\n  const walletName = gWallets[newDefaultWalletId].nickname;\r\n\r\n  document.cookie = \"defaultWallet=\" + walletName;\r\n  $(\".popover__title\").html(walletName);\r\n\r\n  getWebsiteDeployerSC();\r\n}\r\n\r\nasync function getWallets() {\r\n  axios\r\n    .get(\"/mgmt/wallet\")\r\n    .then((resp) => {\r\n      if (resp) {\r\n        gWallets = resp.data;\r\n        feedWallet(gWallets);\r\n      }\r\n    })\r\n    .catch((e) => {\r\n      errorAlert(getErrorMessage(e.response.data.code));\r\n    });\r\n}\r\n\r\nfunction tableInsert(resp, count) {\r\n  const tBody = document\r\n    .getElementById(\"website-deployers-table\")\r\n    .getElementsByTagName(\"tbody\")[0];\r\n  const row = tBody.insertRow(-1);\r\n  const url = \"http://\" + resp.name + \".massa/\";\r\n\r\n  const cell0 = row.insertCell();\r\n  const cell1 = row.insertCell();\r\n  const cell2 = row.insertCell();\r\n  const cell3 = row.insertCell();\r\n\r\n  cell0.innerHTML = resp.name;\r\n  cell1.innerHTML = addressInnerHTML(resp.address);\r\n  cell2.innerHTML = `<a href=\"${url}\" target=\"_blank\" rel=\"noopener noreferrer\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" fill=\"currentColor\" class=\"bi bi-globe2\" viewBox=\"0 0 16 16\"><path d=\"M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm7.5-6.923c-.67.204-1.335.82-1.887 1.855-.143.268-.276.56-.395.872.705.157 1.472.257 2.282.287V1.077zM4.249 3.539c.142-.384.304-.744.481-1.078a6.7 6.7 0 0 1 .597-.933A7.01 7.01 0 0 0 3.051 3.05c.362.184.763.349 1.198.49zM3.509 7.5c.036-1.07.188-2.087.436-3.008a9.124 9.124 0 0 1-1.565-.667A6.964 6.964 0 0 0 1.018 7.5h2.49zm1.4-2.741a12.344 12.344 0 0 0-.4 2.741H7.5V5.091c-.91-.03-1.783-.145-2.591-.332zM8.5 5.09V7.5h2.99a12.342 12.342 0 0 0-.399-2.741c-.808.187-1.681.301-2.591.332zM4.51 8.5c.035.987.176 1.914.399 2.741A13.612 13.612 0 0 1 7.5 10.91V8.5H4.51zm3.99 0v2.409c.91.03 1.783.145 2.591.332.223-.827.364-1.754.4-2.741H8.5zm-3.282 3.696c.12.312.252.604.395.872.552 1.035 1.218 1.65 1.887 1.855V11.91c-.81.03-1.577.13-2.282.287zm.11 2.276a6.696 6.696 0 0 1-.598-.933 8.853 8.853 0 0 1-.481-1.079 8.38 8.38 0 0 0-1.198.49 7.01 7.01 0 0 0 2.276 1.522zm-1.383-2.964A13.36 13.36 0 0 1 3.508 8.5h-2.49a6.963 6.963 0 0 0 1.362 3.675c.47-.258.995-.482 1.565-.667zm6.728 2.964a7.009 7.009 0 0 0 2.275-1.521 8.376 8.376 0 0 0-1.197-.49 8.853 8.853 0 0 1-.481 1.078 6.688 6.688 0 0 1-.597.933zM8.5 11.909v3.014c.67-.204 1.335-.82 1.887-1.855.143-.268.276-.56.395-.872A12.63 12.63 0 0 0 8.5 11.91zm3.555-.401c.57.185 1.095.409 1.565.667A6.963 6.963 0 0 0 14.982 8.5h-2.49a13.36 13.36 0 0 1-.437 3.008zM14.982 7.5a6.963 6.963 0 0 0-1.362-3.675c-.47.258-.995.482-1.565.667.248.92.4 1.938.437 3.008h2.49zM11.27 2.461c.177.334.339.694.482 1.078a8.368 8.368 0 0 0 1.196-.49 7.01 7.01 0 0 0-2.275-1.52c.218.283.418.597.597.932zm-.488 1.343a7.765 7.765 0 0 0-.395-.872C9.835 1.897 9.17 1.282 8.5 1.077V4.09c.81-.03 1.577-.13 2.282-.287z\"/></svg></a>`;\r\n  cell3.innerHTML =\r\n    \"<div><input id='fileid\" +\r\n    count +\r\n    \"' type='file' hidden/><button id='upload-website\" +\r\n    count +\r\n    \"' class='primary-button' type='button' value='Upload MB' >Edit</button><img src='./logo.png' style='display: none' class='massa-logo-spinner loading\" +\r\n    count +\r\n    \" alt='Massa logo' /></span></div>\";\r\n\r\n  document\r\n    .getElementById(`upload-website${count}`)\r\n    .addEventListener(\"click\", function () {\r\n      document.getElementById(`fileid${count}`).value = null;\r\n      document.getElementById(`fileid${count}`).click();\r\n    });\r\n\r\n  document\r\n    .getElementById(`fileid${count}`)\r\n    .addEventListener(\"change\", function (evt) {\r\n      let files = evt.target.files; // get files\r\n      nextFileToUpload = files[0];\r\n\r\n      setTxType(\"uploadWebsiteCreator\" + count);\r\n      $(\"#passwordModal\").modal(\"show\");\r\n    });\r\n}\r\n\r\n$(\"#file-select-button\").click(function () {\r\n  $(\".upload input\").click();\r\n});\r\n\r\n// change button text with file name\r\n$(\".upload input\").on(\"change\", function () {\r\n  let str = $(\".upload input\").val();\r\n\r\n  let n = str.lastIndexOf(\"\\\\\");\r\n\r\n  let result = str.substring(n + 1);\r\n\r\n  $(\"#file-select-button\").html(result);\r\n});\r\n\r\n//check if file is .zip\r\n$(\".upload input\").on(\"change\", function () {\r\n  let str = $(\".upload input\").val();\r\n\r\n  let n = str.lastIndexOf(\".\");\r\n\r\n  let result = str.substring(n + 1);\r\n\r\n  if (result != \"zip\" && $(\".upload input\").val() != \"\") {\r\n    uploadable = false;\r\n\r\n    document.getElementsByClassName(\"fileTypeError\")[0].style.display = \"flex\";\r\n    document.getElementById(\"website-upload\").style.display = \"none\";\r\n    document.getElementById(\"website-upload-refuse\").style.display = \"flex\";\r\n  } else {\r\n    uploadable = true;\r\n    document.getElementsByClassName(\"fileTypeError\")[0].style.display = \"none\";\r\n    document.getElementById(\"website-upload\").style.display = \"flex\";\r\n    document.getElementById(\"website-upload-refuse\").style.display = \"none\";\r\n  }\r\n});\r\n\r\n//check max size file\r\n$(\".upload input\").on(\"change\", function () {\r\n  const fileSize = this.files[0].size / 1024 / 1024; // in MiB\r\n  if (fileSize > 1.5) {\r\n    uploadable = false;\r\n    console.log(fileSize);\r\n    document.getElementsByClassName(\"fileSizeError\")[0].style.display = \"flex\";\r\n    document.getElementById(\"website-upload\").style.display = \"none\";\r\n    document.getElementById(\"website-upload-refuse\").style.display = \"flex\";\r\n  } else {\r\n    uploadable = true;\r\n    document.getElementsByClassName(\"fileSizeError\")[0].style.display = \"none\";\r\n    document.getElementById(\"website-upload\").style.display = \"flex\";\r\n    document.getElementById(\"website-upload-refuse\").style.display = \"none\";\r\n  }\r\n});\r\n\r\n//remove label of input website name on focus\r\n$(\".website-dns input\").on(\"focus\", function () {\r\n  document.getElementById(\"website-info-display\").style.visibility = \"hidden\";\r\n});\r\n\r\n//check if input string is valid\r\n$(\".website-dns input\").on(\"change\", function () {\r\n  let str = $(\".website-dns input\").val();\r\n  let pattern = new RegExp(\"^[a-z0-9]+$\");\r\n  let testPattern = pattern.test(str);\r\n\r\n  if (testPattern == false) {\r\n    uploadable = false;\r\n    document.getElementsByClassName(\"dns-error\")[0].style.display = \"flex\";\r\n    document.getElementById(\"website-upload\").style.display = \"none\";\r\n    document.getElementById(\"website-upload-refuse\").style.display = \"flex\";\r\n  } else {\r\n    uploadable = true;\r\n    document.getElementsByClassName(\"dns-error\")[0].style.display = \"none\";\r\n    document.getElementById(\"website-upload\").style.display = \"flex\";\r\n    document.getElementById(\"website-upload-refuse\").style.display = \"none\";\r\n  }\r\n});\r\n\r\nfunction uploadProcess(file, dnsName, isFullProcess, bodyFormData, callback) {\r\n  document.getElementById(\"wallet-popover\").classList.add(\"popover__disabled\");\r\n  const reader = new FileReader();\r\n  reader.readAsDataURL(file);\r\n  reader.onloadend = (_) => {\r\n    const result = reader.result.length;\r\n\r\n    const chunkSize = Math.floor(result / 260_000) + 1;\r\n\r\n    stepper(dnsName, chunkSize, isFullProcess);\r\n\r\n    callback(bodyFormData);\r\n  };\r\n}\r\n\r\nfunction postUpload(bodyFormData, password) {\r\n  axios({\r\n    url: `/websiteCreator/upload`,\r\n    method: \"POST\",\r\n    data: bodyFormData,\r\n    headers: {\r\n      \"Content-Type\": \"multipart/form-data\",\r\n      Authorization: password,\r\n    },\r\n  })\r\n\r\n    .catch((e) => {\r\n      errorAlert(getErrorMessage(e.response.data.code));\r\n      resetStepper();\r\n    });\r\n}\r\n\r\nfunction putUpload(bodyFormData, password) {\r\n  axios({\r\n    url: `/websiteCreator/prepare`,\r\n    method: \"put\",\r\n    data: bodyFormData,\r\n    headers: {\r\n      \"Content-Type\": \"multipart/form-data\",\r\n      Authorization: password,\r\n    },\r\n  })\r\n    .catch((e) => {\r\n      errorAlert(getErrorMessage(e.response.data.code));\r\n      resetStepper();\r\n    });\r\n}\r\n\r\n// Full deployment process\r\nfunction deployWebsiteAndUpload(password) {\r\n  const dnsNameInputValue = document.getElementById(\"websiteName\").value;\r\n\r\n  const file = document.querySelector(\".upload input\").files[0];\r\n  const bodyFormData = new FormData();\r\n  bodyFormData.append(\"url\", dnsNameInputValue);\r\n  bodyFormData.append(\"nickname\", getDefaultWallet());\r\n  bodyFormData.append(\"zipfile\", file);\r\n\r\n  uploadProcess(file, dnsNameInputValue, true, bodyFormData, (bodyFormData) =>\r\n    putUpload(bodyFormData, password)\r\n  );\r\n}\r\n\r\n// Full deployment process\r\nfunction uploadWebsite(file, count, password) {\r\n  const bodyFormData = new FormData();\r\n\r\n  const address = deployers[count].address;\r\n  bodyFormData.append(\"zipfile\", file);\r\n  bodyFormData.append(\"address\", address);\r\n  bodyFormData.append(\"nickname\", getDefaultWallet());\r\n\r\n  uploadProcess(\r\n    file,\r\n    deployers[count].name,\r\n    false,\r\n    bodyFormData,\r\n    (bodyFormData) => postUpload(bodyFormData, password)\r\n  );\r\n}\r\n\r\nasync function stepper(dnsName, totalChunk, isFullProcess) {\r\n  initStepper(dnsName, totalChunk);\r\n\r\n  if (isFullProcess) {\r\n    step1(dnsName, totalChunk);\r\n  } else {\r\n    $(\".circle\").eq(0).empty();\r\n    $(\".circle\").eq(0).append('<i class=\"bi bi-check\">');\r\n    $(\".title\").eq(0).removeClass(\"loading-dots\");\r\n    step3(getDeployerByDns(dnsName).address, totalChunk);\r\n  }\r\n}\r\n\r\nfunction initStepper(dnsName, totalChunk) {\r\n  $(\".website-card\").hide();\r\n  $(\".stepper\").show();\r\n\r\n  $(\".stepper-title\").html(\"Deployment of \" + dnsName);\r\n  $(\".title\").eq(0).addClass(\"loading-dots\");\r\n\r\n  $(\".title\")\r\n    .eq(2)\r\n    .text(\"Chunk upload \" + 1 + \" on \" + totalChunk);\r\n}\r\n\r\n// Step 1, wait for contract deployment\r\nfunction step1(dnsName, totalChunk) {\r\n  eventManager.subscribe(\r\n    `Website Deployer is deployed at :`,\r\n    getWallet(getDefaultWallet()).address,\r\n    (resp) => step2(dnsName, resp.data.data.split(\":\")[1], totalChunk)\r\n  );\r\n}\r\n\r\n// Step 2, wait for DNS setting\r\nfunction step2(dnsName, contractAddress, totalChunk) {\r\n  eventManager.subscribe(\r\n    `Resolver set to record key : record${dnsName} at address `,\r\n    getWallet(getDefaultWallet()).address,\r\n    (_) => {\r\n      step3(contractAddress, totalChunk);\r\n    }\r\n  );\r\n\r\n  $(\".circle\").eq(0).empty();\r\n  $(\".circle\").eq(0).append('<i class=\"bi bi-check\">');\r\n\r\n  $(\".title\").eq(0).removeClass(\"loading-dots\");\r\n  $(\".title\").eq(1).addClass(\"loading-dots\");\r\n}\r\n// Step 3, Monitor state of chunk uploding\r\nfunction step3(contractAddress, totalChunk) {\r\n  let actualChunk = 1;\r\n\r\n  for (let i = 0; i < totalChunk; i++) {\r\n\r\n    if (i == totalChunk - 1) {\r\n      eventManager.subscribe(\r\n        `Chunk ${i} of Website deployed to ${contractAddress}`,\r\n        getWallet(getDefaultWallet()).address,\r\n        (_) => {\r\n          resetStepper();\r\n        }\r\n      );\r\n    } \r\n    else {\r\n        eventManager.subscribe(\r\n          `Chunk ${i} of Website deployed to ${contractAddress}`,\r\n          getWallet(getDefaultWallet()).address,\r\n          (_) => {\r\n            actualChunk++;\r\n            $(\".title\")\r\n              .eq(2)\r\n              .text(\"Chunk upload \" + actualChunk + \" on \" + totalChunk);\r\n            $(\".title\").eq(2).addClass(\"loading-dots\");\r\n          }\r\n        );\r\n      }\r\n  }\r\n\r\n  $(\".circle\").eq(1).empty();\r\n  $(\".circle\").eq(1).append('<i class=\"bi bi-check\">');\r\n\r\n  $(\".title\").eq(1).removeClass(\"loading-dots\");\r\n  $(\".title\").eq(2).addClass(\"loading-dots\");\r\n}\r\n\r\nfunction resetStepper() {\r\n  $(\".website-card\").show();\r\n  $(\".stepper\").hide();\r\n\r\n  $(\".circle\").empty();\r\n  $(\".circle\").eq(0).html(\"1\");\r\n  $(\".circle\").eq(1).html(\"2\");\r\n  $(\".circle\").eq(2).html(\"3\");\r\n\r\n  $(\".title\").eq(2).html(\"Chunk upload\");\r\n\r\n  $(\".title\").eq(2).removeClass(\"loading-dots\");\r\n  getWebsiteDeployerSC();\r\n  document\r\n    .getElementById(\"wallet-popover\")\r\n    .classList.remove(\"popover__disabled\");\r\n}\r\n"
