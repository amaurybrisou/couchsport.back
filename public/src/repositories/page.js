import repo from './repository';

export default {
    all(){
        return repo.get('/pages');
    },
    upload(payload){
        return repo.post('/images/upload', payload, { headers: { 'Content-Type': 'multipart/form-data' } });
    },
    createOrUpdate(payload){
        return repo.post('/pages/new', payload)
    },
    publish(payload){
        return repo.post('/pages/publish', payload)
    },
    delete(payload){
        return repo.post('/pages/delete', payload)
    }
}