document
  .getElementById("reset-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();
    const url = new URL(window.location.href);
    const params = new URLSearchParams(url.search);
    const code = params.get("code");
    var password = document.getElementById("password").value;
    var message = document.getElementById("response-message");
    const body = {
      password: password,
      code: code,
    };
    try {
      const response = await fetch("http://localhost:3000/updatePassword", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      });
      if (!response.ok) {
        const error = await response.json();
        message.style.color = "#ea9d34";
        message.style.textShadow = "0 0 5px #ea9d34";
        message.textContent = "invalid password";
        throw error;
      } else {
        message.style.color = "#b4637a";
        message.style.textShadow = "0 0 5px #b4637a";
        message.textContent = "Password updated successfully";
      }
    } catch (error) {}
  });
