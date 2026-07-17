(function () {
  var BASE = "/api/v1/admin/movies";
  var editingId = null;

  var grid = document.getElementById("movieGrid");
  var overlay = document.getElementById("overlay");
  var modalTitle = document.getElementById("modalTitle");
  var form = document.getElementById("movieForm");
  var movieId = document.getElementById("movieId");
  var inputTitle = document.getElementById("inputTitle");
  var inputPoster = document.getElementById("inputPoster");
  var inputRows = document.getElementById("inputRows");
  var inputSeatsPerRow = document.getElementById("inputSeatsPerRow");

  document.getElementById("btnAdd").addEventListener("click", openAdd);
  document.getElementById("btnCancel").addEventListener("click", closeModal);
  overlay.addEventListener("click", function (e) {
    if (e.target === overlay) closeModal();
  });
  form.addEventListener("submit", saveMovie);

  loadMovies();

  /* ── CRUD ── */

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
    api("DELETE", BASE + "/" + id).then(function () { loadMovies(); });
  }

  /* ── Render ── */

  function renderGrid(movies) {
    grid.innerHTML = "";
    if (movies.length === 0) {
      grid.innerHTML = '<div class="empty-state">No movies yet. Click "+ Add Movie" to begin.</div>';
      return;
    }
    movies.forEach(function (m) {
      var card = document.createElement("div");
      card.className = "admin-card";

      if (m.poster) {
        var img = document.createElement("img");
        img.className = "admin-poster";
        img.src = m.poster;
        img.alt = m.title;
        img.onerror = function () {
          img.style.display = "none";
          var ph = document.createElement("div");
          ph.className = "admin-poster--empty";
          ph.textContent = "no poster";
          card.insertBefore(ph, card.firstChild);
        };
        card.appendChild(img);
      } else {
        var ph = document.createElement("div");
        ph.className = "admin-poster--empty";
        ph.textContent = "no poster";
        card.appendChild(ph);
      }

      var nameEl = document.createElement("div");
      nameEl.className = "admin-name";
      nameEl.textContent = m.title;
      card.appendChild(nameEl);

      var meta = document.createElement("div");
      meta.className = "admin-meta";
      meta.innerHTML =
        "<span>" + m.rows + " rows &times; " + m.seats_per_row + " seats</span>";
      card.appendChild(meta);

      var actions = document.createElement("div");
      actions.className = "admin-actions";

      var editBtn = document.createElement("button");
      editBtn.className = "btn btn--edit";
      editBtn.textContent = "Edit";
      editBtn.addEventListener("click", function () { openEdit(m); });
      actions.appendChild(editBtn);

      var delBtn = document.createElement("button");
      delBtn.className = "btn btn--delete";
      delBtn.textContent = "Delete";
      delBtn.addEventListener("click", function () { deleteMovie(m.id); });
      actions.appendChild(delBtn);

      card.appendChild(actions);
      grid.appendChild(card);
    });
  }

  /* ── Modal ── */

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
    inputTitle.value = m.title || "";
    inputPoster.value = m.poster || "";
    inputRows.value = m.rows || "";
    inputSeatsPerRow.value = m.seats_per_row || "";
    overlay.classList.remove("hidden");
  }

  function closeModal() {
    overlay.classList.add("hidden");
    editingId = null;
  }

  function saveMovie(e) {
    e.preventDefault();
    var data = {
      title: inputTitle.value.trim(),
      poster: inputPoster.value.trim() || undefined,
      rows: parseInt(inputRows.value, 10),
      seats_per_row: parseInt(inputSeatsPerRow.value, 10),
    };
    if (!data.title) return;
    if (editingId) updateMovie(editingId, data);
    else createMovie(data);
  }
})();
