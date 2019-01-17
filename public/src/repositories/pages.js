import repo from "./repository";

export default {
  all() {
    return repo.get("/pages");
  },
  mines() {
    return repo.get("/profiles/pages");
  },
  get({ id, profile }) {
    return repo.get("/pages?id=" + id + (profile ? "&profile" : ""));
  },
  upload(payload) {
    return repo.post("/images/upload", payload, {
      headers: { "Content-Type": "multipart/form-data" }
    });
  },
  new(payload) {
    return repo.post("/pages/new", payload);
  },
  edit(payload) {
    return repo.post("/pages/update", payload);
  },
  publish(payload) {
    return repo.post("/pages/publish", payload);
  },
  delete(payload) {
    return repo.post("/pages/delete", payload);
  }
};
