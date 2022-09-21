let gWallets = [];
let deployers = [];
let actualTxType = "";
let nextFileToUpload;
let uploadable = false;

// INIT
getWallets();
getWebsiteDeployerSC();
initializeDefaultWallet();
setupModal();

// Write the default wallet text in wallet popover component
async function getWebsiteDeployerSC() {
  let defaultWallet = getDefaultWallet();

  $("#website-deployers-table tbody tr").remove();

  axios
    .get("/my/domains/" + defaultWallet)
    .then((websites) => {
      let count = 0;
      for (const website of websites.data) {
        tableInsert(website, count);
        count++;
      }
      deployers = websites.data;
    })
    .catch((e) => {
      console.error(e);
      errorAlert(getErrorMessage(e.response.data.code));
    });
}

// Write the default wallet text in wallet popover component
function initializeDefaultWallet() {
  let defaultWallet = getDefaultWallet();
  if (defaultWallet === "") {
    defaultWallet = "Connect";
  }
  $(".popover__title").html(defaultWallet);
}

// Retrieve the default wallet nickname in cookies
function getDefaultWallet() {
  let defaultWallet = "";
  const cookies = document.cookie.split(";");
  cookies.forEach((cookie) => {
    const keyValue = cookie.split("=");
    if (keyValue[0] === "defaultWallet") {
      defaultWallet = keyValue[1];
    }
  });
  return defaultWallet;
}

function setupModal() {
  $("#passwordModal").on("shown.bs.modal", function () {
    $("#passwordModal").trigger("focus");
  });
}

function setTxType(txType) {
  actualTxType = txType;
}

async function callTx() {
  const passwordValue = $("#walletPassword").val();

  if (actualTxType === "deployWebsiteCreator") {
    deployWebsiteDeployerSC(nextFileToUpload, address, passwordValue);
  }
  if (actualTxType === "deployWebsiteAndUpload") {
    deployWebsiteAndUpload(passwordValue);
  }
  if (actualTxType.includes("uploadWebsiteCreator")) {
    const websiteIndex = actualTxType.split("uploadWebsiteCreator")[1];
    uploadWebsite(nextFileToUpload, websiteIndex, passwordValue);
  }
}

// open file upload
function openDialog() {
  document.getElementById("fileid0").value = null;
  document.getElementById("fileid0").click();
}

// Handle event on file selecting
function handleFileSelect(evt) {
  let files = evt.target.files; // get files
  let f = files[0];
  const reader = new FileReader();
  reader.onload = (event) => importWallet(JSON.parse(event.target.result)); // desired file content
  reader.onerror = (error) => reject(error);
  reader.readAsText(f);
}

function errorAlert(error) {
  document.getElementsByClassName("alert-danger")[0].style.display = "block";
  document.getElementsByClassName("alert-danger")[0].innerHTML = error;
  setTimeout(function () {
    document.getElementsByClassName("alert-danger")[0].style.display = "none";
  }, 5000);
}

function successMessage(message) {
  document.getElementsByClassName("alert-primary")[0].style.display = "block";
  document.getElementsByClassName("alert-primary")[0].innerHTML = message;
  setTimeout(function () {
    document.getElementsByClassName("alert-primary")[0].style.display = "none";
  }, 5000);
}

// Append wallet accounts in popover component list
async function feedWallet(w) {
  let counter = 0;
  for (const wallet of w) {
    $("#wallet-list").append(
      "<li class='wallet-item'><a class='wallet-link' id='wallet-link-" +
        counter +
        "' onclick='changeDefaultWallet(event)' href='#'>" +
        wallet.nickname +
        "</a></li>"
    );
    counter++;
  }
}

// Handle popover click & update default wallet in cookies
function changeDefaultWallet(event) {
  const idElementClicked = event.target.id;
  const newDefaultWalletId = idElementClicked.split("-")[2];
  const walletName = gWallets[newDefaultWalletId].nickname;

  document.cookie = "defaultWallet=" + walletName;
  $(".popover__title").html(walletName);

  getWebsiteDeployerSC();
}

async function getWallets() {
  axios
    .get("/mgmt/wallet")
    .then((resp) => {
      if (resp) {
        gWallets = resp.data;
        feedWallet(gWallets);
      }
    })
    .catch((e) => {
      errorAlert(getErrorMessage(e.response.data.code));
    });
}

async function deployWebsiteDeployerSC(password) {
  let defaultWallet = getDefaultWallet();
  const dnsNameInputValue = document.getElementById("websiteName").value;

  if (dnsNameInputValue == "") {
    errorAlert("Input a DNS name");
  } else {
    document.getElementsByClassName("loading")[0].style.display =
      "inline-block";
    axios
      .put(
        "/websiteCreator/prepare/",
        { url: dnsNameInputValue, nickname: defaultWallet },
        {
          headers: {
            Authorization: password,
          },
        }
      )
      .then((operation) => {
        document.getElementsByClassName("loading")[0].style.display = "none";
        successMessage(
          "Contract deployed to address " + operation.data.address
        );
        getWebsiteDeployerSC();
      })
      .catch((e) => {
        document.getElementsByClassName("loading")[0].style.display = "none";
        errorAlert(getErrorMessage(e.response.data.code));
      });
  }
}

