import React, { createContext, useState } from 'react'
import Timeline from './pages/Timeline'
import userContext from './utils/userContext';
import { getMessages, exportMessages } from './api/messages'

export default function App() {
    const [username, setUsername] = useState('')
    const [email, setEmail] = useState('')
    const [avatar, setAvatar] = useState('')
    const [follows, setFollows] = useState([])
    const [currentProfile, setCurrentProfile] = useState({
        username: '',
        email: '',
        avatar: '',
        follows: [],
    })
    const [currentMessages, setCurrentMessages] = useState([])

    const setUser = (object) => {
        setUsername(object.username)
        setEmail(object.email)
        setAvatar(object.avatar)
        setFollows(object.follows)
     };
    const getUser = () => {
        return ({
            username: username,
            email: email,
            avatar: avatar,
            follows: follows
        })
    }
    const getCurrentProfile = () => {
        return currentProfile;
    }
    const getCurrentMessages = () => {
        return currentMessages;
    }
    return (
        <userContext.Provider value=
        {{
            username: username,
            email: email,
            avatar: avatar,
            follows: [],
            setUser: setUser, 
            getUser: getUser,
            currentProfile: currentProfile,
            setCurrentProfile: setCurrentProfile,
            getCurrentProfile: getCurrentProfile,
            currentMessages: currentMessages,
            setCurrentMessages: setCurrentMessages,
            getCurrentMessages: getCurrentMessages,
        }}>
            <Timeline />
        </userContext.Provider>
    );
}
