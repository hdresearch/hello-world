const count = document.querySelector("#count");
const button = document.querySelector("#increment");
const status = document.querySelector("#status");

async function request(method = "GET") {
  const response = await fetch("/api/visits", {
    method,
    headers: { Accept: "application/json" },
  });
  const body = await response.json().catch(() => ({}));
  if (!response.ok || typeof body.count !== "number") {
    throw new Error(body.detail || "The stack did not return a visit count.");
  }
  count.textContent = body.count.toLocaleString();
  status.textContent = "Frontend → backend → database is healthy.";
  status.classList.remove("error");
}

async function run(method) {
  button.disabled = true;
  status.textContent = method === "POST" ? "Persisting visit…" : "Connecting to the stack…";
  status.classList.remove("error");
  try {
    await request(method);
  } catch (error) {
    status.textContent = error instanceof Error ? error.message : "The stack is unavailable.";
    status.classList.add("error");
  } finally {
    button.disabled = false;
  }
}

button.addEventListener("click", () => run("POST"));
run("GET");
