(function () {
  var form = document.getElementById("loginForm");
  var username = document.getElementById("inputUsername");
  var password = document.getElementById("inputPassword");
  var errorEl = document.getElementById("loginError");

  form.addEventListener("submit", function (e) {
    e.preventDefault();
    errorEl.classList.add("hidden");

    var body = {
      username: username.value.trim(),
      password: password.value,
    };

    if (!body.username || !body.password) {
      showError("Both fields are required.");
      return;
    }

    fetch("/api/v1/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    }).then(function (r) {
      return r.json().then(function (data) {
        if (!r.ok) throw new Error(data.error || "login failed");
        return data;
      });
    }).then(function (data) {
      if (data.token) localStorage.setItem("token", data.token);
      window.location.href = data.redirect || "/admin/";
    }).catch(function (err) {
      showError(err.message);
    });
  });

  function showError(msg) {
    errorEl.textContent = msg;
    errorEl.classList.remove("hidden");
  }
})();
