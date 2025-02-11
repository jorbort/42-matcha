document
  .getElementById("loginForm")
  .addEventListener("submit", async function (event) {
    event.preventDefault();
    const message = document.getElementById("response-message");
    const formData = {
      username: document.getElementById("username").value,
      password: document.getElementById("password").value,
    };
    try {
      const response = await fetch("http://localhost:3000/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });
      if (!response.ok) {
        const error = await response.json();
        message.textContent = "Invalid username or password";
        message.style.color = "#ea9d34";
        message.style.textShadow = "1px 1px 1px #ea9d34";
        throw error;
      } else if (response.ok) {
        const data = await response.json();
        console.log(data);
        localStorage.setItem("username", data.username);
        localStorage.setItem("profile_completed", data.is_completed);
        window.location.href = "http://localhost:3000/profile";
      }
    } catch (error) {}
  });
