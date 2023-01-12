let gWallets = [];
let deployers = [];
let actualTxType = "";
let nextFileToUpload;
let uploadable = false;

// INIT
const uploadMaxSize = document.getElementById("UploadMaxSize").innerHTML;
setMaxSizeLabel();
getWallets();
getWebsiteDeployerSC();
initializeDefaultWallet();

const eventManager = new EventManager();

async function onSubmitDeploy(txType = "deployWebsiteAndUpload") {
    setTxType(txType);
    callTx();
}

// Write the default wallet text in wallet popover component
async function getWebsiteDeployerSC() {
    let defaultWallet = getDefaultWallet();
    if (defaultWallet === "") {
        errorAlert(getErrorMessage("Wallet-5001"));
        return;
    }

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
        .catch(handleAPIError);
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

function getWallet(nickname) {
    return gWallets.find((w) => w.nickname === nickname);
}

function getDeployerByDns(dns) {
    return deployers.find((c) => c.name === dns);
}

function getWallet(nickname) {
    return gWallets.find((w) => w.nickname === nickname);
}

function getDeployerByDns(dns) {
    return deployers.find((c) => c.name === dns);
}

function setTxType(txType) {
    actualTxType = txType;
}

async function callTx() {
    let wallet = getDefaultWallet();
    if (wallet === "") {
        errorAlert(getErrorMessage("Wallet-5002"));
        return;
    }

    if (actualTxType === "deployWebsiteAndUpload") {
        deployWebsiteAndUpload();
    }
    if (actualTxType.includes("uploadWebsiteCreator")) {
        const websiteIndex = actualTxType.split("uploadWebsiteCreator")[1];
        uploadWebsite(nextFileToUpload, websiteIndex);
    }
    if (actualTxType.includes("uploadMissingChunks")) {
        const websiteIndex = actualTxType.split("uploadMissingChunks")[1];
        uploadMissingChunks(nextFileToUpload, websiteIndex);
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
        .catch(handleAPIError);
}
let missingChunks = [[]];
async function tableInsert(resp, count) {
    const tBody = document
        .getElementById("website-deployers-table")
        .getElementsByTagName("tbody")[0];
    const row = tBody.insertRow(-1);

    const protocol = isExcludedOSAndFirefox() ? "http" : "https";
    const url = `${protocol}://${resp.name}.massa:${window.location.port}/`;

    const cell0 = row.insertCell();
    const cell1 = row.insertCell();
    const cell2 = row.insertCell();
    const cell3 = row.insertCell();
    cell0.innerHTML = resp.name;
    cell1.innerHTML = addressInnerHTML(resp.address);

    if (resp.brokenChunks != null) {
        missingChunks[count] = await resp.brokenChunks;
        cell2.innerHTML = `
            <div class="broken-chunks">
                <div class="tooltip">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle" viewBox="0 0 16 16">
                        <path d="M7.938 2.016A.13.13 0 0 1 8.002 2a.13.13 0 0 1 .063.016.146.146 0 0 1 .054.057l6.857 11.667c.036.06.035.124.002.183a.163.163 0 0 1-.054.06.116.116 0 0 1-.066.017H1.146a.115.115 0 0 1-.066-.017.163.163 0 0 1-.054-.06.176.176 0 0 1 .002-.183L7.884 2.073a.147.147 0 0 1 .054-.057zm1.044-.45a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566z"/>
                        <path d="M7.002 12a1 1 0 1 1 2 0 1 1 0 0 1-2 0zM7.1 5.995a.905.905 0 1 1 1.8 0l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995z"/>
                    </svg>
                    <span class="tooltiptext">Website corrupted. Click on the gear wheel to fix it.</span>
                </div>
                <div class="repair-gear">
                    <input id='repair-fileid${count}' type='file' hidden/>
                    <svg id='fix-broken-website-${count}' xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-gear-fill" viewBox="0 0 16 16">
                        <path d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872l-.1-.34zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z"/>
                    </svg>
                </div>
            </div>`;
        document
            .getElementById(`fix-broken-website-${count}`)
            .addEventListener("click", function () {
                document.getElementById(`repair-fileid${count}`).value = null;
                document.getElementById(`repair-fileid${count}`).click();
            });
        document.getElementById(`repair-fileid${count}`).addEventListener("change", function (evt) {
            let files = evt.target.files;
            nextFileToUpload = files[0];
            onSubmitDeploy("uploadMissingChunks" + count);
        });
    } else {
        cell2.innerHTML = `<a href="${url}" target="_blank" rel="noopener noreferrer"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-globe2" viewBox="0 0 16 16"><path d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm7.5-6.923c-.67.204-1.335.82-1.887 1.855-.143.268-.276.56-.395.872.705.157 1.472.257 2.282.287V1.077zM4.249 3.539c.142-.384.304-.744.481-1.078a6.7 6.7 0 0 1 .597-.933A7.01 7.01 0 0 0 3.051 3.05c.362.184.763.349 1.198.49zM3.509 7.5c.036-1.07.188-2.087.436-3.008a9.124 9.124 0 0 1-1.565-.667A6.964 6.964 0 0 0 1.018 7.5h2.49zm1.4-2.741a12.344 12.344 0 0 0-.4 2.741H7.5V5.091c-.91-.03-1.783-.145-2.591-.332zM8.5 5.09V7.5h2.99a12.342 12.342 0 0 0-.399-2.741c-.808.187-1.681.301-2.591.332zM4.51 8.5c.035.987.176 1.914.399 2.741A13.612 13.612 0 0 1 7.5 10.91V8.5H4.51zm3.99 0v2.409c.91.03 1.783.145 2.591.332.223-.827.364-1.754.4-2.741H8.5zm-3.282 3.696c.12.312.252.604.395.872.552 1.035 1.218 1.65 1.887 1.855V11.91c-.81.03-1.577.13-2.282.287zm.11 2.276a6.696 6.696 0 0 1-.598-.933 8.853 8.853 0 0 1-.481-1.079 8.38 8.38 0 0 0-1.198.49 7.01 7.01 0 0 0 2.276 1.522zm-1.383-2.964A13.36 13.36 0 0 1 3.508 8.5h-2.49a6.963 6.963 0 0 0 1.362 3.675c.47-.258.995-.482 1.565-.667zm6.728 2.964a7.009 7.009 0 0 0 2.275-1.521 8.376 8.376 0 0 0-1.197-.49 8.853 8.853 0 0 1-.481 1.078 6.688 6.688 0 0 1-.597.933zM8.5 11.909v3.014c.67-.204 1.335-.82 1.887-1.855.143-.268.276-.56.395-.872A12.63 12.63 0 0 0 8.5 11.91zm3.555-.401c.57.185 1.095.409 1.565.667A6.963 6.963 0 0 0 14.982 8.5h-2.49a13.36 13.36 0 0 1-.437 3.008zM14.982 7.5a6.963 6.963 0 0 0-1.362-3.675c-.47.258-.995.482-1.565.667.248.92.4 1.938.437 3.008h2.49zM11.27 2.461c.177.334.339.694.482 1.078a8.368 8.368 0 0 0 1.196-.49 7.01 7.01 0 0 0-2.275-1.52c.218.283.418.597.597.932zm-.488 1.343a7.765 7.765 0 0 0-.395-.872C9.835 1.897 9.17 1.282 8.5 1.077V4.09c.81-.03 1.577-.13 2.282-.287z"/></svg></a>`;
    }
    cell3.innerHTML =
        "<div><input id='fileid" +
        count +
        "' type='file' hidden/><button id='upload-website" +
        count +
        "' class='primary-button' type='button' value='Upload MB' >Edit</button><img src='./logo.png' style='display: none' class='massa-logo-spinner loading" +
        count +
        " alt='Massa logo' /></span></div>";

    document.getElementById(`upload-website${count}`).addEventListener("click", function () {
        document.getElementById(`fileid${count}`).value = null;
        document.getElementById(`fileid${count}`).click();
    });

    document.getElementById(`fileid${count}`).addEventListener("change", function (evt) {
        let files = evt.target.files; // get files
        nextFileToUpload = files[0];

        onSubmitDeploy("uploadWebsiteCreator" + count);
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

$(".upload input").on("change", function () {
    const file = this.files[0];
    if (!file) {
        uploadable = false;
        document.getElementsByClassName("fileTypeError")[0].style.display = "none";
        document.getElementsByClassName("fileSizeError")[0].style.display = "none";
        document.getElementById("website-upload").style.display = "none";
        document.getElementById("website-upload-refuse").style.display = "flex";
        $("#file-select-button").html("Import From");
        return;
    }

    const fileSize = file.size;
    const filename = file.name;
    const n = filename.lastIndexOf(".");
    const fileExtension = filename.substring(n + 1);

    if (fileExtension != "zip" && filename != "") {
        uploadable = false;
        document.getElementsByClassName("fileTypeError")[0].style.display = "flex";
        document.getElementsByClassName("fileSizeError")[0].style.display = "none";
        document.getElementById("website-upload").style.display = "none";
        document.getElementById("website-upload-refuse").style.display = "flex";
        $("#file-select-button").html(filename);
    } else if (fileSize > uploadMaxSize) {
        uploadable = false;
        document.getElementsByClassName("fileTypeError")[0].style.display = "none";
        document.getElementsByClassName("fileSizeError")[0].style.display = "flex";
        document.getElementById("website-upload").style.display = "none";
        document.getElementById("website-upload-refuse").style.display = "flex";
    } else {
        uploadable = true;
        document.getElementsByClassName("fileTypeError")[0].style.display = "none";
        document.getElementsByClassName("fileSizeError")[0].style.display = "none";
        document.getElementById("website-upload").style.display = "flex";
        document.getElementById("website-upload-refuse").style.display = "none";
        $("#file-select-button").html(filename);
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
    let testPattern = pattern.test(str);

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

function uploadProcess(file, dnsName, isFullProcess, bodyFormData, callback) {
    document.getElementById("wallet-popover").classList.add("popover__disabled");
    const reader = new FileReader();
    reader.readAsBinaryString(file);
    reader.onloadend = (_) => {
        const result = reader.result.length;

        const chunkSize = Math.floor(result / 260_000) + 1;

        stepper(dnsName, chunkSize, isFullProcess);

        callback(bodyFormData);
    };
}

function postUpload(bodyFormData) {
    axios({
        url: `/websiteCreator/upload`,
        method: "POST",
        data: bodyFormData,
        headers: {
            "Content-Type": "multipart/form-data",
        },
    }).catch((e) => {
        handleAPIError(e);
        resetStepper();
    });
}

function postUploadMissedChunks(bodyFormData) {
    axios({
        url: `/websiteCreator/uploadMissingChunks`,
        method: "POST",
        data: bodyFormData,
        headers: {
            "Content-Type": "multipart/form-data",
        },
    }).catch((e) => {
        errorAlert(getErrorMessage(e.response.data.code));
        resetStepper();
    });
}

function putUpload(bodyFormData) {
    axios({
        url: `/websiteCreator/prepare`,
        method: "put",
        data: bodyFormData,
        headers: {
            "Content-Type": "multipart/form-data",
        },
    }).catch((e) => {
        handleAPIError(e);
        resetStepper();
    });
}

// Full deployment process
function deployWebsiteAndUpload() {
    const dnsNameInputValue = document.getElementById("websiteName").value;

    const file = document.querySelector(".upload input").files[0];
    const bodyFormData = new FormData();
    bodyFormData.append("url", dnsNameInputValue);
    bodyFormData.append("nickname", getDefaultWallet());
    bodyFormData.append("zipfile", file);

    uploadProcess(file, dnsNameInputValue, true, bodyFormData, (bodyFormData) =>
        putUpload(bodyFormData)
    );
}

// Full deployment process
function uploadWebsite(file, count) {
    const bodyFormData = new FormData();
    const address = deployers[count].address;
    bodyFormData.append("zipfile", file);
    bodyFormData.append("address", address);
    bodyFormData.append("nickname", getDefaultWallet());
    uploadProcess(file, deployers[count].name, false, bodyFormData, (bodyFormData) =>
        postUpload(bodyFormData)
    );
}

// Process Missing Chunks
function uploadMissingChunks(file, count) {
    const bodyFormData = new FormData();
    const missedChunks = missingChunks[count];
    const address = deployers[count].address;
    bodyFormData.append("zipfile", file);
    bodyFormData.append("address", address);
    bodyFormData.append("nickname", getDefaultWallet());
    bodyFormData.append("missedChunks", missedChunks);
    uploadProcess(file, deployers[count].name, false, bodyFormData, (bodyFormData) =>
        postUploadMissedChunks(bodyFormData)
    );
}

async function stepper(dnsName, totalChunk, isFullProcess) {
    initStepper(dnsName, totalChunk);

    if (isFullProcess) {
        step1(dnsName, totalChunk);
    } else {
        $(".circle").eq(0).empty();
        $(".circle").eq(0).append('<i class="bi bi-check">');
        $(".title").eq(0).removeClass("loading-dots");
        step3(getDeployerByDns(dnsName).address, totalChunk);
    }
}

function initStepper(dnsName, totalChunk) {
    $(".website-card").hide();
    $(".stepper").show();

    $(".stepper-title").html("Deployment of " + dnsName);
    $(".title").eq(0).addClass("loading-dots");

    $(".title")
        .eq(2)
        .text("Chunk upload " + 1 + " on " + totalChunk);

    eventManager.subscribe(`ERROR :`, getWallet(getDefaultWallet()).address, (resp) => {
        resetStepper();
        console.error(resp.data);
        errorAlert(resp.data.data.split(":")[1]);
    });
}

// Step 1, wait for contract deployment
function step1(dnsName, totalChunk) {
    eventManager.subscribe(
        `Website Deployer is deployed at :`,
        getWallet(getDefaultWallet()).address,
        (resp) => step2(dnsName, resp.data.data.split(":")[1], totalChunk)
    );
}

// Step 2, wait for DNS setting
function step2(dnsName, contractAddress, totalChunk) {
    eventManager.subscribe(
        `Website name ${dnsName} added to DNS at address ${contractAddress}`,
        getWallet(getDefaultWallet()).address,
        (_) => {
            step3(contractAddress, totalChunk);
        }
    );

    $(".circle").eq(0).empty();
    $(".circle").eq(0).append('<i class="bi bi-check">');

    $(".title").eq(0).removeClass("loading-dots");
    $(".title").eq(1).addClass("loading-dots");
}

// Step 3, Monitor state of chunk uploding
function step3(contractAddress, totalChunk) {
    let actualChunk = 1;

    for (let i = 0; i < totalChunk; i++) {
        eventManager.subscribe(
            `Website chunk deployed to ${contractAddress} on key ${i}`,
            getWallet(getDefaultWallet()).address,
            (_) => {
                actualChunk++;
                $(".title")
                    .eq(2)
                    .text("Chunk upload " + actualChunk + " on " + totalChunk);
                $(".title").eq(2).addClass("loading-dots");

                if (actualChunk - 1 == totalChunk) {
                    resetStepper();
                }
            }
        );
    }

    $(".circle").eq(1).empty();
    $(".circle").eq(1).append('<i class="bi bi-check">');

    $(".title").eq(1).removeClass("loading-dots");
    $(".title").eq(2).addClass("loading-dots");
}

function resetStepper() {
    $(".website-card").show();
    $(".stepper").hide();

    $(".circle").empty();
    $(".circle").eq(0).html("1");
    $(".circle").eq(1).html("2");
    $(".circle").eq(2).html("3");

    $(".title").eq(2).html("Chunk upload");

    $(".title").eq(2).removeClass("loading-dots");
    getWebsiteDeployerSC();
    document.getElementById("wallet-popover").classList.remove("popover__disabled");
}

function setMaxSizeLabel() {
    const spans = document.getElementsByClassName("UploadMaxSizeLabel");
    for (let span of spans) {
        span.innerHTML = formatBytes(uploadMaxSize);
    }
}

function formatBytes(bytes, decimals = 2, isBinary = false) {
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"]; // or ['B', 'KB', 'MB', 'GB', 'TB']

    if (!+bytes) {
        return `0 ${sizes[0]}`;
    }

    const inByte = isBinary ? 1024 : 1000;
    const dm = decimals < 0 ? 0 : decimals;

    const pow = Math.floor(Math.log(bytes) / Math.log(inByte));
    const maxPow = Math.min(pow, sizes.length - 1);

    return `${parseFloat((bytes / Math.pow(inByte, maxPow)).toFixed(dm))} ${sizes[maxPow]}`;
}
