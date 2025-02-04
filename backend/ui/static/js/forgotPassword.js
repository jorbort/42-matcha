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
        message.textContent = "invalid email";
        throw error;
      } else {
        message.textContent = "Email sent successfully";
      }
    } catch (error) {
      console.error(error);
    }
  });
