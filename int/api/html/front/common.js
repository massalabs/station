// utils

function ellipsis(str) {
    return str.substr(0, 5) + "..." + str.substr(str.length - 5, str.length);
}

const copyToClipboard = (value, successfully = () => null, failure = () => null) => {
    const clipboard = navigator.clipboard;
    if (clipboard !== undefined && clipboard !== "undefined") {
        navigator.clipboard.writeText(value).then(successfully, failure);
    } else {
        if (document.execCommand) {
            const el = document.createElement("input");
            el.value = value;
            document.body.append(el);

            el.select();
            el.setSelectionRange(0, value.length);

            if (document.execCommand("copy")) {
                successfully();
            }

            el.remove();
        } else {
            failure();
        }
    }
};

const addressInnerHTML = (address) =>
    `${ellipsis(
        address
    )} <svg xmlns="http://www.w3.org/2000/svg" onclick="copyToClipboard('${address}')" style="cursor: pointer;" width="16" height="16" fill="currentColor" class="bi bi-back" viewBox="0 0 16 16"><path d="M0 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v2h2a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-2H2a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H2z"/></svg>`;

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

// isExcludedOSAndFirefox returns true if user uses Windows or Mac and Firefox
function isExcludedOSAndFirefox() {
    // userAgent gets the os and navigator infos.
    // ex : Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0.
    let userAgent = navigator.userAgent;

    // isWindows checks if current os is Windows
    let isWindows = userAgent.indexOf("Windows") != -1;

    // isFirefox checks if current navigator is Firefox
    let isFirefox = userAgent.indexOf("Firefox") != -1;

    return isFirefox && isWindows;
}
window.isExcludedOSAndFirefox = isExcludedOSAndFirefox;
