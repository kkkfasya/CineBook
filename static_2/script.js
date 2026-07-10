(function () {
  var BASE = "/api/v1/admin/movies";

  var editingId = null;

  var grid = document.getElementById("movieGrid");
  var overlay = document.getElementById("overlay");
  var modalTitle = document.getElementById("modalTitle");
  var form = document.getElementById("movieForm");
  var movieId = document.getElementById("movieId");
  var inputName = document.getElementById("inputName");
  var inputPoster = document.getElementById("inputPoster");
  var inputLength = document.getElementById("inputLength");
  var inputRows = document.getElementById("inputRows");
  var inputSeatsPerRow = document.getElementById("inputSeatsPerRow");

  document.getElementById("btnAdd").addEventListener("click", openAdd);
  document.getElementById("btnCancel").addEventListener("click", closeModal);
  overlay.addEventListener("click", function (e) {
    if (e.target === overlay) closeModal();
  });
  form.addEventListener("submit", saveMovie);

  loadMovies();

  // --- API ---

  function api(method, path, body) {
    var opts = {
      method: method,
      headers: { "Content-Type": "application/json" },
    };
    if (body !== undefined) opts.body = JSON.stringify(body);
    return fetch(path, opts).then(function (r) {
      if (r.status === 204) return null;
      return r.json().then(function (data) {
        if (!r.ok) throw new Error(data.error || "request failed");
        return data;
      });
    });
  }

  // --- CRUD ---

  function loadMovies() {
    api("GET", BASE).then(function (movies) {
      renderGrid(movies || []);
    });
  }

  function createMovie(data) {
    return api("POST", BASE, data).then(function () {
      closeModal();
      loadMovies();
    });
  }

  function updateMovie(id, data) {
    return api("PUT", BASE + "/" + id, data).then(function () {
      closeModal();
      loadMovies();
    });
  }

  function deleteMovie(id) {
    if (!confirm("Delete this movie?")) return;
    api("DELETE", BASE + "/" + id).then(function () {
      loadMovies();
    });
  }

  // --- Render ---

  function renderGrid(movies) {
    grid.innerHTML = "";
    if (movies.length === 0) {
      grid.innerHTML = '<div class="empty-state">No movies yet. Click "+ Add Movie" to begin.</div>';
      return;
    }
    movies.forEach(function (m) {
      var card = document.createElement("div");
      card.className = "movie-card";

      if (m.movie_poster) {
        var img = document.createElement("img");
        img.className = "movie-poster";
        img.src = m.movie_poster;
        img.alt = m.movie_name;
        img.onerror = function () {
          img.style.display = "none";
          var placeholder = document.createElement("div");
          placeholder.className = "movie-poster--empty";
          placeholder.textContent = "no poster";
          card.insertBefore(placeholder, card.firstChild);
        };
        card.appendChild(img);
      } else {
        var placeholder = document.createElement("div");
        placeholder.className = "movie-poster--empty";
        placeholder.textContent = "no poster";
        card.appendChild(placeholder);
      }

      var nameEl = document.createElement("div");
      nameEl.className = "movie-name";
      nameEl.textContent = m.movie_name;
      card.appendChild(nameEl);

      var meta = document.createElement("div");
      meta.className = "movie-meta";

      var mins = Math.floor(m.movie_length_sec / 60);
      var secs = m.movie_length_sec % 60;
      meta.innerHTML =
        "<span>" + mins + "m " + secs + "s</span>" +
        "<span>" + m.movie_seat_row + " rows &times; " + m.movie_seat_per_row + " seats</span>";
      card.appendChild(meta);

      var actions = document.createElement("div");
      actions.className = "movie-actions";

      var editBtn = document.createElement("button");
      editBtn.className = "btn btn--edit";
      editBtn.textContent = "Edit";
      editBtn.addEventListener("click", function () {
        openEdit(m);
      });
      actions.appendChild(editBtn);

      var delBtn = document.createElement("button");
      delBtn.className = "btn btn--delete";
      delBtn.textContent = "Delete";
      delBtn.addEventListener("click", function () {
        deleteMovie(m.id);
      });
      actions.appendChild(delBtn);

      card.appendChild(actions);
      grid.appendChild(card);
    });
  }

  // --- Modal ---

  function openAdd() {
    editingId = null;
    modalTitle.textContent = "Add Movie";
    movieId.value = "";
    form.reset();
    overlay.classList.remove("hidden");
  }

  function openEdit(m) {
    editingId = m.id;
    modalTitle.textContent = "Edit Movie";
    movieId.value = m.id;
    inputName.value = m.movie_name || "";
    inputPoster.value = m.movie_poster || "";
    inputLength.value = m.movie_length_sec || "";
    inputRows.value = m.movie_seat_row || "";
    inputSeatsPerRow.value = m.movie_seat_per_row || "";
    overlay.classList.remove("hidden");
  }

  function closeModal() {
    overlay.classList.add("hidden");
    editingId = null;
  }

  function saveMovie(e) {
    e.preventDefault();
    var data = {
      movie_name: inputName.value.trim(),
      movie_poster: inputPoster.value.trim() || undefined,
      movie_length_sec: parseInt(inputLength.value, 10),
      movie_seat_row: parseInt(inputRows.value, 10),
      movie_seat_per_row: parseInt(inputSeatsPerRow.value, 10),
    };

    if (!data.movie_name) return;

    if (editingId) {
      updateMovie(editingId, data);
    } else {
      createMovie(data);
    }
  }
})();
