document.addEventListener("DOMContentLoaded", function () {
  const storedUserId = localStorage.getItem("userID");
  if (storedUserId) {
    document.getElementById("userId").value = storedUserId;
  } else {
    console.warn("User ID not found in localStorage.");
  }
});

if (navigator.geolocation) {
  navigator.geolocation.getCurrentPosition((position) => {
    var latitude = position.coords.latitude;
    var longitude = position.coords.longitude;
    document.cookie = `latitude=${latitude}`;
    document.cookie = `longitude=${longitude}`;
  });
}
