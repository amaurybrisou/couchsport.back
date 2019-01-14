import repo from "./repository";

export default {
  create(payload) {
    return repo.put("/signin", payload);
  },
  login(payload) {
    return repo.post("/login", payload);
  },
  logout() {
    return repo.get("/logout");
  }
};
