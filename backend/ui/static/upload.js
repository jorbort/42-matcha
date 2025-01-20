function getCookie(name) {
  let cookieArr = document.cookie.split(";");
  for (let i = 0; i < cookieArr.length; i++) {
    let cookiePair = cookieArr[i].split("=");
    if (name === cookiePair[0].trim()) {
      return decodeURIComponent(cookiePair[1]);
    }
  }
  return null;
}

document
  .getElementById("uploadForm")
  .addEventListener("submit", function (event) {
    event.preventDefault();

    const formData = new FormData();
    formData.append("user_id", document.getElementById("user_id").value);
    formData.append(
      "picture_number",
      document.getElementById("picture_number").value,
    );
    formData.append("image", document.getElementById("image").files[0]);

    const jwtToken = getCookie("access-token");

    fetch("http://localhost:3000/uploadImg", {
      method: "POST",
      headers: {
        Authorization: "Bearer " + jwtToken,
      },
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        document.getElementById("response").innerText = JSON.stringify(
          data,
          null,
          2,
        );
      })
      .catch((error) => {
        document.getElementById("response").innerText = "Error: " + error;
      });
  });
