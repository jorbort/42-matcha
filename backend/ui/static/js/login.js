document
  .getElementById("loginForm")
  .addEventListener("submit", function (event) {
    event.preventDefault();
    const formData = {
      username: document.getElementById("username").value,
      password: document.getElementById("password").value,
    };

    fetch("http://localhost:3000/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => response.json())
      .then((data) => {
        document.cookie =
          "access-token=" +
          data.access_token +
          "; path=/; secure; samesite=strict";
      })
      .catch((error) => {
        console.log(error);
      });
  });
