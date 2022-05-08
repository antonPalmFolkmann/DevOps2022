import React from 'react';

const userContext = React.createContext({
    username: '',
    email: '',
    avatar: '',
    follows: [],
    setUser: () => {},
    getUser: () => {},
    currentProfile: {
        username: '',
        email: '',
        avatar: '',
        follows: [],
    },
    setCurrentProfile: () => {},
    getCurrentProfile: () => {},
    currentMessages: [],
    setCurrentMessages: () => {},
    getCurrentMessages: () => {},
});


export default userContext;
