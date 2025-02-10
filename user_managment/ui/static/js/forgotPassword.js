document
  .getElementById("reset-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();
    var email = document.getElementById("email").value;
    var message = document.getElementById("response-message");
    const body = {
      email: email,
    };
    try {
      const response = await fetch("http://localhost:3000/SendResetPassword", {
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
        message.textContent = "invalid email";
        throw error;
      } else {
        message.style.color = "#b4637a";
        message.style.textShadow = "0 0 5px #b4637a";
        message.textContent = "Email sent successfully";
      }
    } catch (error) {}
  });
