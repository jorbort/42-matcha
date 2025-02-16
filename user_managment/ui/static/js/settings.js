document.addEventListener("DOMContentLoaded", function () {
  const storedUserId = localStorage.getItem("userID");
  if (storedUserId) {
    document.getElementById("userId").value = storedUserId;
  } else {
    console.warn("User ID not found in localStorage.");
  }
});
