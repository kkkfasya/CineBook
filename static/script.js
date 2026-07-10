function api(method, path, body) {
  var opts = {
    method: method,
    headers: { "Content-Type": "application/json" },
  };
  if (body !== undefined && body !== null) opts.body = JSON.stringify(body);
  return fetch(path, opts).then(function (r) {
    if (r.status === 204) return null;
    return r.json().then(function (data) {
      if (!r.ok) throw new Error(data.error || "request failed");
      return data;
    });
  });
}

function escapeHtml(str) {
  var div = document.createElement("div");
  div.textContent = str;
  return div.innerHTML;
}
