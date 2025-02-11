const message = document.getElementById("message");
if (localStorage.getItem("profile_completed") == "false") {
  message.textContent =
    "Please complete your profile on the settings tab to continue";
}
