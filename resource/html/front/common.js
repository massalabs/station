// utils

function ellipsis(str) {
	return str.substr(0, 5) + '...' + str.substr(str.length - 5, str.length);
}

const copyToClipboard = (
  value,
  successfully = () => null,
  failure = () => null
) => {
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

const addressInnerHTML = (address) => `${ellipsis(address)} <svg xmlns="http://www.w3.org/2000/svg" onclick="copyToClipboard('${address}')" style="cursor: pointer;" width="16" height="16" fill="currentColor" class="bi bi-back" viewBox="0 0 16 16"><path d="M0 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v2h2a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-2H2a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H2z"/></svg>`;
