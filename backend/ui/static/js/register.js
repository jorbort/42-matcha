document
  .getElementById("signup-form")
  .addEventListener("submit", function (event) {
    event.preventDefault();
    let div = document.getElementById("message");
    const formData = {
      username: document.getElementById("username").value,
      email: document.getElementById("email").value,
      password: document.getElementById("password").value,
      first_name: document.getElementById("name").value,
      last_name: document.getElementById("surname").value,
    };
    console.log(formData);
    fetch("http://localhost:3000/create_user", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(response.statusText);
        }
        return response.json();
      })
      .then((data) => {
        div.innerHTML = "User created successfully, sent validation email";
        div.className = "active";
      })
      .catch((error) => {
        div.innerHTML = "Error creating user: " + error.message;
        div.className = "active";
      });
  });
