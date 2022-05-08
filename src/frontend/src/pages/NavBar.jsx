import React, {useState, useContext} from 'react';
import Paper from '@mui/material/Paper';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import FormatListBulletedIcon from '@mui/icons-material/FormatListBulleted';
import ExitToAppIcon from '@mui/icons-material/ExitToApp';
import userContext from '../utils/userContext';
import PublicIcon from '@mui/icons-material/Public';
import { logout } from '../api/logout'
import { getMessages, exportMessages, exportProfile, getMessagesByUser } from '../api/messages'

let emptyProfile = {
    username: '',
    email: '',
    avatar: '',
    follows: [],
}

export default function NavBar() {
    const user = useContext(userContext)

    return (
        <Paper sx={{ position: 'fixed', bottom: 0, left: 0, right: 0, borderColor : 'black'}} elevation={10}>
            <BottomNavigation
            showLabels style={{ 
                backgroundColor : '#a3c9fe',
                borderColor : 'black'
            }}
            >
            <BottomNavigationAction  label='My Profile' icon={<AccountCircleIcon />} onClick={() => { 
                getMessagesByUser(user.username)
                .then((response) => {
                    user.setCurrentProfile(user.getUser());
                    user.setCurrentMessages(exportMessages);
                })
            }}/>
            <BottomNavigationAction  label='My Timeline' icon={<FormatListBulletedIcon />} onClick={() => { 
                getMessages()
                .then((response) => {
                    user.setCurrentProfile(user.getUser());
                    user.setCurrentMessages(exportMessages.filter(
                            (message) => message.authorName == 'Roger Histand'
                        )
                    )
                })
            }}/>
            <BottomNavigationAction label='Public Timeline' icon={<PublicIcon />} onClick={() => { 
                user.setCurrentProfile(emptyProfile);
                getMessages()
                    .then((response) => {
                        user.setCurrentMessages(exportMessages)
                    })
            }}/>
            <BottomNavigationAction label='Logout' icon={<ExitToAppIcon />}  onClick={() => { 
                logout()
                user.setUser(emptyProfile)
                user.setCurrentProfile(emptyProfile)
                getMessages()
                    .then((response) => {
                        user.setCurrentMessages(exportMessages)
                    })
            }}/>
            </BottomNavigation>      
      </Paper>
    )
}