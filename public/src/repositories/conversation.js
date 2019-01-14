import repo from './repository.js';

export default {
    sendMessage: (payload) => {
        return repo.post("/conversations/message/send", payload)
    },
    mines: () => {
        return repo.get("/profile/conversations")
    },
    delete: (payload) => {
        return repo.post("/profile/conversations/delete", payload)
    }
}