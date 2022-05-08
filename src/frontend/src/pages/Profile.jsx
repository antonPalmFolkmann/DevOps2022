import React, {useState, useContext} from 'react';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Stack from '@mui/material/Stack';
import Avatar from '@mui/material/Avatar';
import { Container } from '@mui/material';
import TextField from '@mui/material/TextField';
import userContext from '../utils/userContext';
import FollowButton from './FollowButton'
import { register } from '../api/register'
import { login, loggedInUser } from '../api/login'
import { postMessage } from '../api/messages'


export default function Profile() {
    const user = useContext(userContext)
    const [errorMessage, setErrorMessage] = useState('') 
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [repeatPassword, setRepeatPassword] = useState('')
    const [email, setEmail] = useState('')
    const [message, setMessage] = useState('')
    const [registerClicked, setRegisterClicked] = useState(false)
    return (
        <Container  style={{
            display: 'flex',
            alignItems: 'left',
            justifyContent: 'center',
            paddingTop: 50,
        }}>
            {(() => {
            if (user.username === '' ) {
                if (!registerClicked) {
                // LOGIN //
                    return (
                        <Stack direction="column" spacing={2} style={{
                        display: 'flex',
                        alignItems: 'center',
                        }}>
                            <Typography variant="h4" component="div" sx={{ color: '#a3c9fe' }}>
                            Welcome to minitwit!
                            </Typography>
                            <Typography variant="h5" component="div">
                            Log in 
                            </Typography>
                            <TextField id="outlined-basic" placeholder="Username" variant="outlined" 
                                onChange={(event) => setUsername(event.target.value)}
                            />
                            <TextField id="outlined-basic" type="password" placeholder="Password" variant="outlined" 
                                onChange={(event) => setPassword(event.target.value)}
                             />
                            <Stack direction="row" spacing={2} style={{
                            display: 'flex',
                            alignItems: 'center',
                            }}>
                                <Button variant="outlined" sx={{ borderColor: '#a3c9fe' }} 
                                onClick={() => { 
                                    login(username, password)
                                    .then((response) => {
                                        if (loggedInUser.username !== '') {
                                            user.setUser(loggedInUser)
                                        } else {
                                            setErrorMessage('Error logging in. Please try again.')
                                        }
                                    })
                                }}>
                                Log in
                                </Button>
                                <Button variant="outlined" sx={{ borderColor: '#a3c9fe' }} 
                                onClick={() => { setRegisterClicked(true)}}>
                                Register
                                </Button>
                            </Stack>
                            <Typography variant="p" style={{color: 'red'}}>
                                {errorMessage}
                            </Typography>
                        </Stack>
                    )
                } else {
                    // REGISTER //
                    return (
                        <Stack direction="column" spacing={2} style={{
                        display: 'flex',
                        alignItems: 'center',
                        }}>
                            <Typography variant="h5" component="div">
                            Register 
                            </Typography>
                            <TextField id="outlined-basic" placeholder="Enter username" variant="outlined"
                                onChange={(event) => setUsername(event.target.value)}
                             />
                            <TextField id="outlined-basic" placeholder="Enter e-mail adress" variant="outlined"
                                onChange={(event) => setEmail(event.target.value)}
                             />
                            <TextField id="outlined-basic" type="password" placeholder="Enter password" variant="outlined"
                                onChange={(event) => setPassword(event.target.value)}
                             />
                            <TextField id="outlined-basic" type="password" placeholder="Repeat password" variant="outlined" 
                                onChange={(event) => setRepeatPassword(event.target.value)}
                            />
                            <Stack direction="row" spacing={2} style={{
                            display: 'flex',
                            alignItems: 'center',
                            }}>
                                <Button variant="outlined" sx={{ borderColor: '#a3c9fe' }} 
                                onClick={() => {
                                    if (password === repeatPassword) {
                                        register(username, email, password)
                                        setRegisterClicked(false)
                                    } else {
                                        setErrorMessage('Passwords do not match. Please try again.')
                                    }
                                }}>
                                Register
                                </Button>
                                <Button variant="outlined" sx={{ borderColor: '#a3c9fe' }} 
                                onClick={() => { setRegisterClicked(false)}}>
                                Back
                                </Button>
                            </Stack>
                            <Typography variant="p" style={{color: 'red'}}>
                                {errorMessage}
                            </Typography>
                        </Stack>    
                    )
                }
            } else if (user.username !== '' && user.currentProfile.username !== '') {
                return (
                    <Stack direction="column" spacing={2} style={{
                    display: 'flex',
                    alignItems: 'center',
                    }}>
                        <Avatar
                        sx={{ width: 100, height: 100 }}>
                            {user.currentProfile.username[0]}
                        </Avatar>
                        <Typography variant="h5" component="div">
                            {user.currentProfile.username}
                        </Typography> 
                        {(() => {
                        if (user.currentProfile.username !== user.username) {
                            return (
                                <FollowButton/>
                            )
                        } else {
                            return (
                                <Stack direction="column" spacing={2} style={{
                                display: 'flex',
                                alignItems: 'center',
                                }}>
                                    <TextField
                                    id="outlined-textarea"
                                    label="Enter message..."
                                    multiline
                                    maxRows={4}
                                    sx={{ minWidth: 400, borderColor: '#a3c9fe' }}
                                    onChange={(event) => setMessage(event.target.value)}
                                    />
                                    <Button variant="outlined" value={{message}} sx={{ minWidth: 200, borderColor: '#a3c9fe' }} onClick={() => {
                                        postMessage(user.username, message)
                                        .then((response) => {
                                            setMessage('')
                                        })
                                    }}>
                                        Post message
                                    </Button>
                                </Stack>
                            )
                            }
                        })()}    
                    </Stack>
                )
            }
            })()}
        </Container>
    )
}
