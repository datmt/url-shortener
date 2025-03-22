let authHeader = "";

function login() {
  const u = document.getElementById("username").value;
  const p = document.getElementById("password").value;
  authHeader = "Basic " + btoa(`${u}:${p}`);
  localStorage.setItem("auth", authHeader);
  loadLinks();
}

function logout() {
  localStorage.removeItem("auth");
  authHeader = "";
  document.getElementById("dashboard-section").style.display = "none";
  document.getElementById("login-section").style.display = "block";
}

function createLink() {
  const handle = document.getElementById("handle").value;
  const target = document.getElementById("target").value;

  fetch("/shorten", {
    method: "POST",
    headers: {
      "Authorization": authHeader,
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ handle, target })
  })
  .then(res => {
    if (res.ok) {
      loadLinks();
    } else {
      alert("Failed to create/update link.");
    }
  });
}

function loadLinks() {
  fetch("/shorten", {
    method: "GET",
    headers: { "Authorization": authHeader }
  })
  .then(res => {
    if (!res.ok) throw new Error("Unauthorized");
    return res.json();
  })
  .then(links => {
    document.getElementById("login-section").style.display = "none";
    document.getElementById("dashboard-section").style.display = "block";

    const table = document.getElementById("link-table");
    table.innerHTML = "";
    links.forEach(({ handle, target }) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td><a href="/r/${handle}" target="_blank">${handle}</a></td>
        <td>${target}</td>
        <td><button class="btn btn-sm btn-danger" onclick="deleteLink('${handle}')">Delete</button></td>
      `;
      table.appendChild(row);
    });
  })
  .catch(() => {
    alert("Invalid login or failed to load links.");
    logout();
  });
}

function deleteLink(handle) {
  fetch(`/delete/${handle}`, {
    method: "DELETE",
    headers: { "Authorization": authHeader }
  })
  .then(res => {
    if (res.ok) {
      loadLinks();
    } else {
      alert("Failed to delete link.");
    }
  });
}

window.onload = () => {
  const saved = localStorage.getItem("auth");
  if (saved) {
    authHeader = saved;
    loadLinks();
  }
};
