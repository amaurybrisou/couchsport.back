const AppMessenger = {
    install: function(Vue, options){
        const isDef = v => v !== undefined
        if(!isDef(options.store)) throw 'a vuex store is required';
        
        Vue.prototype.$messenger = {
            namespace: (options.namespace || 'conversations') + '/',
            mutations: options.mutations,
            actions: options.actions,
            store: options.store,
            setMessagesRead(conversationIDX){
                this.store.commit(this.namespace+this.mutations.MESSAGES_READ, conversationIDX)
            },
            sendMessage(m){
                return this.store.dispatch(this.namespace+this.actions.CONVERSATION_SEND_MESSAGE, m);
            }
        }
    }
}

export default AppMessenger.install;

if(typeof window !== 'undefined' && window.Vue){
    Vue.use(AppMessenger);
}