function tableInsert(resp, count) {
  const tBody = document
    .getElementById("website-deployers-table")
    .getElementsByTagName("tbody")[0];
  const row = tBody.insertRow(-1);
  const url = "http://" + resp.name + ".massa/";

  const cell0 = row.insertCell();
  const cell1 = row.insertCell();
  const cell2 = row.insertCell();
  const cell3 = row.insertCell();

  cell0.innerHTML = resp.name;
  cell1.innerHTML = addressInnerHTML(resp.address);
  cell2.innerHTML =
    "<a href='" +
    url +
    "'target = '_blank' rel='noopener noreferrer'>" +
    url +
    "</a>";
  cell3.innerHTML =
    "<div><input id='fileid" +
    count +
    "' type='file' hidden/><button id='upload-website" +
    count +
    "' class='primary-button' id='buttonid' type='button' value='Upload MB' >Edit</button><img src='./logo.png' style='display: none' class='massa-logo-spinner loading" +
    count +
    " alt='Massa logo' /></span></div>";

  document
    .getElementById(`upload-website${count}`)
    .addEventListener("click", function () {
      document.getElementById(`fileid${count}`).value = null;
      document.getElementById(`fileid${count}`).click();
    });

  document
    .getElementById(`fileid${count}`)
    .addEventListener("change", function (evt) {
      let files = evt.target.files; // get files
      nextFileToUpload = files[0];

      setTxType("uploadWebsiteCreator" + count);
      $("#passwordModal").modal("show");
    });
}

function uploadWebsite(file, count, password) {
  let defaultWallet = getDefaultWallet();
  const bodyFormData = new FormData();
  bodyFormData.append("zipfile", file);
  bodyFormData.append("address", deployers[count].address);
  bodyFormData.append("nickname", defaultWallet);
  document.getElementsByClassName("loading" + count)[0].style.display =
    "inline-block";
  axios({
    url: `/websiteCreator/upload`,
    method: "POST",
    data: bodyFormData,
    headers: {
      "Content-Type": "multipart/form-data",
      Authorization: password,
    },
  })
    .then((operation) => {
      document.getElementsByClassName("loading" + count)[0].style.display =
        "none";
      successMessage("Website uploaded to address : " + operation.data.address);
    })

    .catch((e) => {
      document.getElementsByClassName("loading" + count)[0].style.display =
        "none";
      errorAlert(getErrorMessage(e.response.data.code));
    });
}

$("#file-select-button").click(function () {
  $(".upload input").click();
});

// change button text with file name
$(".upload input").on("change", function () {
  let str = $(".upload input").val();

  let n = str.lastIndexOf("\\");

  let result = str.substring(n + 1);

  $("#file-select-button").html(result);
});

//check if file is .zip
$(".upload input").on("change", function () {
  let str = $(".upload input").val();

  let n = str.lastIndexOf(".");

  let result = str.substring(n + 1);

  if (result != "zip" && $(".upload input").val() != "") {
    uploadable = false;

    document.getElementsByClassName("fileError")[0].style.display = "flex";
    document.getElementById("website-upload").style.display = "none";
    document.getElementById("website-upload-refuse").style.display = "flex";
  } else {
    uploadable = true;
    document.getElementsByClassName("fileError")[0].style.display = "none";
    document.getElementById("website-upload").style.display = "flex";
    document.getElementById("website-upload-refuse").style.display = "none";
  }
});

//remove label of input website name on focus
$(".website-dns input").on("focus", function () {
  document.getElementById("website-info-display").style.visibility = "hidden";
});

//check if input string is valid
$(".website-dns input").on("change", function () {
  let str = $(".website-dns input").val();
  let pattern = new RegExp("^[a-z0-9]+$");
  let tesPattern = true;
  testPattern = pattern.test(str);

  if (testPattern == false) {
    uploadable = false;

    document.getElementsByClassName("dns-error")[0].style.display = "flex";
    document.getElementById("website-upload").style.display = "none";
    document.getElementById("website-upload-refuse").style.display = "flex";
  } else {
    uploadable = true;
    document.getElementsByClassName("dns-error")[0].style.display = "none";
    document.getElementById("website-upload").style.display = "flex";
    document.getElementById("website-upload-refuse").style.display = "none";
  }
});

function deployWebsiteAndUpload(password) {
  if (uploadable == true) {
    let defaultWallet = getDefaultWallet();
    const dnsNameInputValue = document.getElementById("websiteName").value;

    const file = document.querySelector(".upload input").files[0];
    const bodyFormData = new FormData();
    bodyFormData.append("url", dnsNameInputValue);
    bodyFormData.append("nickname", defaultWallet);
    bodyFormData.append("zipfile", file);

    document.getElementsByClassName("loading")[0].style.display =
      "inline-block";
    axios({
      url: `/websiteCreator/prepare`,
      method: "put",
      data: bodyFormData,
      headers: {
        "Content-Type": "multipart/form-data",
        Authorization: password,
      },
    })
      .then((operation) => {
        document.getElementsByClassName("loading")[0].style.display = "none";
        successMessage(
          "Website uploaded to address : " + operation.data.address
        );
        getWebsiteDeployerSC();
      })
      .catch((e) => {
        document.getElementsByClassName("loading")[0].style.display = "none";
        errorAlert(getErrorMessage(e.response.data.code));
      });
  } else {
  }
}
